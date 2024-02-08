module github.com/antonybholmes/go-edb-api

go 1.22.0

replace github.com/antonybholmes/go-loctogene => ../go-loctogene

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-annotation => ../go-annotation

require (
	github.com/antonybholmes/go-dna v0.0.0-20240202000110-e2f8eb94d428
	github.com/antonybholmes/go-loctogene v0.0.0-20240201225030-d8a0367b8e0f
	github.com/labstack/echo/v4 v4.11.4
	github.com/labstack/gommon v0.4.2
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
