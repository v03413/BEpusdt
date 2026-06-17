package router

import (
	"html/template"
	"testing"

	"github.com/v03413/bepusdt/static"
)

func TestEmbeddedPaymentTemplatesParse(t *testing.T) {
	if _, err := template.New("payment").ParseFS(static.Payment, "payment/views/*.html"); err != nil {
		t.Fatalf("parse embedded payment templates: %v", err)
	}
}
