package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	"log"

	"github.com/astaxie/beego"
	. "github.com/annchain/explorer/logs"
	"github.com/annchain/explorer/rpc"
	"gopkg.in/mgo.v2/bson"
)

var (
	BLOCK_COLLECT = "block"
	TX_COLLECT    = "transaction"
	MONGO_URL     string
	DB_NAME       string
	ChainID       string
)

type Repo interface {
	Save(interface{})
}

type HTTPResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      string           `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   string           `json:"error"`
}

type Status struct {
	NodeInfo          *NodeInfo `json:"node_info"`
	LatestBlockHeight int       `json:"latest_block_height"`
}

type NodeInfo struct {
	NetWork string `json:"network"`
}

func GetHTTPResp(url string) (bytez []byte, err error) {

	resp, errR := http.Get(url)
	if errR != nil {
		err = errR
		return
	}
	defer resp.Body.Close()
	bytez, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var hr HTTPResponse
	err = json.Unmarshal(bytez, &hr)
	if err != nil {
		return
	}
	if hr.Result == nil {
		err = errors.New(fmt.Sprintf("json.Unmarshal (%s)HTTPResponse wrong ,maybe you need config 'chain_id'", url))
		return
	}
	bytez, err = hr.Result.MarshalJSON()
	if err != nil {
		return
	}
	//str := string(bytez)
	//i := strings.Index(str, ",")
	//bytez = bytez[i+1 : len(bytez)-1]
	return
}

func GetStatus(chainID string) (status Status, err error) {

	url := fmt.Sprintf("%s/status?chainid=\"%s\"", rpc.HTTP_ADDR, chainID)
	bytez, errG := GetHTTPResp(url)
	if errG != nil {
		err = errG
		return
	}

	err = json.Unmarshal(bytez, &status)
	if err != nil {
		log.Fatalf("json.Unmarshal(Status) failed: %s", err.Error())
	}
	return
}

func Init() {

	ChainID = beego.AppConfig.String("chain_id")
	if ChainID == "" {
		status, err := GetStatus("")
		if err != nil {
			log.Fatal(err)
		}
		ChainID = status.NodeInfo.NetWork
	}
	BLOCK_COLLECT += "_" + ChainID
	TX_COLLECT += "_" + ChainID
	BLOCK_COLLECT = strings.Replace(BLOCK_COLLECT, "-", "_", -1)
	TX_COLLECT = strings.Replace(TX_COLLECT, "-", "_", -1)
	fmt.Println("BLOCK_COLLECT = ", BLOCK_COLLECT)
	fmt.Println("TX_COLLECT = ", TX_COLLECT)

	DB_NAME = "block_browser"

	//mongodb://myuser:mypass@localhost:40001
	if beego.AppConfig.String("mogo_addr") != "" {
		if beego.AppConfig.String("mogo_user") != "" {
			MONGO_URL = "mongodb://" +
				beego.AppConfig.String("mogo_user") + ":" +
				beego.AppConfig.String("mogo_pwd") + "@" +
				beego.AppConfig.String("mogo_addr") + "/" +
				DB_NAME
		} else {
			MONGO_URL = beego.AppConfig.String("mogo_addr")
		}

		err := CreateIndex()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := CreateSqlite()
		if err != nil {
			log.Fatal(err)
		}
	}

	return
}

type BlockRepo struct {
	Blocks []Block
	Txs    []Transaction
}

type DisplayItem struct {
	Block
	Tps      int
	Interval float64
}

type Block struct {
	Hash           string
	ChainID        string
	Height         int
	Time           time.Time
	NumTxs         int
	LastCommitHash string
	DataHash       string
	ValidatorsHash string
	AppHash        string
	Reward         uint64
	CoinBase       string
}

type Transaction struct {
	Payload    []byte `json:"-"`
	PayloadHex string
	Hash       string
	From       string
	To         string
	Receipt    string
	Amount     string
	Nonce      uint64
	Gas        string
	Size       int64
	Block      string
	Contract   string
	Time       time.Time
	Height     int
	TxType     string
	Fee        uint64
}

func CreateIndex() (err error) {
	session, err := mgo.Dial(MONGO_URL)
	if err != nil {
		return
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	//block_collect
	c := session.DB(DB_NAME).C(BLOCK_COLLECT)
	index := mgo.Index{
		Key: []string{"height"},
	}
	err = c.EnsureIndex(index)
	return
}

func LatestBlocks(limit int) (displayData []DisplayItem, err error) {

	blocks := []Block{}
	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			err = errD
			return
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(BLOCK_COLLECT)
		query := c.Find(nil).Sort("-height").Limit(limit)
		err = query.All(&blocks)
		if err != nil {
			return
		}

	} else {
		blocks, err = LatestBlocksBySqlite(limit)
		if err != nil {
			return
		}
	}
	if len(blocks) > 0 {
		dur := blocks[0].Time.Sub(blocks[limit-1].Time).Seconds()
		var totalTxs int
		for _, v := range blocks {
			totalTxs += v.NumTxs
		}
		interval := dur / float64(limit)
		tps := int(float64(totalTxs) / dur)
		for _, v := range blocks {
			displayData = append(displayData, DisplayItem{v, tps, interval})
		}
	}
	return
}

func PageBlocks(pageIndex, pageSize int) (displayData []DisplayItem, err error) {
	blocks := []Block{}
	if MONGO_URL != "" {
		session, err := mgo.Dial(MONGO_URL)
		if err != nil {
			log.Fatal("error connect mogo : ", err)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(BLOCK_COLLECT)

		query := c.Find(bson.M{"height": bson.M{"$lte": pageIndex + pageSize, "$gte": pageIndex}}).Sort("-height")
		blocks := []Block{}
		err = query.All(&blocks)
	} else {
		blocks, err = PageBlocksBySqlite(pageIndex, pageSize)
	}
	if len(blocks) > 0 {
		dur := blocks[len(blocks)-1].Time.Sub(blocks[0].Time).Seconds()
		var totalTxs int
		for _, v := range blocks {
			totalTxs += v.NumTxs
		}
		interval := dur / float64(len(blocks))
		tps := int(float64(totalTxs) / dur)
		for _, v := range blocks {
			displayData = append(displayData, DisplayItem{v, tps, interval})
		}
	}
	return

}

func BlocksFromTo(from, to int) (displayData []DisplayItem, err error) {

	blocks := []Block{}
	if MONGO_URL != "" {
		session, err := mgo.Dial(MONGO_URL)
		if err != nil {
			log.Fatal("error connect mogo : ", err)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(BLOCK_COLLECT)

		query := c.Find(bson.M{"height": bson.M{"$lte": to, "$gte": from}}).Sort("-height")
		blocks := []Block{}
		err = query.All(&blocks)
	} else {
		blocks, err = BlocksFromToBySqlite(from, to)
	}
	dur := blocks[len(blocks)-1].Time.Sub(blocks[0].Time).Seconds()
	var totalTxs int
	for _, v := range blocks {
		totalTxs += v.NumTxs
	}
	interval := dur / float64(len(blocks))
	tps := int(float64(totalTxs) / dur)
	for _, v := range blocks {
		displayData = append(displayData, DisplayItem{v, tps, interval})
	}
	return
}

func OneBlock(hash string) (block Block, txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(BLOCK_COLLECT)

		query := c.Find(bson.M{"hash": hash})
		err = query.One(&block)
		if err != nil {
			return
		}
		c2 := session.DB(DB_NAME).C(TX_COLLECT)
		query2 := c2.Find(bson.M{"block": hash})
		err = query2.All(&txs)
		if err != nil {
			return
		}

	} else {
		block, txs, err = OneBlockBySqlite(hash)
	}
	return
}

func BlockByHeight(height int) (block Block, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(BLOCK_COLLECT)

		query := c.Find(bson.M{"height": height})
		err = query.One(&block)

	} else {
		block, err = BlockByHeightBySqlite(height)
	}

	return
}

func TransactionFromTo(from, to int) (txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)

		query := c.Find(nil).Skip(from - 1).Limit(to - from + 1)
		err = query.All(&txs)

	} else {
		log.Fatal("how to use")
		txs, err = TransactionFromToBySqlite(from, to)
	}
	return
}

func TransactionsByBlkhash(hash string) (txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)

		query := c.Find(bson.M{"block": hash})
		err = query.All(&txs)
	} else {
		txs, err = TransactionsByBlkhashBySqlite(hash)
	}
	return
}

func OneTransaction(hash string) (tx Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)

		query := c.Find(bson.M{"hash": hash})
		err = query.One(&tx)
	} else {
		tx, err = OneTransactionBySqlite(hash)
	}
	return
}

func OneContract(hash string) (contract Transaction, txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)

		query := c.Find(bson.M{"contract": hash})
		err = query.One(&contract)
		if err != nil {
			return
		}
		c2 := session.DB(DB_NAME).C(TX_COLLECT)
		query2 := c2.Find(bson.M{"to": hash})
		err = query2.All(&txs)
	} else {
		contract, txs, err = OneContractBySqlite(hash)
	}
	return
}

func Txs(limit int) (txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)
		query := c.Find(bson.M{"contract": bson.M{"$eq": ""}}).Limit(limit)
		err = query.All(&txs)
	} else {
		txs, err = TxsBySqlite(limit)
	}
	return
}

func TxsQuery(fromTo string) (txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)
		query := c.Find(bson.M{"$or": []bson.M{{"from": bson.M{"$eq": fromTo}}, {"to": bson.M{"$eq": fromTo}}}})
		err = query.All(&txs)
	} else {
		txs, err = TxsQueryBySqlite(fromTo)
	}
	return
}

func Contracts(limit int) (txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)
		query := c.Find(bson.M{"contract": bson.M{"$ne": ""}}).Limit(limit)
		err = query.All(&txs)
	} else {
		txs, err = ContractsBySqlite(limit)
	}
	return
}

func Contract(hash string) (tx Transaction, txs []Transaction, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()

		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(TX_COLLECT)

		query := c.Find(bson.M{"hash": hash}).Limit(1)
		err = query.One(&tx)
		if err != nil {
			return
		}
		query2 := c.Find(bson.M{"to": hash})
		err = query2.All(&txs)
	} else {
		tx, txs, err = ContractBySqlite(hash)
	}
	return
}

func Height() (maxHeight int, err error) {
	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(BLOCK_COLLECT)
		//pipe := c.Pipe([]bson.M{bson.M{"$group" : bson.M{"_id":"$height", "height" : bson.M{"$max" : "$height"}, }}})
		query := c.Find(nil).Sort("-height").Limit(1)
		result := Block{}
		err = query.One(&result)
		if err != nil {
			LOG.Info("Info : query Block Max Height , err: %v", err)
		}
		maxHeight = result.Height
	} else {
		status, errS := GetStatus(ChainID)
		if errS != nil {
			err = errS
			return
		}
		maxHeight = status.LatestBlockHeight
	}
	return
}

func CollectionItemNum(collect string) (count int, err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			log.Fatal("error connect mogo : ", errD)
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(collect)
		return c.Count()
	} else {
		count, err = CollectionItemNumBySqlite(collect)
	}
	return
}

func (br *BlockRepo) Save() (err error) {

	if MONGO_URL != "" {
		session, errD := mgo.Dial(MONGO_URL)
		if errD != nil {
			err = errD
			return
		}
		defer session.Close()
		// Optional. Switch the session to a monotonic behavior.
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(DB_NAME).C(BLOCK_COLLECT)
		for _, b := range br.Blocks {
			err = c.Insert(b)
			if err != nil {
				err = errors.New("insert block failed: " + err.Error())
				return
			}
		}

		c2 := session.DB(DB_NAME).C(TX_COLLECT)
		for _, tx := range br.Txs {
			err = c2.Insert(tx)
			if err != nil {
				err = errors.New("insert tx failed: " + err.Error())
				return
			}
		}
	} else {

		err = SaveBlockBySqlite(br.Blocks)
		if err != nil {
			err = errors.New("insert block failed: " + err.Error())
			return
		}
		if len(br.Txs) > 0 {
			err = SaveTxBySqlite(br.Txs)
			if err != nil {
				err = errors.New("insert tx failed: " + err.Error())
				return
			}
		}
	}
	return
}

func TestSave() []byte {
	session, err := mgo.Dial(MONGO_URL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(DB_NAME).C("mytest")

	var bs []byte
	bs = make([]byte, 5)
	for i := 0; i < 5; i++ {
		bs[i] = byte(i)
	}
	fmt.Println("bs: ", bs)
	c.Insert(bs)

	var bs2 []byte
	query := c.Find(nil)
	query.All(&bs2)
	fmt.Println("bs2: ", bs)

	return bs2
}
