package SM2

var ECC_WORDSIZE uint32 = 8
var SM2_WORDSIZE int = 8
var SM2_NUMBITS int = 256
var SM2_NUMWORD int = 32
var ERR_INFINITY_POINT uint32 = 0x00000001
var ERR_NOT_VALID_ELEMENT uint32 = 0x00000002
var ERR_NOT_VALID_POINT uint32 = 0x00000003
var ERR_ORDER uint32 = 0x00000004
var ERR_ARRAY_NULL uint32 = 0x00000005
var ERR_C3_MATCH uint32 = 0x00000006
var ERR_ECURVE_INIT uint32 = 0x00000007
var ERR_SELFTEST_KG uint32 = 0x00000008
var ERR_SELFTEST_ENC uint32 = 0x00000009
var ERR_SELFTEST_DEC uint32 = 0x0000000A
var ERR_PUBKEY_INIT uint32 = 0x0000000B
var ERR_DATA_MEMCMP uint32 = 0x0000000C
var ERR_GENERATE_R uint32 = 0x0000000D
var ERR_GENERATE_S uint32 = 0x0000000E
var ERR_OUTRANGE_R uint32 = 0x0000000F
var ERR_OUTRANGE_S uint32 = 0x00000010
var ERR_GENERATE_T uint32 = 0x00000011

/* 椭圆曲线方程为： y2 = x3 + ax + b     >>>   Fp-256 */
var SM2_p  =[32]uint8  {0x85,0x42,0xd6,0x9e,0x4c,0x04,0x4F,0x18,0xe8,0xb9,0x24,0x35,0xbf,0x6f,0xf7,0xde,
			0x45,0x72,0x83,0x91,0x5c,0x45,0x51,0x7d,0x72,0x2e,0xdb,0x8b,0x08,0xF1,0xdF,0xc3}
var SM2_a = [32]uint8 {0x78,0x79,0x68,0xb4,0xfa,0x32,0xc3,0xFd,0x24,0x17,0x84,0x2e,0x73,0xbb,0xFe,0xFF,
			0x2F,0x3c,0x84,0x8b,0x68,0x31,0xd7,0xe0,0xec,0x65,0x22,0x8b,0x39,0x37,0xe4,0x98}
var SM2_b = [32]uint8 {0x63,0xE4,0xc6,0xd3,0xb2,0x3b,0x0c,0x84,0x9c,0xf8,0x42,0x41,0x48,0x4b,0xfe,0x48,
			0xF6,0x1d,0x59,0xa5,0xb1,0x6B,0xa0,0x6e,0x6e,0x12,0xD1,0xda,0x27,0xc5,0x24,0x9a}
var SM2_n = [32]uint8 {0x85,0x42,0xD6,0x9E,0x4C,0x04,0x4F,0x18,0xE8,0xB9,0x24,0x35,0xBF,0x6F,0xF7,0xDD,
			0x29,0x77,0x20,0x63,0x04,0x85,0x62,0x8D,0x5A,0xE7,0x4E,0xE7,0xC3,0x2E,0x79,0xB7}
var SM2_Gx = [32]uint8 {0x42,0x1D,0xEB,0xD6,0x1B,0x62,0xEA,0xB6,0x74,0x64,0x34,0xEB,0xC3,0xCC,0x31,0x5E,
			0x32,0x22,0x0B,0x3B,0xAD,0xD5,0x0B,0xDC,0x4C,0x4E,0x6C,0x14,0x7F,0xED,0xD4,0x3D}
var SM2_Gy = [32]uint8 {0x06,0x80,0x51,0x2B,0xCB,0xB4,0x2C,0x07,0xD4,0x73,0x49,0xD2,0x15,0x3B,0x70,0xC4,
			0xE5,0xD7,0xFD,0xFC,0xBF,0xA3,0x6E,0xA1,0xA8,0x58,0x41,0xB9,0xE4,0x6E,0x09,0xA2}
var SM2_h = [32]uint8 {0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
			0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x01}
var para_p, para_a, para_b, para_n, para_Gx, para_Gy, para_h Big
var Gx, Gy, p, a, b, n Big  ////
var G, nG *Epoint
var mip *Miracl

