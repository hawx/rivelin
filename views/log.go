package views

const log = pre + `<ul class="blocks">
  {{ range . }}
  <li class="block">
    <header class="block-title">
      <h1><a href="#">{{.Header}}</a></h1>
    </header>
    <ul class="items">
      {{ range .Items }}
      <li class="item">
        <h2><a href="{{.URI}}">{{.URI}}</a> <span class="code {{.Status}}">{{.Code}}</span></h2>
      </li>
      {{ end }}
    </ul>
  </li>
  {{ end }}
</ul> ` + post
