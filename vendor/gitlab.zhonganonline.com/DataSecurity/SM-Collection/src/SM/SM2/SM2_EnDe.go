package SM2

import "fmt"

func SM2_Encrypt(randK []uint8,pubKey *Epoint,M[] uint8,klen int,C[] uint8) uint32 {
	var C1x, C1y, x2, y2, rand Big
	var C1, kP, S *Epoint
	var i int = 0
	var x2y2 = [64]uint8{0}
	var md SM3_STATE
	C1x = Mirvar(0)
	C1y = Mirvar(0)
	x2 = Mirvar(0)
	y2 = Mirvar(0)
	rand = Mirvar(0)
	C1 = Epoint_init()
	kP = Epoint_init()
	S = Epoint_init()
	//Step2. calculate C1=[k]G=(rGx,rGy)
	Bytes_to_big(SM2_NUMWORD, randK, rand)
	Ecurve_mult(rand, G, C1) //C1=[k]G
	Epoint_get(C1, C1x, C1y)
	Big_to_bytes(SM2_NUMWORD, C1x, C, true)
	Big_to_bytes(SM2_NUMWORD, C1y, C[SM2_NUMWORD:], true)
	//Step3. test if S=[h]pubKey if the point at infinity
	Ecurve_mult(para_h, pubKey, S)
	if Point_at_infinity(S)!=0 { // if S is point at infinity, return error;
		return ERR_INFINITY_POINT
	}
	//Step4. calculate [k]PB=(x2,y2)
	Ecurve_mult(rand, pubKey, kP) //kP=[k]P
	Epoint_get(kP, x2, y2)
	//Step5. KDF(x2||y2,klen)
	Big_to_bytes(SM2_NUMWORD, x2, x2y2[:], true)
	Big_to_bytes(SM2_NUMWORD, y2, x2y2[SM2_NUMWORD:], true)
	SM3_KDF(x2y2[:], 64, uint32(klen), C[SM2_NUMWORD*2:])
	if Test_Null(C[SM2_NUMWORD*2:], klen) != 0 {
		return ERR_ARRAY_NULL
	}
	//Step6. C2=M^t
	for i = 0; i < klen; i++ {
		C[SM2_NUMWORD*2+i] = M[i] ^ C[SM2_NUMWORD*2+i]
	}
	//Step7. C3=hash(x2,M,y2)
	SM3_init(&md)
	SM3_process(&md, x2y2[:], SM2_NUMWORD)
	SM3_process(&md, M, klen)
	SM3_process(&md, x2y2[SM2_NUMWORD:], SM2_NUMWORD)
	SM3_done(&md, C[SM2_NUMWORD*2+klen:])
	return 0
}

func SM2_Decrypt(dB Big,C[] uint8,Clen int,klen int,M[] uint8)uint32 {
	var md SM3_STATE
	var i uint32 = 0
	var x2y2 = [64]uint8{0}
	var hash = [32]uint8{0}
	var C1x, C1y, x2, y2 Big
	var C1, S, dBC1 *Epoint
	C1x = Mirvar(0)
	C1y = Mirvar(0)
	x2 = Mirvar(0)
	y2 = Mirvar(0)
	C1 = Epoint_init()
	S = Epoint_init()
	dBC1 = Epoint_init()
	//Step1. test if C1 fits the curve
	Bytes_to_big(SM2_NUMWORD, C, C1x)
	Bytes_to_big(SM2_NUMWORD, C[SM2_NUMWORD:], C1y)
	Epoint_set(C1x, C1y, 0, C1)
	i = Test_Point(C1)
	if i != 0 {
		return i
	}
	//Step2. S=[h]C1 and test if S is the point at infinity
	Ecurve_mult(para_h, C1, S)
	if Point_at_infinity(S)!=0 { // if S is point at infinity, return error;
		return ERR_INFINITY_POINT
	}
	//Step3. [dB]C1=(x2,y2)
	Ecurve_mult(dB, C1, dBC1)
	Epoint_get(dBC1, x2, y2)
	Big_to_bytes(SM2_NUMWORD, x2, x2y2[:], true)
	Big_to_bytes(SM2_NUMWORD, y2, x2y2[SM2_NUMWORD:], true)
	//Step4. t=KDF(x2||y2,klen)
	SM3_KDF(x2y2[:], uint32(SM2_NUMWORD*2), uint32(klen), M)
	if Test_Null(M, klen) != 0 {
		return ERR_ARRAY_NULL
	}
	//Step5. M=C2^t
	for i = 0; i < uint32(klen); i++ {
		M[i] = M[i] ^ C[uint32(SM2_NUMWORD*2)+i]
	}
	//Step6. hash(x2,m,y2)
	SM3_init(&md)
	SM3_process(&md, x2y2[:], SM2_NUMWORD)
	SM3_process(&md, M, klen)
	SM3_process(&md, x2y2[SM2_NUMWORD:], SM2_NUMWORD)
	SM3_done(&md, hash[:])
	if memcmp(hash[:],C[SM2_NUMWORD*2+klen:],SM2_NUMWORD)!=0 {
		return ERR_C3_MATCH
	} else {
		return 0}
}

/*
	SM2_En
	@input: 
		PriKey [32]uint8
		Message []uint8
*/
func SM2_En(user_priKey []uint8, Message []uint8) ([]uint8, uint32) {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	SM2_rand := Rand_Gen( SM2_n[:] )

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

	//fmt.Printf("PubKeyMerge:: %x\n", PubKeyMerge)

	//encrypt data and compare the result with the standard data	
	Msg_Len := len(Message)
	var Cipher = make([]uint8, Msg_Len+64+32)
	tmp = SM2_Encrypt(SM2_rand[:], PubKey, Message[:], Msg_Len, Cipher[:])
	fmt.Printf("cipher1:%x\n",Cipher[:64])
	fmt.Printf("cipher2:%x\n",Cipher[64: Msg_Len+64])
	fmt.Printf("cipher3:%x\n",Cipher[Msg_Len+64:])
	if tmp != 0 {
		return nil, tmp
	}

	return Cipher, 0
}

/*
	SM2_De
	@input: 
		PriKey [32]uint8
		CyperText []uint8
*/
func SM2_De(user_priKey []uint8, Cipher []uint8) ([]uint8, uint32) {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var tmp uint32 = 0
	var PubKeyMerge = [64]uint8{0}

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

	//fmt.Printf("PubKeyMerge:: %x\n", PubKeyMerge)

	//encrypt data and compare the result with the standard data	
	Cipher_Len := len(Cipher)
	var Message = make([]uint8, Cipher_Len-64-32)

	tmp = SM2_Decrypt(PriKey, Cipher[:], Cipher_Len, Cipher_Len-64-32 ,Message)
	if tmp != 0 {
		return nil, tmp
	}

	return Message, 0
}
