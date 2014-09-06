package main

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

var indexTemplateRaw = `
<!DOCTYPE html>
<html lang="en">
<head>
	<title>aliases</title>
	{{range .stylesheets}}
	<link rel="stylesheet" href={{.}}>
	{{end}}
</head>
<body>
	<ul>
		{{range .aliases}}
		<li>
		<a href={{aliasHRef .}}>{{.Name}}</a>
		<span>{{.Description}}</span>
		</li>
		{{end}}
	</ul>
</body>
</html>
`

type index struct {
	sites    []*Site
	css      []string
	template *template.Template
	p        []byte
}

func aliasHRef(s *Site) (string, error) {
	u, err := url.Parse(s.Bind)
	if err != nil {
		return "", err
	}
	if u.Scheme != "" {
		return s.Bind, nil
	}
	if u.Host != "" {
		return s.Bind, nil
	}
	ustr := strings.TrimPrefix(s.Bind, "/")
	return ustr, nil
}

func (x *index) compile() ([]byte, error) {
	var err error
	x.template, err = template.New("alias-index").
		Funcs(template.FuncMap{"aliasHRef": aliasHRef}).
		Parse(indexTemplateRaw)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	x.p = nil
	err = x.template.Execute(&buf, map[string]interface{}{
		"stylesheets": x.css,
		"aliases":     x.sites,
	})
	if err != nil {
		return nil, err
	}
	x.p = buf.Bytes()
	return x.p, nil
}

func Index(s []*Site, css []string) (http.Handler, error) {
	x := new(index)
	x.sites = s
	x.css = css
	p, err := x.compile()
	if err != nil {
		return nil, err
	}
	h := func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		w.Write(p)
	}
	return http.HandlerFunc(h), nil
}
