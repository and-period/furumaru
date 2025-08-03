package mailer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridError - エラー構造体
type SendGridError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
	Help    string `json:"help"`
}

func (e SendGridError) Error() string {
	return fmt.Sprintf("mailer: failed to SendGrid. code=%d, %s=%s", e.Code, e.Field, e.Message)
}

func (c *client) SendFromInfo(
	ctx context.Context, emailID, toName, toAddress string, substitutions map[string]interface{},
) error {
	ps := []*Personalization{
		{
			Type:          AddressTypeTo,
			Name:          toName,
			Address:       toAddress,
			Substitutions: substitutions,
		},
	}
	return c.sendEmail(ctx, emailID, c.fromName, c.fromAddress, "", nil, ps)
}

func (c *client) MultiSend(
	ctx context.Context, emailID, fromName, fromAddress string, ps []*Personalization,
) error {
	return c.sendEmail(ctx, emailID, fromName, fromAddress, "", nil, ps)
}

func (c *client) MultiSendFromInfo(ctx context.Context, emailID string, ps []*Personalization) error {
	return c.sendEmail(ctx, emailID, c.fromName, c.fromAddress, "", nil, ps)
}

func (c *client) sendEmail(
	ctx context.Context, emailID, fromName, fromAddress, subject string, cs []*Content, ps []*Personalization,
) error {
	msg := c.newMessage(emailID, fromName, fromAddress, subject, cs, ps)
	resp, err := c.client.SendWithContext(ctx, msg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to send email", log.Error(err))
		return c.mailError(err)
	}
	if resp.StatusCode < 400 {
		return nil
	}
	var out *SendGridError
	if err = json.Unmarshal([]byte(resp.Body), out); err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal response", slog.String("body", resp.Body), log.Error(err))
		return c.mailError(err)
	}
	return c.mailError(err)
}

func (c *client) newMessage(
	emailID, fromName, fromAddress, subject string, cs []*Content, ps []*Personalization,
) *mail.SGMailV3 {
	now := c.now()
	from := &mail.Email{
		Name:    fromName,
		Address: fromAddress,
	}
	msg := mail.NewV3Mail()
	msg.Subject = subject
	msg.SetTemplateID(c.templateMap[emailID])
	msg.SetFrom(from)
	msg.SetSendAt(int(now.Unix()))
	msg.AddContent(c.newContents(cs)...)
	msg.AddPersonalizations(c.newPersonalizations(ps)...)
	msg.AddCategories(emailID)
	return msg
}

func (c *client) newContents(contents []*Content) []*mail.Content {
	cs := make([]*mail.Content, len(contents))
	for i := range contents {
		c := &mail.Content{
			Type:  contents[i].ContentType,
			Value: contents[i].Value,
		}
		cs[i] = c
	}
	return cs
}

func (c *client) newPersonalizations(personalizations []*Personalization) []*mail.Personalization {
	ps := make([]*mail.Personalization, len(personalizations))
	for i := range personalizations {
		email := &mail.Email{
			Name:    personalizations[i].Name,
			Address: personalizations[i].Address,
		}
		p := mail.NewPersonalization()
		switch personalizations[i].Type {
		case AddressTypeTo:
			p.AddTos(email)
		case AddressTypeCC:
			p.AddCCs(email)
		case AddressTypeBCC:
			p.AddBCCs(email)
		}
		if personalizations[i].Substitutions != nil {
			p.DynamicTemplateData = personalizations[i].Substitutions
		}
		ps[i] = p
	}
	return ps
}
