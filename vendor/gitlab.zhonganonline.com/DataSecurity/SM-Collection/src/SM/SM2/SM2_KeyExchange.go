package SM2

func CheckR(R *Epoint) bool {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

 	if Test_Point(R) != 0 {
		return false  ////
	}
	return true
}

func Responser__Get_R_S02_S03(Responser__PriKey []uint8, Requester__R *Epoint, Requester__pubkey []uint8, Requester__Z []uint8, Responser__Z []uint8, klen int) (Big, []uint8, []uint8, []uint8, uint32 ) {
 	if !SM2_INIT_FLAG {
		SM2_Init()
	}

 	/* Step */
 	Responser__r := Cal_r()
 	var Responser__R *Epoint
 	Responser__R = Epoint_init()
	SM2_KeyGeneration(Responser__r, Responser__R)

	var Requester__Pubkey *Epoint
 	Requester__Pubkey = Epoint_init()
 	var Requester__Pubkey_x Big
 	Requester__Pubkey_x = Mirvar(0)
 	var Requester__Pubkey_y Big
 	Requester__Pubkey_y = Mirvar(0)
 	Bytes_to_big(32, Requester__pubkey[0:], Requester__Pubkey_x)
	Bytes_to_big(32, Requester__pubkey[32:], Requester__Pubkey_y)
	Epoint_set(Requester__Pubkey_x, Requester__Pubkey_y, 0, Requester__Pubkey)

 	/* Step */
 	var Responser__x Big
 	Responser__x = Mirvar(0)
 	var Responser__y Big
 	Responser__y = Mirvar(0)
 	Epoint_get(Responser__R, Responser__x, Responser__y)
 	Responser__x_ := Cal_x_(Responser__x)

 	/* Step */
 	var Responser__d Big
 	Responser__d = Mirvar(0)
 	Bytes_to_big(32, Responser__PriKey, Responser__d)
 	Responser__t := Cal_t(Responser__d, Responser__r, Responser__x_)

 	/* Step */
 	var Requester__x Big
 	Requester__x = Mirvar(0)
 	var Requester__y Big
 	Requester__y = Mirvar(0)
 	Epoint_get(Requester__R, Requester__x, Requester__y)
 	Requester__x_ := Cal_x_(Requester__x)

 	/* Step */
 	Responser__UV := Cal_UV(Requester__Pubkey, Requester__R, Requester__x_, Responser__t)
 	if Point_at_infinity(Responser__UV) == 1 {
		return nil, nil, nil, nil, ERR_INFINITY_POINT
	}

 	/* Step */
 	Responser__K := Cal_K( Responser__UV, Requester__Z, Responser__Z, klen)

 	/* Step */
 	var Head uint8 = 0x02
 	Responser__S02 := Cal_S( Head, Responser__UV, Requester__Z, Responser__Z, Requester__R, Responser__R)

 	/* Step */
 	Head = 0x03
 	Responser__S03 := Cal_S( Head, Responser__UV, Requester__Z, Responser__Z, Requester__R, Responser__R)

 	/* Step */
 	return Responser__r, Responser__S02, Responser__S03, Responser__K, 0
}

