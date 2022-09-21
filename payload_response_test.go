//go:build !integration
// +build !integration

package valuefirst

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseMessageAck_GetError(t *testing.T) {
	type fields struct {
		GUID  interface{}
		Error *ResponseMessageAckError
	}
	tests := []struct {
		name    string
		fields  func() fields
		wantErr bool
	}{
		{
			name: "error DTD",
			fields: func() fields {
				return fields{
					GUID: nil,
					Error: &ResponseMessageAckError{
						Description: "The specified message does not conform to DTD",
						Code:        ErrCodeGeneralNotConformDTD,
					},
				}
			},
			wantErr: true,
		},
		{
			name: "error GUID",
			fields: func() fields {
				return fields{
					GUID: &ResponseMessageAckGUID{
						GUID:       "KI2EC395172653F410004TXOENC9KDEFAULT",
						SubmitDate: "2018-02-14 12:39:51",
						Error: &ResponseMessageAckGUIDError{
							Code:     ErrCodeMessageInvalidSenderID,
							Sequence: "1",
						},
						ID: 1,
					},
					Error: nil,
				}
			},
			wantErr: true,
		},
		{
			name: "error GUIDs",
			fields: func() fields {
				return fields{
					GUID: &ResponseMessageAckGUID{
						GUID:       "KI2EC395172653F410004TXOENC9KDEFAULT",
						SubmitDate: "2018-02-14 12:39:51",
						Error: &ResponseMessageAckGUIDErrors{
							{
								Code:     ErrCodeMessageInvalidSenderID,
								Sequence: "1",
							},
							{
								Code:     ErrCodeMessageInvalidMessage,
								Sequence: "2",
							},
						},
						ID: 1,
					},
					Error: nil,
				}
			},
			wantErr: true,
		},
		{
			name: "ignored error GUID",
			fields: func() fields {
				return fields{
					GUID: &ResponseMessageAckGUID{
						GUID:       "KI2EC395172653F410004TXOENC9KDEFAULT",
						SubmitDate: "2018-02-14 12:39:51",
						Error: &ResponseMessageAckGUIDError{
							Code:     ErrCodeSchedulerSuccess,
							Sequence: "1",
						},
						ID: 1,
					},
					Error: nil,
				}
			},
			wantErr: false,
		},
		{
			name: "ignored error GUIDs",
			fields: func() fields {
				return fields{
					GUID: &ResponseMessageAckGUID{
						GUID:       "KI2EC395172653F410004TXOENC9KDEFAULT",
						SubmitDate: "2018-02-14 12:39:51",
						Error: &ResponseMessageAckGUIDErrors{
							{
								Code:     ErrCodeSchedulerSuccess,
								Sequence: "1",
							},
							{
								Code:     ErrCodeGeneralSuccess,
								Sequence: "2",
							},
							{
								Code:     ErrCodeStatusRequestSuccess,
								Sequence: "3",
							},
						},
						ID: 1,
					},
					Error: nil,
				}
			},
			wantErr: false,
		},
		{
			name: "nil error",
			fields: func() fields {
				return fields{
					GUID:  nil,
					Error: nil,
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			r := &ResponseMessageAck{
				GUID:  fields.GUID,
				Error: fields.Error,
			}
			err := r.GetError()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
