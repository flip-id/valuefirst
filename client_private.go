package valuefirst

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/fairyhunter13/pool"
	"github.com/fairyhunter13/reflecthelper/v5"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (c *client) isError(resp *http.Response) bool {
	return resp.StatusCode >= http.StatusBadRequest
}

func (c *client) getFullURL(subPath string) string {
	return c.opt.BaseURL + subPath
}

func (c *client) setCustomIPs(req *http.Request) *client {
	if len(c.opt.CustomIPs) <= 0 {
		return c
	}

	for _, ip := range c.opt.CustomIPs {
		req.Header.Add(fiber.HeaderXForwardedFor, ip)
	}
	return c
}

func (c *client) formatUnknownError(resp *http.Response) (err error) {
	if !c.isError(resp) {
		return
	}

	err = formatUnknown(resp)
	return
}

func (c *client) getRequestTokenManagement(action string) (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodPost, c.getFullURL(action), nil)
	return
}

func (c *client) doRequestTokenManagement(ctx context.Context, req *http.Request, decode interface{}) (err error) {
	if !reflecthelper.IsPtr(decode) {
		err = ErrDecodeVarIsNotPtr
		return
	}

	req.SetBasicAuth(c.opt.BasicAuth.User, c.opt.BasicAuth.Password)
	req = req.WithContext(ctx)
	resp, err := c.
		setCustomIPs(req).
		opt.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	err = c.formatUnknownError(resp)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(decode)
	return
}

func (c *client) putBuffer(buff *bytes.Buffer) {
	pool.Put(buff)
}

func (c *client) getRequestSendSMS(action string, in interface{}) (req *http.Request, buff *bytes.Buffer, err error) {
	buff = pool.GetBuffer()
	err = json.NewEncoder(buff).Encode(in)
	if err != nil {
		return
	}

	req, err = http.NewRequest(http.MethodPost, c.getFullURL(action), buff)
	return
}

func (c *client) doRequestSendSMS(ctx context.Context, req *http.Request, decode interface{}) (err error) {
	if !reflecthelper.IsPtr(decode) {
		err = ErrDecodeVarIsNotPtr
		return
	}

	token, err := c.opt.TokenManager.Get(ctx)
	if err != nil {
		return
	}

	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req = req.WithContext(ctx)
	resp, err := c.
		setCustomIPs(req).
		opt.client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	err = c.formatUnknownError(resp)
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(decode)
	return
}
