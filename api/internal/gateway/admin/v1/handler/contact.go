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
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) contactRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListContacts)
	arg.GET("/:contactId", h.GetContact)
	arg.POST("", h.CreateContact)
	arg.PATCH("/:contactId", h.UpdateContact)
	arg.DELETE("/:contactId", h.DeleteContact)
}

func (h *handler) ListContacts(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	in := &messenger.ListContactsInput{
		Limit:  limit,
		Offset: offset,
	}

	contacts, total, err := h.messenger.ListContacts(ctx, in)
	if err != nil {
		httpError(ctx, err)
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
			httpError(ctx, err)
			return
		}
		threads = append(threads, thread...)
	}
	var (
		contactCategories service.ContactCategories
		users             uentity.Users
		responders        service.Admins
	)

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		usersIn := &user.MultiGetUsersInput{
			UserIDs: contacts.UserIDs(),
		}
		users, err = h.user.MultiGetUsers(ectx, usersIn)
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
		httpError(ctx, err)
		return
	}
	res := &response.ContactsResponse{
		Contacts:   service.NewContacts(contacts).Response(),
		Threads:    service.NewThreads(threads).Response(),
		Categories: contactCategories.Response(),
		Users:      service.NewUsers(users).Response(),
		Responders: responders.Response(),
		Total:      total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateContact(ctx *gin.Context) {
	req := &request.CreateContactRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	var (
		category  *service.ContactCategory
		uuser     *uentity.User
		responder *uentity.Admin
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		if req.UserID == "" {
			return nil
		}
		in := &user.GetUserInput{
			UserID: req.UserID,
		}
		uuser, err = h.user.GetUser(ectx, in)
		return err
	})
	eg.Go(func() (err error) {
		if req.ResponderID == "" {
			return nil
		}
		in := &user.GetAdminInput{
			AdminID: req.ResponderID,
		}
		responder, err = h.user.GetAdmin(ectx, in)
		return err
	})
	eg.Go(func() (err error) {
		category, err = h.getContactCategory(ectx, req.CategoryID)
		return err
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		badRequest(ctx, err)
		return
	}
	if err != nil {
		httpError(ctx, err)
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
		httpError(ctx, err)
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
		httpError(ctx, err)
		return
	}
	res := &response.ContactResponse{
		Contact:   service.NewContact(scontact).Response(),
		Category:  category.Response(),
		Threads:   service.NewThreads(entity.Threads{sthread}).Response(),
		User:      service.NewUser(uuser).Response(),
		Responder: service.NewAdmin(responder).Response(),
	}

	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetContact(ctx *gin.Context) {
	contactID := util.GetParam(ctx, "contactId")
	contact, err := h.getContact(ctx, contactID)
	if err != nil {
		httpError(ctx, err)
		return
	}

	var (
		category  *service.ContactCategory
		threads   service.Threads
		sender    *uentity.User
		responder *uentity.Admin
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
		in := &user.GetUserInput{
			UserID: contact.UserID,
		}
		sender, err = h.user.GetUser(ectx, in)
		return err
	})
	eg.Go(func() (err error) {
		if contact.ResponderID == "" {
			return nil
		}
		in := &user.GetAdminInput{
			AdminID: contact.ResponderID,
		}
		responder, err = h.user.GetAdmin(ectx, in)
		return err
	})
	if err := eg.Wait(); err != nil {
		httpError(ctx, err)
		return
	}

	res := response.ContactResponse{
		Contact:   contact.Response(),
		Category:  category.Response(),
		Threads:   threads.Response(),
		User:      service.NewUser(sender).Response(),
		Responder: service.NewAdmin(responder).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateContact(ctx *gin.Context) {
	req := &request.UpdateContactRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
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
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteContact(ctx *gin.Context) {
	in := &messenger.DeleteContactInput{
		ContactID: util.GetParam(ctx, "contactId"),
	}
	if err := h.messenger.DeleteContact(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
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

func (h *handler) getContactDetailsByContactID(ctx context.Context, contactID string) (service.Threads, int64, error) {
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
