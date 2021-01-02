package exinswap

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Endpoint = "https://app.exinswap.com/api/v1"
)

type Error struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Msg)
}

var httpClient = resty.New().
	SetHeader("Accept", "application/json").
	SetHostURL(Endpoint).
	SetTimeout(300 * time.Millisecond)

func SetTimeout(timeout time.Duration) {
	httpClient.SetTimeout(timeout)
}

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, int64, error) {
	var body struct {
		Error
		Data        json.RawMessage `json:"data,omitempty"`
		Code        int             `json:"code,omitempty"`
		Success     bool            `json:"success,omitempty"`
		Message     string          `json:"message,omitempty"`
		TimestampMs int64           `json:"timestampMs,omitempty"`
	}

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		if resp.IsError() {
			return nil, 0, &Error{
				Code: resp.StatusCode(),
				Msg:  resp.Status(),
			}
		}
		return nil, 0, err
	}

	if body.Error.Code > 0 {
		return nil, 0, &body.Error
	}
	return body.Data, body.TimestampMs, nil
}

func UnmarshalResponse(resp *resty.Response, v interface{}) (timestampMs int64, err error) {
	data, timestampMs, err := DecodeResponse(resp)
	if err != nil {
		return timestampMs, err
	}

	if v != nil {
		return timestampMs, json.Unmarshal(data, v)
	}

	return timestampMs, nil
}
