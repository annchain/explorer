package controllers

/*
import (
	"github.com/astaxie/beego"
	"github.com/annchain/explorer/rpc"
)

type AccountController struct {
	beego.Controller
}

func (this *AccountController) List() {
	defer this.ServeJSON()
	callResult := rpc.TCall("accounts", []interface{}{})
	this.Data["json"] = callResult
}
*/

/* reason

curl 'http://10.253.4.248:9090/v1/account/list' |json_reformat
{
    "code": 500,
    "payload": "Response error: RPC method unknown: accounts"
}

*/
