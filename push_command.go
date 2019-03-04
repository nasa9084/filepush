package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nasa9084/errors"
)

type PushCommand struct {
	Filename string `short:"f" long:"file" required:"yes"`
	To       string `short:"t" long:"to" required:"yes"`
}

func (cmd *PushCommand) Execute([]string) error {
	log.Printf("push %s to %s", cmd.Filename, cmd.To)
	f, err := os.Open(cmd.Filename)
	if err != nil {
		return errors.Wrap(err, "opening file")
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return errors.Wrap(err, "buffering file content")
	}
	f.Close()
	if !strings.HasPrefix(cmd.To, "http://") && !strings.HasPrefix(cmd.To, "https://") {
		cmd.To = "http://" + cmd.To
	}
	req, err := http.NewRequest(
		http.MethodPost,
		cmd.To,
		&buf,
	)
	if err != nil {
		return errors.Wrap(err, "")
	}
	req.Header.Set("X-FILENAME", cmd.Filename)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "requesting")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}
	log.Print("done")
	return nil
}
