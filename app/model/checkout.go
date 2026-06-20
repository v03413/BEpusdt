package model

type Checkout struct {
	Name   string `json:"name"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
	Link   string `json:"link"`
}

var checkoutMap = make(map[string]Checkout)

// RegisterCheckout 注册收银台模板
func RegisterCheckout(name string, checkout Checkout) {
	checkoutMap[name] = checkout
}

func CheckoutList() map[string]Checkout {
	return checkoutMap
}
