package main

import (
	"hawx.me/code/rivelin/models"
	"hawx.me/code/rivelin/views"

	"hawx.me/code/serve"

	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	port     = flag.String("port", "8080", "")
	socket   = flag.String("socket", "", "")
	river    = flag.String("river", "", "")
	timezone = flag.String("timezone", "UTC", "")
	withLog  = flag.Bool("with-log", false, "")
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

	loc, err := time.LoadLocation(*timezone)
	if err != nil || loc == nil {
		println("err: --timezone invalid")
		return
	}

	riverURL, err := url.Parse(*river)
	if err != nil {
		println("err: --river invalid")
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

		views.List.Execute(w, river.SetLocation(*loc))
	})

	if *withLog {
		http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
			logURL, _ := riverURL.Parse("log")
			resp, err := http.Get(logURL.String())
			if err != nil {
				w.WriteHeader(500)
			}

			var logList []models.LogLine
			err = json.NewDecoder(resp.Body).Decode(&logList)

			views.Log.Execute(w, models.MakeLogBlocks(logList))
		})
	}

	serve.Serve(*port, *socket, http.DefaultServeMux)
}
