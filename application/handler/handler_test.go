package handler

import (
	"errors"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr/cmd"

	"github.com/stretchr/testify/require"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"testing"
)

func TestHandler_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppAdder := NewMockApplicationAdder(ctrl)

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
				mockAppAdder.EXPECT().AddApplication(gomock.Any(), "test-app").Return(errors.New("internal error")),
			},
			expected: nil,
			expErr:   errors.New("internal error"),
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
