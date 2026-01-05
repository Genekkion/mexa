package test

import (
	"encoding/json"
	"mexa/internal/utils"
)

func JsonB(v any) []byte {
	return utils.Must(func() ([]byte, error) {
		return json.Marshal(v)
	})
}

func JsonS(v any) string {
	return string(JsonB(v))
}
