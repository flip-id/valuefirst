package valuefirst

import (
	"github.com/fairyhunter13/phone"
	"github.com/fairyhunter13/pool"
	"strconv"
)

// RequestSendSMS is the request of SendSMS.
type RequestSendSMS struct {
	Version string `json:"@VER"`
	User    struct {
		Username      string `json:"@USERNAME,omitempty"`
		Password      string `json:"@PASSWORD,omitempty"`
		ChannelType   string `json:"@CH_TYPE,omitempty"`
		UnixTimestamp string `json:"@UNIXTIMESTAMP,omitempty"`
	} `json:"USER"`
	DLR RequestSendSMSDLR      `json:"DLR"`
	SMS RequestSendSMSMessages `json:"SMS"`
}

// SetVersion sets the version of the RequestSendSMS.
func (r *RequestSendSMS) SetVersion(version string) *RequestSendSMS {
	r.Version = version
	return r
}

// SetUnixTimestamp sets the UnixTimestamp of the RequestSendSMS.
func (r *RequestSendSMS) SetUnixTimestamp(unixTimestamp int64) *RequestSendSMS {
	r.User.UnixTimestamp = strconv.FormatInt(unixTimestamp, 10)
	return r
}

// SetChannelType sets the channel type of the RequestSendSMS.
func (r *RequestSendSMS) SetChannelType(chType string) *RequestSendSMS {
	r.User.ChannelType = chType
	return r
}

// SetPassword sets the password of the RequestSendSMS.
func (r *RequestSendSMS) SetPassword(password string) *RequestSendSMS {
	r.User.Password = password
	return r
}

// SetUsername sets the username of the RequestSendSMS.
func (r *RequestSendSMS) SetUsername(username string) *RequestSendSMS {
	r.User.Username = username
	return r
}

// SetDLRURL sets the DLR URL of the RequestSendSMS.
func (r *RequestSendSMS) SetDLRURL(urlStr string) *RequestSendSMS {
	r.DLR.URL = urlStr
	return r
}

// AddMessage adds a message to the RequestSendSMS.
func (r *RequestSendSMS) AddMessage(req *RequestSendSMSMessage) *RequestSendSMS {
	if req == nil {
		return r
	}

	r.SMS = append(r.SMS, *req)
	return r
}

// SetTypeWhatsapp sets the type of RequestSendSMS to Whatsapp.
func (r *RequestSendSMS) SetTypeWhatsapp() *RequestSendSMS {
	r.User.ChannelType = ChannelTypeWhatsapp
	return r
}

// Default sets the default values for the RequestSendSMS.
func (r *RequestSendSMS) Default() *RequestSendSMS {
	if r.Version == "" {
		r.Version = DefaultPayloadVersion
	}

	r.SMS = r.SMS.Default()
	return r
}

// RequestSendSMSDLR is the DLR part of the RequestSendSMS.
type RequestSendSMSDLR struct {
	URL string `json:"@URL"`
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

	// 	Embed Whatsapp Request
	RequestSendWhatsappMessage
}

// SetUDH sets the UDH of the message.
func (r *RequestSendSMSMessage) SetUDH(udh string) *RequestSendSMSMessage {
	r.UDH = udh
	return r
}

// SetCoding sets the coding of the message.
func (r *RequestSendSMSMessage) SetCoding(coding string) *RequestSendSMSMessage {
	r.Coding = coding
	return r
}

// SetText sets the text of the message.
func (r *RequestSendSMSMessage) SetText(text string) *RequestSendSMSMessage {
	r.Text = text
	return r
}

// SetTemplate sets the template of the RequestSendSMSMessage.
func (r *RequestSendSMSMessage) SetTemplate(templateID string, parameters ...string) *RequestSendSMSMessage {
	if templateID == "" {
		return r
	}

	builder := pool.GetStrBuilder()
	defer pool.Put(builder)

	_, _ = builder.WriteString(templateID)
	for _, param := range parameters {
		_, _ = builder.WriteString(DefaultTemplateSeparator)
		_, _ = builder.WriteString(param)
	}

	r.TemplateInfo = builder.String()
	return r
}

