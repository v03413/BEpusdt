package task

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/v03413/bepusdt/app/model"
)

func TestAmountMatchHasPrefix(t *testing.T) {
	initTaskTestEnv(t, map[model.ConfKey]string{model.PaymentMatchMode: string(model.HasPrefix)})

	cases := []struct {
		paid, order string
		want        bool
	}{
		{"1.5", "1.5", true},
		{"1.50123", "1.5", true},
		{"1.509", "1.5", true},
		{"1.51", "1.5", false},
		{"1.4999", "1.5", false},
		{"2", "2", true},
		{"2.009", "2", true},
		{"2.01", "2", false},
		{"2.5", "2", false},
		{"10.379", "10.37", true},
		{"10.38", "10.37", false},
		{"20", "2", false},
		{"100.5", "100", false},
		{"123.456789", "123.45", true},
		// 订单金额小数位多于原子精度时按订单金额的小数位截断
		{"1.5059", "1.505", true},
		{"1.506", "1.505", false},
	}
	for _, c := range cases {
		got := amountMatch(decimal.RequireFromString(c.paid), c.order, string(model.UsdtTrc20))
		if got != c.want {
			t.Errorf("paid %s order %s: got %v, want %v", c.paid, c.order, got, c.want)
		}
	}
}

// 同一地址并发订单按原子精度递增金额，截断后每笔付款最多对应其中一笔
func TestAmountMatchHasPrefixConcurrentOrders(t *testing.T) {
	initTaskTestEnv(t, map[model.ConfKey]string{model.PaymentMatchMode: string(model.HasPrefix)})

	orders := []string{"1.5", "1.51", "1.52"}
	payments := map[string]string{"1.5": "1.5", "1.50123": "1.5", "1.51": "1.51", "1.519": "1.51", "1.52": "1.52"}
	for paid, want := range payments {
		var matched []string
		for _, o := range orders {
			if amountMatch(decimal.RequireFromString(paid), o, string(model.UsdtTrc20)) {
				matched = append(matched, o)
			}
		}
		if len(matched) != 1 || matched[0] != want {
			t.Errorf("paid %s matched %v, want only %s", paid, matched, want)
		}
	}
}
