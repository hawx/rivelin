package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"hawx.me/code/rivelin/models"
	"hawx.me/code/serve"
)

func printHelp() {
	fmt.Println(`Usage: rivelin [options]

  Rivelin is a web reader for riverjs (http://riverjs.org) files.

   --river URL
      Full URL to the riverjs file to read from.

   --timezone TZ='UTC'
      Display datetimes using this timezone.

   --with-log
      Also serve a log of feed reading activity at '/log'. Will probably only
      work when reading from a riviera generated feed.

   --web PATH
      Path to the 'web' directory.

   --port PORT='8080'
      Serve on given port.

   --socket SOCK
      Serve at given socket, instead.`)
}

func main() {
	var (
		port     = flag.String("port", "8080", "")
		socket   = flag.String("socket", "", "")
		river    = flag.String("river", "", "")
		timezone = flag.String("timezone", "UTC", "")
		withLog  = flag.Bool("with-log", false, "")
		webPath  = flag.String("web", "web", "")
	)
	flag.Usage = printHelp
	flag.Parse()

	if *river == "" {
		fmt.Println("err: --river must be given")
		return
	}

	loc, err := time.LoadLocation(*timezone)
	if err != nil || loc == nil {
		fmt.Println("err: --timezone invalid")
		return
	}

	riverURL, err := url.Parse(*river)
	if err != nil {
		fmt.Println("err: --river invalid")
		return
	}

	templates, err := template.ParseGlob(*webPath + "/template/*.gotmpl")
	if err != nil {
		fmt.Println(err)
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
			return
		}

		if err = templates.ExecuteTemplate(w, "list.gotmpl", river.SetLocation(*loc)); err != nil {
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

			if err = templates.ExecuteTemplate(w, "log.gotmpl", models.MakeLogBlocks(logList)); err != nil {
				log.Println("/log", err)
			}
		})
	}

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir(*webPath+"/static"))))

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      http.DefaultServeMux,
	}

	serve.Server(*port, *socket, srv)
}
