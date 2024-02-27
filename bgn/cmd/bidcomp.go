package main

import (
	"fmt"
	"math/big"
	//"reflect"
	//  "math/rand"
	// "log"
	// "log"
	"github.com/sachaservan/bgn"
	//"time"
	"math/rand"
	"time"
	//"miracl/core/BN254"
	//"os"
	//  "io"
	"log"
)

const KEYBITS = 128
const POLYBASE = 5
const MSGSPACE = 1000000000 // message space for polynomial coefficients
const FPSCALEBASE = 3
const FPPREC = 0.0001
const DET = true // deterministic ops

func setup() (*bgn.PublicKey, *bgn.SecretKey, error) {

	start := time.Now()
	pk, sk, err := bgn.NewKeyGen(KEYBITS, big.NewInt(MSGSPACE), POLYBASE, FPSCALEBASE, FPPREC, DET)
	elapsed := time.Since(start)
	log.Printf("time %s", elapsed)
	return pk, sk, err

}

func bid_generate(bid int, pk *bgn.PublicKey, r1a int, r2a int) (*bgn.Ciphertext, *bgn.Ciphertext, *bgn.Ciphertext) {

	// start := time.Now()
	c1 := pk.Encrypt(big.NewInt(int64(bid)))
	//encrypting r_1(A)
	c2 := pk.Encrypt(big.NewInt(int64(r1a)))
	//encrypting r_2(A)
	c3 := pk.Encrypt(big.NewInt(int64(r2a)))
	//    elapsed := time.Since(start)
	//log.Printf("time %s", elapsed)
	return c1, c2, c3

}

func bid_random(a *bgn.Ciphertext, b *bgn.Ciphertext, c *bgn.Ciphertext, pk *bgn.PublicKey, r1b int, r2b int) *bgn.Ciphertext {

	//start := time.Now()
	c4 := pk.Encrypt(big.NewInt(int64(r1b)))

	b1 := pk.Add(b, c4)

	d := pk.Encrypt(big.NewInt(int64(r2b)))

	c5 := pk.Mult(a, b1)

	c6 := pk.Mult(c, d)

	enc := pk.Add(c5, c6)
	// elapsed := time.Since(start)
	// log.Printf("time %s", elapsed)
	return enc

}

func bid_comp(encr *bgn.Ciphertext, pk *bgn.PublicKey, sk *bgn.SecretKey) *big.Int {
	// start := time.Now()
	bgn.ComputeDecryptionPreprocessing(pk, sk)
	dec := sk.DecryptFailSafe(encr, pk)
	//   elapsed := time.Since(start)
	//   log.Printf("time %s", elapsed)
	return dec

}

func encbidgenerate1(r1a int, r2a int, r1b int, r2b int) big.Int {

	pk_i, sk_i, err_i := setup()
	if err_i != nil {
		panic(err_i)
	}

	rand.Seed(time.Now().UnixNano())
	bid := rand.Intn(1000)
	fmt.Println("===plain bid value")
	fmt.Println(bid)

	/* p_min := 0
	   p_max := 10000

	   bid_array_i :=make([]int, param)
	   for i:=0; i<param; i++ {
	   	bid_array_i[i] = rand.Intn(p_max - p_min + 1) + p_min
	   }
	   fmt.Print("\n")
	   fmt.Println("===bids generated w.r.t to various parameters===")
	   fmt.Println(bid_array_i)
	*/

	fmt.Print("\n")
	fmt.Println("===encrypted bid generation")
	//encbid_array_i :=make([]big.Int, param)

	a, b, c := bid_generate(bid, pk_i, r1a, r2a)
	encr := bid_random(a, b, c, pk_i, r1b, r2b)
	fmt.Println("===encrypted bid value")
	fmt.Println(encr)

	decr := bid_comp(encr, pk_i, sk_i)
	// encbid_array_i[i]= *(decr)
	// fmt.Println(reflect.TypeOf(encbid_array_i[i]))

	fmt.Println(decr)

	/*
	   fmt.Print("\n")
	   fmt.Println("===weighted enc array w.r.t to various parameters===")
	   weightedenc_bid_array_i :=make([]big.Int, param)
	   for i:=0; i<param; i++{
	   a := &encbid_array_i[i]
	   b := big.NewInt(int64(WEIGHT[i]))
	   weightedenc_bid_array_i[i]= *big.NewInt(0).Mul(a,b)

	   }
	   fmt.Println(weightedenc_bid_array_i)

	   //generating final score of the domain
	   fmt.Print("\n")
	   fmt.Println("===final score of the domain===")
	   Final_score := big.NewInt(0)
	   for i:=0; i<param; i++{
	   c := &weightedenc_bid_array_i[i]
	   Final_score.Add(Final_score,c)
	   }*/

	return *decr

}

