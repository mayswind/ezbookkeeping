package utils

import (
	"bytes"
	"encoding/gob"
)

// Clone deep-clones src object to dst object
func Clone(src, dst interface{}) error {
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}

	err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)

	if err != nil {
		return err
	}

	return nil
}
