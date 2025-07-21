package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) contactRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/contacts", h.authentication)

	r.GET("", h.ListContacts)
	r.GET("/:contactId", h.GetContact)
	r.POST("", h.CreateContact)
	r.PATCH("/:contactId", h.UpdateContact)
	r.DELETE("/:contactId", h.DeleteContact)
}

func (h *handler) ListContacts(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	in := &messenger.ListContactsInput{
		Limit:  limit,
		Offset: offset,
	}

	contacts, total, err := h.messenger.ListContacts(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	contactIDs := contacts.IDs()
	threads := make([]*entity.Thread, 0, len(contactIDs))
	for _, contact := range contacts {
		in := &messenger.ListThreadsInput{
			ContactID: contact.ID,
		}
		thread, _, err := h.messenger.ListThreads(ctx, in)
		if err != nil {
			h.httpError(ctx, err)
			return
		}
		threads = append(threads, thread...)
	}
	var (
		contactCategories service.ContactCategories
		users             service.Users
		responders        service.Admins
	)

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		users, err = h.multiGetUsers(ectx, contacts.UserIDs())
		return
	})
	eg.Go(func() (err error) {
		responders, err = h.multiGetAdmins(ectx, contacts.ResponderIDs())
		return
	})
	eg.Go(func() (err error) {
		contactCategories, err = h.multiGetContactCategories(ctx, contacts.CategoryIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ContactsResponse{
		Contacts:   service.NewContacts(contacts).Response(),
		Threads:    service.NewThreads(threads).Response(),
		Categories: contactCategories.Response(),
		Users:      users.Response(),
		Responders: responders.Response(),
		Total:      total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateContact(ctx *gin.Context) {
	req := &request.CreateContactRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	var (
		category  *service.ContactCategory
		user      *service.User
		responder *service.Admin
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if req.UserID == "" {
			return nil
		}
		user, err = h.getUser(ectx, req.UserID)
		return err
	})
	eg.Go(func() (err error) {
		if req.ResponderID == "" {
			return nil
		}
		responder, err = h.getAdmin(ectx, req.ResponderID)
		return err
	})
	eg.Go(func() (err error) {
		category, err = h.getContactCategory(ectx, req.CategoryID)
		return err
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		h.badRequest(ctx, err)
		return
	}
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	in := &messenger.CreateContactInput{
		Title:       req.Title,
		Content:     req.Content,
		CategoryID:  req.CategoryID,
		Username:    req.Username,
		UserID:      req.UserID,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		ResponderID: req.ResponderID,
		Note:        req.Note,
	}
	scontact, err := h.messenger.CreateContact(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	threadIn := &messenger.CreateThreadInput{
		ContactID: scontact.ID,
		UserID:    req.UserID,
		UserType:  2,
		Content:   req.Content,
	}
	sthread, err := h.messenger.CreateThread(ctx, threadIn)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.ContactResponse{
		Contact:   service.NewContact(scontact).Response(),
		Category:  category.Response(),
		Threads:   service.NewThreads(entity.Threads{sthread}).Response(),
		User:      user.Response(),
		Responder: responder.Response(),
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetContact(ctx *gin.Context) {
	contactID := util.GetParam(ctx, "contactId")
	contact, err := h.getContact(ctx, contactID)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		category  *service.ContactCategory
		threads   service.Threads
		sender    *service.User
		responder *service.Admin
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		threads, _, err = h.getContactDetailsByContactID(ectx, contact.ID)
		return
	})
	eg.Go(func() (err error) {
		category, err = h.getContactCategory(ectx, contact.CategoryID)
		return
	})
	eg.Go(func() (err error) {
		if contact.UserID == "" {
			return nil
		}
		sender, err = h.getUser(ectx, contact.UserID)
		return err
	})
	eg.Go(func() (err error) {
		if contact.ResponderID == "" {
			return nil
		}
		responder, err = h.getAdmin(ectx, contact.ResponderID)
		return err
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	res := response.ContactResponse{
		Contact:   contact.Response(),
		Category:  category.Response(),
		Threads:   threads.Response(),
		User:      sender.Response(),
		Responder: responder.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateContact(ctx *gin.Context) {
	req := &request.UpdateContactRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.UpdateContactInput{
		ContactID:   util.GetParam(ctx, "contactId"),
		Title:       req.Title,
		Content:     req.Content,
		Username:    req.Username,
		UserID:      req.UserID,
		CategoryID:  req.CategoryID,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Status:      entity.ContactStatus(req.Status),
		ResponderID: req.ResponderID,
		Note:        req.Note,
	}

	if err := h.messenger.UpdateContact(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteContact(ctx *gin.Context) {
	in := &messenger.DeleteContactInput{
		ContactID: util.GetParam(ctx, "contactId"),
	}
	if err := h.messenger.DeleteContact(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) getContact(ctx context.Context, contactID string) (*service.Contact, error) {
	in := &messenger.GetContactInput{
		ContactID: contactID,
	}
	scontact, err := h.messenger.GetContact(ctx, in)
	if err != nil {
		return nil, err
	}
	contact := service.NewContact(scontact)
	return contact, nil
}

func (h *handler) getContactDetailsByContactID(
	ctx context.Context,
	contactID string,
) (service.Threads, int64, error) {
	in := &messenger.ListThreadsInput{
		ContactID: contactID,
	}
	sthreads, total, err := h.messenger.ListThreads(ctx, in)
	if err != nil {
		return nil, 0, err
	}
	threads := service.NewThreads(sthreads)
	return threads, total, nil
}
