package lib

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
)

var client *tikv.RawKVClient

// var encoder transform.Transformer = encoding.Encoder{}
// var decoder transform.Transformer = encoding.Decoder{}
var encoder transform.Transformer = simplifiedchinese.GBK.NewEncoder()
var decoder transform.Transformer = simplifiedchinese.GBK.NewDecoder()

func Connect(pd string) error {
	cli, err := tikv.NewRawKVClient([]string{pd}, config.Security{})
	if err != nil {
		return err
	}
	client = cli
	return nil
}

func SetNames(name string) error {
	switch name {
	case "utf8":
		encoder = encoding.Encoder{}
		decoder = encoding.Decoder{}
	case "gbk":
		encoder = simplifiedchinese.GBK.NewEncoder()
		decoder = simplifiedchinese.GBK.NewDecoder()
	default:
		return errors.New(fmt.Sprintf("no found names :%s", name))
	}
	return nil
}

func strToBytes(str string) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(str)), encoder)
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
	transWriter := transform.NewWriter(writer, decoder)
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
	transReader := transform.NewReader(reader, encoder)
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
