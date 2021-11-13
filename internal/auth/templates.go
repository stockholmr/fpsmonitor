package auth

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
				<link rel="stylesheet" href="/bootstrap?ver=5.1.3" type="text/css" />
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

func (t *Templates) Login(w http.ResponseWriter, data TemplateData) error {
	tmpl, err := t.page().New("content").Delims(t.LeftDelim, t.RightDelim).Parse(`
		<div id="login">
			<form action="/login" method="POST">

				<< if (ne .Error "") >>
					<div class="alert alert-danger" role="alert">
						<< .Error >>
					</div>
				<< end >>

				<div class="form-group">
					<label class="sr-only" for="username">Username</label>
					<input type="text" class="form-control form-control-sm" id="username" placeholder="Username" name="username" />
				</div>
				<div class="form-group mb-2">
					<label class="sr-only" for="password">Password</label>
					<input type="password" class="form-control form-control-sm" id="password" placeholder="Password" name="password" />
				</div>

				<div class="checkbox mb-3">
					<label id="remember">
						<input type="checkbox" name="remember" value="true"> Remember me
					</label>
				</div>

				<button type="submit" class="btn btn-primary btn-sm">Login</button>

			</form>
		</div>
	`)

	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "page", t.mergeData(t.Data, data))
}

func (t *Templates) Register(w http.ResponseWriter, data TemplateData) error {
	tmpl, err := t.page().New("content").Delims(t.LeftDelim, t.RightDelim).Parse(`
		<div id="register">
			<form action="/register" method="POST">

				<div class="form-group">
					<label class="sr-only" for="username">Username</label>
					<input type="text" class="form-control form-control-sm" id="username" placeholder="Username" name="username" />
				</div>
				<div class="form-group mb-3">
					<label class="sr-only" for="password">Password</label>
					<input type="password" class="form-control form-control-sm" id="password" placeholder="Password" name="password" />
				</div>

				<button type="submit" class="btn btn-primary btn-sm">Register</button>

			</form>
		</div>
	`)

	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "page", t.mergeData(t.Data, data))
}
