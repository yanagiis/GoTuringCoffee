package lib

import (
	"strings"

	jsoniter "github.com/json-iterator/go"
	nats "github.com/nats-io/nats.go"
)

// JsonEncoder is a JSON Encoder implementation for EncodedConn.
// This encoder will use the builtin encoding/json to Marshal
// and Unmarshal most types, including structs.
type JsoniterEncoder struct {
	// Empty
}

func init() {
	nats.RegisterEncoder("jsoniter", &JsoniterEncoder{})
}

// Encode
func (je *JsoniterEncoder) Encode(subject string, v interface{}) ([]byte, error) {
	b, err := jsoniter.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Decode
func (je *JsoniterEncoder) Decode(subject string, data []byte, vPtr interface{}) (err error) {
	switch arg := vPtr.(type) {
	case *string:
		// If they want a string and it is a JSON string, strip quotes
		// This allows someone to send a struct but receive as a plain string
		// This cast should be efficient for Go 1.3 and beyond.
		str := string(data)
		if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
			*arg = str[1 : len(str)-1]
		} else {
			*arg = str
		}
	case *[]byte:
		*arg = data
	default:
		err = jsoniter.Unmarshal(data, arg)
	}
	return
}
