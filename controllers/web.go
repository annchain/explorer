package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.zhonganonline.com/ann/angine/types"
	go_common "gitlab.zhonganonline.com/ann/ann-module/lib/go-common"
	"gitlab.zhonganonline.com/ann/ann-module/lib/go-rpc/client"
	"github.com/annchain/explorer/repository"
)

const (
	DisplayNum   = 25
	BlockHashLen = 40
	AccPubkeyLen = 64
)

type WebController struct {
	beego.Controller
}

type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type IndexShow struct {
	LatestBlockHeight int `json:""`
	Transactions      int
}

func (wc *WebController) Index() {
	wc.Layout = "layout.html"
	wc.TplName = "web.tpl"
}

func (wc *WebController) IndexShow() {
	defer wc.ServeJSON()
	status, err := repository.GetStatus(repository.ChainID)
	if err != nil {
		log.Printf("repository.GetStatus failed:  %v\n", err)
		return
	}
	countTx, err := repository.CollectionItemNum("transaction_t")
	if err != nil {
		log.Printf("repository.CollectionItemNum failed:  %v\n", err)
		return
	}
	data := IndexShow{
		status.LatestBlockHeight,
		countTx,
	}

	wc.Data["json"] = &Result{
		Success: true,
		Data:    data,
	}

}

type PageControl struct {
	FirstPage   int
	CurrentPage int
	PrevPage    int
	NextPage    int
	LastPage    int
	Items       int
	Pages       int
}

func (wc *WebController) Blocks() {
	wc.Layout = "layout.html"
	wc.TplName = "blocks.tpl"
	count, errC := repository.CollectionItemNum("block_t")
	if errC != nil {
		log.Println(errC)
		return
	}
	page := wc.GetString(":page")
	var (
		pageNum int
		err     error
	)
	if page == "latest" {
		pageNum = 1
	} else {
		pageNum, err = strconv.Atoi(page)
		if err != nil {
			log.Printf("Invalid PageNum %s\n", page)
			return
		}
	}

	data, err := repository.PageBlocks(pageNum, DisplayNum)
	if err != nil {
		log.Fatal(err)
	}

	wc.Data["Blocks"] = data
	pageControl := PageControl{
		FirstPage:   1,
		Items:       count,
		LastPage:    int(math.Ceil(float64(count) / float64(DisplayNum))),
		CurrentPage: pageNum,
	}
	pageControl.PrevPage = pageControl.CurrentPage - 1
	pageControl.NextPage = pageControl.CurrentPage + 1
	if pageControl.PrevPage < 1 {
		pageControl.PrevPage = 1
	}
	if pageControl.NextPage > pageControl.LastPage {
		pageControl.NextPage = pageControl.LastPage
	}
	wc.Data["Page"] = pageControl
}

func (wc *WebController) Block() {
	hash := wc.GetString(":hash")
	block, txs, err := repository.OneBlock(hash)
	if err != nil {
		wc.TplName = "error.tpl"
	} else {
		wc.Layout = "layout.html"
		wc.TplName = "block.tpl"
		block.Time = block.Time.Local()
		wc.Data["Block"] = block
		wc.Data["Transactions"] = txs
	}

}

func (wc *WebController) TxsPage() {
	wc.Layout = "layout.html"
	wc.TplName = "txs.tpl"
}

func (wc *WebController) Txs() {
	defer wc.ServeJSON()
	data, _ := repository.Txs(DisplayNum)
	wc.Data["json"] = &Result{
		Success: true,
		Data:    data,
	}
}

func (wc *WebController) ContractPage() {
	wc.Layout = "layout.html"
	wc.TplName = "contract.tpl"
}

func (wc *WebController) Contracts() {
	defer wc.ServeJSON()
	data, _ := repository.Contracts(DisplayNum)
	wc.Data["json"] = &Result{
		Success: true,
		Data:    data,
	}
}

func (wc *WebController) TxByHash() {
	hash := wc.GetString(":hash")
	wc.Layout = "layout.html"
	wc.TplName = "tx_view.tpl"
	tx, _ := repository.OneTransaction(hash)

	wc.Data["Transaction"] = tx

}

