module github.com/antonybholmes/go-edb-server

go 1.23

replace github.com/antonybholmes/go-auth => ../go-auth

replace github.com/antonybholmes/go-dna => ../go-dna

replace github.com/antonybholmes/go-genes => ../go-genes

replace github.com/antonybholmes/go-mutations => ../go-mutations

replace github.com/antonybholmes/go-basemath => ../go-basemath

replace github.com/antonybholmes/go-sys => ../go-sys

replace github.com/antonybholmes/go-mailer => ../go-mailer

replace github.com/antonybholmes/go-geneconv => ../go-geneconv

replace github.com/antonybholmes/go-motiftogene => ../go-motiftogene

replace github.com/antonybholmes/go-pathway => ../go-pathway

replace github.com/antonybholmes/go-gex => ../go-gex

require (
	github.com/antonybholmes/go-basemath v0.0.0-20240825181410-a6174a39116c // indirect
	github.com/antonybholmes/go-dna v0.0.0-20240830030422-fdb7452d202d
	github.com/labstack/echo/v4 v4.12.0
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20240901220333-c6210ff22d37
	github.com/antonybholmes/go-genes v0.0.0-20240901185659-8b17b98bbf71
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.2.0
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20240901222928-242cea64757a
	github.com/antonybholmes/go-sys v0.0.0-20240901191116-4f230479c4a8
	github.com/gorilla/sessions v1.4.0
	github.com/labstack/echo-contrib v0.17.1
)

require (
	github.com/antonybholmes/go-geneconv v0.0.0-20240901185654-d75dd1745c01
	github.com/antonybholmes/go-math v0.0.0-20240825181410-a6174a39116c
	github.com/antonybholmes/go-motiftogene v0.0.0-20240901185652-b6d4aa5d3c35
	github.com/antonybholmes/go-mutations v0.0.0-20240831050902-7ba8652704c1
	github.com/antonybholmes/go-pathway v0.0.0-20240901185656-7497cb6ea997
	github.com/redis/go-redis/v9 v9.6.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/antonybholmes/go-gex v0.0.0-20240825181414-7343636e387b
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
