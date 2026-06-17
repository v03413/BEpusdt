package epusdt

import (
	"strings"
	"testing"

	"github.com/v03413/bepusdt/app/model"
)

func TestCheckoutCashierTemplateName(t *testing.T) {
	tests := []struct {
		name string
		mode model.PaymentTemplateMode
		want string
	}{
		{name: "official uses upstream cashier", mode: model.PaymentTemplateOfficial, want: "cashier.html"},
		{name: "wolf uses LangGe cashier", mode: model.PaymentTemplateWolf, want: "wolf.cashier.html"},
		{name: "custom uses custom cashier filename", mode: model.PaymentTemplateCustom, want: "cashier.html"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkoutCashierTemplateName(tt.mode); got != tt.want {
				t.Fatalf("checkoutCashierTemplateName(%q) = %q, want %q", tt.mode, got, tt.want)
			}
		})
	}
}

func TestJSONScriptEscapesHTMLSensitiveCharacters(t *testing.T) {
	got := string(jsonScript(map[string]string{
		"name": "</script><img src=x onerror=alert(1)>&",
	}))

	for _, unsafe := range []string{"</script>", "<img", ">", "&"} {
		if strings.Contains(got, unsafe) {
			t.Fatalf("jsonScript output contains unsafe substring %q: %s", unsafe, got)
		}
	}
	if !strings.Contains(got, `\u003c/script\u003e`) {
		t.Fatalf("jsonScript output does not contain escaped script close tag: %s", got)
	}
}

func TestCheckoutCounterTemplateName(t *testing.T) {
	tests := []struct {
		name      string
		mode      model.PaymentTemplateMode
		tradeType model.TradeType
		want      string
	}{
		{name: "official uses trade template", mode: model.PaymentTemplateOfficial, tradeType: model.UsdtTrc20, want: "usdt.trc20.html"},
		{name: "wolf uses universal LangGe checkout", mode: model.PaymentTemplateWolf, tradeType: model.UsdtTrc20, want: "wolf.checkout.html"},
		{name: "custom uses trade template", mode: model.PaymentTemplateCustom, tradeType: model.UsdtTrc20, want: "usdt.trc20.html"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkoutCounterTemplateName(tt.mode, tt.tradeType); got != tt.want {
				t.Fatalf("checkoutCounterTemplateName(%q, %q) = %q, want %q", tt.mode, tt.tradeType, got, tt.want)
			}
		})
	}
}
