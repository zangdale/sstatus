package sstatus

import (
	"bytes"
	"encoding/gob"
)

func Unmarshal[T any](data []byte) (res *T, err error) {
	b := new(bytes.Buffer)
	b.Write(data)
	res = new(T)
	err = gob.NewDecoder(b).Decode(res)
	return res, err
}

func Marshal(data any) (res []byte, err error) {
	b := new(bytes.Buffer)
	err = gob.NewEncoder(b).Encode(data)
	return b.Bytes(), err
}