func encbidgenerate2(r1a int, r2a int, r1b int, r2b int) big.Int {

	pk_i, sk_i, err_i := setup()
	if err_i != nil {
		panic(err_i)
	}

	rand.Seed(time.Now().UnixNano())
	bid := rand.Intn(10000)
	fmt.Print("\n")
	fmt.Println("===plain bid value")
	fmt.Println(bid)

	/* p_min := 0
	   p_max := 10000

	   bid_array_i :=make([]int, param)
	   for i:=0; i<param; i++ {
	   	bid_array_i[i] = rand.Intn(p_max - p_min + 1) + p_min
	   }
	   fmt.Print("\n")
	   fmt.Println("===bids generated w.r.t to various parameters===")
	   fmt.Println(bid_array_i)
	*/

	fmt.Println("===encrypted bids generation")
	//encbid_array_i :=make([]big.Int, param)

	a, b, c := bid_generate(bid, pk_i, r1b, r2b)
	encr := bid_random(a, b, c, pk_i, r1a, r2a)
	fmt.Println("===encrypted bid value")
	fmt.Println(encr)
	decr := bid_comp(encr, pk_i, sk_i)
	// encbid_array_i[i]= *(decr)
	// fmt.Println(reflect.TypeOf(encbid_array_i[i]))

	fmt.Println(decr)

	/*
	   fmt.Print("\n")
	   fmt.Println("===weighted enc array w.r.t to various parameters===")
	   weightedenc_bid_array_i :=make([]big.Int, param)
	   for i:=0; i<param; i++{
	   a := &encbid_array_i[i]
	   b := big.NewInt(int64(WEIGHT[i]))
	   weightedenc_bid_array_i[i]= *big.NewInt(0).Mul(a,b)

	   }
	   fmt.Println(weightedenc_bid_array_i)

	   //generating final score of the domain
	   fmt.Print("\n")
	   fmt.Println("===final score of the domain===")
	   Final_score := big.NewInt(0)
	   for i:=0; i<param; i++{
	   c := &weightedenc_bid_array_i[i]
	   Final_score.Add(Final_score,c)
	   }*/

	return *decr

}

func bid_final_score(user_i big.Int, user_j big.Int) (big.Int, big.Int) {

	rand.Seed(time.Now().UnixNano())
	r1a := rand.Intn(10000)
	r2a := rand.Intn(10000)
	r1b := rand.Intn(10000)
	r2b := rand.Intn(10000)

	fmt.Println(r1a, r2a, r1b, r2b)
	//     param := PARAM
	//     WEIGHT :=weight

	user_i = encbidgenerate1(r1a, r2a, r1b, r2b)
	user_j = encbidgenerate2(r1a, r2a, r1b, r2b)
	//elapsed := time.Since(start)
	// log.Printf("time %s", elapsed)
	return user_i, user_j
}

func main() {
	start := time.Now()
	fmt.Println("******Adding weight to the bids of the domains******")
	fmt.Print("\n")
	var user_i big.Int
	var user_j big.Int
	PARAM := 4
	fmt.Println("===weights w.r.t to various parameters===")
	weight := make([]int, PARAM)
	for i := 0; i < 100; i++ {
		rand.Seed(time.Now().UnixNano())
		weight[rand.Intn(PARAM)]++
	}
	fmt.Println(weight)
	n := 4

	domain := make([]big.Int, n)
	for i := 0; i < n; i = i + 2 {
		domain[i], domain[i+1] = bid_final_score(user_i, user_j)

	}


	fmt.Println("finalscore:", domain)
	var k int
	k = 0
	for j := 0; j < n; j++ {

		// if domain[0]<domain[j]
		// if domain[0]<domain[j]
		r := domain[0].Cmp(&domain[j])
		if r == -1 {
			domain[0] = domain[j]
			k = j
		}

	}
	fmt.Println("winning bid:domain", k, "with the value", domain[0])
	elapsed := time.Since(start)
	log.Printf("time %s", elapsed)
}
