package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"hawx.me/code/rivelin/models"
	"hawx.me/code/rivelin/views"
	"hawx.me/code/serve"
)

func printHelp() {
	fmt.Println(`Usage: rivelin [options]

  Rivelin is a web reader for riverjs (http://riverjs.org) files.

 DISPLAY
   --river URL
      Full URL to the riverjs file to read from.

   --timezone TZ='UTC'
      Display datetimes using this timezone.

   --with-log
      Also serve a log of feed reading activity at '/log'. Will probably only
      work when reading from a riviera generated feed.

 SERVE
   --port PORT='8080'
      Serve on given port.

   --socket SOCK
      Serve at given socket, instead.
`)
}

func main() {
	var (
		port     = flag.String("port", "8080", "")
		socket   = flag.String("socket", "", "")
		river    = flag.String("river", "", "")
		timezone = flag.String("timezone", "UTC", "")
		withLog  = flag.Bool("with-log", false, "")
	)

	flag.Usage = printHelp
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

	// use default values from DefaultTransport
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := httpClient.Get(*river)
		if err != nil {
			log.Println("/", err)
			w.WriteHeader(502)
			return
		}
		defer resp.Body.Close()

		bufferedBody := bufio.NewReader(resp.Body)
		bufferedBody.ReadBytes('(')

		var river models.River
		if err = json.NewDecoder(bufferedBody).Decode(&river); err != nil {
			log.Println("/", err)
		}

		if err = views.List.Execute(w, river.SetLocation(*loc)); err != nil {
			log.Println("/", err)
		}
	})

	if *withLog {
		http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
			logURL, _ := riverURL.Parse("log")

			resp, err := httpClient.Get(logURL.String())
			if err != nil {
				log.Println("/log", err)
				w.WriteHeader(502)
				return
			}
			defer resp.Body.Close()

			var logList []models.LogLine
			if err = json.NewDecoder(resp.Body).Decode(&logList); err != nil {
				log.Println("/log", err)
				w.WriteHeader(500)
				return
			}

			views.Log.Execute(w, models.MakeLogBlocks(logList))
		})
	}

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      http.DefaultServeMux,
	}

	serve.Server(*port, *socket, srv)
}
