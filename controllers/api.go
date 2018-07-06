package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"github.com/annchain/explorer/repository"
)

type ApiController struct {
	beego.Controller
}

type BlockListResp struct {
	Total  int
	Blocks []repository.DisplayItem
}

type TxListResp struct {
	Total int
	Txs   []repository.Transaction
}

func (p *ApiController) ContractQuery2() {

	p.ServeJSON()
	contractId := p.Ctx.Input.Query("contract_id")
	method := p.Ctx.Input.Query("method")
	postUrl := fmt.Sprintf("http://%s/contract/query", beego.AppConfig.String("api_addr"))
	data := make(url.Values)
	data["privkey"] = []string{"770152ec65cc029a45402b91a7c0d888b3c329f924674b470f926545747ca096"}
	data["id"] = []string{contractId}
	data["method"] = []string{method}
	data["chain_id"] = []string{repository.ChainID}
	fmt.Println(data, postUrl, "==============", contractId, method)

	/*	resp, err := http.PostForm(postUrl, data)
				if err != nil {
					p.Data["json"] = &Result{
						Success: false,
						Data:    "PostData failed: " + err.Error(),
					}
					return
				}
				defer resp.Body.Close()
			bytez, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.Data["json"] = &Result{
				Success: false,
				Data:    "ReadAll data failed: " + err.Error(),
			}
			return
		}
		fmt.Println(string(bytez))
	*/
	bytez := []byte("just test")
	p.Data["json"] = &Result{
		Success: true,
		Data:    bytez,
	}
}

func (p *ApiController) QueryBlock() {
	param := p.Ctx.Input.Query("param")
	defer p.ServeJSON()

	var err error
	var bk repository.Block
	if strings.HasPrefix(param, "0x") {
		bk, _, err = repository.OneBlock(param)
		if err != nil {
			p.Data["json"] = &Result{
				Success: false,
				Data:    "get block from repo failed",
			}
			return
		}
	} else {
		height, err := strconv.Atoi(param)
		if err != nil {
			p.Data["json"] = &Result{
				Success: false,
				Data:    "input format error.",
			}
			return
		}
		bk, _ = repository.BlockByHeight(height)
	}

	if bk.AppHash == "" {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "can not find block:" + param,
		}
		return
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    bk,
	}
}

func (p *ApiController) QueryBlockList() {
	from := p.Ctx.Input.Query("from")
	to := p.Ctx.Input.Query("to")
	defer p.ServeJSON()

	htFrom, err := strconv.Atoi(from)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	htTo, err := strconv.Atoi(to)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	if htTo-htFrom <= 0 || htFrom < 1 {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "to should be bigger than from",
		}
		return
	}

	var resp BlockListResp
	resp.Total, err = repository.Height()
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get height from repo failed",
		}
		return
	}
	resp.Blocks, err = repository.BlocksFromTo(htFrom, htTo)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get blocks from repo failed",
		}
		return
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    resp,
	}
}

func (p *ApiController) QueryTx() {
	hash := p.Ctx.Input.Query("hash")
	defer p.ServeJSON()

	tx, err := repository.OneTransaction(hash)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get tx from repo failed",
		}
		return
	}
	tx.PayloadHex = common.ToHex(tx.Payload)
	p.Data["json"] = &Result{
		Success: true,
		Data:    tx,
	}
}

func (p *ApiController) QueryTxsList() {
	from := p.Ctx.Input.Query("from")
	to := p.Ctx.Input.Query("to")
	defer p.ServeJSON()

	htFrom, err := strconv.Atoi(from)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	htTo, err := strconv.Atoi(to)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "input format error.",
		}
		return
	}
	if htTo-htFrom <= 0 || htFrom < 1 {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "to should be bigger than from",
		}
		return
	}

	var resp TxListResp
	resp.Total, err = repository.CollectionItemNum(repository.TX_COLLECT)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get collection num error.",
		}
		return
	}
	resp.Txs, err = repository.TransactionFromTo(htFrom, htTo)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get txs from repo failed",
		}
		return
	}
	for i := 0; i < len(resp.Txs); i++ {
		resp.Txs[i].PayloadHex = common.ToHex(resp.Txs[i].Payload)
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    resp,
	}
}

func (p *ApiController) QueryTxsListByBlk() {
	hash := p.Ctx.Input.Query("hash")
	defer p.ServeJSON()

	var err error
	var resp TxListResp
	resp.Txs, err = repository.TransactionsByBlkhash(hash)
	if err != nil {
		p.Data["json"] = &Result{
			Success: false,
			Data:    "get txs by block hash from repo failed",
		}
		return
	}
	resp.Total = len(resp.Txs)
	for i := 0; i < len(resp.Txs); i++ {
		resp.Txs[i].PayloadHex = common.ToHex(resp.Txs[i].Payload)
	}

	p.Data["json"] = &Result{
		Success: true,
		Data:    resp,
	}
}

func (p *ApiController) Test() {
	_ = p.Ctx.Input.Query("hash")
	defer p.ServeJSON()

	p.Data["json"] = &Result{
		Success: true,
		Data:    repository.TestSave(),
	}
}
func (p *ApiController) Health() {
	_ = p.Ctx.Input.Query("")
	defer p.ServeJSON()

	p.Data["json"] = &Result{
		Success: true,
		Data:    repository.TestSave(),
	}
}
