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

func Request(ctx context.Context) *resty.Request {
	return httpClient.R().SetContext(ctx)
}

func RequestWithTimeout(ctx context.Context, timeout time.Duration) *resty.Request {
	return httpClient.SetTimeout(timeout).R().SetContext(ctx)
}

func DecodeResponse(resp *resty.Response) ([]byte, error) {
	var body struct {
		Error
		Data    json.RawMessage `json:"data,omitempty"`
		Code    int             `json:"code,omitempty"`
		Success bool            `json:"success,omitempty"`
		Message string          `json:"message,omitempty"`
	}

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		if resp.IsError() {
			return nil, &Error{
				Code: resp.StatusCode(),
				Msg:  resp.Status(),
			}
		}
		return nil, err
	}

	if body.Error.Code > 0 {
		return nil, &body.Error
	}
	return body.Data, nil
}

func UnmarshalResponse(resp *resty.Response, v interface{}) error {
	data, err := DecodeResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		return json.Unmarshal(data, v)
	}

	return nil
}
