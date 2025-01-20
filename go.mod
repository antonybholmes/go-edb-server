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

replace github.com/antonybholmes/go-seqs => ../go-seqs

replace github.com/antonybholmes/go-cytobands => ../go-cytobands

replace github.com/antonybholmes/go-beds => ../go-beds

require (
	github.com/antonybholmes/go-basemath v0.0.0-20250107213632-9971295f8456 // indirect
	github.com/antonybholmes/go-dna v0.0.0-20250110222441-27b549fda20d
	github.com/labstack/echo/v4 v4.13.3
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20250113143741-d6cf1d634ada
	github.com/antonybholmes/go-genes v0.0.0-20250113143746-9246da82bc47
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.3.0
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20250110222437-d062d791b918
	github.com/antonybholmes/go-sys v0.0.0-20250113143747-03c4e3605208
	github.com/gorilla/sessions v1.4.0
	github.com/labstack/echo-contrib v0.17.2
)

require (
	github.com/antonybholmes/go-beds v0.0.0-20250113143739-af286db2a33a
	github.com/antonybholmes/go-geneconv v0.0.0-20250113143747-5a766d6d259a
	github.com/antonybholmes/go-math v0.0.0-20250107213632-9971295f8456
	github.com/antonybholmes/go-motifs v0.0.0-20250113143744-ecede76bcf68
	github.com/antonybholmes/go-mutations v0.0.0-20250113143742-0863f2deb9f3
	github.com/antonybholmes/go-pathway v0.0.0-20250113143738-be266628ab1a
	github.com/antonybholmes/go-seqs v0.0.0-20250114235934-62d6d351e63f
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
	github.com/xuri/nfp v0.0.0-20250111060730-82a408b9aa71 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/antonybholmes/go-cytobands v0.0.0-20250113143738-d358dd8ef3b8
	github.com/antonybholmes/go-gex v0.0.0-20250110222442-95fc36ceefe0
	github.com/aws/aws-sdk-go v1.55.6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xyproto/randomstring v1.2.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
