package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	mock_stripe "github.com/and-period/furumaru/api/mock/pkg/stripe"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	signaturemock = "signature"
	errmock       = errors.New("some error")
)

type mocks struct {
	receiver *mock_stripe.MockReceiver
}

type testResponse struct {
	code int
	body interface{}
}

type testOptions struct {
	now func() time.Time
}

type testOption func(opt *testOptions)

func withNow(now time.Time) testOption {
	return func(opts *testOptions) {
		opts.now = func() time.Time {
			return now
		}
	}
}

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		receiver: mock_stripe.NewMockReceiver(ctrl),
	}
}

func newHandler(mocks *mocks, opts *testOptions) Handler {
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Receiver:  mocks.receiver,
	}
	handler := NewHandler(params).(*handler)
	handler.now = opts.now
	return handler
}

func newRoutes(h Handler, r *gin.Engine) {
	h.Routes(r.Group(""))
}

func testPost(
	t *testing.T,
	setup func(*testing.T, *mocks, *gomock.Controller),
	expect *testResponse,
	path string,
	body interface{},
	opts ...testOption,
) {
	testHTTP(t, setup, expect, newHTTPRequest(t, http.MethodPost, path, body), opts...)
}

/**
 * testHTTP - HTTPハンドラのテストを実行
 */
func testHTTP(
	t *testing.T,
	setup func(*testing.T, *mocks, *gomock.Controller),
	expect *testResponse,
	req *http.Request,
	opts ...testOption,
) {
	t.Parallel()

	// setup
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newMocks(ctrl)
	dopts := &testOptions{
		now: jst.Now,
	}
	for i := range opts {
		opts[i](dopts)
	}
	h := newHandler(mocks, dopts)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	newRoutes(h, r)
	setup(t, mocks, ctrl)

	// test
	r.ServeHTTP(w, req)
	require.Equal(t, expect.code, w.Code)
	if isError(w) || expect.body == nil {
		return
	}

	body, err := json.Marshal(expect.body)
	require.NoError(t, err, err)
	require.JSONEq(t, string(body), w.Body.String())
}

func isError(res *httptest.ResponseRecorder) bool {
	return res.Code < 200 || 300 <= res.Code
}

/**
 * newHTTPRequest - HTTP Request(application/json)を生成
 */
func newHTTPRequest(t *testing.T, method, path string, body interface{}) *http.Request {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var buf []byte
	if body != nil {
		var err error
		buf, err = json.Marshal(body)
		require.NoError(t, err, err)
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(buf))
	require.NoError(t, err, err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Stripe-Signature", signaturemock)
	return req
}

func TestHandler(t *testing.T) {
	t.Parallel()
	h := NewHandler(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, h)
}
