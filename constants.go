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
)
