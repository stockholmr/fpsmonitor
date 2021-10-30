module fpsmonitor

go 1.17

require (
	github.com/gorilla/mux v1.8.0
	github.com/jcelliott/lumber v0.0.0-20160324203708-dd349441af25
	github.com/jmoiron/sqlx v1.3.4
	github.com/justinas/alice v1.2.0
	github.com/mattn/go-sqlite3 v1.14.9
	gopkg.in/guregu/null.v3 v3.5.0
)

replace fpsmonitor => ./
