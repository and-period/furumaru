//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package geolocation

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/text/width"
	"googlemaps.github.io/maps"
)

var ErrNotFound = errors.New("geocoding: not found")

type Client interface {
	// GetGeolocation - 住所から緯度経度を取得
	GetGeolocation(ctx context.Context, in *GetGeolocationInput) (*GetGeolocationOutput, error)
}

type client struct {
	client *maps.Client
	logger *zap.Logger
}

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type Params struct {
	APIKey string
}

func NewClient(params *Params, opts ...Option) (Client, error) {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	c, err := maps.NewClient(maps.WithAPIKey(params.APIKey))
	if err != nil {
		return nil, err
	}
	return &client{
		client: c,
		logger: dopts.logger,
	}, nil
}

type Address struct {
	PostalCode   string
	Prefecture   string
	City         string
	AddressLine1 string
	AddressLine2 string
}

type GetGeolocationInput struct {
	*Address
}

func (a *Address) String() string {
	if a == nil {
		return ""
	}
	// 空白の削除 & 半角文字を全角文字に変換
	if a.AddressLine1 != "" {
		a.AddressLine1 = width.Narrow.String(strings.Replace(a.AddressLine1, " ", "", -1))
	}
	if a.AddressLine2 != "" {
		a.AddressLine2 = width.Narrow.String(strings.Replace(a.AddressLine2, " ", "", -1))
	}

	var builder strings.Builder
	if len(a.PostalCode) == 7 {
		builder.WriteString(a.PostalCode[:3])
		builder.WriteString("-")
		builder.WriteString(a.PostalCode[3:])
	}
	if builder.Len() > 0 {
		builder.WriteString(" ")
	}
	builder.WriteString(a.Prefecture)
	builder.WriteString(a.City)
	builder.WriteString(a.AddressLine1)
	builder.WriteString(a.AddressLine2)

	return builder.String()
}

type GetGeolocationOutput struct {
	Latitude  float64
	Longitude float64
}

func (c *client) GetGeolocation(ctx context.Context, in *GetGeolocationInput) (*GetGeolocationOutput, error) {
	req := &maps.GeocodingRequest{
		Language: "ja",
		Region:   "JP",
		Address:  in.Address.String(),
	}
	res, err := c.client.Geocode(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	out := &GetGeolocationOutput{
		Latitude:  res[0].Geometry.Location.Lat,
		Longitude: res[0].Geometry.Location.Lng,
	}
	return out, nil
}
