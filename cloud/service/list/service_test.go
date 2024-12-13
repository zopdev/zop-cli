package list

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd/terminal"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/service"
)

var (
	errFailFetch  = errors.New("failed to fetch")
	errUnmarshall = errors.New("invalid character 'i' looking for beginning of value")
)

func Test_Service_GetAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCont, mocks := container.NewMockContainer(t, func(_ *container.Container, ctrl *gomock.Controller) any {
		return service.NewMockHTTP(ctrl)
	})
	mockCont.Services["api-service"] = mocks.HTTPService
	ctx := &gofr.Context{Container: mockCont, Out: terminal.New()}

	testCases := []struct {
		name      string
		mockCalls []*gomock.Call
		expResult []*CloudAccountResponse
		expError  error
	}{
		{
			name: "successful GET call",
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().Get(ctx, "/cloud-accounts", nil).
					Return(&http.Response{
						Body: io.NopCloser(bytes.NewBufferString(
							`{"data": [{"Name": "Account1", "Provider": "AWS", "ProviderID": "12345", "UpdatedAt": "2024-01-01", "CreatedAt": "2023-01-01"}]}`)),
					}, nil),
			},
			expResult: []*CloudAccountResponse{
				{Name: "Account1", Provider: "AWS", ProviderID: "12345", UpdatedAt: "2024-01-01", CreatedAt: "2023-01-01"},
			},
			expError: nil,
		},
		{
			name: "error during GET call",
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().Get(ctx, "/cloud-accounts", nil).
					Return(nil, errFailFetch),
			},
			expResult: nil,
			expError:  errFailFetch,
		},
		{
			name: "invalid JSON response",
			mockCalls: []*gomock.Call{
				mocks.HTTPService.EXPECT().Get(ctx, "/cloud-accounts", nil).
					Return(&http.Response{
						Body: io.NopCloser(bytes.NewBufferString("invalid-json")),
					}, nil),
			},
			expResult: nil,
			expError:  errUnmarshall,
		},
	}

	svc := New()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.GetAccounts(ctx)

			if tt.expError != nil {
				require.EqualError(t, err, tt.expError.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expResult, result)
			}
		})
	}
}
