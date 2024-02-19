module github.com/antonybholmes/go-edb-api

go 1.22.0

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-genes => ../go-genes

replace github.com/antonybholmes/go-math => ../go-math

replace github.com/antonybholmes/go-auth => ../go-auth

replace github.com/antonybholmes/go-env => ../go-env

replace github.com/antonybholmes/go-mailer => ../go-mailer

require (
	github.com/antonybholmes/go-dna v0.0.0-20240215223821-4bcce26db858
	github.com/antonybholmes/go-math v0.0.0-20240215163921-12bb7e52185c
	github.com/labstack/echo/v4 v4.11.4
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20240219043615-26013b3f5c61
	github.com/antonybholmes/go-env v0.0.0-20240216174519-d83d9222e5a7
	github.com/antonybholmes/go-gene v0.0.0-20240219040039-0e816bfeef5b
	github.com/antonybholmes/go-genes v0.0.0-20240219040039-0e816bfeef5b
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rs/zerolog v1.32.0
)

require github.com/antonybholmes/go-mailer v0.0.0-00010101000000-000000000000

require (
	github.com/gofrs/uuid/v5 v5.0.0 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xyproto/randomstring v1.0.5 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
