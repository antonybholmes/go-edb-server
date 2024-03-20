package consts

import (
	"os"

	"github.com/antonybholmes/go-sys/env"
)

var NAME string
var VERSION string
var COPYRIGHT string
var JWT_SECRET string
var SESSION_SECRET string

func init() {
	env.Load()

	NAME = os.Getenv("NAME")
	VERSION = os.Getenv("VERSION")
	COPYRIGHT = os.Getenv("COPYRIGHT")
	JWT_SECRET = os.Getenv("JWT_SECRET")
	SESSION_SECRET = os.Getenv("SESSION_SECRET")
}
