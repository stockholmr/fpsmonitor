package computer

import (
	"html/template"
	"log"
	"net/http"
)

type TemplateData map[string]interface{}

type Templates struct {
	LeftDelim  string
	RightDelim string
	Data       TemplateData
}

func InitTemplates() *Templates {
	return &Templates{
		LeftDelim:  "<<",
		RightDelim: ">>",
		Data:       make(TemplateData),
	}
}

func (t *Templates) SetData(key string, value interface{}) {
	t.Data[key] = value
}

func (t *Templates) mergeData(data ...TemplateData) TemplateData {
	newData := make(TemplateData)

	for i := range data {
		for k, v := range data[i] {
			newData[k] = v
		}
	}

	return newData
}

func (t *Templates) page() *template.Template {
	tmpl, err := template.New("page").Delims(t.LeftDelim, t.RightDelim).Parse(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8" />
				<meta language="english" />
				<meta http-equiv="X-UA-Compatible" content="IE=edge">
				<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" />
				<title><< .Title >></title>
				<link rel="stylesheet" href="/bootstrap" type="text/css" />
				<link rel="stylesheet" href="/vuejsdev" type="text/css" />
				<style type="text/css"><< .Styles >></style>
			</head>

			<body>
				<<template "content" .>>
			</body>
		</html>
	`)

	if err != nil {
		log.Panic(err)
	}

	return tmpl
}

func (t *Templates) List(w http.ResponseWriter, data TemplateData) error {
	tmpl, err := t.page().New("content").Delims(t.LeftDelim, t.RightDelim).Parse(`
		<< print .Computers >>
	`)

	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "page", t.mergeData(t.Data, data))
}
