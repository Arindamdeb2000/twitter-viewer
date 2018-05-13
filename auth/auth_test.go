package auth

import (
	"bytes"
	"encoding/base64"
	"log"
	"reflect"
	"testing"
)

func TestNewAuthorizer(t *testing.T) {
	logger := log.New(bytes.NewBuffer(make([]byte, 0)), `test`, log.Flags())
	expected := &AuthorizerImpl{
		ApiUrl:      `testUrl`,
		ApiKey:      `testKey`,
		ApiSecret:   `testSecret`,
		Logger:      logger,
		UrlEncoding: base64.URLEncoding,
		Ecoding:     base64.StdEncoding,
	}
	a := NewAuthorizer(`testUrl`, `testKey`, `testSecret`, logger)
	if a == nil {
		t.Errorf("NewAuthorizer error: AuthorizerImpl expected but got nil\n")
	}
	if !reflect.DeepEqual(a, expected) {
		t.Errorf(`NewAuthorizer error: %#v expected but got #%v`, expected, a)
	}
}

func TestAuthorizerImpl_Authorize(t *testing.T) {
	//TODO(h.lazar) write more tests
}

func TestEncodeApiKey(t *testing.T) {
	a := NewAuthorizer(`testUrl`, `testKey`, `testSecret`, log.New(bytes.NewBuffer(make([]byte, 0)), `test`, log.Flags()))
	token := a.encodeApiKey()
	expected := `dGVzdEtleTp0ZXN0U2VjcmV0`
	if token != expected {
		t.Errorf("encodeApiKey wrong result! %v expected but got %#\n", expected, token)
	}
}
