package model

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func newConfTestDB(t *testing.T) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "conf-test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	if err := db.AutoMigrate(&Conf{}); err != nil {
		t.Fatalf("migrate conf: %v", err)
	}

	Db = db
	confCache = sync.Map{}

	t.Cleanup(func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	})
}

func upsertConfForTest(t *testing.T, key ConfKey, value string) {
	t.Helper()

	if err := Db.Where("k = ?", key).Delete(&Conf{}).Error; err != nil {
		t.Fatalf("delete conf %s: %v", key, err)
	}
	if err := Db.Create(&Conf{K: key, V: value}).Error; err != nil {
		t.Fatalf("create conf %s: %v", key, err)
	}
}

func TestNormalizePaymentTemplateMode(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  PaymentTemplateMode
	}{
		{name: "empty defaults to official", value: "", want: PaymentTemplateOfficial},
		{name: "official", value: "official", want: PaymentTemplateOfficial},
		{name: "wolf", value: "wolf", want: PaymentTemplateWolf},
		{name: "custom", value: "custom", want: PaymentTemplateCustom},
		{name: "trim and lower", value: " WOLF ", want: PaymentTemplateWolf},
		{name: "unknown defaults to official", value: "modern", want: PaymentTemplateOfficial},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePaymentTemplateMode(tt.value); got != tt.want {
				t.Fatalf("NormalizePaymentTemplateMode(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsValidPaymentTemplateMode(t *testing.T) {
	valid := []string{"", "official", "wolf", "custom", " WOLF "}
	for _, value := range valid {
		if !IsValidPaymentTemplateMode(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}

	invalid := []string{"modern", "enhanced", "default"}
	for _, value := range invalid {
		if IsValidPaymentTemplateMode(value) {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
}

func TestNormalizePaymentTemplateLanguage(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  PaymentTemplateLanguageMode
	}{
		{name: "empty defaults to auto", value: "", want: PaymentTemplateLanguageAuto},
		{name: "auto", value: "auto", want: PaymentTemplateLanguageAuto},
		{name: "simplified chinese", value: "zh", want: PaymentTemplateLanguageZh},
		{name: "traditional chinese preserves script tag", value: "zh-Hant", want: PaymentTemplateLanguageZhHant},
		{name: "case insensitive traditional chinese", value: " ZH-HANT ", want: PaymentTemplateLanguageZhHant},
		{name: "japanese", value: "ja", want: PaymentTemplateLanguageJa},
		{name: "unknown defaults to auto", value: "de", want: PaymentTemplateLanguageAuto},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePaymentTemplateLanguage(tt.value); got != tt.want {
				t.Fatalf("NormalizePaymentTemplateLanguage(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsValidPaymentTemplateLanguage(t *testing.T) {
	valid := []string{"", "auto", "zh", "zh-Hant", "en", "ru", "vi", "tr", "ja", "ko", " ZH-HANT "}
	for _, value := range valid {
		if !IsValidPaymentTemplateLanguage(value) {
			t.Fatalf("expected %q to be valid", value)
		}
	}

	invalid := []string{"de", "zh-TW", "browser", "default"}
	for _, value := range invalid {
		if IsValidPaymentTemplateLanguage(value) {
			t.Fatalf("expected %q to be invalid", value)
		}
	}
}

func TestUseCustomPaymentAssets(t *testing.T) {
	tests := []struct {
		name       string
		mode       string
		staticPath string
		want       bool
	}{
		{name: "no config uses embedded assets", want: false},
		{name: "legacy static path uses custom assets", staticPath: "/tmp/payment", want: true},
		{name: "custom mode with static path uses custom assets", mode: "custom", staticPath: "/tmp/payment", want: true},
		{name: "custom mode without static path uses embedded assets", mode: "custom", want: false},
		{name: "official ignores legacy static path", mode: "official", staticPath: "/tmp/payment", want: false},
		{name: "wolf ignores legacy static path", mode: "wolf", staticPath: "/tmp/payment", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newConfTestDB(t)

			if tt.mode != "" {
				upsertConfForTest(t, PaymentTemplate, tt.mode)
			}
			if tt.staticPath != "" {
				upsertConfForTest(t, PaymentStaticPath, tt.staticPath)
			}

			if got := UseCustomPaymentAssets(); got != tt.want {
				t.Fatalf("UseCustomPaymentAssets() = %v, want %v", got, tt.want)
			}
		})
	}
}
