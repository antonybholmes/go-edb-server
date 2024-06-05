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
	github.com/antonybholmes/go-dna v0.0.0-20240528182114-dd434dbd0c49
	github.com/antonybholmes/go-math v0.0.0-20240215163921-12bb7e52185c
	github.com/labstack/echo/v4 v4.12.0
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20240604033739-7068bb278fb0
	github.com/antonybholmes/go-genes v0.0.0-20240604033740-7a140d97cfdd
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20240603221758-2f95021dab54
	github.com/antonybholmes/go-sys v0.0.0-20240505052557-9f8864ac77aa
	github.com/gorilla/sessions v1.2.2
	github.com/labstack/echo-contrib v0.17.1
	github.com/michaeljs1990/sqlitestore v0.0.0-20210507162135-8585425bc864
)

require (
	github.com/antonybholmes/go-gene-conversion v0.0.0-20240604180412-159f13b8d8da
	github.com/antonybholmes/go-mutations v0.0.0-20240604033738-9a3814392967
)

require (
	github.com/antonybholmes/go-microarray v0.0.0-20240504032631-9fb6b43a10d4
	github.com/aws/aws-sdk-go v1.53.16 // indirect
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
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/time v0.5.0 // indirect
)
