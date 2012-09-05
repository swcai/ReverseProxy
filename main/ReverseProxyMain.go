package main

import (
	"flag"
	"fmt"
	"stanleycai.com/rproxy"
	"runtime"
	"strings"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var port *int = flag.Int("p", 8080, "The port for ReverseProxy.")
var clientStr *string = flag.String("c", "", "The client list for ReverseProxy")

// This is an oversimplified client list implementation. Only in 0.x version

func parseClientString(cs string) []string {
	strings.Fields(";")
	return nil
}

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	addr := fmt.Sprintf("127.0.0.1:%v", *port)
	// start main loop
	if cs := strings.Split(*clientStr, ";"); len(cs) != 0 {
		rproxy.ListenAndServe(addr, cs)
	}
}
