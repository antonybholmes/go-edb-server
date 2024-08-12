module github.com/antonybholmes/go-edb-api

go 1.22.5

replace github.com/antonybholmes/go-auth => ../go-auth

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-genes => ../go-genes

replace github.com/antonybholmes/go-microarray => ../go-microarray

replace github.com/antonybholmes/go-mutations => ../go-mutations

replace github.com/antonybholmes/go-basemath => ../go-basemath

replace github.com/antonybholmes/go-sys => ../go-sys

replace github.com/antonybholmes/go-mailer => ../go-mailer

replace github.com/antonybholmes/go-geneconv => ../go-geneconv

replace github.com/antonybholmes/go-motiftogene => ../go-motiftogene

replace github.com/antonybholmes/go-pathway => ../go-pathway

require (
	github.com/antonybholmes/go-basemath v0.0.0-20240802221548-7773050a8f2f // indirect
	github.com/antonybholmes/go-dna v0.0.0-20240805202906-f91775279da8
	github.com/labstack/echo/v4 v4.12.0
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20240805202905-b7f199898342
	github.com/antonybholmes/go-genes v0.0.0-20240805202902-bf341a39aa7c
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20240805202858-c5a47009e8bf
	github.com/antonybholmes/go-sys v0.0.0-20240801224521-3bed2c519a83
	github.com/gorilla/sessions v1.3.0
	github.com/labstack/echo-contrib v0.17.1
	github.com/michaeljs1990/sqlitestore v0.0.0-20210507162135-8585425bc864
)

require (
	github.com/antonybholmes/go-geneconv v0.0.0-20240805202859-92730bef5f4f
	github.com/antonybholmes/go-math v0.0.0-20240802221548-7773050a8f2f
	github.com/antonybholmes/go-motiftogene v0.0.0-20240805202904-08e20e1f17c9
	github.com/antonybholmes/go-mutations v0.0.0-20240805202902-cdf90fd60d59
	github.com/antonybholmes/go-pathway v0.0.0-20240805202900-f2186667883b
)

require gopkg.in/yaml.v2 v2.4.0 // indirect

require (
	github.com/antonybholmes/go-microarray v0.0.0-20240504032631-9fb6b43a10d4
	github.com/aws/aws-sdk-go v1.55.5 // indirect
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
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/time v0.6.0 // indirect
)
