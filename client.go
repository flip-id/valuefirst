package valuefirst

import (
	"context"
	"github.com/flip-id/valuefirst/manager"
)

var _ Client = new(client)

// Client is the interface for the ValueFirst API client.
type Client interface {
	manager.TokenClient
	// SendSMS sends SMS messages to a list of addresses.
	SendSMS(ctx context.Context, req *RequestSendSMS) (resp *ResponseMessage, err error)
}

type client struct {
	opt *Option
}

// New initialize the ValueFirst API client.
func New(opts ...FnOption) (c Client, err error) {
	opt := (new(Option)).Assign(opts...)
	err = opt.Validate()
	if err != nil {
		return
	}

	newClient := new(client)
	newClient.opt = opt.setValueFirstClient(newClient).
		Default().
		Clone()
	c = newClient
	return
}

// GenerateToken generate a token for the user.
func (c *client) GenerateToken(ctx context.Context) (res manager.ResponseGenerateToken, err error) {
	req, err := c.getRequestTokenManagement(URLActionGenerateToken)
	if err != nil {
		return
	}

	err = c.doRequestTokenManagement(ctx, req, &res)
	return
}

// EnableToken enables the token.
func (c *client) EnableToken(ctx context.Context, token string) (res manager.ResponseEnableToken, err error) {
	req, err := c.getRequestTokenManagement(URLActionEnableToken)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add(QueryParamToken, token)
	req.URL.RawQuery = q.Encode()
	err = c.doRequestTokenManagement(ctx, req, &res)
	return
}

// DisableToken disables the token.
func (c *client) DisableToken(ctx context.Context, token string) (res manager.ResponseEnableToken, err error) {
	req, err := c.getRequestTokenManagement(URLActionDisableToken)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add(QueryParamToken, token)
	req.URL.RawQuery = q.Encode()
	err = c.doRequestTokenManagement(ctx, req, &res)
	return
}

// DeleteToken disables the token.
func (c *client) DeleteToken(ctx context.Context, token string) (res manager.ResponseEnableToken, err error) {
	req, err := c.getRequestTokenManagement(URLActionDeleteToken)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add(QueryParamToken, token)
	req.URL.RawQuery = q.Encode()
	err = c.doRequestTokenManagement(ctx, req, &res)
	return
}

// SendSMS sends SMS messages to a list of addresses.
func (c *client) SendSMS(ctx context.Context, req *RequestSendSMS) (resp *ResponseMessage, err error) {
	if req == nil {
		err = ErrNilRequest
		return
	}

	r, buff, err := c.getRequestSendSMS(URLActionSendSMS, req.Default())
	if err != nil {
		return
	}
	defer c.putBuffer(buff)

	err = c.doRequestSendSMS(ctx, r, &resp)
	if err != nil {
		return
	}

	err = resp.GetError()
	return
}
