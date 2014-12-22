package main

import (
	"github.com/hawx/rivelin/models"

	"github.com/hawx/serve"

	"encoding/json"
	"flag"
	"html/template"
	"bytes"
	"io/ioutil"
	"net/http"
)

var (
	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
	river  = flag.String("river", "", "")
)

// what!
var tmpl = template.Must(template.New("tmpl").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Rivelin</title>
    <style>
      html, body {
          margin: 0;
          padding: 0;
      }

      body {
          font: 100%/1.2rem Helvetica, sans-serif;
          color: #333;
          background: #fff;
      }

      a, a:visited {
          text-decoration: none;
          color: #1E68A6;
      }
      a:hover, a:focus, a:active {
          text-decoration: underline;
      }

      body > header h1, body > .updated, body > ul {
          max-width: 36rem;
          margin: 0 auto;
      }

      body > header {
          background: #1E68A6;
          padding: .75rem 0;
          margin: 0 0 3rem;
      }
      body > header h1 {
          font-size: 1rem;
          color: #fff;
      }

      ul { list-style: none; padding: 0; }

      .block {
          clear: both;
          padding: .5rem 0 0;
          border-top: 1px solid #ccc;
          margin: 1.5rem 0 0;
      }
      .block > header h1, .block > header time {
          float: left;
          padding: 0 .5rem 0 0;
          margin: -1.3rem 0 0;
          font-size: .75rem;
          font-weight: normal;
          background: #fff;
      }
      .block > header .icon {
          position: relative;
          float: left;
          margin: 0 .5rem 0 -1.5rem;
          border: 0 none;
          vertical-align: middle;
      }
      .block > header time {
          float: right;
          padding: 0 0 0 .5rem;
          color: #777;
      }

      .item {
          clear: both;
          position: relative;
          padding: 1rem 0;
          margin: 0;
      }
      .item > header {
          margin: 0 0 .3rem;
      }
      .item > h2 {
          font-size: 1rem;
          margin: 0;
      }
      .item > p {
          font-size: 0.875rem;
          margin: .2rem 0;
      }
      .item > time {
          clear: both;
          margin: 0 1.5rem 0 0;
          font-size: .6875rem;
          color: #666;
      }
    </style>
  </head>
  <body>
    <header>
      <h1>Rivelin</h1>
    </header>

    <ul>
      {{range .UpdatedFeeds.UpdatedFeed}}
        <li class="block">
          <header>
            <h1>
              <img class="icon" src="http://www.google.com/s2/favicons?domain={{.WebsiteUrl}}" alt="">
              <a href="{{.WebsiteUrl}}">{{.FeedTitle}}</a>
              (<a href="{{.FeedUrl}}">Feed</a>)
            </h1>
            {{.WhenLastUpdate.HtmlFormat}}
          </header>
          <ul>
            {{range .Item}}
              <li class="item" id="{{.Id}}">
                <h2><a rel="external" href="{{.Link}}">{{.Title}}</a></h2>
                <p>{{.Body}}</p>
                {{.PubDate.HtmlFormat}}
              </li>
            {{end}}
          </ul>
        </li>
      {{end}}
    </ul>
  </body>
</html>
`))

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

		tmpl.Execute(w, river)
	})

	serve.Serve(*port, *socket, http.DefaultServeMux)
}
