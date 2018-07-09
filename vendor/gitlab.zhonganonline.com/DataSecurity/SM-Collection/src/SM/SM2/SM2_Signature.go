package SM2

func SM2_Sign(message []uint8,len int,ZA []uint8,rand []uint8,d []uint8,R []uint8,S []uint8) uint32 {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	a, b, n, p, Gx, Gy = para_a, para_b, para_n, para_p, para_Gx, para_Gy

	var hash [32]uint8
	var M_len int = len + int(SM3_len)/8
	var M []uint8 = nil
	//var i int
	var dA, r, s, e, k, KGx, KGy Big
	var rem, rk, z1, z2 Big
	var KG *Epoint

	dA = Mirvar(0)
	e = Mirvar(0)
	k = Mirvar(0)
	KGx = Mirvar(0)
	KGy = Mirvar(0)
	r = Mirvar(0)
	s = Mirvar(0)
	rem = Mirvar(0)
	rk = Mirvar(0)
	z1 = Mirvar(0)
	z2 = Mirvar(0)
	Bytes_to_big(SM2_NUMWORD, d, dA) //cinstr(dA,d);
	KG = Epoint_init()
	//step1,set M=ZA||M
	//M = (char *)malloc(sizeof(char) * (M_len + 1)) 记号
	M = make([]uint8,M_len+1)
	memcpy(M, ZA, int(SM3_len)/8)
	memcpy(M[SM3_len/8:], message, len)
	//step2,generate e=H(M)
	SM3_256(M, M_len, hash[:])
	Bytes_to_big(int(SM3_len)/8, hash[:], e)
	//step3:generate k
	Bytes_to_big(int(SM3_len)/8, rand, k)
	//step4:calculate kG
	ecurve_mult(k, G, KG)
	//step5:calculate r
	Epoint_get(KG, KGx, KGy)
	Add(e, KGx, r)

	Divide(r, n, rem)
	//judge r=0 or n+k=n?
	Add(r, k, rk)
	if Test_Zero(r)!=0 || Test_n(rk)!=0 {
		return ERR_GENERATE_R
	}
	//step6:generate s
	incr(dA, 1, z1)
	xgcd(z1, n, z1, z1, z1)
	Multiply(r, dA, z2)
	Divide(z2, n, rem)
	subtract(k, z2, z2)
	Add(z2, n, z2)
	Multiply(z1, z2, s)
	Divide(s, n, rem)
	//judge s=0?
	if Test_Zero(s)!=0 {
		return ERR_GENERATE_S
	}
	Big_to_bytes(SM2_NUMWORD, r, R, true)
	Big_to_bytes(SM2_NUMWORD, s, S, true)
	//free(M);
	return 0
}

func SM2_Verify(message []uint8, len int, ZA[]uint8, Px[]uint8, Py[]uint8, R[]uint8, S[]uint8) uint32 {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var hash [32]uint8
	var M_len int = len + int(SM3_len)/8
	var M []uint8 = nil
	var PAx, PAy, r, s, e, t, rem, x1, y1, RR Big
	var PA, sG, tPA *Epoint

	PAx = Mirvar(0)
	PAy = Mirvar(0)
	r = Mirvar(0)
	s = Mirvar(0)
	e = Mirvar(0)
	t = Mirvar(0)
	x1 = Mirvar(0)
	y1 = Mirvar(0)
	rem = Mirvar(0)
	RR = Mirvar(0)
	PA = Epoint_init()
	sG = Epoint_init()
	tPA = Epoint_init()
	Bytes_to_big(SM2_NUMWORD, Px, PAx)
	Bytes_to_big(SM2_NUMWORD, Py, PAy)
	Bytes_to_big(SM2_NUMWORD, R, r)
	Bytes_to_big(SM2_NUMWORD, S, s)
	if Epoint_set(PAx, PAy, 0, PA) == 0 { //initialise public key
		return ERR_PUBKEY_INIT
	}
	//step1: test if r belong to [1,n-1]
	if Test_Range(r)!=0 {
		return ERR_OUTRANGE_R
	}
	//step2: test if s belong to [1,n-1]
	if Test_Range(s)!=0 {
		return ERR_OUTRANGE_S
	}
	//step3,generate M
	M = make([]uint8,M_len+1)
	memcpy(M, ZA, 32)
	memcpy(M[32:], message, len)
	//step4,generate e=H(M)
	SM3_256(M, M_len, hash[:])
	Bytes_to_big(int(SM3_len)/8, hash[:], e)
	//step5:generate t
	Add(r, s, t)
	Divide(t, n, rem)
	if Test_Zero(t)!=0 {
		return ERR_GENERATE_T
	}
	//step 6: generate(x1,y1)
	ecurve_mult(s, G, sG)
	ecurve_mult(t, PA, tPA)
	ecurve_add(sG, tPA)
	Epoint_get(tPA, x1, y1)
	//step7:generate RR
	Add(e, x1, RR)
	Divide(RR, n, rem)
	//free(M);
	if compare(RR, r) == 0 {
		return 0
	} else {
		return ERR_DATA_MEMCMP
	}
}

