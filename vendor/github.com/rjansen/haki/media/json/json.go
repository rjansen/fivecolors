package json

import (
	"encoding/json"
	"github.com/rjansen/l"
	"io"
)

const (
	//ContentType is a constant to hold the json content type value
	ContentType = "application/json"
	//ContentTypeUTF8 is a constant to hold the utf8 json content type value
	ContentTypeUTF8 = "application/json; charset=utf-8"
)

//Marshal writes a json representation of the struct instance
func Marshal(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(&data)
}

//Unmarshal reads a json representation into the struct instance
func Unmarshal(r io.Reader, result interface{}) error {
	return json.NewDecoder(r).Decode(&result)
}

//MarshalBytes writes a json representation of the struct instance
func MarshalBytes(data interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	l.Debug("json.MarshalBytes",
		l.Int("len", len(jsonBytes)),
		l.Err(err),
	)
	return jsonBytes, err
}

//UnmarshalBytes reads a json representation into the struct instance
func UnmarshalBytes(raw []byte, result interface{}) error {
	err := json.Unmarshal(raw, &result)
	l.Debug("json.UnmarshalBytes",
		l.Bool("nilResult", result == nil),
		l.Err(err),
	)
	return err
}

//Media is a struct to helps writes and reads of a json representation
type Media struct {
}

//Marshal writes a json representation of the struct instance
func (Media) Marshal(writer io.Writer, val interface{}) error {
	return Marshal(writer, &val)
}

//Unmarshal reads a json representation into the struct instance
func (Media) Unmarshal(reader io.Reader, ref interface{}) error {
	return Unmarshal(reader, &ref)
}

//MarshalBytes writes a json representation of the struct instance
func (Media) MarshalBytes(val interface{}) ([]byte, error) {
	return MarshalBytes(&val)
}

//UnmarshalBytes reads a json representation into the struct instance
func (Media) UnmarshalBytes(raw []byte, ref interface{}) error {
	return UnmarshalBytes(raw, &ref)
}
