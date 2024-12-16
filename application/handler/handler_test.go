package handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/cmd"
	"gofr.dev/pkg/gofr/cmd/terminal"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/testutil"

	svc "zop.dev/cli/zop/application/service"
)

var errAPICall = errors.New("error in API call")

func TestHandler_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppAdder := NewMockApplicationService(ctrl)

	testCases := []struct {
		name      string
		appName   string
		mockCalls []*gomock.Call
		expected  any
		expErr    error
	}{
		{
			name:    "success",
			appName: "test-app",
			mockCalls: []*gomock.Call{
				mockAppAdder.EXPECT().AddApplication(gomock.Any(), "test-app").Return(nil),
			},
			expected: "Application test-app added successfully!",
			expErr:   nil,
		},
		{
			name:     "missing name parameter",
			appName:  "",
			expected: nil,
			expErr:   ErrorApplicationNameNotProvided,
		},
		{
			name:    "error adding application",
			appName: "test-app",
			mockCalls: []*gomock.Call{
				mockAppAdder.EXPECT().AddApplication(gomock.Any(), "test-app").Return(errAPICall),
			},
			expected: nil,
			expErr:   errAPICall,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCont, _ := container.NewMockContainer(t)
			ctx := &gofr.Context{
				Container: mockCont,
				Request:   cmd.NewRequest([]string{"", "-name=" + tc.appName}),
			}

			h := New(mockAppAdder)
			res, err := h.Add(ctx)

			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expected, res)
		})
	}
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockApplicationService(ctrl)

	testCases := []struct {
		name      string
		mockCalls []*gomock.Call
		expected  string
		expErr    error
	}{
		{
			name: "success",
			mockCalls: []*gomock.Call{
				mockSvc.EXPECT().GetApplications(gomock.Any()).
					Return([]svc.Application{
						{ID: 1, Name: "app1",
							Envs: []svc.Environment{{Name: "env1", Order: 1}, {Name: "env2", Order: 2}}},
						{ID: 2, Name: "app2",
							Envs: []svc.Environment{{Name: "dev", Order: 1}, {Name: "prod", Order: 2}}},
					}, nil),
			},
			expected: "Applications and their environments:\n\n1.\x1b[38;5;6m app1 " +
				"\n\t\x1b[0m\x1b[38;5;2menv1 > env2 \n\x1b[0m2.\x1b[38;5;6m app2 " +
				"\n\t\x1b[0m\x1b[38;5;2mdev > prod \n\x1b[0m",
		},
		{
			name: "failure",
			mockCalls: []*gomock.Call{
				mockSvc.EXPECT().GetApplications(gomock.Any()).
					Return(nil, errAPICall),
			},
			expErr:   errAPICall,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := New(mockSvc)
			out := testutil.StdoutOutputForFunc(func() {
				ctx := &gofr.Context{
					Request: cmd.NewRequest([]string{""}),
					Out:     terminal.New(),
				}

				_, err := h.List(ctx)

				require.Equal(t, tc.expErr, err)
			})

			require.Equal(t, tc.expected, out)
		})
	}
}
