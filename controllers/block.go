package controllers

/*
import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/rlp"
	. "github.com/annchain/explorer/logs"
	"github.com/annchain/explorer/rpc"
)

type BlockController struct {
	beego.Controller
}

func (bc *BlockController) BlockChain() {
	maxHeight := bc.Ctx.Input.Param(":maxHeight")
	iMaxHeight, err := strconv.Atoi(maxHeight)
	defer bc.ServeJSON()
	if err != nil {
		LOG.Error("request parameter  maxHeight  input error. %v", err)
		bc.Data["json"] = &rpc.Result{
			Code:    400,
			Payload: "input parameter error.",
		}
		return
	}
	callResult := rpc.TCall("blockchain", []interface{}{0, iMaxHeight})
	bc.Data["json"] = callResult
}

func (bc *BlockController) Block() {
	height := bc.Ctx.Input.Param(":height")
	iHeight, err := strconv.Atoi(height)
	defer bc.ServeJSON()
	if err != nil {
		LOG.Error("request parameter height  input error. %v", err)
		bc.Data["json"] = &rpc.Result{
			Code:    400,
			Payload: "input parameter error.",
		}
		return
	}
	callResult := rpc.TCall("block", []interface{}{iHeight})

	resultBlock := (callResult.Payload).(*anntypes.ResultBlock)

	for _, o := range resultBlock.Block.Txs {
		tx := new(ethtypes.Transaction)
		error := rlp.DecodeBytes(o, tx)
		if error != nil {
			LOG.Error("rlp tx error  %v", err)
		}

		fmt.Println("tx :", tx.Hash().Hex(), " , gas : ", tx.Gas())
	}

	callResult.Payload = resultBlock

	bc.Data["json"] = callResult

}
*/

/* reson(use blcokchain api ,and the result of blockchain is different with the result of browser)

curl 'http://10.253.4.248:9090/v1/block/blockchain/1'|json_reformat |grep -e 'time' -e 'validators_hash'
                    "time": "2017-09-07T03:36:38.036Z",
                    "validators_hash": "crQYuixsB29qatZc5czHn3Xi15s=",

curl 'http://10.253.4.248:9090/v1/block/info/1' |json_reformat |grep -e 'time' -e 'validators_hash'
                "time": "2017-09-07T03:36:38.036Z",
                "validators_hash": "crQYuixsB29qatZc5czHn3Xi15s=",


curl 'http://10.253.4.248:5657/block?height=1'|json_reformat |grep -e 'time' -e 'validators_hash'
                    "time": "2017-09-07T03:36:38.036Z",
                    "validators_hash": "72B418BA2C6C076F6A6AD65CE5CCC79F75E2D79B",


*/