func Requester__Get_S02_S03(Requester__PriKey []uint8, Requester__r Big, Responser__R *Epoint, Responser__pubkey []uint8, Requester__Z []uint8, Responser__Z []uint8, klen int) ([]uint8, []uint8, []uint8, uint32 ) {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var Requester__R *Epoint
 	Requester__R = Epoint_init()
	SM2_KeyGeneration(Requester__r, Requester__R)

	var Responser__Pubkey *Epoint
 	Responser__Pubkey = Epoint_init()
 	var Responser__Pubkey_x Big
 	Responser__Pubkey_x = Mirvar(0)
 	var Responser__Pubkey_y Big
 	Responser__Pubkey_y = Mirvar(0)
 	Bytes_to_big(32, Responser__pubkey[0:], Responser__Pubkey_x)
	Bytes_to_big(32, Responser__pubkey[32:], Responser__Pubkey_y)
	Epoint_set(Responser__Pubkey_x, Responser__Pubkey_y, 0, Responser__Pubkey)


 	/* Step */
 	var Requester__x Big
 	Requester__x = Mirvar(0)
 	var Requester__y Big
 	Requester__y = Mirvar(0)
 	Epoint_get(Requester__R, Requester__x, Requester__y)
 	Requester__x_ := Cal_x_(Requester__x)

 	/* Step */
 	var Requester__d Big
 	Requester__d = Mirvar(0)
 	Bytes_to_big(32, Requester__PriKey, Requester__d)
 	Requester__t := Cal_t(Requester__d, Requester__r, Requester__x_)

 	/* Step */
 	var Responser__x Big
 	Responser__x = Mirvar(0)
 	var Responser__y Big
 	Responser__y = Mirvar(0)
 	Epoint_get(Responser__R, Responser__x, Responser__y)
 	Responser__x_ := Cal_x_(Responser__x)

 	/* Step */
 	Requester__UV := Cal_UV(Responser__Pubkey, Responser__R, Responser__x_, Requester__t)
 	if Point_at_infinity(Requester__UV) == 1 {
		return nil, nil, nil, ERR_INFINITY_POINT
	}

 	/* Step */
 	Requester__K := Cal_K( Requester__UV, Responser__Z, Requester__Z, klen)

 	/* Step */
 	var Head uint8 = 0x02
 	Requester__S02 := Cal_S( Head, Requester__UV, Responser__Z, Requester__Z, Responser__R, Requester__R)

 	/* Step */
 	Head = 0x03
 	Requester__S03 := Cal_S( Head, Requester__UV, Responser__Z, Requester__Z, Responser__R, Requester__R)

 	/* Step */
 	return Requester__S02, Requester__S03, Requester__K, 0
}

//--------B2: RB=[rb]G=(x2,y2)--------
func Cal_r() Big {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var r_Big Big
	r_Big = Mirvar(0)
	r := Rand_Gen( SM2_n[:] )
	Bytes_to_big(32, r, r_Big)

	return r_Big
}

func Cal_R() *Epoint {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var r_Big Big
	r_Big = Mirvar(0)
	r := Rand_Gen( SM2_n[:] )
	Bytes_to_big(32, r, r_Big)

	// R_Epoint
	var R_Epoint *Epoint
	R_Epoint = Epoint_init()

	//--------B2: RB=[rb]G=(x2,y2)--------
	SM2_KeyGeneration(r_Big, R_Epoint)

	return R_Epoint
}

//--------B3: x_=2^w+x & (2^w-1)--------
func Cal_x_(x Big) Big {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var x_ Big
	x_ = Mirvar(0)

	var tmp Big
	tmp = Mirvar(0)

	var w int = 0
	w = SM2_W(para_n)

	//--------B3: x_=2^w+x & (2^w-1)--------
	expb2(w, x_)        // X_=2^w
	Divide(x, x_, tmp) // x=x mod x_=x & (2^w-1)
	Add(x_, x, x_)

	return x_
}

//--------B4: tB=(dB+x_*rB)mod n--------
func Cal_t(d Big, r Big, x_ Big) Big {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var tmp Big
	tmp = Mirvar(0)

	//--------B4: tB=(dB+x_*rB)mod n--------
	Multiply(x_, r, x_)
	Add(d, x_, x_)
	Divide(x_, para_n, tmp)

	return x_
}

//--------B5: x1_=2^w+x1 & (2^w-1)--------
//	func B3

//--------B6: V=[h*t](OP+[x1_]OR)--------
func Cal_UV(OP *Epoint, OR *Epoint, Ox_ Big, t Big) *Epoint {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	// V_Epoint
	var V *Epoint
	V = Epoint_init()


	//--------B6: V=[h*t](OP+[Ox_]OR)--------
	ecurve_mult(Ox_, OR, V) // v=[Ox_]OR
	ecurve_add(OP, V) // V=OP+V
	Multiply(para_h, t, t) // t=t*h
	ecurve_mult(t, V, V)

	return V
}

