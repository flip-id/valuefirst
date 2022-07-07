//go:build integration
// +build integration

package valuefirst

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fairyhunter13/dotenv"
	"github.com/flip-id/valuefirst/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	c    Client
	once sync.Once
)

func setupClient() {
	once.Do(func() {
		err := dotenv.Load2(
			dotenv.WithPaths(".env"),
		)
		if err != nil {
			log.Fatalln(err)
		}

		c, err = New(
			WithBasicAuth(os.Getenv("VALUEFIRST_USERNAME"), os.Getenv("VALUEFIRST_PASSWORD")),
			WithCustomIPs(os.Getenv("VALUEFIRST_IP")),
		)
		if err != nil {
			log.Fatalln(err)
		}
	})
}

// Run integration tests.
// Notes: Run this test only on local, not on CI/CD.
func TestMain(m *testing.M) {
	flag.Parse()
	setupClient()

	os.Exit(m.Run())
}

func formatTime(in time.Time) string {
	return in.Format(storage.TimeFormatExpiredDate)
}

func TestSendSMSMultipleSuccess(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	future := now.Add(24 * time.Hour)
	req := &RequestSendSMS{
		SMS: RequestSendSMSMessages{
			RequestSendSMSMessage{
				Text: fmt.Sprintf("Hello, message are coming at %s!", formatTime(now)),
				Address: []RequestSendSMSMessageAddress{
					{
						From: os.Getenv("SENDER"),
						To:   os.Getenv("PHONE_NUMBER"),
					},
				},
			},
			RequestSendSMSMessage{
				Text: fmt.Sprintf("Hello, message are coming at %s!", formatTime(future)),
				Address: []RequestSendSMSMessageAddress{
					{
						From: os.Getenv("SENDER"),
						To:   os.Getenv("PHONE_NUMBER"),
					},
				},
			},
		},
	}
	resp, err := c.SendSMS(ctx, req)
	respB, _ := json.Marshal(resp)
	errB, _ := json.Marshal(err)
	log.Printf("resp: %s", respB)
	log.Printf("err: %s", errB)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUIDs)) == reflect.TypeOf(resp.MessageAck.GUID))

	guids := resp.MessageAck.GUID.(*ResponseMessageAckGUIDs)
	for idx, guid := range *guids {
		assert.NotEmpty(t, guid.GUID)
		assert.NotEmpty(t, guid.SubmitDate)
		assert.Nil(t, guid.Error)
		assert.Equal(t, idx+1, guid.ID)
	}
}

func TestSendSingleSMSSuccess(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	req := &RequestSendSMS{
		SMS: RequestSendSMSMessages{
			{
				Text: fmt.Sprintf("This single SMS are coming at %s!", formatTime(now)),
				Address: []RequestSendSMSMessageAddress{
					{
						From: os.Getenv("SENDER"),
						To:   os.Getenv("PHONE_NUMBER"),
					},
				},
			},
		},
	}
	resp, err := c.SendSMS(ctx, req)
	respB, _ := json.Marshal(resp)
	errB, _ := json.Marshal(err)
	log.Printf("resp: %s", respB)
	log.Printf("err: %s", errB)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUID)) == reflect.TypeOf(resp.MessageAck.GUID))

	guid := resp.MessageAck.GUID.(*ResponseMessageAckGUID)
	assert.NotEmpty(t, guid.GUID)
	assert.NotEmpty(t, guid.SubmitDate)
	assert.Nil(t, guid.Error)
	assert.Equal(t, 1, guid.ID)
}

