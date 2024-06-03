module github.com/antonybholmes/go-edb-api

go 1.22.2

replace github.com/antonybholmes/go-auth => ../go-auth

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-genes => ../go-genes

replace github.com/antonybholmes/go-microarray => ../go-microarray

replace github.com/antonybholmes/go-mutations => ../go-mutations

replace github.com/antonybholmes/go-math => ../go-math

replace github.com/antonybholmes/go-sys => ../go-sys

replace github.com/antonybholmes/go-mailer => ../go-mailer

replace github.com/antonybholmes/go-gene-conversion => ../go-gene-conversion

require (
	github.com/antonybholmes/go-dna v0.0.0-20240503021126-08c3c39059f5
	github.com/antonybholmes/go-math v0.0.0-20240215163921-12bb7e52185c
	github.com/labstack/echo/v4 v4.12.0
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20240502232413-63b01c2be949
	github.com/antonybholmes/go-genes v0.0.0-20240505003605-c7f9b24cc7f3
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20240504032634-6b3ed4b8d58f
	github.com/antonybholmes/go-sys v0.0.0-20240501232923-152b6e4cc204
	github.com/gorilla/sessions v1.2.2
	github.com/labstack/echo-contrib v0.17.1
	github.com/michaeljs1990/sqlitestore v0.0.0-20210507162135-8585425bc864
)

require (
	github.com/antonybholmes/go-gene-conversion v0.0.0-20240601171451-8624e6ee9dcb
	github.com/antonybholmes/go-mutations v0.0.0-20240505003604-87e9d6ee2b76
)

require (
	github.com/antonybholmes/go-microarray v0.0.0-20240504032631-9fb6b43a10d4
	github.com/aws/aws-sdk-go v1.53.15 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xyproto/randomstring v1.0.5 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
