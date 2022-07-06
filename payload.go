package valuefirst

import (
	"strconv"
)

// RequestSendSMS is the request of SendSMS.
type RequestSendSMS struct {
	Version string                 `json:"@VER"`
	User    struct{}               `json:"USER"`
	DLR     RequestSendSMSDLR      `json:"DLR"`
	SMS     RequestSendSMSMessages `json:"SMS"`
}

// RequestSendSMSDLR is the DLR part of the RequestSendSMS.
type RequestSendSMSDLR struct {
	URL string `json:"@URL"`
}

// Default sets the default values for the RequestSendSMS.
func (r *RequestSendSMS) Default() *RequestSendSMS {
	if r.Version == "" {
		r.Version = DefaultPayloadVersion
	}

	r.SMS = r.SMS.Default()
	return r
}

// RequestSendSMSMessages is a collection of RequestSendSMSMessage.
type RequestSendSMSMessages []RequestSendSMSMessage

// Default sets the default values for each RequestSendSMSMessage.
func (r RequestSendSMSMessages) Default() RequestSendSMSMessages {
	for idx, val := range r {
		r[idx] = *val.Default(idx)
	}
	return r
}

// RequestSendSMSMessage is the message part of the RequestSendSMS.
type RequestSendSMSMessage struct {
	// UDH is used for sending binary messages. For text message the value should be 0.
	UDH string `json:"@UDH"`

	// Extended type of messages. For text message the value should be 1.
	Coding string `json:"@CODING"`

	// This field describe the message text to be sent to receiver.
	// SMS can contain up to 160 characters in Message Text.
	// API allows user to send Message text of more than 160 characters.
	// Credits will be deducted in the multiple of 160 characters according to the length of SMS.
	Text string `json:"@TEXT"`

	// New Parameter TEMPLATEINFO has been added for above functionality
	// which will contain the template id and variables value to be replaced in template text.
	// Template info parameter will have ~ separated values.
	//
	// If both TEXT and TEMPLATEINFO is given then priority will be given to Text.
	//
	// New error code i.e. INVALID_TEMPLATEINFO = 28694; has been created for error if occurred any,
	// related to TEMPLATEINFO parameter, which include like invalid templateid is provided,
	// variables count mismatch than the template Text variables count,
	// template text not found for the given template id.
	TemplateInfo string `json:"@TEMPLATEINFO"`

	// Unique property of message. Default value is 0. For sending Flash SMS the value should be 1.
	Property string `json:"@PROPERTY"`

	// Unique ID of message. The client sends this value.
	// In future communication, server sends this value back to the client.
	// This value is used in future to check status of the message.
	ID string `json:"@ID"`

	// It is now possible to schedule a message.
	// To schedule message to go at a later time,
	// user can specify “SEND_ON” date as attribute of SMS tag.
	// Only absolute date is supported.
	// The value should be given in “YYYY-MM-DD HH:MM:SS TIMEZONE” format.
	// Time zone is difference w.r.t. to GMT.
	// Please refer Scheduling Support for more information on this feature.
	SendOn string `json:"@SEND_ON"`

	// Describe the Sender as well as Receiver address.
	Address RequestSendSMSMessageAddresses `json:"ADDRESS"`
}

// Default sets	the default values for the RequestSendSMSMessage.
func (r *RequestSendSMSMessage) Default(id int) *RequestSendSMSMessage {
	if r.UDH == "" {
		r.UDH = DefaultUDHText
	}

	if r.Coding == "" {
		r.Coding = DefaultCodingText
	}

	if r.Property == "" {
		r.Property = DefaultPropertyText
	}

	if r.ID == "" {
		r.ID = strconv.Itoa(id + 1)
	}

	r.Address = r.Address.Default()
	return r
}

// RequestSendSMSMessageAddresses is a collection of RequestSendSMSMessageAddress.
type RequestSendSMSMessageAddresses []RequestSendSMSMessageAddress

// Default sets the default values for each RequestSendSMSMessageAddress.
func (r RequestSendSMSMessageAddresses) Default() RequestSendSMSMessageAddresses {
	for idx, val := range r {
		r[idx] = *val.Default(idx)
	}
	return r
}

// RequestSendSMSMessageAddress is the address part of the RequestSendSMSMessage.
type RequestSendSMSMessageAddress struct {
	// The Sender of the message.
	// This field should conform to Sender Phone Number guidelines.
	From string `json:"@FROM"`

	// Person receiving the SMS, should confirm to Receiver Phone Number guidelines.
	To string `json:"@TO"`

	// Unique Sequence ID.
	// Must be an integer and must be unique to each SMS.
	// While checking message status, you must send this value.
	Sequence string `json:"@SEQ"`

	// A text that identify message.
	// This is an optional parameter.
	Tag string `json:"@TAG"`
}

// Default sets the default values for the RequestSendSMSMessageAddress.
func (r *RequestSendSMSMessageAddress) Default(id int) *RequestSendSMSMessageAddress {
	if r.Sequence == "" {
		r.Sequence = strconv.Itoa(id + 1)
	}
	return r
}

// GetMessageGUID return the GUID of the ResponseMessageAck.
func (r *ResponseMessageAck) GetMessageGUID() (res *ResponseMessageAckGUID, ok bool) {
	res, ok = r.GUID.(*ResponseMessageAckGUID)
	return
}

// GetMessageGUIDs return the GUIDs of the ResponseMessageAck.
func (r *ResponseMessageAck) GetMessageGUIDs() (res *ResponseMessageAckGUIDs, ok bool) {
	res, ok = r.GUID.(*ResponseMessageAckGUIDs)
	return
}

// GetMessageError return the Error of the ResponseMessageAckGUID.
func (r *ResponseMessageAckGUID) GetMessageError() (res *ResponseMessageAckGUIDError, ok bool) {
	res, ok = r.Error.(*ResponseMessageAckGUIDError)
	return
}

// GetMessageErrors return the Errors of the ResponseMessageAckGUID.
func (r *ResponseMessageAckGUID) GetMessageErrors() (res *ResponseMessageAckGUIDErrors, ok bool) {
	res, ok = r.Error.(*ResponseMessageAckGUIDErrors)
	return
}
