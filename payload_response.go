package valuefirst

import "fmt"

// ResponseMessage is the message response of SendSMS.
type ResponseMessage struct {
	MessageAck ResponseMessageAck `json:"MESSAGEACK"`
}

// GetError returns the error if any.
func (r *ResponseMessage) GetError() (err error) {
	err = r.MessageAck.GetError()
	return
}

// ResponseMessageAck is the message part of the ResponseMessage.
type ResponseMessageAck struct {
	GUID interface{} `json:"GUID"`

	// "Err": {
	//            "Desc": "The Specified message does not conform to DTD",
	//            "Code": 65535
	//        }
	Error *ResponseMessageAckError `json:"Err,omitempty"`
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

// ResponseMessageAckError is the error part of the ResponseMessageAck.
type ResponseMessageAckError struct {
	Description string `json:"Desc"`
	Code        int    `json:"Code"`
}

// Error implements the error interface.
func (r *ResponseMessageAckError) Error() string {
	return fmt.Sprintf(
		"error message ACK: %s (code: %d)",
		r.Description,
		r.Code,
	)
}

// GetError returns the error if any.
func (r *ResponseMessageAck) GetError() (err error) {
	if r.Error != nil {
		err = r.Error
		return
	}

	switch newType := r.GUID.(type) {
	case *ResponseMessageAckGUID:
		err = newType.GetError()
	case *ResponseMessageAckGUIDs:
		err = newType.GetError()
	}
	return
}

// ResponseMessageAckGUIDs is a collection of ResponseMessageAckGUID.
type ResponseMessageAckGUIDs []ResponseMessageAckGUID

// GetError returns the error if any.
func (r ResponseMessageAckGUIDs) GetError() (err error) {
	var errList = new(ErrorList)
	defer func() {
		if len(*errList) <= 0 {
			return
		}

		err = errList
	}()

	for _, val := range r {
		switch newType := val.Error.(type) {
		case *ResponseMessageAckGUIDError:
			errList.Append(filterError(newType))
		case *ResponseMessageAckGUIDErrors:
			errList.Append(filterErrors(newType))
		}
	}
	return
}

// ResponseMessageAckGUID is the GUID part of the ResponseMessageAck.
type ResponseMessageAckGUID struct {
	// A globally unique message ID that is generated for each <SMS> tag.
	// Note that, in future to check the status of the message you must save this code.
	GUID string `json:"GUID"`

	// The date and time when the transaction was completed.
	SubmitDate string `json:"SUBMITDATE"`

	// (In case of any error)
	// To conserve bandwidth utilization ValueFirst JSON API only sends
	// Sequence information of messages that has either some error or were rejected
	// because of some error. If there are no errors in a particular message,
	// you shall not receive any confirmation of each address SEQ.
	// For instance, in the above example in message ID 1 (of client)
	// the TO number "My company" was rejected as non-numeric.
	// The second message does not have any error,
	// and hence there was no error information for the second part.
	// SEQ: The Sequence ID (provided by client) that has error.
	// CODE: Reason why the message wasn’t accepted.
	// The table shown next describes these error conditions.
	Error interface{} `json:"ERROR,omitempty"`

	// Unique SMS ID sent by the customer. For each message a unique GUID is generated.
	// The Server sends the SMS ID so that
	// the client application can map the GUID to the SMS ID provided by them.
	ID int `json:"ID"`
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

// GetError returns the Error of the ResponseMessageAckGUID.
func (r *ResponseMessageAckGUID) GetError() (err error) {
	switch newType := r.Error.(type) {
	case *ResponseMessageAckGUIDError:
		err = filterError(newType)
	case *ResponseMessageAckGUIDErrors:
		err = filterErrors(newType)
	}
	return
}

// ResponseMessageAckGUIDErrors is the collection of the ResponseMessageAckGUIDError.
type ResponseMessageAckGUIDErrors []ResponseMessageAckGUIDError

// Error implements the error interface.
func (r ResponseMessageAckGUIDErrors) Error() (res string) {
	var errList ErrorList
	for _, err := range r {
		errList.Append(&err)
	}

	res = errList.Error()
	return
}

// ResponseMessageAckGUIDError is the error part of the ResponseMessageAckGUID.
//
//	"ERROR": {
//			"CODE": 28675,
//			"SEQ": 1
//			// OR
//			"SEQ": "1"
//	},
type ResponseMessageAckGUIDError struct {
	Code     int    `json:"CODE"`
	Sequence string `json:"SEQ"`
}

// Error implements the error interface.
func (r *ResponseMessageAckGUIDError) Error() (res string) {
	res = fmt.Sprintf(
		"error ValueFirst: CODE: %d, SEQ: %s",
		r.Code,
		r.Sequence,
	)
	return
}
