package docs

import (
	"bytes"
	"io/ioutil"
	"sort"
	"text/template"
)

func buildMd(patterns []pattern, filename string) (string, error) {
	parse, err := template.New("md").Parse(mdTemplateStr)
	if err != nil {
		return "", err
	}
	var outputFile = filename
	if outputFile == "" {
		outputFile = defaultFilename
	}
	o := bytes.NewBufferString("")
	sort.Slice(patterns, func(i, j int) bool {
		a := patterns[i]
		b := patterns[j]
		if a.Kind > b.Kind {
			return true
		} else if a.Kind == b.Kind {
			return a.Name < b.Name
		} else {
			return false
		}
	})
	if err := parse.Execute(o, patterns); err != nil {
		return "", err
	}
	all, err := ioutil.ReadAll(o)
	if err != nil {
		return "", err
	}
	return string(all), ioutil.WriteFile(outputFile+".md", all, 0666)
}

const mdTemplateStr = `
| Name         | Kind  |   Type  | Description |
|:------------:|:-----:| :-----: | ----------- |
{{- range . }}
| {{ .Name }}   | {{ .Kind }} | {{ .Typ }} | {{ .Description }} |
{{- end -}}
`
