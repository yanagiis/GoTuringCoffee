package lib

***REMOVED***
	"strings"

	jsoniter "github.com/json-iterator/go"
	nats "github.com/nats-io/go-nats"
***REMOVED***

// JsonEncoder is a JSON Encoder implementation for EncodedConn.
// This encoder will use the builtin encoding/json to Marshal
// and Unmarshal most types, including structs.
type JsoniterEncoder struct {
	// Empty
***REMOVED***

func init(***REMOVED*** {
	nats.RegisterEncoder("jsoniter", &JsoniterEncoder{***REMOVED******REMOVED***
***REMOVED***

// Encode
func (je *JsoniterEncoder***REMOVED*** Encode(subject string, v interface{***REMOVED******REMOVED*** ([]byte, error***REMOVED*** {
	b, err := jsoniter.Marshal(v***REMOVED***
***REMOVED***
		return nil, err
***REMOVED***
	return b, nil
***REMOVED***

// Decode
func (je *JsoniterEncoder***REMOVED*** Decode(subject string, data []byte, vPtr interface{***REMOVED******REMOVED*** (err error***REMOVED*** {
	switch arg := vPtr.(type***REMOVED*** {
	case *string:
		// If they want a string and it is a JSON string, strip quotes
		// This allows someone to send a struct but receive as a plain string
		// This cast should be efficient for Go 1.3 and beyond.
		str := string(data***REMOVED***
		if strings.HasPrefix(str, `"`***REMOVED*** && strings.HasSuffix(str, `"`***REMOVED*** {
			*arg = str[1 : len(str***REMOVED***-1]
	***REMOVED*** else {
			*arg = str
	***REMOVED***
	case *[]byte:
		*arg = data
	default:
		err = jsoniter.Unmarshal(data, arg***REMOVED***
***REMOVED***
	return
***REMOVED***
