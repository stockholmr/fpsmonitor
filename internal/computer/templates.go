package computer

import "html/template"

func editorPage() *template.Template {
	return template.Must(template.New("page").Delims("<<", ">>").Parse(`
		<!DOCTYPE html>
		<html lang="en">
			<head>
				<meta charset="utf-8" />
				<meta language="english" />
				<meta http-equiv="X-UA-Compatible" content="IE=edge">
				<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" />
				<title><<.Title>></title>
                <link rel="stylesheet" href="/bootstrap" type="text/css" />
				<link rel="stylesheet" href="/computers/stylesheet" type="text/css" />
			</head>

			<body>
				<div class="container">
					<table class="table table-dark">
						<thead>
							<tr>
								<th scope="col">#</th>
								<th scope="col">Date</th>
								<th scope="col">ComputerName</th>
								<th scope="col">Username</th>
							</tr>
						</thead>
						<tbody>
							<<range .Records>>
								<tr>
									<td></td>
									<td>
										<< .Created.String >>
									</td>
									<td>
										<< .ComputerName.String >>
									</td>
									<td>
										<< .Username.String >>
									</td>
								</tr>
							<<end>>
						</tbody>
					</table>
				</div>
			</body>
		</html>
	`))
}