//------------B7:KB=KDF(VX,VY,OZ,Z,KLEN)----------
func Cal_K( UV *Epoint, OZ []uint8, Z []uint8, klen int) []uint8 {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var K = make([]uint8, 128/8)
	var HashMerge = make([]uint8, uint32(32*2) + 256/4)
	var UVx Big
	UVx = Mirvar(0)
	var UVy Big
	UVy = Mirvar(0)

	//------------B7:KB=KDF(VX,VY,OZ,Z,KLEN)----------
	Epoint_get(UV, UVx, UVy)
	Big_to_bytes(32, UVx, HashMerge[0 : 32], true)
	Big_to_bytes(32, UVy, HashMerge[32 : 64], true)
	memcpy(HashMerge[32*2 : 32*2+256/8], OZ, 256 / 8);
	memcpy(HashMerge[32*2+256/8 : ], Z, 256 / 8);
	SM3_KDF(HashMerge[0:], uint32(32*2)+256/4, uint32(klen)/8, K)

	return K
}

//--------------- B8:(optional) -------------
// SB=hash(0x02||Vy||HASH(Vx||ZA||ZB||x1||y1||x2||y2
func Cal_S( Head uint8, V *Epoint, ZA_String []uint8, ZB_String []uint8, RA *Epoint, RB *Epoint) []uint8 {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}
	var md SM3_STATE

	var Vx Big
	Vx = Mirvar(0)
	var Vy Big
	Vy = Mirvar(0)
	Epoint_get(V, Vx, Vy)

	var RAx Big
	RAx = Mirvar(0)
	var RAy Big
	RAy = Mirvar(0)
	Epoint_get(RA, RAx, RAy)

	var RBx Big
	RBx = Mirvar(0)
	var RBy Big
	RBy = Mirvar(0)
	Epoint_get(RB, RBx, RBy)

	var S = make([]uint8, 32)
	var HashMerge = make([]uint8, 32 + int(256)/8 + int(256)/8 + 32*2 + 32*2)
	var Hash_x = make([]uint8, 32)
	var flag_left int = 0
	var flag_right int = 32
	Big_to_bytes(flag_right - flag_left, Vx, HashMerge[flag_left:flag_right], true)

	flag_left = flag_right
	flag_right += int(256)/8
	memcpy(ZA_String, HashMerge[flag_left:flag_right], flag_right - flag_left)

	flag_left = flag_right
	flag_right += int(256)/8
	memcpy(ZB_String, HashMerge[flag_left:flag_right], flag_right - flag_left)

	flag_left = flag_right
	flag_right += 32
	Big_to_bytes(flag_right - flag_left, RAx, HashMerge[flag_left:flag_right], true)

	flag_left = flag_right
	flag_right += 32
	Big_to_bytes(flag_right - flag_left, RAy, HashMerge[flag_left:flag_right], true)

	flag_left = flag_right
	flag_right += 32
	Big_to_bytes(flag_right - flag_left, RBx, HashMerge[flag_left:flag_right], true)

	flag_left = flag_right
	flag_right += 32
	Big_to_bytes(flag_right - flag_left, RBy, HashMerge[flag_left:flag_right], true)

	SM3_init(&md)
	SM3_process(&md, HashMerge, flag_right)
	SM3_done(&md, Hash_x)

	var tmpArr = []uint8{Head}
	HashMerge = make([]uint8, 1 + 32 + int(256)/8)
	flag_left = 0
	flag_right = 1
	memcpy(tmpArr, HashMerge[flag_left:flag_right], flag_right - flag_left)

	flag_left = flag_right
	flag_right += 32
	Big_to_bytes(flag_right - flag_left, Vy, HashMerge[flag_left:flag_right], true)

	flag_left = flag_right
	flag_right += int(256)/8
	memcpy(Hash_x, HashMerge[flag_left:flag_right], flag_right - flag_left)

	SM3_init(&md)
	SM3_process(&md, HashMerge, flag_right)
	SM3_done(&md, S)

	return S
}

