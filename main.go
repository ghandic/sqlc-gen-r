package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"strings"
	"text/template"

	"github.com/sqlc-dev/plugin-sdk-go/codegen"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

//go:embed r.go.tmpl
var templateFS embed.FS

func main() {
	codegen.Run(generate)
}

type Options struct {
	Filename string `json:"filename" yaml:"filename"`
	Out      string `json:"out" yaml:"out"`
}

func parseOpts(req *plugin.GenerateRequest) (*Options, error) {
	var options Options
	if len(req.PluginOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.PluginOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling plugin options: %w", err)
	}

	return &options, nil
}

func generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	// fmt.Println(req)
	options, _ := parseOpts(req)

	pluginOptions := make(map[string]any)
	err := json.Unmarshal(req.PluginOptions, &pluginOptions)
	if err != nil {
		log.Fatal("failed to unmarshal plugin options: ", err)
	}

	funcMap := template.FuncMap{
		"Contains": strings.Contains,
		// https://stackoverflow.com/a/18276968/1149933
		"Dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dict call")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"GetPluginOption": func(name string) any {
			option, ok := pluginOptions[name]
			if !ok {
				return ""
			}
			return option
		},
		"Split":   strings.Split,
		"ToLower": strings.ToLower,
	}

	// Read the template from the embedded file system
	templateContent, err := fs.ReadFile(templateFS, "r.go.tmpl")
	if err != nil {
		log.Fatalf("Error reading embedded template file: %v", err)
	}

	// Parse the template content
	tmpl, err := template.New("embedded_template").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		log.Fatalf("Error parsing embedded template: %v", err)
	}

	resp := plugin.GenerateResponse{}
	for i := range req.Queries {
		paramMap := make(map[string]int)
		for j := range req.Queries[i].Params {
			colName := req.Queries[i].Params[j].Column.Name
			val, ok := paramMap[colName]
			if !ok {
				paramMap[colName] = 1
				continue
			}
			paramMap[colName] = val + 1
			req.Queries[i].Params[j].Column.Name = colName + fmt.Sprintf("%v", val)
		}
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, req)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	resp.Files = append(resp.Files, &plugin.File{
		Name:     options.Filename,
		Contents: buf.Bytes(),
	})

	return &resp, nil
}