var SM2_INIT_FLAG bool = false

// x3
func SM2_Init() uint32 {
	// ??? 1000 or 10000 ?
	mip = Mirsys(10000, 16)
	mip.IOBASE = 16

	para_p = Mirvar(0)
	para_a = Mirvar(0)
	para_b = Mirvar(0)
	para_n = Mirvar(0)
	para_Gx = Mirvar(0)
	para_Gy = Mirvar(0)
	para_h = Mirvar(0)
	G = Epoint_init()
	nG = Epoint_init()
	Bytes_to_big(SM2_NUMWORD, SM2_p[:], para_p)
	Bytes_to_big(SM2_NUMWORD, SM2_a[:], para_a)
	Bytes_to_big(SM2_NUMWORD, SM2_b[:], para_b)
	Bytes_to_big(SM2_NUMWORD, SM2_n[:], para_n)
	Bytes_to_big(SM2_NUMWORD, SM2_Gx[:], para_Gx)
	Bytes_to_big(SM2_NUMWORD, SM2_Gy[:], para_Gy)
	Bytes_to_big(SM2_NUMWORD, SM2_h[:], para_h)
	ecurve_init(para_a, para_b, para_p, 0)
	//Initialises GF(p) elliptic curve.
	//MR_PROJECTIVE specifying projective coordinates
	if Epoint_set(para_Gx, para_Gy, 0, G) == 0 { //initialise point G
		return ERR_ECURVE_INIT
	}
	Ecurve_mult(para_n, G, nG)
	if Point_at_infinity(nG)==0 { //test if the order of the point is n
		return ERR_ORDER
	}
	SM2_INIT_FLAG = true
	return 0
}

// x3
func SM2_KeyGeneration(priKey Big, pubKey *Epoint) uint32 {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var x, y Big
	x = Mirvar(0)
	y = Mirvar(0)
	Ecurve_mult(priKey, G, pubKey) //通过大数和基点产生公钥
	Epoint_get(pubKey, x, y)
	if Test_PubKey(pubKey) != 0 {
		return ERR_PUBKEY_INIT
	} else {
		return 0
	}
}

func SM2_GetPubKey(user_priKey []uint8) ([]uint8, uint32) {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var tmp uint32 = 0
	var PubKeyMerge = make([]uint8, SM2_NUMWORD*2)

	//generate key pair
	var PriKey, PubKey_x, PubKey_y Big
	var PubKey *Epoint


	PubKey_x = Mirvar(0)
	PubKey_y = Mirvar(0)
	PriKey = Mirvar(0)
	PubKey = Epoint_init()
	Bytes_to_big(len(user_priKey), user_priKey[:], PriKey) //PriKey is the standard private key's Big format

	tmp = SM2_KeyGeneration(PriKey, PubKey)
	if tmp != 0 {return nil, tmp}

	Epoint_get(PubKey, PubKey_x, PubKey_y)
	Big_to_bytes(SM2_NUMWORD, PubKey_x, PubKeyMerge[:], true)
	Big_to_bytes(SM2_NUMWORD, PubKey_y, PubKeyMerge[SM2_NUMWORD:], true)

	return PubKeyMerge, 0
}

// for SM2_EnDe
func Test_Null(array []uint8,len int) int {
	var i int = 0
	for i = 0; i < len; i++ {
		if array[i] != 0x00 {
			return 0
		}
	}
	return 1
}

// for SM2_Signature
func Test_Zero(x Big) int {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var zero Big
	zero = Mirvar(0)
	if compare(x, zero) == 0 {
		return 1
	} else {
		return 0
	}
}

func Test_n(x Big)int {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	// Bytes_to_big(32,SM2_n,n);
	if compare(x, n) == 0 {
		return 1
	} else {
		return 0
	}
}

func Test_Range(x Big)int {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var one, decr_n Big
	one = Mirvar(0)
	decr_n = Mirvar(0)
	convert(1, one)
	decr(n, 1, decr_n)
	if (compare(x, one) < 0) || (compare(x, decr_n) > 0) { //这里原本是(compare(x, one) < 0) | (compare(x, decr_n) > 0)
		return 1
	}
	return 0
}
