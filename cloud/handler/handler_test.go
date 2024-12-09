package handler

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
)

var (
	errTest = errors.New("import error")
)

func TestImport_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountImporter := NewMockAccountImporter(ctrl)
	mockAccountImporter.EXPECT().PostAccounts(gomock.Any()).Return(nil)

	handler := New(mockAccountImporter)
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

	handler := New(mockAccountImporter)
	ctx := &gofr.Context{}

	result, err := handler.Import(ctx)

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}

	if errors.Is(err, errTest) {
		t.Errorf("expected 'import error', got %v", err.Error())
	}
}
