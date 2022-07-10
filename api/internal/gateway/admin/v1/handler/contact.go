package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/gin-gonic/gin"
)

func (h *handler) contactRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListContacts)
	arg.GET("/:contactId", h.GetContact)
	arg.PATCH("/:contactId", h.UpdateContact)
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
	contacts, err := h.messenger.ListContacts(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ContactsResponse{
		Contacts: service.NewContacts(contacts).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetContact(ctx *gin.Context) {
	in := &messenger.GetContactInput{
		ContactID: util.GetParam(ctx, "contactId"),
	}
	contact, err := h.messenger.GetContact(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.ContactResponse{
		Contact: service.NewContact(contact).Response(),
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
		ContactID: util.GetParam(ctx, "contactId"),
		Status:    service.ContactStatus(req.Status).MessengerEntity(),
		Priority:  service.ContactPriority(req.Priority).MessengerEntity(),
		Note:      req.Note,
	}
	if err := h.messenger.UpdateContact(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
