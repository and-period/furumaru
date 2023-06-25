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
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) contactRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("/:contactId", h.GetContact)
	arg.POST("", h.CreateContact)
	arg.PATCH("/:contactId", h.UpdateContact)
	arg.DELETE("/:contactId", h.DeleteContact)
}

func (h *handler) CreateContact(ctx *gin.Context) {
	req := &request.CreateContactRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if req.UserID == "" {
			return nil
		}
		in := &user.GetUserInput{
			UserID: req.UserID,
		}
		_, err := h.user.GetUser(ectx, in)
		return err
	})
	eg.Go(func() error {
		if req.ResponderID == "" {
			return nil
		}
		in := &user.GetAdminInput{
			AdminID: req.ResponderID,
		}
		_, err := h.user.GetAdmin(ectx, in)
		return err
	})
	eg.Go(func() error {
		_, err := h.getContactCategory(ectx, req.CategoryID)
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
		Contact: service.NewContact(scontact).Response(),
		Threads: service.NewThreads(entity.Threads{sthread}).Response(),
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
	threads, _, err := h.getContactDetailsByContactID(ctx, contact)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := response.ContactResponse{
		Contact: contact.Response(),
		Threads: threads.Response(),
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

func (h *handler) getContactDetailsByContactID(ctx context.Context, contact *service.Contact) (service.Threads, int64, error) {
	in := &messenger.ListThreadsByContactIDInput{
		ContactID: contact.ID,
	}
	sthreads, total, err := h.messenger.ListThreadsByContactID(ctx, in)
	if err != nil {
		return nil, 0, err
	}
	threads := service.NewThreads(sthreads)

	return threads, total, nil
}
