package consts

import (
	"os"

	"github.com/antonybholmes/go-sys/env"
)

var NAME string
var VERSION string
var COPYRIGHT string
var JWT_PRIVATE_KEY []byte
var JWT_PUBLIC_KEY []byte
var SESSION_SECRET string

func init() {
	env.Load()

	NAME = os.Getenv("NAME")
	VERSION = os.Getenv("VERSION")
	COPYRIGHT = os.Getenv("COPYRIGHT")
	JWT_PRIVATE_KEY = []byte(os.Getenv("JWT_SECRET"))
	JWT_PUBLIC_KEY = []byte(os.Getenv("JWT_SECRET"))
	SESSION_SECRET = os.Getenv("SESSION_SECRET")
}
