package SM2

import (
	"fmt"
	"math/rand"
	"math/big"
	"time"
	"strconv"
)

/* 
*  <随机数生成函数> *
*  @input:椭圆方程阶数 n ( Big->[]uint8 )
*  @output:产生的随机数 rand ( Big->[]uint8)
*  使用方法: 
*  SM2_rand_SelfGen := Rand_Gen( SM2_n[:] )
*/

func Rand_Gen(SM2_n_Bytes []uint8) []uint8{
	SM2_n_Str := Bytes_to_String(SM2_n_Bytes)

	SM2_n_Big := new(big.Int)
	SM2_rand_Big := new(big.Int)
	rand_Seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	SM2_n_Big.SetString(SM2_n_Str, 16)
	SM2_rand_Big.Rand(rand_Seed, SM2_n_Big)
	SM2_rand_Str := SM2_rand_Big.Text(16)

	return String_to_Bytes(SM2_rand_Str)
}

func Bytes_to_String(target []uint8) string{
	var n = len(target)
	var temp string = ""
	var str string = ""
	for i := 0; i < n; i++ {
		temp = strconv.FormatUint(uint64(target[i]), 16)
		for ; len(temp) <2 ; {
			temp = "0"+temp
		}
		str += temp
	}
	return str
}

func String_to_Bytes(target string) []uint8{
	var substr string
	var result_slice []uint8 = make([]uint8, 32)

	for ; len(target) <2*32 ; {
		target = "0"+target
	}

	for i, n := 0, 32; i < n; i++ {
		substr = target[i*2:i*2+2]
		val, err := strconv.ParseUint(substr, 16, 32)
		if err==nil {
			result_slice[i] = uint8(val)
		}else{
			fmt.Printf("Error:: strconv.ParseUint(substr, 16, 32), i=%d, n=%d \n", i, n)
		}
	}
	return result_slice
}

/* </随机数生成函数> */
