package views

const list = pre + `<ul class="blocks">
  {{range .UpdatedFeeds.UpdatedFeed}}
  <li class="block">
    <header class="block-title">
      <h1>
        <img class="icon" src="//www.google.com/s2/favicons?domain={{.WebsiteUrl}}" alt="">
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
</ul>` + post
