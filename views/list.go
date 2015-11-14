package views

import "html/template"

var List = template.Must(template.New("list").Parse(list))

const list = `<!DOCTYPE html>
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
          font: 14px/1.3 Verdana, Geneva, sans-serif;
          color: #000;
          background: #fff;
      }

      a, a:visited {
          text-decoration: none;
          color: #365da9;
      }
      a:hover, a:focus, a:active {
          text-decoration: underline;
          color: #2a6497;
      }

      .container {
          max-width: 40em;
          margin: 0 auto;
          padding: 0 1em;
      }
      .container:before, .container:after {
          clear: both;
          content: " ";
          display: table;
      }

      .page-title {
          background: #eee;
          border-bottom: 1px solid #ddd;
          padding: 0;
          margin: 0;
      }
      .page-title h1 {
          font-size: 1.5em;
          padding: 1.3rem;
          margin: 0;
          height: 1.3rem;
          line-height: 1.3rem;
          display: inline-block;
          padding-left: 0;
          font-weight: bold;
      }

      ul { list-style: none; padding: 0; }

      .blocks {
          width: auto;
          margin: 2.6rem 0;
      }

      .block {
          clear: both;
          padding: .5rem 0 0;
          border-top: 1px solid #ddd;
          margin: 1.1rem 0 0;
      }
      .block-title h1, .block-title time {
          float: left;
          padding: 0 .5rem 0 0;
          margin: -1.1rem 0 0;
          font-size: .75rem;
          font-weight: normal;
          background: #fff;
      }
      .block-title .icon {
          position: relative;
          float: left;
          margin: 0 .5rem 0 -1.5rem;
          border: 0 none;
          vertical-align: middle;
      }
      .block-title time {
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
      .item header {
          margin: 0 0 .3rem;
      }
      .item h2 {
          font-size: 1rem;
          margin: 0;
      }
      .item p {
          font-size: 0.875rem;
          margin: .2rem 0;
      }
      .item .timea {
          clear: both;
          margin: 0 1.5rem 0 0;
          font-size: .6875rem;
          color: #666;
      }

      footer {
          text-align: center;
          padding-bottom: 3rem;
          font-size: .6875rem;
          color: #bbb;
      }
      footer a, footer a:hover, footer a:visited, footer a:focus, footer a:active {
          color: #bbb;
          text-decoration: underline;
      }

      @media screen and (max-width: 40rem) {
          .block-title .icon, .block-title .feed { display: none; }
      }
    </style>
  </head>
  <body>
    <div class="container">
      <ul class="blocks">
        {{range .UpdatedFeeds.UpdatedFeed}}
        <li class="block">
          <header class="block-title">
            <h1>
              <img class="icon" src="http://www.google.com/s2/favicons?domain={{.WebsiteUrl}}" alt="">
              <a href="{{.WebsiteUrl}}">{{.FeedTitle}}</a>
              <span class="feed">(<a href="{{.FeedUrl}}">Feed</a>)</span>
            </h1>
            {{.WhenLastUpdate.HtmlFormat}}
          </header>
          <ul class="items">
            {{range .Item}}
            <li class="item" id="{{.Id}}">
              <h2><a rel="external" href="{{.Link}}">{{.Title}}</a></h2>
              <p>{{.FilteredBody}}</p>
              <a class="timea" rel="external" href="{{.Link}}">{{.PubDate.HtmlFormat}}</a>
            </li>
            {{end}}
          </ul>
        </li>
        {{end}}
      </ul>

      <footer>
        <a href="http://hawx.me/code/rivelin">rivelin</a> + <a href="http://hawx.me/code/riviera">riviera</a>
      </footer>
    </div>
  </body>
</html>`
