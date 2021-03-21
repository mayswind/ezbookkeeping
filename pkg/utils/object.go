package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
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

// PrintObjectFields prints all fields in specified object
func PrintObjectFields(obj interface{}) {
	if obj == nil {
		return
	}

	elem := reflect.ValueOf(obj).Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fmt.Printf("[%s] %v\n", typ.Field(i).Name, field.Interface())
	}
}