// SetProperty sets the property of the message.
func (r *RequestSendSMSMessage) SetProperty(property string) *RequestSendSMSMessage {
	r.Property = property
	return r
}

// SetID sets the ID of the message.
func (r *RequestSendSMSMessage) SetID(id string) *RequestSendSMSMessage {
	r.ID = id
	return r
}

// SetSendOn sets the scheduled sent time of the message.
func (r *RequestSendSMSMessage) SetSendOn(sendOn string) *RequestSendSMSMessage {
	r.SendOn = sendOn
	return r
}

func (r *RequestSendSMSMessage) AddAddress(addr *RequestSendSMSMessageAddress) *RequestSendSMSMessage {
	if addr == nil {
		return r
	}

	r.Address = append(r.Address, *addr)
	return r
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

	r.Text = Encode(r.Text)
	r.TemplateInfo = Encode(r.TemplateInfo)
	r.Address = r.Address.Default()
	return r
}

// RequestSendWhatsappMessage is the Whatsapp message part of the RequestSendSMSMessage.
type RequestSendWhatsappMessage struct {
	Caption     string `json:"@CAPTION,omitempty"`
	ContentType string `json:"@CONTENTTYPE,omitempty"`
	Type        string `json:"@TYPE,omitempty"`
	MessageType string `json:"@MSGTYPE,omitempty"`
	MediaData   string `json:"@MEDIADATA,omitempty"`
	BURLInfo    string `json:"@B_URLINFO,omitempty"`
}

// SetCaption sets the caption of Whatsapp message.
func (r *RequestSendSMSMessage) SetCaption(caption string) *RequestSendSMSMessage {
	r.Caption = caption
	return r
}

// SetContentType sets the content type of Whatsapp message.
func (r *RequestSendSMSMessage) SetContentType(contentType string) *RequestSendSMSMessage {
	r.ContentType = contentType
	return r
}

// SetType sets the type of Whatsapp message.
func (r *RequestSendSMSMessage) SetType(typ string) *RequestSendSMSMessage {
	r.Type = typ
	return r
}

// SetMessageType sets the message type of Whatsapp message.
func (r *RequestSendSMSMessage) SetMessageType(msgType MessageType) *RequestSendSMSMessage {
	r.MessageType = msgType.String()
	return r
}

// SetMediaData sets the media data of Whatsapp message.
func (r *RequestSendSMSMessage) SetMediaData(mediaData string) *RequestSendSMSMessage {
	r.MediaData = mediaData
	return r
}

// SetBURLInfo sets the dynamic URL information of Whatsapp message.
func (r *RequestSendSMSMessage) SetBURLInfo(bURLInfo string) *RequestSendSMSMessage {
	r.BURLInfo = bURLInfo
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

// SetFrom sets the sender address of the message.
func (r *RequestSendSMSMessageAddress) SetFrom(from string) *RequestSendSMSMessageAddress {
	r.From = from
	return r
}

// SetTo sets the destination address of the message.
func (r *RequestSendSMSMessageAddress) SetTo(to string) *RequestSendSMSMessageAddress {
	r.To = to
	return r
}

// SetSequence sets the sequence address of the message.
func (r *RequestSendSMSMessageAddress) SetSequence(seq string) *RequestSendSMSMessageAddress {
	r.Sequence = seq
	return r
}

// SetTag sets the tag address of the message.
func (r *RequestSendSMSMessageAddress) SetTag(tag string) *RequestSendSMSMessageAddress {
	r.Tag = tag
	return r
}

// Normalize normalizes the RequestSendSMSMessageAddress.
func (r *RequestSendSMSMessageAddress) Normalize() *RequestSendSMSMessageAddress {
	r.From = phone.NormalizeID(r.From, 0)
	r.To = phone.NormalizeID(r.To, 0)
	return r
}

// Default sets the default values for the RequestSendSMSMessageAddress.
func (r *RequestSendSMSMessageAddress) Default(id int) *RequestSendSMSMessageAddress {
	if r.Sequence == "" {
		r.Sequence = strconv.Itoa(id + 1)
	}
	return r.Normalize()
}
