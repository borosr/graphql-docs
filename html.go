package docs

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strings"
	"text/template"
)

func buildHtml(templates []pattern, filename string) (string, error) {
	parse, err := template.New("html").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).Parse(htmlTemplateStr)
	if err != nil {
		return "", err
	}
	var outputFile = filename
	if outputFile == "" {
		outputFile = defaultFilename
	}
	o := bytes.NewBufferString("")
	sort.Slice(templates, func(i, j int) bool {
		a := templates[i]
		b := templates[j]
		if a.Kind > b.Kind {
			return true
		} else if a.Kind == b.Kind {
			return a.Name < b.Name
		} else {
			return false
		}
	})
	if err := parse.Execute(o, templates); err != nil {
		return "", err
	}
	all, err := ioutil.ReadAll(o)
	if err != nil {
		return "", err
	}
	return string(all), ioutil.WriteFile(outputFile+".html", all, 0666)
}

const htmlTemplateStr = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>GraphQL Documentation</title>
  <meta name="description" content="GraphQL Documentation">
  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Open+Sans">
  <style>
  body {
	font-family: 'Open Sans', serif;
    color: darkslategray;
  }
  ul {
    list-style-type: none;
  }
  hr {
    opacity: .5;
    width: 90vw;
  }
  a, a:link, a:visited, a:active {
    color: darkslategray;
    text-decoration: none;
  }
  .kind {
	font-size: 1.2rem;
    padding: 0 .6rem;
    border-radius: .7rem;
	margin: 0 1rem;
    color: darkslategray;
    text-transform: uppercase;
  }
  .kind.field {
	background-color: orange;
  }
  .kind.argument {
    background-color: yellow;
  }
  .type {
    text-transform: uppercase;
	font-size: 1.2rem;
    padding: 0 .6rem;
    border-radius: .7rem;
	margin: 0 1rem;
    color: rgba(246, 247, 233, 1);
    background-color: rgba(106, 195, 96, 1);
  }
  .type.String {
	background-color: rgba(96, 140, 195, 1);
  }
  .type.Int {
	background-color: rgba(117, 194, 185, 1);
  }
  .type.Float {
	background-color: rgba(171, 117, 194, 1);
  }
  .type.Boolean {
	background-color: rgba(193, 206, 131, 1);
  }
  .description {
    
  }
  </style>
</head>
<body>
	<ul>
	{{- range . }}
		<li {{ if eq .Kind "field" -}}id="field-{{ .Name | ToLower }}"{{- end -}}>
			<h2>{{- .Name -}}
			{{- if .Parent -}}
			| <a href="#field-{{- .Parent -}}">#{{- .Parent -}}</a>
			{{- end -}}
			</h2>
			<div>
				<span class="kind {{ .Kind }}">{{- .Kind -}}</span>
				{{ if .Typ -}}<span class="type {{ .Typ }}">{{- .Typ -}}</span>{{- end}}
			</div>
			<p class="description">{{ .Description }}</p>
			<hr>
		</li>
	{{- end -}}
	</ul>
</body>
</html>
`
