// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypto

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"gitlab.zhonganonline.com/DataSecurity/SM-Collection/src/SM/SM2"
	"gitlab.zhonganonline.com/ann/gemmill/ed25519"
	"gitlab.zhonganonline.com/ann/gemmill/ed25519/extra25519"
	"gitlab.zhonganonline.com/ann/gemmill/go-wire"
	gcmn "gitlab.zhonganonline.com/ann/gemmill/modules/go-common"
)

const (
	CryptoTypeED25519   = "ed25519"
	CryptoTypeSecp256k1 = "secp256k1"
	CryptoTypeGMSM2     = "gmsm2"
)

// PrivKey is part of PrivAccount and state.PrivValidator.
type PrivKey interface {
	Bytes() []byte
	Sign(msg []byte) Signature
	PubKey() PubKey
	Equals(PrivKey) bool
	KeyString() string
}

// Types of PrivKey implementations
const (
	PrivKeyTypeEd25519   = byte(0x01)
	PrivKeyTypeSecp256k1 = byte(0x02)
	PrivKeyTypeGmsm2     = byte(0x03)
)

// for wire.readReflect
var _ = wire.RegisterInterface(
	struct{ PrivKey }{},
	wire.ConcreteType{PrivKeyEd25519{}, PrivKeyTypeEd25519},
	wire.ConcreteType{PrivKeySecp256k1{}, PrivKeyTypeSecp256k1},
	wire.ConcreteType{PrivKeyGmsm2{}, PrivKeyTypeGmsm2},
)

func PrivKeyFromBytes(privKeyBytes []byte) (privKey PrivKey, err error) {
	err = wire.ReadBinaryBytes(privKeyBytes, &privKey)
	return
}

//-------------------------------------

// Implements PrivKey
type PrivKeyEd25519 [64]byte

func (privKey PrivKeyEd25519) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeyEd25519) Sign(msg []byte) Signature {
	privKeyBytes := [64]byte(privKey)
	signatureBytes := ed25519.Sign(&privKeyBytes, msg)
	return SignatureEd25519(*signatureBytes)
}

func (privKey PrivKeyEd25519) PubKey() PubKey {
	privKeyBytes := [64]byte(privKey)
	return PubKeyEd25519(*ed25519.MakePublicKey(&privKeyBytes))
}

func (privKey PrivKeyEd25519) Equals(other PrivKey) bool {
	if otherEd, ok := other.(PrivKeyEd25519); ok {
		return bytes.Equal(privKey[:], otherEd[:])
	} else {
		return false
	}
}

func (privKey PrivKeyEd25519) KeyString() string {
	return gcmn.Fmt("%X", privKey[:])
}

func (privKey PrivKeyEd25519) ToCurve25519() *[32]byte {
	keyCurve25519 := new([32]byte)
	privKeyBytes := [64]byte(privKey)
	extra25519.PrivateKeyToCurve25519(keyCurve25519, &privKeyBytes)
	return keyCurve25519
}

func (privKey PrivKeyEd25519) String() string {
	return gcmn.Fmt("PrivKeyEd25519{*****}")
}

// Deterministically generates new priv-key bytes from key.
func (privKey PrivKeyEd25519) Generate(index int) PrivKeyEd25519 {
	newBytes := wire.BinarySha256(struct {
		PrivKey [64]byte
		Index   int
	}{privKey, index})
	var newKey [64]byte
	copy(newKey[:], newBytes)
	return PrivKeyEd25519(newKey)
}

func GenPrivKeyEd25519() PrivKeyEd25519 {
	privKeyBytes := new([64]byte)
	copy(privKeyBytes[:32], CRandBytes(32))
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyEd25519FromSecret(secret []byte) PrivKeyEd25519 {
	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	privKeyBytes := new([64]byte)
	copy(privKeyBytes[:32], privKey32)
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

//-------------------------------------

// PrivKeySecp256k1 Implements PrivKey
type PrivKeySecp256k1 [32]byte

func (privKey PrivKeySecp256k1) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeySecp256k1) Sign(msg []byte) Signature {
	priv__, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	sig__, err := priv__.Sign(Sha256(msg))
	if err != nil {
		gcmn.PanicSanity(err)
	}

	return SignatureSecp256k1(sig__.Serialize())
}

func (privKey PrivKeySecp256k1) PubKey() PubKey {
	_, pub__ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	pub := [64]byte{}
	copy(pub[:], pub__.SerializeUncompressed()[1:])
	return PubKeySecp256k1(pub)
}

func (privKey PrivKeySecp256k1) Equals(other PrivKey) bool {
	if otherSecp, ok := other.(PrivKeySecp256k1); ok {
		return bytes.Equal(privKey[:], otherSecp[:])
	} else {
		return false
	}
}

func (privKey PrivKeySecp256k1) String() string {
	return gcmn.Fmt("PrivKeySecp256k1{*****}")
}

func (privKey PrivKeySecp256k1) KeyString() string {
	return gcmn.Fmt("%X", privKey[:])
}

/*
// Deterministically generates new priv-key bytes from key.
func (key PrivKeySecp256k1) Generate(index int) PrivKeySecp256k1 {
	newBytes := wire.BinarySha256(struct {
		PrivKey [64]byte
		Index   int
	}{key, index})
	var newKey [64]byte
	copy(newKey[:], newBytes)
	return PrivKeySecp256k1(newKey)
}
*/

func GenPrivKeySecp256k1() PrivKeySecp256k1 {
	privKeyBytes := [32]byte{}
	copy(privKeyBytes[:], CRandBytes(32))
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKeyBytes[:])
	copy(privKeyBytes[:], priv.Serialize())
	return PrivKeySecp256k1(privKeyBytes)
}

// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeySecp256k1FromSecret(secret []byte) PrivKeySecp256k1 {
	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey32)
	privKeyBytes := [32]byte{}
	copy(privKeyBytes[:], priv.Serialize())
	return PrivKeySecp256k1(privKeyBytes)
}

//-------------------------------------

var (
	one = new(big.Int).SetInt64(1)

	gmA = new(big.Int).SetBytes(SM2.SM2_a[:])
	gmB = new(big.Int).SetBytes(SM2.SM2_b[:])
	gmN = new(big.Int).SetBytes(SM2.SM2_n[:])
	gmP = new(big.Int).SetBytes(SM2.SM2_p[:])
)

// PrivKeyGmsm2 Implements PrivKey
type PrivKeyGmsm2 [32]byte

func (privKey PrivKeyGmsm2) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeyGmsm2) Sign(msg []byte) Signature {
	r, s, _, errCode := SM2.SM2_Si(privKey[:], nil, msg)
	if errCode != 0 {
		gcmn.PanicSanity(fmt.Errorf("Sign failed: %d", errCode))
	}

	var sig SignatureGmsm2
	copy(sig[:32], r)
	copy(sig[32:], s)

	return sig
}

func (privKey PrivKeyGmsm2) PubKey() PubKey {
	pk, errCode := SM2.SM2_GetPubKey(privKey[:])
	if errCode != 0 {
		gcmn.PanicSanity(fmt.Errorf("Get public key failed: %d", errCode))
	}

	var pubKey PubKeyGmsm2
	copy(pubKey[:], pk)
	return pubKey
}

func (privKey PrivKeyGmsm2) Equals(other PrivKey) bool {
	if otherSecp, ok := other.(PrivKeyGmsm2); ok {
		return bytes.Equal(privKey[:], otherSecp[:])
	} else {
		return false
	}
}

func (privKey PrivKeyGmsm2) String() string {
	return gcmn.Fmt("PrivKeyGmsm2{*****}")
}

func (privKey PrivKeyGmsm2) KeyString() string {
	return gcmn.Fmt("%X", privKey[:])
}

func GenPrivKeyGmsm2() PrivKeyGmsm2 {
	b := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		gcmn.PanicSanity(err)
	}

	k := new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(gmN, one)
	k.Mod(k, n)
	k.Add(k, one)

	var privKey PrivKeyGmsm2
	copy(privKey[:], k.Bytes())
	return privKey
}

func GenPrivkeyByBytes(cryptoType string, data []byte) (PrivKey, error) {
	var privkey PrivKey
	switch cryptoType {
	case CryptoTypeED25519:
		var ed PrivKeyEd25519
		copy(ed[:], data)
		privkey = ed
	case CryptoTypeSecp256k1:
		var sp PrivKeySecp256k1
		copy(sp[:], data)
		privkey = sp
	case CryptoTypeGMSM2:
		var gm PrivKeyGmsm2
		copy(gm[:], data)
		privkey = gm
	default:
		return nil, fmt.Errorf("Unknow crypto type")
	}
	return privkey, nil
}

func GenPrivkeyByType(cryptoType string) (PrivKey, error) {
	var privkey PrivKey
	switch cryptoType {
	case CryptoTypeED25519:
		privkey = GenPrivKeyEd25519()
	case CryptoTypeSecp256k1:
		privkey = GenPrivKeySecp256k1()
	case CryptoTypeGMSM2:
		privkey = GenPrivKeyGmsm2()
	default:
		return nil, fmt.Errorf("Unknow crypto type")
	}
	return privkey, nil
}

// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
// func GenPrivKeyGmsm2FromSecret(secret []byte) PrivKeySecp256k1 {
// 	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
// 	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey32)
// 	privKeyBytes := [32]byte{}
// 	copy(privKeyBytes[:], priv.Serialize())
// 	return PrivKeySecp256k1(privKeyBytes)
// }
