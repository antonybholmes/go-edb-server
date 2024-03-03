module github.com/antonybholmes/go-edb-api

go 1.22.0

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-genes => ../go-genes

replace github.com/antonybholmes/go-math => ../go-math

replace github.com/antonybholmes/go-auth => ../go-auth

replace github.com/antonybholmes/go-sys => ../go-sys

replace github.com/antonybholmes/go-mailer => ../go-mailer

require (
	github.com/antonybholmes/go-dna v0.0.0-20240220233159-e0a18c04f799
	github.com/antonybholmes/go-math v0.0.0-20240215163921-12bb7e52185c
	github.com/labstack/echo/v4 v4.11.4
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20240220233159-3c0eb135aebe
	github.com/antonybholmes/go-genes v0.0.0-20240220233158-15b6002680ec
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rs/zerolog v1.32.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20240219230546-a8e3e894e0d7
	github.com/antonybholmes/go-sys v0.0.0-20240219230548-9ab0febd5fc5
	github.com/gorilla/sessions v1.2.2
	github.com/labstack/echo-contrib v0.15.0
	github.com/michaeljs1990/sqlitestore v0.0.0-20210507162135-8585425bc864
)

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xyproto/randomstring v1.0.5 // indirect
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
