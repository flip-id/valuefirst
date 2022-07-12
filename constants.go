package valuefirst

// List of URLs used in this package.
const (
	URLToken               = "/api/messages/token"
	URLActionGenerateToken = URLToken + "?action=generate"
	URLActionEnableToken   = URLToken + "?action=enable"
	URLActionDisableToken  = URLToken + "?action=disable"
	URLActionDeleteToken   = URLToken + "?action=delete"
	URLActionSendSMS       = "/servlet/psms.JsonEservice"
)

const (
	// QueryParamToken is the key for adding the token to the query params of URL.
	QueryParamToken = "token"

	// DefaultPayloadVersion is the default version of Payload.
	DefaultPayloadVersion = "1.2"

	// DefaultUDHText is the default UDH used for the text message.
	// UDH is used for sending binary messages.
	// For text message the value should be 0.
	DefaultUDHText = "0"

	// DefaultCodingText is the default coding used for the text message.
	// Extended type of messages.
	// For text message the value should be 1.
	DefaultCodingText = "1"

	// DefaultPropertyText is the default property used for the text message.
	// Unique property of message.
	// Default value is 0.
	// For sending Flash SMS the value should be 1.
	DefaultPropertyText = "0"
	// DefaultPropertyFlash is the default property used for the flash message.
	DefaultPropertyFlash = "1"
	// ChannelTypeWhatsapp is the channel type for Whatsapp messages.
	ChannelTypeWhatsapp = "4"
	// 	DefaultTemplateSeparator is the default separator used for TEMPLATEINFO in ValueFirst.
	DefaultTemplateSeparator = "~"
)

// MessageType is the type of the message for Whatsapp.
type MessageType string

// String returns the string representation of the MessageType.
func (m MessageType) String() string {
	return string(m)
}

const (
	// MessageTypePlain is the plain type of Whatsapp message.
	MessageTypePlain MessageType = "1"
	// MessageTypePlainTwoWay is the plain type of Whatsapp message with two-way communication.
	MessageTypePlainTwoWay MessageType = "2"
	// MessageTypeRich is the rich type of Whatsapp message.
	MessageTypeRich MessageType = "3"
	// MessageTypeRichTwoWay is the rich type of Whatsapp message with two-way communication.
	MessageTypeRichTwoWay MessageType = "4"
	// MessageTypeBusinessCardSharingTwoWay is the business card sharing type of Whatsapp message with two-way communication.
	MessageTypeBusinessCardSharingTwoWay MessageType = "5"
	// MessageTypeLocationSharingTwoWay is the location sharing type of Whatsapp message with two-way communication.
	MessageTypeLocationSharingTwoWay MessageType = "6"
	// MessageTypeListTwoWay is the list type of Whatsapp message with two-way communication.
	MessageTypeListTwoWay MessageType = "7"
	// MessageTypeReplyButtonTwoWay is the reply button type of Whatsapp message with two-way communication.
	MessageTypeReplyButtonTwoWay MessageType = "8"
)
