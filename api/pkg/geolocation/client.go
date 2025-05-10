//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package geolocation

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/text/width"
	"googlemaps.github.io/maps"
)

var (
	ErrNotFound = errors.New("geocoding: not found")

	errInvalidAddress = errors.New("geocoding: invalid address")
)

type Client interface {
	// GetGeolocation - 住所から緯度経度を取得
	GetGeolocation(ctx context.Context, in *GetGeolocationInput) (*GetGeolocationOutput, error)
	// GetAddress - 経度緯度から住所を取得
	GetAddress(ctx context.Context, in *GetAddressInput) (*GetAddressOutput, error)
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

func newAddress(components []maps.AddressComponent) (*Address, error) {
	if len(components) == 0 {
		return nil, errInvalidAddress
	}
	res := &Address{}
	for i := len(components) - 1; i >= 0; i-- {
		component := components[i]
		switch {
		case slices.Contains(component.Types, "postal_code"):
			res.PostalCode = strings.Replace(component.LongName, "-", "", -1)
		case slices.Contains(component.Types, "administrative_area_level_1"):
			res.Prefecture = component.LongName
		case slices.Contains(component.Types, "locality"):
			res.City = component.LongName
		case slices.Contains(component.Types, "sublocality"):
			res.AddressLine1 += component.LongName
		case slices.Contains(component.Types, "premise"), slices.Contains(component.Types, "subpremise"):
			if len(res.AddressLine2) > 0 {
				res.AddressLine2 += " "
			}
			res.AddressLine2 += component.LongName
		}
	}
	return res, nil
}

func (a *Address) FormatWiden() {
	a.AddressLine1 = width.Widen.String(strings.TrimSpace(a.AddressLine1))
	a.AddressLine2 = width.Widen.String(strings.TrimSpace(a.AddressLine2))
}

func (a *Address) FormatShorten() {
	a.AddressLine1 = width.Narrow.String(strings.TrimSpace(a.AddressLine1))
	a.AddressLine2 = width.Narrow.String(strings.TrimSpace(a.AddressLine2))
}

func (a *Address) String() string {
	if a == nil {
		return ""
	}
	a.FormatWiden()

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
	if a.AddressLine1 == "" || a.AddressLine2 == "" {
		return builder.String()
	}
	if isNumber(a.AddressLine1[len(a.AddressLine1)-1:]) && isNumber(a.AddressLine2[:1]) {
		builder.WriteString("−")
	}
	builder.WriteString(a.AddressLine2)

	return builder.String()
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

type GetGeolocationInput struct {
	*Address
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
	c.logger.Debug("Request geocoding by address", zap.Any("request", req))
	res, err := c.client.Geocode(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	c.logger.Debug("Received geocoding by address", zap.Any("response", res))
	out := &GetGeolocationOutput{
		Latitude:  res[0].Geometry.Location.Lat,
		Longitude: res[0].Geometry.Location.Lng,
	}
	return out, nil
}

type GetAddressInput struct {
	Latitude  float64
	Longitude float64
}

type GetAddressOutput struct {
	*Address
}

func (c *client) GetAddress(ctx context.Context, in *GetAddressInput) (*GetAddressOutput, error) {
	req := &maps.GeocodingRequest{
		Language: "ja",
		Region:   "JP",
		LatLng: &maps.LatLng{
			Lat: in.Latitude,
			Lng: in.Longitude,
		},
	}
	c.logger.Debug("Request geocoding by geolocation", zap.Any("request", req))
	res, err := c.client.Geocode(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	c.logger.Debug("Received geocoding by geolocation", zap.Any("response", res))
	address, err := newAddress(res[0].AddressComponents)
	if err != nil {
		return nil, fmt.Errorf("failed to parse address: %w", err)
	}
	out := &GetAddressOutput{
		Address: address,
	}
	return out, nil
}
