package mockery

import "encoding/json"

type RequesterArgSameAsNamedImport interface {
	Get(json string) *json.RawMessage
}
