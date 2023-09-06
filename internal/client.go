package internal

import (
	"bytes"
	"io"

	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
)

var client *tikv.RawKVClient

var names = "utf-8"

func Connect(pd string) error {
	cli, err := tikv.NewRawKVClient([]string{pd}, config.Security{})
	if err != nil {
		return err
	}
	client = cli
	return nil
}

func SetNames(name string) error {

	if name == "utf8" {
		name = "utf-8"
	}
	if _, err := ianaindex.MIB.Encoding(names); err != nil {
		return err
	}

	names = name
	return nil
}

func strToBytes(str string) ([]byte, error) {
	e, err := ianaindex.MIB.Encoding(names)
	if err != nil {
		return []byte{}, err
	}
	reader := transform.NewReader(bytes.NewReader([]byte(str)), e.NewEncoder())
	return io.ReadAll(reader)
}

func RawGet(key string, writer io.Writer) (n int, err error) {
	var bsKey, bsVal []byte
	n = 0
	if bsKey, err = strToBytes(key); err != nil {
		return
	}
	if bsVal, err = client.Get(bsKey); err != nil {
		return
	}
	n = len(bsVal)
	if _, err = writer.Write(bsVal); err != nil {
		return
	}
	return
}

func Get(key string, writer io.Writer) (int, error) {
	e, err := ianaindex.MIB.Encoding(names)
	if err != nil {
		return 0, err
	}
	transWriter := transform.NewWriter(writer, e.NewDecoder())
	return RawGet(key, transWriter)
}

func RawPut(key string, reader io.Reader) error {
	bsKey, err := strToBytes(key)
	if err != nil {
		return err
	}
	var bsVal []byte
	if bsVal, err = io.ReadAll(reader); err != nil {
		return err
	}

	if err = client.Put(bsKey, bsVal); err != nil {
		return err
	}
	return nil
}

func Put(key string, reader io.Reader) error {
	e, err := ianaindex.MIB.Encoding(names)
	if err != nil {
		return err
	}
	transReader := transform.NewReader(reader, e.NewEncoder())
	return RawPut(key, transReader)
}

func Delete(key string) error {
	bsKey, err := strToBytes(key)
	if err != nil {
		return err
	}
	if err = client.Delete(bsKey); err != nil {
		return err
	}
	return nil
}

func DeleteRange(key1, key2 string) (err error) {
	var bsKey1, bsKey2 []byte
	if bsKey1, err = strToBytes(key1); err != nil {
		return
	}
	if bsKey2, err = strToBytes(key2); err != nil {
		return
	}
	if err = client.DeleteRange(bsKey1, bsKey2); err != nil {
		return err
	}
	return nil
}
