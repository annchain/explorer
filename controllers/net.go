package controllers

/*
import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/annchain/explorer/rpc"
)

type NetController struct {
	beego.Controller
}

func (this *NetController) Info() {
	tmResult := new(anntypes.TMResult)
	_, err := rpc.Call("net_info", []interface{}{}, tmResult)
	if err != nil || *tmResult == nil {
		fmt.Errorf("jsonRPC call: %v", err)
		this.ServeJSON()
		return
	}
	resultNetInfo := (*tmResult).(*anntypes.ResultNetInfo)

	//latest_block_height

	this.Data["json"] = resultNetInfo

	this.ServeJSON()
}

*/

/* reason

curl 'http://10.253.4.248:9090/v1/net/info'|json_reformat
null

*/
