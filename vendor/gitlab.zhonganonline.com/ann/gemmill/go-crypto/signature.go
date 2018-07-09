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

	"gitlab.zhonganonline.com/ann/gemmill/go-wire"
	gcmn "gitlab.zhonganonline.com/ann/gemmill/modules/go-common"
)

// Signature is a part of Txs and consensus Votes.
type Signature interface {
	Bytes() []byte
	IsZero() bool
	String() string
	Equals(Signature) bool
	KeyString() string
}

// Types of Signature implementations
const (
	SignatureTypeEd25519   = byte(0x01)
	SignatureTypeSecp256k1 = byte(0x02)
	SignatureTypeGmsm2     = byte(0x03)
)

// for wire.readReflect
var _ = wire.RegisterInterface(
	struct{ Signature }{},
	wire.ConcreteType{SignatureEd25519{}, SignatureTypeEd25519},
	wire.ConcreteType{SignatureSecp256k1{}, SignatureTypeSecp256k1},
	wire.ConcreteType{SignatureGmsm2{}, SignatureTypeGmsm2},
)

func SignatureFromBytes(sigBytes []byte) (sig Signature, err error) {
	err = wire.ReadBinaryBytes(sigBytes, &sig)
	return
}

//-------------------------------------

// Implements Signature
type SignatureEd25519 [64]byte

func (sig SignatureEd25519) Bytes() []byte {
	return wire.BinaryBytes(struct{ Signature }{sig})
}

func (sig SignatureEd25519) IsZero() bool { return len(sig) == 0 }

func (sig SignatureEd25519) String() string { return fmt.Sprintf("/%X.../", gcmn.Fingerprint(sig[:])) }

func (sig SignatureEd25519) Equals(other Signature) bool {
	if otherEd, ok := other.(SignatureEd25519); ok {
		return bytes.Equal(sig[:], otherEd[:])
	} else {
		return false
	}
}

func (sig SignatureEd25519) KeyString() string {
	return fmt.Sprintf("%X", sig[:])
}

//-------------------------------------

// Implements Signature
type SignatureSecp256k1 []byte

func (sig SignatureSecp256k1) Bytes() []byte {
	return wire.BinaryBytes(struct{ Signature }{sig})
}

func (sig SignatureSecp256k1) IsZero() bool { return len(sig) == 0 }

func (sig SignatureSecp256k1) String() string { return fmt.Sprintf("/%X.../", gcmn.Fingerprint(sig[:])) }

func (sig SignatureSecp256k1) Equals(other Signature) bool {
	if otherEd, ok := other.(SignatureSecp256k1); ok {
		return bytes.Equal(sig[:], otherEd[:])
	} else {
		return false
	}
}

func (sig SignatureSecp256k1) KeyString() string {
	return fmt.Sprintf("%X", sig[:])
}

//-------------------------------------

// SignatureGmsm2 Implements Signature
type SignatureGmsm2 [64]byte

func (sig SignatureGmsm2) Bytes() []byte {
	return wire.BinaryBytes(struct{ Signature }{sig})
}

func (sig SignatureGmsm2) IsZero() bool { return len(sig) == 0 }

func (sig SignatureGmsm2) String() string { return fmt.Sprintf("/%X.../", gcmn.Fingerprint(sig[:])) }

func (sig SignatureGmsm2) Equals(other Signature) bool {
	if otherEd, ok := other.(SignatureGmsm2); ok {
		return bytes.Equal(sig[:], otherEd[:])
	}
	return false
}

func (sig SignatureGmsm2) KeyString() string {
	return fmt.Sprintf("%X", sig[:])
}
