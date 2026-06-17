package admin

import (
	"path/filepath"
	"testing"

	"github.com/v03413/bepusdt/app/model"
)

func TestValidateConfSetsRejectsInvalidPaymentTemplate(t *testing.T) {
	err := validateConfSets([]model.Conf{
		{K: model.PaymentTemplate, V: "modern"},
	})
	if err == nil {
		t.Fatal("expected invalid payment template to be rejected")
	}
}

func TestValidateConfSetsAllowsKnownPaymentTemplates(t *testing.T) {
	for _, mode := range []string{"official", "wolf", "custom"} {
		t.Run(mode, func(t *testing.T) {
			err := validateConfSets([]model.Conf{
				{K: model.PaymentTemplate, V: mode},
			})
			if err != nil {
				t.Fatalf("validateConfSets(%q) returned error: %v", mode, err)
			}
		})
	}
}

func TestValidateConfSetsRejectsInvalidPaymentTemplateLanguage(t *testing.T) {
	err := validateConfSets([]model.Conf{
		{K: model.PaymentTemplateLanguage, V: "de"},
	})
	if err == nil {
		t.Fatal("expected invalid payment template language to be rejected")
	}
}

func TestValidateConfSetsAllowsKnownPaymentTemplateLanguages(t *testing.T) {
	for _, language := range []string{"auto", "zh", "zh-Hant", "en", "ru", "vi", "tr", "ja", "ko"} {
		t.Run(language, func(t *testing.T) {
			err := validateConfSets([]model.Conf{
				{K: model.PaymentTemplateLanguage, V: language},
			})
			if err != nil {
				t.Fatalf("validateConfSets(%q) returned error: %v", language, err)
			}
		})
	}
}

func TestValidateConfSetsKeepsExistingGuards(t *testing.T) {
	err := validateConfSets([]model.Conf{
		{K: model.ApiAuthToken, V: "manual-token"},
	})
	if err == nil {
		t.Fatal("expected api token changes to be rejected")
	}

	missing := filepath.Join(t.TempDir(), "missing-payment-assets")
	err = validateConfSets([]model.Conf{
		{K: model.PaymentStaticPath, V: missing},
	})
	if err == nil {
		t.Fatal("expected missing static path to be rejected")
	}
}
