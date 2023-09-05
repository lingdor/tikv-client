package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/lingdor/goscan"
	"github.com/lingdor/tikv-client/lib"
	"os"
	"strings"
)

/*
*
tikv-client is simple command line program

tikv -pd 127.0.0.1:2379,127.0.0.2:2379
or
tikv --node 127.0.0.1:20160
*/
var exitErr error = errors.New("exit")

func main() {

	var param_pd string
	var param_exec string
	//var param_client string
	flag.StringVar(&param_pd, "pd", "", "The tikv cluster pd endpoint")
	flag.StringVar(&param_exec, "exec", "", "execute for parameter once")
	//flag.StringVar(&param_client, "node", "", "The tikv node endpoint")
	flag.Parse()
	ctx := context.Background()

	if err := lib.Connect(param_pd); err != nil {
		mainErr(err)
	}

	if param_exec != "" {
		reader := bytes.NewReader([]byte(param_exec + "\r"))
		scanner := goscan.NewFScanner(reader)
		if err := exec(ctx, scanner); err != nil {
			mainErr(err)
		}
	}

	fmt.Println("tikv connected success!")

	for {
		scanner := goscan.NewScanner()
		if err := exec(ctx, scanner); err != nil {
			if errors.Is(err, exitErr) {
				return
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
	var isEndParam bool
	cmd, isEndParam, err := scanner.Scan()
	if err != nil {
		mainErr(err)
	} else if isEndParam {
		return nil
	}

	lowerCmd := strings.ToLower(cmd)

	if lowerCmd == "get" {

		if words, err := scanner.ScanWords(); err != nil {
			mainErr(err)
		} else {
			for _, word := range words {
				var valLen int
				fmt.Printf("%s: ", word)
				if valLen, err = lib.Get(word, os.Stdout); err != nil {
					mainErr(err)
				}
				fmt.Printf("(length:%d)\n", valLen)
			}
		}
	} else if lowerCmd == "rawget" {

		var getKey string

		if getKey, isEndParam, err = scanner.Scan(); err != nil {
			mainErr(err)
		} else if !isEndParam {
			if check, err := scanner.CheckToEnd(); err != nil {
				mainErr(err)
			} else if !check {
				return errors.New("wrong parameter")
			}
		}
		if _, err := lib.RawGet(getKey, os.Stdout); err != nil {
			mainErr(err)
		}

	} else if lowerCmd == "rawput" {

		var putKey string
		if putKey, isEndParam, err = scanner.Scan(); err != nil {
			mainErr(err)
		} else if !isEndParam {
			if check, err := scanner.CheckToEnd(); err != nil {
				mainErr(err)
			} else if !check {
				return errors.New("wrong parameter")
			}
		}
		if err := lib.RawPut(putKey, os.Stdin); err != nil {
			mainErr(err)
		}
	} else if lowerCmd == "put" {

		var putKey, putVal string

		if putKey, isEndParam, err = scanner.Scan(); err != nil {
			mainErr(err)
		} else if isEndParam {
			return errors.New("wrong parameter")
		}
		if putVal, isEndParam, err = scanner.Scan(); err != nil {
			mainErr(err)
		} else if !isEndParam {
			if putVal, isEndParam, err = scanner.Scan(); err != nil {
				mainErr(err)
			} else if !isEndParam {
				if check, err := scanner.CheckToEnd(); err != nil {
					mainErr(err)
				} else if !check {
					return errors.New("wrong parameter")
				}
			}
		}

		putReader := bytes.NewReader([]byte(putVal))
		if err = lib.Put(putKey, putReader); err != nil {
			mainErr(err)
		}
	} else if lowerCmd == "exit" {
		return exitErr
	}
	fmt.Println("done!")
	return nil
}
