package main

import (
	"github.com/hawx/rivelin/models"
	"github.com/hawx/rivelin/views"

	"github.com/hawx/serve"

	"encoding/json"
	"flag"
	"bytes"
	"io/ioutil"
	"net/http"
)

var (
	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
	river  = flag.String("river", "", "")
)

const (
	PREFIX = "onGetRiverStream("
	SUFFIX = ")"
)

func main() {
	flag.Parse()
	if *river == "" {
		println("err: --river must be given")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(*river)
		if err != nil {
			w.WriteHeader(500) // could not connect
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(500) // could not parse, should have nice msg!?!
		}

		// hate this. maybe, func TrimPrefix(io.Reader, string) io.Reader ???
		data = bytes.TrimSuffix(bytes.TrimPrefix(data, []byte(PREFIX)), []byte(SUFFIX))

		var river models.River
		json.Unmarshal(data, &river)

		views.List.Execute(w, river)
	})

	serve.Serve(*port, *socket, http.DefaultServeMux)
}
