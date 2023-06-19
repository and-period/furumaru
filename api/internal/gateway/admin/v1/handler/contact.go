package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/gin-gonic/gin"
)

func (h *handler) contactRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("/:contactId", h.GetContact)
}

func (h *handler) CreateContact(ctx *gin.Context) {
	req := &request.CreateContactRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	if req.UserID != "" {
		in := &user.GetUserInput{
			UserID: req.UserID,
		}
		if _, err := h.user.GetUser(ctx, in); err != nil {
			badRequest(ctx, err)
			return
		}
	}
	if req.ResponderID != "" {
		in := &user.GetAdminInput{
			AdminID: req.ResponderID,
		}
		if _, err := h.user.GetAdmin(ctx, in); err != nil {
			badRequest(ctx, err)
			return
		}
	}
	if _, err := h.getContactCategory(ctx, req.CategoryID); err != nil {
		badRequest(ctx, err)
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
	thread := service.NewThread(sthread)
	threads := make([]*response.Thread, 1)
	threads[0] = thread.Response()
	contact := service.NewContact(scontact)
	res := &response.ContactResponse{
		Contact: contact.Response(),
		Threads: threads,
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
