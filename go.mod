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
	github.com/antonybholmes/go-basemath v0.0.0-20241223034309-5f70a342a0e3 // indirect
	github.com/antonybholmes/go-dna v0.0.0-20241231004053-571f6c9d6eb6
	github.com/labstack/echo/v4 v4.13.3
	github.com/labstack/gommon v0.4.2 // indirect
)

require (
	github.com/antonybholmes/go-auth v0.0.0-20250101192611-ce4e7f6d596b
	github.com/antonybholmes/go-genes v0.0.0-20241225054554-d10c7e194b23
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/labstack/echo-jwt/v4 v4.3.0
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/antonybholmes/go-mailer v0.0.0-20241204182257-78d1094d458e
	github.com/antonybholmes/go-sys v0.0.0-20241219150701-64d5bd9623f7
	github.com/gorilla/sessions v1.4.0
	github.com/labstack/echo-contrib v0.17.2
)

require (
	github.com/antonybholmes/go-beds v0.0.0-20241230042333-b0ba8926c8ae
	github.com/antonybholmes/go-geneconv v0.0.0-20241219150701-97a6b2b1d896
	github.com/antonybholmes/go-math v0.0.0-20241223034309-5f70a342a0e3
	github.com/antonybholmes/go-motifs v0.0.0-20241227053706-4a94b275f707
	github.com/antonybholmes/go-mutations v0.0.0-20241204182259-2078e660a775
	github.com/antonybholmes/go-pathway v0.0.0-20241219150653-d23de5e61a0d
	github.com/antonybholmes/go-seqs v0.0.0-20241228025617-9d6717e605ca
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
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	github.com/antonybholmes/go-cytobands v0.0.0-20241225054645-47ddf66b3bfd
	github.com/antonybholmes/go-gex v0.0.0-20241204182259-8c7006a4366e
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
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)
