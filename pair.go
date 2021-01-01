package exinswap

import (
	"context"

	"github.com/shopspring/decimal"
)

type ChainAsset struct {
	UUID    string `json:"uuid"`
	Symbol  string `json:"symbol"`
	Name    string `json:"name"`
	IconURL string `json:"iconUrl"`
}

type Asset struct {
	UUID       string     `json:"uuid"`
	Symbol     string     `json:"symbol"`
	Name       string     `json:"name"`
	RouteId    string     `json:"routeId"`
	IconURL    string     `json:"iconUrl"`
	ChainAsset ChainAsset `json:"chainAsset"`
}

type Pair struct {
	Asset0                  *Asset          `json:"asset0,omitempty"`
	Asset1                  *Asset          `json:"asset1,omitempty"`
	Asset0Balance           decimal.Decimal `json:"asset0Balance,omitempty"`
	Asset1Balance           decimal.Decimal `json:"asset1Balance,omitempty"`
	LpAsset                 *Asset          `json:"lpAsset,omitempty"`
	LpAssetSupply           string          `json:"lpAssetSupply,omitempty"`
	PriceRate24hours        decimal.Decimal `json:"priceRate24hours,omitempty"`
	PriceReverseRate24hours decimal.Decimal `json:"priceReverseRate24hours,omitempty"`
	UsdtTradeVolume24hours  decimal.Decimal `json:"usdtTradeVolume24hours,omitempty"`
	CreatedAt               int             `json:"createdAt,omitempty"`
	UpdatedAt               int             `json:"updatedAt,omitempty"`
}

func ReadPairs() (pairs []*Pair, timestampMs int64, err error) {
	uri := "/pairs"
	resp, err := Request(context.Background()).Get(uri)
	if err != nil {
		return nil, 0, err
	}

	if timestampMs, err = UnmarshalResponse(resp, &pairs); err != nil {
		return nil, timestampMs, err
	}
	return pairs, timestampMs, nil
}
