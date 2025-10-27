package tests

import (
	"testing"

	"github.com/venomous-maker/go-mpesa/Services"
)

func TestNewBusinessBuyGoodsService(t *testing.T) {
	svc := Services.NewBusinessBuyGoodsService(nil, nil)
	if svc == nil {
		t.Fatal("expected service not nil")
	}
}
