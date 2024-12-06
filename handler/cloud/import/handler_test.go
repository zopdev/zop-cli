package _import

import (
	"errors"
	"go.uber.org/mock/gomock"

	"gofr.dev/pkg/gofr"
	"testing"
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
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "Successfully Imported!" {
		t.Fatalf("expected 'Successfully Imported!', got %v", result)
	}
}

func TestImport_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAccountImporter := NewMockAccountImporter(ctrl)
	mockAccountImporter.EXPECT().PostAccounts(gomock.Any()).Return(errors.New("import error"))

	handler := New(mockAccountImporter)
	ctx := &gofr.Context{}

	result, err := handler.Import(ctx)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
	if err.Error() != "import error" {
		t.Fatalf("expected 'import error', got %v", err.Error())
	}
}
