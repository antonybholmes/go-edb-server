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

replace github.com/antonybholmes/go-motifs => ../go-motifs

replace github.com/antonybholmes/go-pathway => ../go-pathway

replace github.com/antonybholmes/go-gex => ../go-gex

require (
	github.com/antonybholmes/go-basemath v0.0.0-20240825181410-a6174a39116c // indirect
	github.com/antonybholmes/go-dna v0.0.0-20241007150544-1b58eb1162ce
	github.com/labstack/echo/v4 v4.13.3
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20241022154234-65b9f487c031
	github.com/antonybholmes/go-genes v0.0.0-20241018214057-74bdb6aa6946
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.3.0
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20241007150543-6b1570c6da90
	github.com/antonybholmes/go-sys v0.0.0-20241021142915-13ff92cbbcdc
	github.com/gorilla/sessions v1.4.0
	github.com/labstack/echo-contrib v0.17.1
)

require (
	github.com/antonybholmes/go-geneconv v0.0.0-20241018214059-3aa39a68396d
	github.com/antonybholmes/go-math v0.0.0-20240825181410-a6174a39116c
	github.com/antonybholmes/go-motifs v0.0.0-20241019032547-082676e3625a
	github.com/antonybholmes/go-mutations v0.0.0-20241008211559-3b0ab30c604d
	github.com/antonybholmes/go-pathway v0.0.0-20241022154233-27225f703cfd
	github.com/redis/go-redis/v9 v9.7.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/matoous/go-nanoid/v2 v2.1.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/xuri/efp v0.0.0-20241211021726-c4e992084aa6 // indirect
	github.com/xuri/excelize/v2 v2.9.0 // indirect
	github.com/xuri/nfp v0.0.0-20240318013403-ab9948c2c4a7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/antonybholmes/go-gex v0.0.0-20241009210106-9c16df39c4d0
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xyproto/randomstring v1.2.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
