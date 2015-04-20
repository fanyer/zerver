package components

import (
	"io"
	"io/ioutil"

	"github.com/cosiner/zerver"
	"github.com/pquerna/ffjson/ffjson"
)

// Ffjson is a implementation of zerver.ResourceMaster
type Ffjson struct{}

func NewFfjsonResource() zerver.ResourceMaster {
	return Ffjson{}
}

func (Ffjson) Init(zerver.Enviroment) error          { return nil }
func (Ffjson) Destroy()                              {}
func (Ffjson) Marshal(v interface{}) ([]byte, error) { return ffjson.Marshal(v) }
func (Ffjson) Pool(data []byte)                      { ffjson.Pool(data) }

func (Ffjson) Unmarshal(data []byte, v interface{}) error { return ffjson.Unmarshal(data, v) }

func (Ffjson) Send(w io.Writer, key string, value interface{}) error {
	if key == "" {
		data, err := ffjson.Marshal(value)
		if err == nil {
			_, err = w.Write(data)
			ffjson.Pool(data)
		}
		return err
	}

	w.Write(zerver.JSONObjStart)
	w.Write(zerver.Bytes(key))

	if s, is := value.(string); is {
		w.Write(zerver.JSONQuoteMid)
		w.Write(zerver.Bytes(s))
		_, err := w.Write(zerver.JSONQuoteEnd)
		return err
	}

	data, err := ffjson.Marshal(value)
	if err == nil {
		w.Write(zerver.JSONObjMid)
		w.Write(data)
		_, err = w.Write(zerver.JSONObjEnd)
		ffjson.Pool(data)
	}
	return err
}

func (Ffjson) Receive(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err == nil {
		err = ffjson.Unmarshal(data, v)
	}
	return err
}