func TestSendSMSMultipleError(t *testing.T) {
	t.Parallel()
	var idx time.Duration
	t.Run("no sender", func(t *testing.T) {
		idx := idx
		now := time.Now().UTC().Add(idx * time.Hour)
		future := now.Add(24 * time.Hour)
		ctx := context.Background()
		req := &RequestSendSMS{
			SMS: RequestSendSMSMessages{
				RequestSendSMSMessage{
					Text: fmt.Sprintf("Hello, message are coming at %s!", formatTime(now)),
					Address: []RequestSendSMSMessageAddress{
						{
							From: "",
							To:   os.Getenv("PHONE_NUMBER"),
						},
					},
				},
				RequestSendSMSMessage{
					Text: fmt.Sprintf("Hello, message are coming at %s!", formatTime(future)),
					Address: []RequestSendSMSMessageAddress{
						{
							From: "",
							To:   os.Getenv("PHONE_NUMBER"),
						},
					},
				},
			},
		}
		resp, err := c.SendSMS(ctx, req)
		respB, _ := json.Marshal(resp)
		errB, _ := json.Marshal(err)
		log.Printf("resp: %s", respB)
		log.Printf("err: %s", errB)
		assert.NotNil(t, err)
		assert.NotNil(t, resp)
		assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUIDs)) == reflect.TypeOf(resp.MessageAck.GUID))

		guids := resp.MessageAck.GUID.(*ResponseMessageAckGUIDs)
		for idx, guid := range *guids {
			assert.NotEmpty(t, guid.GUID)
			assert.NotEmpty(t, guid.SubmitDate)
			assert.NotNil(t, guid.Error)
			assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUIDError)) == reflect.TypeOf(guid.Error))
			assert.Equal(t, idx+1, guid.ID)
		}
	})
	idx++
	t.Run("no destination", func(t *testing.T) {
		idx := idx
		now := time.Now().UTC().Add(idx * time.Hour)
		future := now.Add(24 * time.Hour)
		ctx := context.Background()
		req := &RequestSendSMS{
			SMS: RequestSendSMSMessages{
				RequestSendSMSMessage{
					Text: fmt.Sprintf("Hello, message are coming at %s!", formatTime(now)),
					Address: []RequestSendSMSMessageAddress{
						{
							From: os.Getenv("SENDER"),
							To:   "",
						},
					},
				},
				RequestSendSMSMessage{
					Text: fmt.Sprintf("Hello, message are coming at %s!", formatTime(future)),
					Address: []RequestSendSMSMessageAddress{
						{
							From: os.Getenv("SENDER"),
							To:   "",
						},
					},
				},
			},
		}
		resp, err := c.SendSMS(ctx, req)
		respB, _ := json.Marshal(resp)
		errB, _ := json.Marshal(err)
		log.Printf("resp: %s", respB)
		log.Printf("err: %s", errB)
		assert.NotNil(t, err)
		assert.NotNil(t, resp)
		assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUIDs)) == reflect.TypeOf(resp.MessageAck.GUID))

		guids := resp.MessageAck.GUID.(*ResponseMessageAckGUIDs)
		for idx, guid := range *guids {
			assert.NotEmpty(t, guid.GUID)
			assert.NotEmpty(t, guid.SubmitDate)
			assert.NotNil(t, guid.Error)
			assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUIDError)) == reflect.TypeOf(guid.Error))
			assert.Equal(t, idx+1, guid.ID)
		}
	})
}

func TestSendSMSSingleError(t *testing.T) {
	t.Parallel()
	var idx time.Duration
	t.Run("no sender", func(t *testing.T) {
		idx := idx
		now := time.Now().UTC().Add(idx * time.Hour)
		ctx := context.Background()
		req := &RequestSendSMS{
			SMS: RequestSendSMSMessages{
				{
					Text: fmt.Sprintf("This single SMS are coming at %s!", formatTime(now)),
					Address: []RequestSendSMSMessageAddress{
						{
							From: "",
							To:   os.Getenv("PHONE_NUMBER"),
						},
					},
				},
			},
		}
		resp, err := c.SendSMS(ctx, req)
		respB, _ := json.Marshal(resp)
		errB, _ := json.Marshal(err)
		log.Printf("resp: %s", respB)
		log.Printf("err: %s", errB)
		assert.NotNil(t, err)
		assert.NotNil(t, resp)
		assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUID)) == reflect.TypeOf(resp.MessageAck.GUID))

		guid := resp.MessageAck.GUID.(*ResponseMessageAckGUID)
		assert.NotEmpty(t, guid.GUID)
		assert.NotEmpty(t, guid.SubmitDate)
		assert.NotNil(t, guid.Error)
		assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUIDError)) == reflect.TypeOf(guid.Error))
		assert.Equal(t, 1, guid.ID)
	})
	idx++
	t.Run("no destination", func(t *testing.T) {
		idx := idx
		now := time.Now().UTC().Add(idx * time.Hour)
		ctx := context.Background()
		req := &RequestSendSMS{
			SMS: RequestSendSMSMessages{
				{
					Text: fmt.Sprintf("This single SMS are coming at %s!", formatTime(now)),
					Address: []RequestSendSMSMessageAddress{
						{
							From: os.Getenv("SENDER"),
							To:   "",
						},
					},
				},
			},
		}
		resp, err := c.SendSMS(ctx, req)
		respB, _ := json.Marshal(resp)
		errB, _ := json.Marshal(err)
		log.Printf("resp: %s", respB)
		log.Printf("err: %s", errB)
		assert.NotNil(t, err)
		assert.NotNil(t, resp)
		assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUID)) == reflect.TypeOf(resp.MessageAck.GUID))

		guid := resp.MessageAck.GUID.(*ResponseMessageAckGUID)
		assert.NotEmpty(t, guid.GUID)
		assert.NotEmpty(t, guid.SubmitDate)
		assert.NotNil(t, guid.Error)
		assert.True(t, reflect.TypeOf(new(ResponseMessageAckGUIDError)) == reflect.TypeOf(guid.Error))
		assert.Equal(t, 1, guid.ID)
	})
}
