package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nasa9084/errors"
)

type ListenCommand struct {
	Port      int    `short:"p" long:"port" default:"18080"`
	Directory string `short:"d" long:"directory" default:"."`
}

func (cmd *ListenCommand) Execute([]string) error {
	http.HandleFunc(`/`, cmd.listen)
	log.Printf("listen on port %d", cmd.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cmd.Port), nil); err != nil {
		return errors.Wrap(err, "listening")
	}
	return nil
}

const openFileFlg = os.O_CREATE | os.O_RDWR | os.O_TRUNC

func (cmd *ListenCommand) listen(w http.ResponseWriter, r *http.Request) {
	filename := r.Header.Get("X-FILENAME")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("receive file: %s", filename)
	f, err := os.OpenFile(filepath.Join(cmd.Directory, filename), openFileFlg, 0644)
	log.Printf("save file to %s", filepath.Join(cmd.Directory, filename))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Print(err.Error())
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		log.Print(err.Error())
		return
	}
	log.Printf("finished saving")
	w.WriteHeader(http.StatusOK)
}
