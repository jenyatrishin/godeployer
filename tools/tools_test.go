package tools

import (
	"reflect"
	"testing"
)

func TestUserError (t *testing.T) {
	testError := UserError("Test error message")
	testErrorType := reflect.TypeOf(testError).Kind()
	if testErrorType != reflect.Ptr {
		t.Error("Expected receive Ptr data type")
	}
}
