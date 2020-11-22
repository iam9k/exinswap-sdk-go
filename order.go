package exinswap

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	OrderStateTrading  = "Trading"
	OrderStateRejected = "Rejected"
	OrderStateDone     = "Done"
)

type Order struct {
	ID          string          `json:"id,omitempty"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	State       string          `json:"state,omitempty"`
	PayAssetID  string          `json:"pay_asset_id,omitempty"`
	Funds       decimal.Decimal `json:"funds,omitempty"`
	FillAssetID string          `json:"fill_asset_id,omitempty"`
	Amount      decimal.Decimal `json:"amount,omitempty"`
	PriceImpact decimal.Decimal `json:"price_impact,omitempty"`
	RouteAssets []string        `json:"route_assets,omitempty"`
	// route id
	Routes string `json:"routes,omitempty"`
}

type PreOrderReq struct {
	PayAssetID  string          `json:"pay_asset_id,omitempty"`
	FillAssetID string          `json:"fill_asset_id,omitempty"`
	Funds       decimal.Decimal `json:"funds,omitempty"`
	Amount      decimal.Decimal `json:"amount,omitempty"`
}
