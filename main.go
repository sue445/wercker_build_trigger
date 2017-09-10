package main

import (
	"flag"
	"fmt"
)

var (
	Version  string
	Revision string
)
var configFile, token string
var isPrintVersion bool

func init() {
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.StringVar(&token, "token", "", "API token")
	flag.BoolVar(&isPrintVersion, "version", false, "Whether showing version")

	flag.Parse()
}

func main() {
	if isPrintVersion {
		printVersion()
		return
	}

	if len(configFile) == 0 || len(token) == 0 {
		flag.PrintDefaults()
		return
	}

}

func printVersion() {
	fmt.Printf("zatsu_monitor v%s, build %s\n", Version, Revision)
}
