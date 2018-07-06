package rpc

import (
	"fmt"
	"log"

	"github.com/astaxie/beego"
)

var HTTP_ADDR string

func Init() {

	HTTP_ADDR = fmt.Sprintf("http://%s", beego.AppConfig.String("rpc_addr"))
	if HTTP_ADDR == "" {
		log.Panic("HTTP_ADDR is nil")
	}
}

type Result struct {
	Code    uint32      `json:"code"`
	Payload interface{} `json:"payload"`
}
