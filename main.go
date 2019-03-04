package main

import (
	"log"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	Push   PushCommand   `command:"push"`
	Listen ListenCommand `command:"listen"`
}

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		if fe, ok := err.(*flags.Error); !ok || fe.Type != flags.ErrHelp {
			log.Printf("error: %v", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
