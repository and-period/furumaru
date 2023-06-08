package handler

import (
	"context"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
)

func (h *handler) contactRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("/:contactId")
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