func (wc *WebController) H5TxByHash() {
	hash := wc.GetString(":hash")
	wc.TplName = "h5.tpl"
	if tx, err := repository.OneTransaction(hash); err == nil {

		if tx.Height > 0 {
			block, _ := repository.BlockByHeight(tx.Height)
			wc.Data["Block"] = block
		}
		wc.Data["Transaction"] = tx

	}
}

type Account struct {
	Pubkey  string
	Balance string
}

func (wc *WebController) Account() {
	pubkey := wc.GetString(":pubkey")
	wc.Layout = "layout.html"
	wc.TplName = "account.tpl"
	tcpServer := fmt.Sprintf("tcp://%s", beego.AppConfig.String("api_addr"))
	var (
		acc = Account{
			Pubkey: pubkey,
		}
		err error
	)
	acc.Balance, err = queryBalance(pubkey, tcpServer)
	if err != nil {
		log.Println("queryBalance failed: ", err)
	}
	wc.Data["Account"] = acc
}

func queryBalance(pubkey, tcpServer string) (balance string, err error) {
	clientJSON := rpcclient.NewClientJSONRPC(nil, tcpServer)
	tmResult := new(types.RPCResult)

	addrHex := go_common.SanitizeHex(pubkey)
	addr := common.Hex2Bytes(addrHex)
	query := append([]byte{0x41}, addr...)

	_, err = clientJSON.Call("query", []interface{}{beego.AppConfig.String("chain_id"), query}, tmResult)
	if err != nil {
		return
	}

	res := (*tmResult).(*types.ResultQuery)
	balance = string(res.Result.Data)
	return
}

func (wc *WebController) ContractQuery() {
	contractId := wc.Ctx.Input.Query("contract_id")
	method := wc.Ctx.Input.Query("method")
	wc.Layout = "layout.html"
	wc.TplName = "bubuji.tpl"
	postUrl := fmt.Sprintf("http://%s/contract/query", beego.AppConfig.String("api_addr"))
	data := make(url.Values)
	data["privkey"] = []string{"770152ec65cc029a45402b91a7c0d888b3c329f924674b470f926545747ca096"}
	data["id"] = []string{contractId}
	data["method"] = []string{method}
	data["chain_id"] = []string{repository.ChainID}

	resp, err := http.PostForm(postUrl, data)
	if err != nil {
		log.Println("http.PostForm failed: ", err)
		return
	}
	defer resp.Body.Close()
	bytez, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll failed: ", err)
		return
	}
	fmt.Println(string(bytez))

	wc.Data["Data"] = string(bytez)
}

func (wc *WebController) Search() {
	hash := wc.GetString(":hash")

	reg := regexp.MustCompile(`^[1-9][0-9]*$`)

	page := 0

	if reg.MatchString(hash) {
		fmt.Printf("hash is block height ")
		height, err := strconv.Atoi(hash)
		if err != nil {
			fmt.Println("err is : %v", err)
		}
		block, err := repository.BlockByHeight(height)
		if err == nil {
			page = 1
			hash = block.Hash
		} else {
			//err = job.BlockChain(height)
			//if err != nil {
			//	log.Println(fmt.Sprintf("job.BlockChain(%d) failed: %v", height, err))
			//	return
			//}
			block, err := repository.BlockByHeight(height)
			if err == nil {
				page = 1
				hash = block.Hash
			}
		}
	} else {
		switch len(hash) {
		case BlockHashLen: //block  or contract
			_, _, err := repository.OneBlock(hash)
			if err == nil {
				page = 1
			} else {
				_, err2 := repository.OneTransaction(hash)
				if err2 == nil {
					page = 2
				}
			}
			break

		case AccPubkeyLen:
			page = 3
		}
	}

	fmt.Printf("hash :%s", hash)

	switch page {
	case 0:
		wc.Redirect("/", 302)
		break
	case 1:
		wc.Redirect("/view/blocks/hash/"+hash, 302)
		break
	case 2:
		wc.Redirect("/view/txs/hash/"+hash, 302)
		break
	case 3:
		wc.Redirect("/view/account/pubkey/"+hash, 302)
	}

}
