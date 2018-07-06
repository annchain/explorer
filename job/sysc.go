package job

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.zhonganonline.com/ann/angine/types"
	"github.com/annchain/explorer/repository"
	"github.com/annchain/explorer/rpc"
	"sync"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
)

const (
	Interval        = 1 * time.Second
	Step            = 500
	TagPrefixLength = 4
	gape            = 1
)

var EthSigner = HomesteadSigner{}

func SyncTimingTask() {
	ticker := time.NewTicker(Interval)
	for range ticker.C {
		fmt.Println(time.Now())
		// 获取db 最高height
		heightOfDb := 1
		data, err := repository.LatestBlocks(1)
		if len(data) > 0 {
			heightOfDb = data[0].Height + 1
		} else {
			heightOfDb = 1
		}
		if err != nil {
			log.Println("SyncTimingTask repository.Height failed: ", err)
			continue
		}

		heightOfNode, err := repository.Height()

		if err != nil {
			log.Println("SyncTimingTask repository.Height failed: ", err)
			continue
		}

		if heightOfNode-heightOfDb > 0 {
			for ; heightOfDb+gape <= heightOfNode; heightOfDb = heightOfDb + gape {
				var group sync.WaitGroup
				log.Println("current height ", heightOfDb, "gape=", gape)
				for i := 0; i < gape; i++ {
					group.Add(1)
					go worker(heightOfDb+i, &group)
				}
				group.Wait()
			}
		}
		err = BlockChain(heightOfDb)
	}

}

func worker(height int, group *sync.WaitGroup) {
	defer group.Done()
	err := BlockChain(height)
	if err != nil {
		log.Println(err)
	}
}

type Metas struct {
	BlockMetas []BlockMeta `json:"block_metas"`
}

