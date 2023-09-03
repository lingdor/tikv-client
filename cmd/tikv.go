package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/lingdor/goscan"
	"github.com/lingdor/tikv-client/lib"
	"os"
	"regexp"
	"strings"
)

/*
*
tikv-client is simple command line program

tikv -pd 127.0.0.1:2379,127.0.0.2:2379
or
tikv --node 127.0.0.1:20160
*/
func main() {

	var param_pd string
	var param_exec string
	//var param_client string
	flag.StringVar(&param_pd, "pd", "", "The tikv cluster pd endpoint")
	flag.StringVar(&param_exec, "exec", "", "execute for parameter once")
	//flag.StringVar(&param_client, "node", "", "The tikv node endpoint")
	flag.Parse()

	lib.Connect(param_pd)

	if param_exec != "" {
		exec(param_exec)
		return
	}

	fmt.Println("tikv connected success!")

	for {
		if scanner, err := goscan.NewScanStd(); err != nil {
			mainErr(err)
		} else {
			cmd, _, err := scanner.Scan()
			if err != nil {
				mainErr(err)
			}

			lowerCmd := strings.ToLower(cmd)
			if lowerCmd == "get" {

				var getKey string
				if words, err := scanner.ScanWords(); err != nil {
					mainErr(err)
				} else {
					for _, word := range words {
						var valLen int
						fmt.Printf("%s: ", word)
						if valLen, err = lib.Get(getKey, os.Stdout); err != nil {
							mainErr(err)
						}
						fmt.Printf("(length:%d)\n", valLen)
					}
				}
			} else if lowerCmd == "put" {
				var putKey, putVal string

				if putKey, _, err = scanner.Scan(); err != nil {
					mainErr(err)
				}
				if putVal, _, err = scanner.Scan(); err != nil {
					mainErr(err)
				}

				putReader := bytes.NewReader([]byte(putVal))
				if err = lib.Put(putKey, putReader); err != nil {
					mainErr(err)
				}
			}
			fmt.Println("done!")
		}
	}

}

func mainErr(err error) {
	panic(err)

}

var CommandRegexp *regexp.Regexp

func init() {
	CommandRegexp = regexp.MustCompile("^(\\w+)\\s*(.*)")
}

func exec(line string) error {
	//fmt.Printf("|%s|\n", line)
	allString := CommandRegexp.FindAllString(line, 0)
	fmt.Printf("result:%+v\n", allString)
	return nil
}
