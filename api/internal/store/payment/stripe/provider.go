package stripe

import (
	"github.com/and-period/furumaru/api/internal/store/payment"
	pkgstripe "github.com/and-period/furumaru/api/pkg/stripe"
)

type provider struct {
	client pkgstripe.Client
}

type Params struct {
	SecretKey string
}

func NewProvider(params *Params, opts ...pkgstripe.Option) payment.Provider {
	client := pkgstripe.NewClient(&pkgstripe.Params{
		SecretKey: params.SecretKey,
	}, opts...)
	return &provider{
		client: client,
	}
}