type BlockMeta struct {
	Hash   string  `json:"Hash"`   // The block hash
	Header *Header `json:"Header"` // The block's Header
}
type BlockID struct {
	Hash []byte `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
}
type Header struct {
	ChainID            string   `protobuf:"bytes,1,opt,name=ChainID,proto3" json:"chain_id,omitempty"`
	Height             int      `protobuf:"varint,2,opt,name=Height,proto3" json:"height,omitempty"`
	Time               int64    `protobuf:"varint,3,opt,name=Time,proto3" json:"time,omitempty"`
	NumTxs             int64    `protobuf:"varint,4,opt,name=NumTxs,proto3" json:"num_txs,omitempty"`
	Maker              []byte   `protobuf:"bytes,5,opt,name=Maker,proto3" json:"maker,omitempty"`
	LastBlockID        *BlockID `protobuf:"bytes,6,opt,name=LastBlockID" json:"last_block_id,omitempty"`
	LastCommitHash     string   `protobuf:"bytes,7,opt,name=LastCommitHash,proto3" json:"last_commit_hash,omitempty"`
	DataHash           string   `protobuf:"bytes,8,opt,name=DataHash,proto3" json:"data_hash,omitempty"`
	ValidatorsHash     string   `protobuf:"bytes,9,opt,name=ValidatorsHash,proto3" json:"validators_hash,omitempty"`
	AppHash            string   `protobuf:"bytes,10,opt,name=AppHash,proto3" json:"app_hash,omitempty"`
	ReceiptsHash       string   `protobuf:"bytes,11,opt,name=ReceiptsHash,proto3" json:"receipts_hash,omitempty"`
	LastNonEmptyHeight int64    `protobuf:"varint,12,opt,name=LastNonEmptyHeight,proto3" json:"last_non_empty_height,omitempty"`
	CoinBase           string   `json:"coin_base,omitempty"`
	BlockRewards       uint64   `protobuf:"varint,14,opt,name=BlockRewards,proto3" json:"block_rewards,omitempty"`
}

func BlockChain(h int) (err error) {
	fmt.Println("current block height : ", h)

	block, errG := GetBlock(h)
	if errG != nil {
		err = errors.New(fmt.Sprintf("GetBlock(height:%d) Error :%v", h, errG))
		return
	}
	//save block
	br := &repository.BlockRepo{
		Blocks: []repository.Block{},
		Txs:    []repository.Transaction{},
	}

	rb := repository.Block{
		Hash:           strings.ToLower(block.BlockMeta.Hash),
		ChainID:        block.Block.Header.ChainID,
		Height:         block.Block.Header.Height,
		Time:           time.Unix(0, block.Block.Header.Time),
		LastCommitHash: strings.ToLower(block.Block.Header.LastCommitHash),
		DataHash:       strings.ToLower(block.Block.Header.DataHash),
		ValidatorsHash: strings.ToLower(block.Block.Header.ValidatorsHash),
		AppHash:        strings.ToLower(block.Block.Header.AppHash),
		Reward:         128,
		CoinBase:       block.Block.Header.CoinBase,
	}
	if len(block.Block.Data.Txs) == 0 && len(block.Block.Data.ExTxs) == 0 {
		br.Blocks = append(br.Blocks, rb)
	} else {
		TxsData := make(map[string]int)
		//unique
		for _, v := range block.Block.Data.Txs {
			TxsData[v] = 0
		}
		//rb.NumTxs = len(block.Block.Data.Txs) + len(block.Block.Data.ExTxs)
		rb.NumTxs = len(TxsData)
		br.Blocks = append(br.Blocks, rb)
		for k, _ := range TxsData {

			rtx, errP := parseTransaction(&rb, k)
			if errP != nil {
				err = errP
				return
			}
			br.Txs = append(br.Txs, rtx)
		}
	}
	err = br.Save()
	if err != nil {
		err = errors.New(fmt.Sprintf("[save failed] %s", err.Error()))
		return
	}

	fmt.Println("  save blocks len:", len(br.Blocks))
	fmt.Println("  save txs len:", len(br.Txs))
	return
}

type BlockTransaction struct {
	GasLimit  *big.Int
	GasPrice  *big.Int
	Nonce     uint64
	Sender    []byte
	Payload   []byte
	Signature []byte
}

func parseTransaction(rb *repository.Block, v string) (rtx repository.Transaction, err error) {
	rtx = repository.Transaction{
		Block:  rb.Hash,
		Time:   rb.Time,
		Height: rb.Height,
	}

	tx := &BlockTransaction{}
	decoder := gob.NewDecoder(bytes.NewReader(common.FromHex(v)))
	err = decoder.Decode(tx)
	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(tx); err != nil {
		return rtx, err
	}

	hash := types.Tx(buf.Bytes()).Hash()
	rtx = repository.Transaction{
		Hash:       hex.EncodeToString(hash),
		PayloadHex: string(tx.Payload),
		Block:      rb.Hash,
		Time:       rb.Time,
		Height:     rb.Height,
	}
	return
}

type ResultBlock struct {
	BlockMeta struct {
		Hash string `json:"hash"`
		PartsHeader struct {
			Total int    `json:"total"`
			Hash  string `json:"hash"`
		} `json:"parts_header"`
	} `json:"block_meta"`
	Block struct {
		Header struct {
			ChainID string `json:"chain_id"`
			Height  int    `json:"height"`
			Time    int64  `json:"time"`
			NumTxs  int    `json:"num_txs"`
			Maker   string `json:"maker"`
			LastBlockID struct {
				Hash        string      `json:"hash"`
				PartsHeader interface{} `json:"parts_header"`
			} `json:"last_block_id"`
			LastCommitHash     string `json:"last_commit_hash"`
			DataHash           string `json:"data_hash"`
			ValidatorsHash     string `json:"validators_hash"`
			AppHash            string `json:"app_hash"`
			ReceiptsHash       string `json:"receipts_hash"`
			LastNonEmptyHeight int    `json:"last_non_empty_height"`
			CoinBase           string `json:"coin_base"`
		} `json:"header"`
		Data struct {
			Txs   []string `json:"txs"`
			ExTxs []string `json:"ex_txs"`
		} `json:"data"`
		LastCommit struct {
			BlockID    interface{}   `json:"block_id"`
			PreCommits []interface{} `json:"pre_commits"`
		} `json:"last_commit"`
		ValidatorSet struct {
			Validators []struct {
				Address     string `json:"address"`
				PubKey      string `json:"pub_key"`
				VotingPower int    `json:"voting_power"`
				Accum       int    `json:"accum"`
				IsCa        bool   `json:"is_ca"`
			} `json:"validators"`
			Proposer         interface{} `json:"proposer"`
			TotalVotingPower int         `json:"total_voting_power"`
		} `json:"validator_set"`
	} `json:"block"`
}

func GetBlock(height int) (result ResultBlock, err error) {
	url := fmt.Sprintf("%s/block?height=%d&chainid=\"%s\"", rpc.HTTP_ADDR, height, repository.ChainID)
	bytez, errB := repository.GetHTTPResp(url)
	if errB != nil {
		err = errB
		return
	}
	err = json.Unmarshal(bytez, &result)
	return
}
