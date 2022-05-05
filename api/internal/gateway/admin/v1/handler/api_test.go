package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/service"
	mock_user "github.com/and-period/marche/api/mock/user/service"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var (
	idmock    = "kSByoE6FetnPs5Byk3a9Zx"
	tokenmock = "eyJraWQiOiJXOWxyODBzODRUVXQ3eWdyZ"
	errmock   = errors.New("some error")
)

type mocks struct {
	user *mock_user.MockUserService
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
		user: mock_user.NewMockUserService(ctrl),
	}
}

func newAPIV1Handler(mocks *mocks, opts *testOptions) APIV1Handler {
	return &apiV1Handler{
		now:         opts.now,
		logger:      zap.NewNop(),
		sharedGroup: &singleflight.Group{},
		waitGroup:   &sync.WaitGroup{},
		user:        mocks.user,
	}
}

func newRoutes(h APIV1Handler, r *gin.Engine) {
	h.Routes(r.Group(""))
}

func testHTTP(
	t *testing.T,
	setup func(*testing.T, *mocks, *gomock.Controller),
	expect *testResponse,
	req *http.Request,
	opts ...testOption,
) {
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
	h := newAPIV1Handler(mocks, dopts)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	newRoutes(h, r)
	setup(t, mocks, ctrl)

	// auth := &uentity.AdminAuth{AdminID: idmock, Role: uentity.AdminRoleAdministrator}
	// mocks.user.EXPECT().GetAdminAuth(gomock.Any(), gomock.Any()).Return(auth, nil).MaxTimes(1)

	// test
	r.ServeHTTP(w, req)
	require.Equal(t, expect.code, w.Code)
	if isError(w) || expect.body == nil {
		return
	}

	body, err := json.Marshal(expect.body)
	require.NoError(t, err, err)
	require.Equal(t, string(body), w.Body.String())
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
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenmock))
	req.Header.Add("adminId", idmock)
	return req
}

/**
 * newMultipartRequset - HTTP Request(multipart/form-data)を生成
 */
func newMultipartRequest(t *testing.T, method, path, field string) *http.Request {
	const filename = "calmato.png"

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	filepath := getFilepath(t)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)
	defer file.Close()

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, field, filename))
	header.Set("Content-Type", "multipart/form-data")
	part, err := writer.CreatePart(header)
	require.NoError(t, err, err)

	_, err = io.Copy(part, file)
	require.NoError(t, err, err)

	req, err := http.NewRequest(method, path, buf)
	require.NoError(t, err, err)

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenmock))
	req.Header.Add("adminId", idmock)
	return req
}

func getFilepath(t *testing.T) string {
	const filename = "and-period.png"

	dir, err := os.Getwd()
	assert.NoError(t, err)

	strs := strings.Split(dir, "api/internal")
	if len(strs) == 0 {
		t.Fatal("test: invalid file path")
	}
	return filepath.Join(strs[0], "/api/tmp", filename)
}

func TestAPIV1Handler(t *testing.T) {
	t.Parallel()
	h := NewAPIV1Handler(&Params{}, WithLogger(zap.NewNop()))
	assert.NotNil(t, h)
}

func TestSetAuth(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{Header: http.Header{}}
	setAuth(ctx, "admin-id", service.AdminRoleDeveloper)
	assert.Equal(t, "admin-id", getAdminID(ctx))
}
