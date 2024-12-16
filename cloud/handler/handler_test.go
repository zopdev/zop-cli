package handler

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/logging"

	"zop.dev/cli/zop/cloud/service/list"
)

var (
	errTest                  = errors.New("import error")
	errFailedToFetchAccounts = errors.New("failed to fetch accounts")
)

func TestImport_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountImporter := NewMockAccountImporter(ctrl)
	mockAccountImporter.EXPECT().PostAccounts(gomock.Any()).Return(nil)

	handler := New(mockAccountImporter, nil)
	ctx := &gofr.Context{}
	result, err := handler.Import(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if result != successMessage {
		t.Errorf("expected 'Successfully Imported!', got %v", result)
	}
}

func TestImport_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountImporter := NewMockAccountImporter(ctrl)
	mockAccountImporter.EXPECT().PostAccounts(gomock.Any()).Return(errTest)

	ctx := &gofr.Context{Container: &container.Container{Logger: logging.NewMockLogger(logging.INFO)}}
	handler := New(mockAccountImporter, nil)

	result, err := handler.Import(ctx)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if !errors.Is(err, errTest) {
		t.Errorf("expected 'import error', got %v", err.Error())
	}
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccGetter := NewMockAccountGetter(ctrl)

	tests := []struct {
		name         string
		accounts     []*list.CloudAccountResponse
		expectedResp string
		expectedErr  error
		mocks        []*gomock.Call
	}{
		{
			name:         "No accounts",
			accounts:     []*list.CloudAccountResponse{},
			expectedResp: "No accounts found",
			expectedErr:  nil,
			mocks: []*gomock.Call{
				mockAccGetter.EXPECT().GetAccounts(gomock.Any()).Return([]*list.CloudAccountResponse{}, nil),
			},
		},
		{
			name: "Single account",
			accounts: []*list.CloudAccountResponse{
				{Name: "Account1", Provider: "AWS", ProviderID: "12345", UpdatedAt: "2024-01-01", CreatedAt: "2023-01-01"},
			},
			expectedResp: "Name                 Provider             ProviderID           UpdateAt             CreatedAt           \n" +
				"Account1             AWS                  12345                2024-01-01           2023-01-01          \n",
			expectedErr: nil,
			mocks: []*gomock.Call{
				mockAccGetter.EXPECT().GetAccounts(gomock.Any()).Return([]*list.CloudAccountResponse{
					{Name: "Account1", Provider: "AWS", ProviderID: "12345", UpdatedAt: "2024-01-01", CreatedAt: "2023-01-01"},
				}, nil),
			},
		},
		{
			name: "Multiple accounts with truncation",
			accounts: []*list.CloudAccountResponse{
				{Name: "ThisIsAVeryLongAccountNameThatShouldBeTruncated",
					Provider: "GCP", ProviderID: "67890", UpdatedAt: "2024-02-02", CreatedAt: "2023-02-02"},
				{Name: "Account2", Provider: "Azure",
					ProviderID: "11111", UpdatedAt: "2024-03-03", CreatedAt: "2023-03-03"},
			},
			expectedResp: "Name                 Provider             ProviderID           UpdateAt             CreatedAt           \n" +
				"ThisIsAVeryLongAc... GCP                  67890                2024-02-02           2023-02-02          \n" +
				"Account2             Azure                11111                2024-03-03           2023-03-03          \n",
			expectedErr: nil,
			mocks: []*gomock.Call{
				mockAccGetter.EXPECT().GetAccounts(gomock.Any()).Return([]*list.CloudAccountResponse{
					{Name: "ThisIsAVeryLongAccountNameThatShouldBeTruncated",
						Provider: "GCP", ProviderID: "67890", UpdatedAt: "2024-02-02", CreatedAt: "2023-02-02"},
					{Name: "Account2", Provider: "Azure", ProviderID: "11111", UpdatedAt: "2024-03-03", CreatedAt: "2023-03-03"},
				}, nil),
			},
		},
		{
			name:         "Error from GetAccounts",
			accounts:     nil,
			expectedResp: "",
			expectedErr:  errFailedToFetchAccounts,
			mocks: []*gomock.Call{
				mockAccGetter.EXPECT().GetAccounts(gomock.Any()).Return(nil, errFailedToFetchAccounts),
			},
		},
	}

	handler := New(nil, mockAccGetter)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &gofr.Context{}

			resp, err := handler.List(ctx)

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