func SM2_Si(user_priKey []uint8, IDA []uint8, Message []uint8) ([]uint8, []uint8, []uint8, uint32) {
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
	if tmp != 0 {return nil, nil, nil, tmp}

	Epoint_get(PubKey, PubKey_x, PubKey_y)
	Big_to_bytes(SM2_NUMWORD, PubKey_x, PubKeyMerge[:], true)
	Big_to_bytes(SM2_NUMWORD, PubKey_y, PubKeyMerge[SM2_NUMWORD:], true)


	var r = make([]uint8, 32)
	var s = make([]uint8, 32)// Signature424C 49 43 45 31 32 33 40 59 41 48 4F 4F 2E 43 4F 11

	var IDA_len = len(IDA)
	var ENTLA = [2]uint8{ uint8((IDA_len*8)>>8), uint8(IDA_len*8)}      //the length of userA's identification,presentation in ASCII code
	
	var Msg_len int = len(Message)	//the length of Message
	var ZA = make([]uint8, 32)		//ZA=Hash(ENTLA|| IDA|| a|| b|| Gx || Gy || xA|| yA)
	
	N := IDA_len+2+SM2_NUMWORD*6
	var Msg = make([]uint8,N,N)	//210=IDA_len+2+SM2_NUMWORD*6
	// ENTLA || IDA || a || b || Gx || Gy || xA || yA
	memcpy(Msg[:], ENTLA[:], 2)
	memcpy(Msg[2:N], IDA[:], IDA_len)
	memcpy(Msg[2+IDA_len:N], SM2_a[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD:N], SM2_b[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD*2:N], SM2_Gx[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD*3:N], SM2_Gy[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD*4:N], PubKeyMerge[:], SM2_NUMWORD*2)	
	SM3_256(Msg[:], N, ZA[:])

	tmp = SM2_Sign(Message, Msg_len, ZA[:], SM2_rand, user_priKey[:], r[:], s[:])
	if tmp != 0 {return nil, nil, nil, tmp}

	return r, s, PubKeyMerge, 0
}

func SM2_Ve(user_pubKey []uint8, IDA []uint8, Message []uint8, R []uint8, S []uint8) (bool, uint32) {
	if !SM2_INIT_FLAG {
		SM2_Init()
	}

	var tmp uint32 = 0

	var IDA_len = len(IDA)
	var ENTLA = [2]uint8{ uint8((IDA_len*8)>>8), uint8(IDA_len*8)}      //the length of userA's identification,presentation in ASCII code
	
	var Msg_len int = len(Message)	//the length of Message
	var ZA [32]uint8		//ZA=Hash(ENTLA|| IDA|| a|| b|| Gx || Gy || xA|| yA)
	
	N := IDA_len+2+SM2_NUMWORD*6
	Msg := make([]uint8,N,N)	//210=IDA_len+2+SM2_NUMWORD*6
	// ENTLA || IDA || a || b || Gx || Gy || xA || yA
	memcpy(Msg[:], ENTLA[:], 2)
	memcpy(Msg[2:N], IDA[:], IDA_len)
	memcpy(Msg[2+IDA_len:N], SM2_a[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD:N], SM2_b[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD*2:N], SM2_Gx[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD*3:N], SM2_Gy[:], SM2_NUMWORD)
	memcpy(Msg[2+IDA_len+SM2_NUMWORD*4:N], user_pubKey[:], SM2_NUMWORD*2)	
	SM3_256(Msg[:], N, ZA[:])

	tmp = SM2_Verify(Message, Msg_len, ZA[:], user_pubKey[:SM2_NUMWORD], user_pubKey[SM2_NUMWORD:], R, S)
	if tmp != 0 {return false, tmp}

	return true, 0
}
