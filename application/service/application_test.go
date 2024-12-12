package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd/terminal"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/service"
)

var errAPICall = errors.New("error in API call")

func Test_AddApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCont, mocks := container.NewMockContainer(t, func(_ *container.Container, ctrl *gomock.Controller) any {
		return service.NewMockHTTP(ctrl)
	})
	mockCont.Services["api-service"] = mocks.HTTPService
	ctx := &gofr.Context{Container: mockCont, Out: terminal.New()}

	b, err := json.Marshal(MockErrorResponse{Error: "Something went wrong"})
	if err != nil {
		t.Fatalf("Failed to marshal test response body: %v", err)
	}

	testCases := []struct {
		name      string
		mockCalls []*gomock.Call
		expError  error
	}{
		{
			name: "success Post call",
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().PostWithHeaders(ctx, "application", nil, gomock.Any(), nil).
					Return(&http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(&errorReader{})}, nil),
			},
			expError: nil,
		},
		{
			name: "error in Post call",
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().PostWithHeaders(ctx, "application", nil, gomock.Any(), nil).
					Return(nil, errAPICall),
			},
			expError: errAPICall,
		},
		{
			name: "unexpected response",
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().PostWithHeaders(ctx, "application", nil, gomock.Any(), nil).
					Return(&http.Response{StatusCode: http.StatusInternalServerError, Body: io.NopCloser(bytes.NewBuffer(b))}, nil),
			},
			expError: &ErrAPIService{StatusCode: http.StatusInternalServerError, Message: "Something went wrong"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := New()

			errSvc := s.AddApplication(ctx, "test")

			require.Equal(t, tt.expError, errSvc)
		})
	}
}

func Test_AddApplication_WithEnvs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCont, mocks := container.NewMockContainer(t, func(_ *container.Container, ctrl *gomock.Controller) any {
		return service.NewMockHTTP(ctrl)
	})

	mockCont.Services["api-service"] = mocks.HTTPService
	ctx := &gofr.Context{Container: mockCont, Out: terminal.New()}

	testCases := []struct {
		name         string
		mockCalls    []*gomock.Call
		userInput    string
		expectedEnvs []Environment
		expError     error
	}{
		{
			name:      "success with environments",
			userInput: "y\nprod\ny\ndev\nn\n",
			expectedEnvs: []Environment{
				{Name: "prod", Order: 1},
				{Name: "dev", Order: 2},
			},
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().PostWithHeaders(ctx, "application", nil, gomock.Any(), nil).
					DoAndReturn(func(_ *gofr.Context, _ string, _, body, _ interface{}) (*http.Response, error) {
						var app Application
						_ = json.Unmarshal(body.([]byte), &app)
						require.Equal(t, "test", app.Name)
						require.Equal(t, []Environment{
							{Name: "prod", Order: 1},
							{Name: "dev", Order: 2},
						}, app.Envs)
						return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(bytes.NewBuffer(nil))}, nil
					}),
			},
			expError: nil,
		},
		{
			name:         "no environments added",
			userInput:    "n\n",
			expectedEnvs: []Environment{},
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().PostWithHeaders(ctx, "application", nil, gomock.Any(), nil).
					DoAndReturn(func(_ *gofr.Context, _ string, _, body, _ interface{}) (*http.Response, error) {
						var app Application
						_ = json.Unmarshal(body.([]byte), &app)
						require.Equal(t, "test", app.Name)
						require.Empty(t, app.Envs)
						return &http.Response{StatusCode: http.StatusCreated, Body: io.NopCloser(bytes.NewBuffer(nil))}, nil
					}),
			},
			expError: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := New()

			// Mock user input
			r, w, _ := os.Pipe()
			_, _ = w.WriteString(tt.userInput)

			oldStdin := os.Stdin
			os.Stdin = r

			defer func() { os.Stdin = oldStdin }()

			errSvc := s.AddApplication(ctx, "test")
			require.Equal(t, tt.expError, errSvc)
		})
	}
}
