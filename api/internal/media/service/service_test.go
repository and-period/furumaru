package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/and-period/furumaru/api/internal/media/database"
	mock_database "github.com/and-period/furumaru/api/mock/media/database"
	mock_sqs "github.com/and-period/furumaru/api/mock/pkg/sqs"
	mock_storage "github.com/and-period/furumaru/api/mock/pkg/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	tmpURL     = "http://tmp.and-period.jp"
	storageURL = "http://and-period.jp"
	unknownURL = "http://example.com"
)

type mocks struct {
	db       *dbMocks
	tmp      *mock_storage.MockBucket
	storage  *mock_storage.MockBucket
	producer *mock_sqs.MockProducer
}

type dbMocks struct {
	Broadcast *mock_database.MockBroadcast
}

type testCaller func(ctx context.Context, t *testing.T, service *service)

func newMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		db:       newDBMocks(ctrl),
		tmp:      mock_storage.NewMockBucket(ctrl),
		storage:  mock_storage.NewMockBucket(ctrl),
		producer: mock_sqs.NewMockProducer(ctrl),
	}
}

func newDBMocks(ctrl *gomock.Controller) *dbMocks {
	return &dbMocks{
		Broadcast: mock_database.NewMockBroadcast(ctrl),
	}
}

func newService(mocks *mocks) *service {
	params := &Params{
		WaitGroup: &sync.WaitGroup{},
		Database: &database.Database{
			Broadcast: mocks.db.Broadcast,
		},
		Tmp:      mocks.tmp,
		Storage:  mocks.storage,
		Producer: mocks.producer,
	}
	tmpHost, _ := url.Parse(tmpURL)
	storageHost, _ := url.Parse(storageURL)
	mocks.tmp.EXPECT().GetHost().Return(tmpHost, nil)
	mocks.storage.EXPECT().GetHost().Return(storageHost, nil)
	srv, _ := NewService(params)
	return srv.(*service)
}

func testService(setup func(ctx context.Context, mocks *mocks), testFunc testCaller) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := newMocks(ctrl)

		srv := newService(mocks)
		setup(ctx, mocks)

		testFunc(ctx, t, srv)
		srv.waitGroup.Wait()
	}
}

func testImageFile(t *testing.T) (io.Reader, *multipart.FileHeader) {
	const filename, format = "and-period.png", "image"

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, format, filename))
	header.Set("Content-Type", "multipart/form-data")
	part := &multipart.FileHeader{
		Filename: filepath,
		Header:   header,
		Size:     3 << 20, // 3MB
	}

	return file, part
}

func testVideoFile(t *testing.T) (io.Reader, *multipart.FileHeader) {
	const filename, format = "and-period.mp4", "video"

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	defer writer.Close()

	filepath := testFilePath(t, filename)
	file, err := os.Open(filepath)
	require.NoError(t, err, err)

	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, format, filename))
	header.Set("Content-Type", "multipart/form-data")
	part := &multipart.FileHeader{
		Filename: filepath,
		Header:   header,
		Size:     10 << 20, // 10MB
	}

	return file, part
}

func testFilePath(t *testing.T, filename string) string {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	strs := strings.Split(dir, "api/internal")
	if len(strs) == 0 {
		t.Fatal("test: invalid file path")
	}
	return filepath.Join(strs[0], "/api/tmp", filename)
}

func TestService(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mocks := newMocks(ctrl)
	params := &Params{Storage: mocks.storage, Tmp: mocks.tmp}
	surl, err := url.Parse(storageURL)
	require.NoError(t, err)
	turl, err := url.Parse(tmpURL)
	require.NoError(t, err)
	mocks.storage.EXPECT().GetHost().Return(surl, nil)
	mocks.tmp.EXPECT().GetHost().Return(turl, nil)
	srv, err := NewService(params, WithLogger(zap.NewNop()))
	assert.NoError(t, err)
	assert.NotNil(t, srv)
}
