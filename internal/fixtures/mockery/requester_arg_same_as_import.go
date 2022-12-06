package mockery

import "encoding/json"

type RequesterArgSameAsImport interface {
	Get(json string) *json.RawMessage
}
