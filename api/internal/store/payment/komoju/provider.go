package komoju

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/store/payment"
)

type provider struct {
	client      *apiClient
	host        string
	captureMode CaptureMode
}

// Params contains parameters for creating a new KOMOJU provider.
type Params struct {
	Host         string      // KOMOJU接続用URL
	ClientID     string      // KOMOJU接続時のBasic認証ユーザー名
	ClientSecret string      // KOMOJU接続時のBasic認証パスワード
	CaptureMode  CaptureMode // 売上処理方法
}

// NewProvider creates a new KOMOJU payment provider.
func NewProvider(cli *http.Client, params *Params, opts ...Option) payment.Provider {
	return &provider{
		client:      newAPIClient(cli, params.ClientID, params.ClientSecret, opts...),
		host:        params.Host,
		captureMode: params.CaptureMode,
	}
}
