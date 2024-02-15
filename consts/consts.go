package consts

import "os"

var NAME string = os.Getenv("NAME")
var VERSION string = os.Getenv("VERSION")
var COPYRIGHT string = os.Getenv("COPYRIGHT")
