//go:build !integration
// +build !integration

package valuefirst

import (
	"context"
	"github.com/flip-id/valuefirst/manager"
	"github.com/flip-id/valuefirst/mocks"
	"github.com/flip-id/valuefirst/storage"
	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type TestStruct struct {
	Manager manager.TokenManager
	Hub     storage.Hub
	Mock    struct {
		Manager *mocks.MockTokenManager
		Hub     *mocks.MockHub
	}
	HTTPClient *http.Client
	Ctrl       *gomock.Controller
}

func (s *TestStruct) Finish() *TestStruct {
	s.Ctrl.Finish()
	httpmock.DeactivateAndReset()
	return s
}

func (s *TestStruct) Setup(t *testing.T) *TestStruct {
	s.Ctrl = gomock.NewController(t)

	// Mocks
	s.Mock.Manager = mocks.NewMockTokenManager(s.Ctrl)
	s.Mock.Hub = mocks.NewMockHub(s.Ctrl)

	// Implementations
	s.Manager = s.Mock.Manager
	s.Hub = s.Mock.Hub

	s.HTTPClient = new(http.Client)
	httpmock.ActivateNonDefault(s.HTTPClient)
	return s
}

func Test_client_SendSMS(t *testing.T) {
	t.Parallel()
	type fields struct {
		opt *Option
	}
	type args struct {
		ctx context.Context
		req *RequestSendSMS
	}
	tests := []struct {
		name     string
		fields   func(s *TestStruct) fields
		args     func() args
		wantResp func() *ResponseMessage
		wantErr  bool
	}{
		{
			name: "success sending the SMS",
			fields: func(s *TestStruct) fields {
				s.Mock.Manager.EXPECT().Get(gomock.Any()).
					Return(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIzNDU2Nzg5LCJuYW1lIjoiSm9zZXBoIn0.OpOSSw7e485LOP5PrzScxHb7SR6sAOMRckfFwi4rp7o",
					nil,
				)
				httpmock.RegisterResponder(
					http.MethodPost,
					BaseURL+URLActionSendSMS,
					httpmock.NewJsonResponderOrPanic(http.StatusOK, ResponseMessage{
						MessageAck: ResponseMessageAck{
							GUID: ResponseMessageAckGUID{
								GUID:       "km8vd0730664x1f440e00aflkjFLIPTECHWA",
								SubmitDate: "2022-08-31 13:07:30",
								ID:         1,
							},
						},
					}),
				)
				return fields{
					opt: &Option{
						Client:       s.HTTPClient,
						TokenManager: s.Manager,
					},
				}
			},
			args: func() args {
				return args{
					ctx: context.Background(),
					req: (new(RequestSendSMS)).
						AddMessage(
						(new(RequestSendSMSMessage)).
							SetText("Hello World").
							AddAddress(
							(new(RequestSendSMSMessageAddress)).
								SetFrom("Flip").
								SetTo("6281301933595"),
						),
					),
				}
			},
			wantResp: func() *ResponseMessage {
				return &ResponseMessage{
					MessageAck: ResponseMessageAck{
						GUID: &ResponseMessageAckGUID{
							GUID:       "km8vd0730664x1f440e00aflkjFLIPTECHWA",
							SubmitDate: "2022-08-31 13:07:30",
							ID:         1,
						},
					},
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct := (new(TestStruct)).Setup(t)
			defer testStruct.Finish()
			c := &client{
				opt: tt.fields(testStruct).opt.Default(),
			}
			args := tt.args()
			gotResp, err := c.SendSMS(args.ctx, args.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			wantResp := tt.wantResp()
			assert.EqualValuesf(t, wantResp, gotResp, "SendSMS(%v, %v)", args.ctx, args.req)
		})
	}
}
