package exinswap

import "github.com/shopspring/decimal"
import "errors"

var (
	SwapFee = "0.003"

	ErrInsufficientLiquiditySwapped = errors.New("insufficient liquidity swapped")
)

// Result represent Swap Result
type Result struct {
	PayAssetID  string
	PayAmount   decimal.Decimal
	FillAssetID string
	FillAmount  decimal.Decimal
	FeeAssetID  string
	FeeAmount   decimal.Decimal
	RouteID     string
}

// Swap trade in a pair
func Swap(pair *Pair, payAssetID string, payAmount decimal.Decimal) (*Result, error) {
	K := pair.Asset0Balance.Mul(pair.Asset1Balance)
	if !K.IsPositive() {
		return nil, ErrInsufficientLiquiditySwapped
	}

	payAmount = payAmount.Truncate(8)

	r := &Result{
		PayAssetID: payAssetID,
		PayAmount:  payAmount,
		FeeAssetID: payAssetID,
		FeeAmount:  payAmount.Mul(Decimal(SwapFee)).Truncate(8),
		RouteID:    pair.LpAsset.RouteId,
	}

	funds := payAmount.Sub(r.FeeAmount)
	if !funds.IsPositive() {
		return nil, errors.New("pay amount must be positive")
	}

	switch payAssetID {
	case pair.Asset0.UUID:
		newBase := pair.Asset0Balance.Add(funds)
		newQuote := K.Div(newBase)
		r.FillAssetID = pair.Asset1.UUID
		r.FillAmount = pair.Asset1Balance.Sub(newQuote).Truncate(8)
	case pair.Asset1.UUID:
		newQuote := pair.Asset1Balance.Add(funds)
		newBase := K.Div(newQuote)
		r.FillAssetID = pair.Asset0.UUID
		r.FillAmount = pair.Asset0Balance.Sub(newBase).Truncate(8)
	default:
		return nil, errors.New("invalid pay asset id")
	}

	return r, nil
}

// ReverseSwap is a Reverse version of Swap
func ReverseSwap(pair *Pair, fillAssetID string, fillAmount decimal.Decimal) (*Result, error) {
	K := pair.Asset0Balance.Mul(pair.Asset1Balance)
	if !K.IsPositive() {
		return nil, ErrInsufficientLiquiditySwapped
	}

	fillAmount = fillAmount.Truncate(8)
	if !fillAmount.IsPositive() {
		return nil, errors.New("invalid fill amount")
	}

	r := &Result{
		FillAssetID: fillAssetID,
		FillAmount:  fillAmount,
	}

	switch fillAssetID {
	case pair.Asset0.UUID:
		newBase := pair.Asset0Balance.Sub(fillAmount)
		if !newBase.IsPositive() {
			return nil, ErrInsufficientLiquiditySwapped
		}

		newQuote := K.Div(newBase)
		r.PayAssetID = pair.Asset1.UUID
		r.PayAmount = newQuote.Sub(pair.Asset1Balance)
	case pair.Asset1.UUID:
		newQuote := pair.Asset1Balance.Sub(fillAmount)
		if !newQuote.IsPositive() {
			return nil, ErrInsufficientLiquiditySwapped
		}

		newBase := K.Div(newQuote)
		r.PayAssetID = pair.Asset0.UUID
		r.PayAmount = newBase.Sub(pair.Asset0Balance)
	default:
		return nil, errors.New("invalid fill asset id")
	}

	r.PayAmount = r.PayAmount.Div(decimal.NewFromInt(1).Sub(Decimal(SwapFee)))
	r.PayAmount = Ceil(r.PayAmount, 8)
	r.FeeAssetID = r.PayAssetID
	r.FeeAmount = r.PayAmount.Mul(Decimal(SwapFee)).Truncate(8)

	return r, nil
}
