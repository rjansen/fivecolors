package model

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalDateTime(d time.Time) graphql.Marshaler {
	return graphql.WriterFunc(
		func(w io.Writer) {
			if !d.IsZero() {
				w.Write(
					[]byte(
						fmt.Sprintf(`"%s"`, d.Format(time.RFC3339)),
					),
				)
			}
		},
	)
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {
	switch v := v.(type) {
	case string:
		return time.Parse(v, time.RFC3339)
	case []byte:
		return time.Parse(string(v), time.RFC3339)
	default:
		return time.Time{}, fmt.Errorf("err_invalid_datetimetype{type: %T}", v)
	}
}

type Object map[string]interface{}

// UnmarshalGQL implements the graphql.Marshaler interface
func (j *Object) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		return json.Unmarshal([]byte(v), j)
	case []byte:
		return json.Unmarshal(v, j)
	default:
		return fmt.Errorf("err_invalid_jsontype{type: %T}", v)
	}
}

// MarshalGQL implements the graphql.Marshaler interface
func (j Object) MarshalGQL(w io.Writer) {
	err := json.NewEncoder(w).Encode(j)
	if err != nil {
		panic(err)
	}
}
