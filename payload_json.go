package valuefirst

import (
	"encoding/json"
)

var (
	_ json.Unmarshaler = new(ResponseMessageAck)
	_ json.Unmarshaler = new(ResponseMessageAckGUID)
)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *ResponseMessageAck) UnmarshalJSON(data []byte) (err error) {
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	switch result["GUID"].(type) {
	case map[string]interface{}:
		r.GUID = new(ResponseMessageAckGUID)
	case []interface{}:
		r.GUID = new(ResponseMessageAckGUIDs)
	}

	type tmp ResponseMessageAck
	err = json.Unmarshal(data, (*tmp)(r))
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (r *ResponseMessageAckGUID) UnmarshalJSON(data []byte) (err error) {
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	switch result["ERROR"].(type) {
	case map[string]interface{}:
		r.Error = new(ResponseMessageAckGUIDError)
	case []interface{}:
		r.Error = new(ResponseMessageAckGUIDErrors)
	}

	type tmp ResponseMessageAckGUID
	err = json.Unmarshal(data, (*tmp)(r))
	return
}