func SM2_W(n Big) int {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var n1 Big
	var w int = 0
	n1 = Mirvar(0)
	w = logb2(para_n) //approximate integer log to the base 2 of para_n
	expb2(w, n1)      //n1=2^w
	if compare(para_n, n1) == 1 {
		w++
	}
	if (w % 2) == 0 {
		w = w/2 - 1
	} else {
		w = (w+1)/2 - 1
	}
	return w
}

func SM3_Z(ID []uint8, ELAN uint16, pubKey *Epoint, hash []uint8) {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var Px = [32]uint8{0}
	var Py = [32]uint8{0}
	var IDlen = [2]uint8{0}
	var x, y Big
	var md *SM3_STATE=&SM3_STATE{}
	x = Mirvar(0)
	y = Mirvar(0)
	Epoint_get(pubKey, x, y)
	Big_to_bytes(32, x, Px[0:], true)
	Big_to_bytes(32, y, Py[0:], true)
	var di = []uint8{}
	var gao = []uint8{}

	di = append(di, uint8(ELAN%(uint16(1<<8))))
	gao = append(gao, uint8(ELAN/(uint16(1<<8))))
	memcpy(IDlen[0:], gao[0:], 1);
	memcpy(IDlen[1:], di[0:], 1);
	SM3_init(md)
	SM3_process(md, IDlen[0:], 2)
	SM3_process(md, ID, int(ELAN)/8)
	SM3_process(md, SM2_a[0:], 32)
	SM3_process(md, SM2_b[0:], 32)
	SM3_process(md, SM2_Gx[0:], 32)
	SM3_process(md, SM2_Gy[0:], 32)
	SM3_process(md, Px[0:], 32)
	SM3_process(md, Py[0:], 32)
	SM3_done(md, hash)
	return
}

func SM2_Gen_Z(ID []uint8, user_PriKey []uint8) ([]uint8, uint32) {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}
	
	var tmp uint32 = 0
	var PubKeyMerge = make([]uint8, SM2_NUMWORD*2)

	var PriKey, PubKey_x, PubKey_y Big
	PubKey_x = Mirvar(0)
	PubKey_y = Mirvar(0)
	PriKey = Mirvar(0)

	var PubKey *Epoint
	PubKey = Epoint_init()

	Bytes_to_big(len(user_PriKey), user_PriKey[:], PriKey)

	tmp = SM2_KeyGeneration(PriKey, PubKey)
	if tmp != 0 {return nil, tmp}

	Epoint_get(PubKey, PubKey_x, PubKey_y)
	Big_to_bytes(SM2_NUMWORD, PubKey_x, PubKeyMerge[:], true)
	Big_to_bytes(SM2_NUMWORD, PubKey_y, PubKeyMerge[SM2_NUMWORD:], true)

	var ID_len = len(ID)
	var ENTL = [2]uint8{ uint8((ID_len*8)>>8), uint8(ID_len*8)}      //the length of userA's identification,presentation in ASCII code

	var Z = make([]uint8, 32)	//Z=Hash(ENTL|| ID|| a|| b|| Gx || Gy || xA|| yA)
	
	N := ID_len+2+SM2_NUMWORD*6
	var Msg = make([]uint8,N,N)	//210=ID_len+2+SM2_NUMWORD*6
	// ENTL || ID || a || b || Gx || Gy || xA || yA
	memcpy(Msg[:], ENTL[:], 2)
	memcpy(Msg[2:N], ID[:], ID_len)
	memcpy(Msg[2+ID_len:N], SM2_a[:], SM2_NUMWORD)
	memcpy(Msg[2+ID_len+SM2_NUMWORD:N], SM2_b[:], SM2_NUMWORD)
	memcpy(Msg[2+ID_len+SM2_NUMWORD*2:N], SM2_Gx[:], SM2_NUMWORD)
	memcpy(Msg[2+ID_len+SM2_NUMWORD*3:N], SM2_Gy[:], SM2_NUMWORD)
	memcpy(Msg[2+ID_len+SM2_NUMWORD*4:N], PubKeyMerge[:], SM2_NUMWORD*2)	
	SM3_256(Msg[:], N, Z[:])

	return Z, 0
}
