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
	"fmt"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"gitlab.zhonganonline.com/DataSecurity/SM-Collection/src/SM/SM2"
	"gitlab.zhonganonline.com/ann/gemmill/ed25519"
	"gitlab.zhonganonline.com/ann/gemmill/ed25519/extra25519"
	"gitlab.zhonganonline.com/ann/gemmill/go-wire"
	gcmn "gitlab.zhonganonline.com/ann/gemmill/modules/go-common"
	"golang.org/x/crypto/ripemd160"
)

// PubKey is part of Account and Validator.
type PubKey interface {
	Address() []byte
	Bytes() []byte
	KeyString() string
	VerifyBytes(msg []byte, sig Signature) bool
	Equals(PubKey) bool
}

// Types of PubKey implementations
const (
	PubKeyTypeEd25519   = byte(0x01)
	PubKeyTypeSecp256k1 = byte(0x02)
	PubKeyTypeGmsm2     = byte(0x03)
)

// for wire.readReflect
var _ = wire.RegisterInterface(
	struct{ PubKey }{},
	wire.ConcreteType{PubKeyEd25519{}, PubKeyTypeEd25519},
	wire.ConcreteType{PubKeySecp256k1{}, PubKeyTypeSecp256k1},
	wire.ConcreteType{PubKeyGmsm2{}, PubKeyTypeGmsm2},
)

func PubKeyFromBytes(pubKeyBytes []byte) (pubKey PubKey, err error) {
	err = wire.ReadBinaryBytes(pubKeyBytes, &pubKey)
	return
}

//-------------------------------------

// Implements PubKey
type PubKeyEd25519 [32]byte

func (pubKey PubKeyEd25519) Address() []byte {
	w, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(pubKey[:], w, n, err)
	if *err != nil {
		gcmn.PanicCrisis(*err)
	}
	// append type byte
	encodedPubkey := append([]byte{PubKeyTypeEd25519}, w.Bytes()...)
	hasher := ripemd160.New()
	hasher.Write(encodedPubkey) // does not error
	return hasher.Sum(nil)
}

func (pubKey PubKeyEd25519) Bytes() []byte {
	return wire.BinaryBytes(struct{ PubKey }{pubKey})
}

func (pubKey PubKeyEd25519) VerifyBytes(msg []byte, sig_ Signature) bool {
	sig, ok := sig_.(SignatureEd25519)
	if !ok {
		return false
	}
	pubKeyBytes := [32]byte(pubKey)
	sigBytes := [64]byte(sig)
	return ed25519.Verify(&pubKeyBytes, msg, &sigBytes)
}

// For use with golang/crypto/nacl/box
// If error, returns nil.
func (pubKey PubKeyEd25519) ToCurve25519() *[32]byte {
	keyCurve25519, pubKeyBytes := new([32]byte), [32]byte(pubKey)
	ok := extra25519.PublicKeyToCurve25519(keyCurve25519, &pubKeyBytes)
	if !ok {
		return nil
	}
	return keyCurve25519
}

func (pubKey PubKeyEd25519) String() string {
	return gcmn.Fmt("PubKeyEd25519{%X}", pubKey[:])
}

// Must return the full bytes in hex.
// Used for map keying, etc.
func (pubKey PubKeyEd25519) KeyString() string {
	return gcmn.Fmt("%X", pubKey[:])
}

func (pubKey PubKeyEd25519) Equals(other PubKey) bool {
	if otherEd, ok := other.(PubKeyEd25519); ok {
		return bytes.Equal(pubKey[:], otherEd[:])
	} else {
		return false
	}
}

//-------------------------------------

// Implements PubKey
type PubKeySecp256k1 [64]byte

func (pubKey PubKeySecp256k1) Address() []byte {
	w, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(pubKey[:], w, n, err)
	if *err != nil {
		gcmn.PanicCrisis(*err)
	}
	// append type byte
	encodedPubkey := append([]byte{PubKeyTypeSecp256k1}, w.Bytes()...)
	hasher := ripemd160.New()
	hasher.Write(encodedPubkey) // does not error
	return hasher.Sum(nil)
}

func (pubKey PubKeySecp256k1) Bytes() []byte {
	return wire.BinaryBytes(struct{ PubKey }{pubKey})
}

func (pubKey PubKeySecp256k1) VerifyBytes(msg []byte, sig_ Signature) bool {
	pub__, err := secp256k1.ParsePubKey(append([]byte{0x04}, pubKey[:]...), secp256k1.S256())
	if err != nil {
		return false
	}
	sig, ok := sig_.(SignatureSecp256k1)
	if !ok {
		return false
	}
	sig__, err := secp256k1.ParseDERSignature(sig[:], secp256k1.S256())
	if err != nil {
		return false
	}
	return sig__.Verify(Sha256(msg), pub__)
}

func (pubKey PubKeySecp256k1) String() string {
	return gcmn.Fmt("PubKeySecp256k1{%X}", pubKey[:])
}

// Must return the full bytes in hex.
// Used for map keying, etc.
func (pubKey PubKeySecp256k1) KeyString() string {
	return gcmn.Fmt("%X", pubKey[:])
}

func (pubKey PubKeySecp256k1) Equals(other PubKey) bool {
	if otherSecp, ok := other.(PubKeySecp256k1); ok {
		return bytes.Equal(pubKey[:], otherSecp[:])
	} else {
		return false
	}
}

//-------------------------------------

// PubKeyGmsm2 Implements PubKey
type PubKeyGmsm2 [64]byte

func (pubKey PubKeyGmsm2) Address() []byte {
	w, n, err := new(bytes.Buffer), new(int), new(error)
	wire.WriteBinary(pubKey[:], w, n, err)
	if *err != nil {
		gcmn.PanicCrisis(*err)
	}
	// append type byte
	encodedPubkey := append([]byte{PubKeyTypeGmsm2}, w.Bytes()...)
	hasher := ripemd160.New()
	hasher.Write(encodedPubkey) // does not error
	return hasher.Sum(nil)
}

func (pubKey PubKeyGmsm2) Bytes() []byte {
	return wire.BinaryBytes(struct{ PubKey }{pubKey})
}

func (pubKey PubKeyGmsm2) VerifyBytes(msg []byte, sig_ Signature) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("there is a bug in the SM2.SM2_Ve")
			// SM2.SM2_Ve has a bug, running verify only will cause panic
		}
	}()

	sig, ok := sig_.(SignatureGmsm2)
	if !ok {
		return false
	}

	r := make([]byte, 32)
	s := make([]byte, 32)
	copy(r, sig[:32])
	copy(s, sig[32:])

	valid, errCode := SM2.SM2_Ve(pubKey[:], nil, msg, r, s)
	if errCode != 0 || !valid {
		return false
	}

	return true
}

func (pubKey PubKeyGmsm2) String() string {
	return gcmn.Fmt("PubKeyGmsm2{%X}", pubKey[:])
}

// Must return the full bytes in hex.
// Used for map keying, etc.
func (pubKey PubKeyGmsm2) KeyString() string {
	return gcmn.Fmt("%X", pubKey[:])
}

func (pubKey PubKeyGmsm2) Equals(other PubKey) bool {
	if otherSecp, ok := other.(PubKeyGmsm2); ok {
		return bytes.Equal(pubKey[:], otherSecp[:])
	} else {
		return false
	}
}
