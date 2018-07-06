package routers

import (
	"github.com/astaxie/beego"
	. "github.com/annchain/explorer/controllers"
)

func init() {

	beego.Router("/", &WebController{}, "get:Index")
	beego.Router("/health", &ApiController{}, "head:Health")
	beego.Router("/health", &ApiController{}, "get:Health")
	webNs := beego.NewNamespace("/view",
		beego.NSRouter("/", &WebController{}, "get:IndexShow"),
		beego.NSRouter("/blocks/:page", &WebController{}, "get:Blocks"),
		beego.NSRouter("/blocks/hash/:hash", &WebController{}, "get:Block"),
		beego.NSRouter("/account/pubkey/:pubkey", &WebController{}, "get:Account"),

		beego.NSRouter("/txs/page", &WebController{}, "get:TxsPage"),
		beego.NSRouter("/txs/latest", &WebController{}, "get:Txs"),

		beego.NSRouter("/contracts/page", &WebController{}, "get:ContractPage"),
		beego.NSRouter("/contracts/latest", &WebController{}, "get:Contracts"),
		beego.NSRouter("/txs/hash/:hash", &WebController{}, "get:TxByHash"),
		beego.NSRouter("/h5/txs/hash/:hash", &WebController{}, "get:H5TxByHash"),
		beego.NSRouter("/search/:hash", &WebController{}, "get:Search"),
		beego.NSRouter("/contract/query", &WebController{}, "get:ContractQuery"),
	)

	apiNs := beego.NewNamespace("/v1",
		beego.NSNamespace("/account"), //beego.NSRouter("/list", &AccountController{}, "get:List"),

		beego.NSNamespace("/txs",
			beego.NSRouter("/query/:fromTo", &TxsController{}, "get:Query")),

		beego.NSNamespace("/info"), //beego.NSRouter("/status", &InfoController{}, "get:Status"),

		beego.NSNamespace("/net"), //beego.NSRouter("/info", &NetController{}, "get:Info"),

		beego.NSNamespace("/block"), //			beego.NSRouter("/blockchain/:maxHeight", &BlockController{}, "get:BlockChain"),
		//			beego.NSRouter("/info/:height", &BlockController{}, "get:Block"),

	)

	apiNsV2 := beego.NewNamespace("/v2",
		beego.NSNamespace("/query",
			beego.NSRouter("/block", &ApiController{}, "get:QueryBlock"),
			beego.NSRouter("/blocklist", &ApiController{}, "get:QueryBlockList"),
			beego.NSRouter("/tx", &ApiController{}, "get:QueryTx"),
			beego.NSRouter("/txlist", &ApiController{}, "get:QueryTxsList"),
			beego.NSRouter("/txlistByBlock", &ApiController{}, "get:QueryTxsListByBlk"),
			beego.NSRouter("/test", &ApiController{}, "get:Test"),
		),
	)
	beego.AddNamespace(webNs)
	beego.AddNamespace(apiNs)
	beego.AddNamespace(apiNsV2)

	beego.SetStaticPath("/assets", "static")

}
