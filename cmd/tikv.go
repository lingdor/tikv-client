package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lingdor/goscan"
	"github.com/lingdor/tikv-client/internal"
)

/*
*
tikv-client is simple command line program

tikv -pd 127.0.0.1:2379,127.0.0.2:2379
or
tikv --node 127.0.0.1:20160
*/
var errExit = errors.New("exit")
var errWrongParam = errors.New("wrong parameter")

func main() {

	var param_pd string
	var param_exec string
	//var param_client string
	flag.StringVar(&param_pd, "pd", "", "The tikv cluster pd endpoint")
	flag.StringVar(&param_exec, "exec", "", "execute for parameter once")
	//flag.StringVar(&param_client, "node", "", "The tikv node endpoint")
	flag.Parse()
	ctx := context.Background()

	if err := internal.Connect(param_pd); err != nil {
		mainErr(err)
	}

	if param_exec != "" {
		reader := bytes.NewReader([]byte(param_exec + "\r"))
		scanner := goscan.NewFLineScanner(reader)
		if err := exec(ctx, scanner); err != nil {
			mainErr(err)
		}
	}

	fmt.Println("tikv connected success!")

	for {
		scanner := goscan.NewLineScanner()
		if err := exec(ctx, scanner); err != nil {
			if errors.Is(err, errExit) {
				return
			} else if err == errWrongParam {
				fmt.Println(err.Error())
				continue
			} else if err == io.EOF {
				fmt.Println("param uncheck")
			} else {
				mainErr(err)
			}
		}
	}

}

func mainErr(err error) {
	panic(err)

}

func exec(ctx context.Context, scanner goscan.Scanner) error {
	cmd, err := scanner.Scan()
	if err != nil {
		return err
	}

	lowerCmd := strings.ToLower(cmd)

	if lowerCmd == "get" {

		if words, err := scanner.ScanWords(); err != nil {
			return err
		} else {
			for _, word := range words {
				var valLen int
				fmt.Printf("%s: ", word)
				if valLen, err = internal.Get(word, os.Stdout); err != nil {
					return err
				}
				fmt.Printf(" (len:%d)\n", valLen)
			}
		}
	} else if lowerCmd == "rawget" {

		var getKey string

		if getKey, err = scanner.Scan(); err != nil {
			return err
		}
		if err = scanner.ToEnd(); err != nil {
			return err
		}
		if _, err := internal.RawGet(getKey, os.Stdout); err != nil {
			return err
		}

	} else if lowerCmd == "rawput" {

		var putKey string
		if putKey, err = scanner.Scan(); err != nil {
			return err
		}
		if err = scanner.ToEnd(); err != nil {
			return err
		}
		if err := internal.RawPut(putKey, os.Stdin); err != nil {
			mainErr(err)
		}
	} else if lowerCmd == "put" {

		var putKey, putVal string

		if putKey, err = scanner.Scan(); err != nil {
			return err
		}
		if putVal, err = scanner.Scan(); err != nil {
			return err
		}
		if err = scanner.ToEnd(); err != nil {
			return err
		}

		putReader := bytes.NewReader([]byte(putVal))
		if err = internal.Put(putKey, putReader); err != nil {
			return err
		}
	} else if lowerCmd == "exit" {
		return errExit
	} else if lowerCmd == "set" {

		if words, err := scanner.ScanWords(); err != nil {
			return err
		} else {
			if len(words) != 2 {
				return errWrongParam
			}
			if strings.ToLower(words[0]) == "names" {
				internal.SetNames(words[1])
			} else {
				return errWrongParam
			}
		}
	}
	fmt.Println("done!")
	return nil
}
