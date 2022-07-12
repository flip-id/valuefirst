//go:build integration
// +build integration

package valuefirst

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	waC    Client
	onceWa sync.Once
)

func setupClientWhatsapp() {
	onceWa.Do(func() {
		var err error
		waC, err = New(
			WithBasicAuth(os.Getenv("VALUEFIRST_WA_USERNAME"), os.Getenv("VALUEFIRST_WA_PASSWORD")),
			WithCustomIPs(os.Getenv("VALUEFIRST_IP")),
		)
		if err != nil {
			log.Fatalln(err)
		}
	})
}

func TestSendWhatsappPlainSuccess(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	req := (new(RequestSendSMS)).
		SetTypeWhatsapp().
		AddMessage((new(RequestSendSMSMessage)).
		SetTemplate(os.Getenv("VALUEFIRST_WA_TEMPLATE_PLAIN"), "test "+now.Format(time.RFC3339), "https://flip.id").
		SetMessageType(MessageTypePlain).
		AddAddress((new(RequestSendSMSMessageAddress)).
		SetTo(os.Getenv("PHONE_NUMBER")).
		SetFrom(os.Getenv("VALUEFIRST_WA_SENDER")),
	),
	)

	resp, err := waC.SendSMS(ctx, req)
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
