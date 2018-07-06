package repository

import (
	"fmt"
	"testing"
	"time"
)

func TestHeight(t *testing.T) {
	Height()
}

func TestBlockRepo_Save(t *testing.T) {

	br := &BlockRepo{
		Blocks: []Block{},
		Txs:     []Transaction{},
	}

	var i int
	for i = 0; i < 2000000; i++ {
		bl := Block{
			Hash:           "0x1234",
			ChainID:        "cc-id",
			Height:         i + 1,
			Time:           time.Now(),
			NumTxs:         0,
			LastCommitHash: "0x12345",
			DataHash:       "0x12346",
			ValidatorsHash: "0x12347",
			AppHash:        "0x12348",
		}
		br.Blocks = append(br.Blocks, bl)
	}

	tm1 := time.Now().UnixNano() / 1000000
	br.Save()
	tm2 := time.Now().UnixNano() / 1000000

	fmt.Println("time cost:", tm2-tm1)

	CreateIndex()
	tm3 := time.Now().UnixNano() / 1000000
	h, e := Height()
	if e != nil {
		t.Errorf("get height error")
	}
	tm4 := time.Now().UnixNano() / 1000000
	fmt.Println("time cost:", tm4-tm3)
	fmt.Println("height:", h)
}
