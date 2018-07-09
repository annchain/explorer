package job

import (
"github.com/annchain/annchain/eth/rlp"
)

func ToBytes(tx interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(tx)
}

func FromBytes(bs []byte, tx interface{}) error {
	return rlp.DecodeBytes(bs, tx)
}










