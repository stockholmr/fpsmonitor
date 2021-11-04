package auth

import "html/template"

func Page() *template.Template {
	return template.Must(template.New("page").Delims("<<", ">>").Parse(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8" />
				<meta language="english" />
				<meta http-equiv="X-UA-Compatible" content="IE=edge">
				<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" />
				<title><<.title>></title>
                <link rel="stylesheet" href="/bootstrap" type="text/css" />
			</head>

			<body>
                <<template "content" .>>
			</body>
		</html>
	`))
}

func Register() *template.Template {
	tmpl := template.Must(Page().New("content").Delims("<<", ">>").Parse(`
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
	`))
	return tmpl
}

func login() *template.Template {
	tmpl := template.Must(Page().New("content").Delims("<<", ">>").Parse(`
	<div id="login">
        <form action="/login" method="POST">

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
	`))
	return tmpl
}
