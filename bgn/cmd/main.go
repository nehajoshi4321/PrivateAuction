package main

import (
	"fmt"
	"math/big"
	//  "math/rand"
	// "log"
	"github.com/sachaservan/bgn"
	"log"
	"math/rand"
	"time"
	//"time"
	//"miracl/core/BN254"
	"io"
	"os"
)

const KEYBITS = 512
const POLYBASE = 3
const MSGSPACE = 10000 // message space for polynomial coefficients
const FPSCALEBASE = 3
const FPPREC = 0.0001
const DET = true // deterministic ops

/*func main() {

	printWelcome()

	//keyBits := 512 // length of q1 and q2
	//messageSpace := big.NewInt(1021)
	//polyBase := 3 // base for the ciphertext polynomial
	//fpScaleBase := 3
	//fpPrecision := 0.0001
	bids :=30
	r1aa :=40
	r2aa :=30
	r1bb :=3
	r2bb :=8


	//runSimpleCheck(keyBits, polyBase)
	//runPolyArithmeticCheck(keyBits, messageSpace, polyBase, fpScaleBase, fpPrecision)
	auction_value :=auction_bid(bids, r1aa, r2aa, r1bb, r2bb)

	fmt.Println(auction_value)
}

*/

func printWelcome() {
	fmt.Println("====================================")
	fmt.Println(" ____   _____ _   _ ")
	fmt.Println("|  _ \\ / ____| \\ | |")
	fmt.Println("| |_) | |  __ |  \\| |")
	fmt.Println("|  _ <| | |_  | . `  |")
	fmt.Println("| |_) | |__|  | |\\  |")
	fmt.Println("|____/ \\_____|_| \\_|")

	fmt.Println("Boneh Goh Nissim Cryptosystem in Go")
	fmt.Println("====================================")

}
func auction_bid(bid int, r1a int, r2a int, r1b int, r2b int) big.Int {

	start := time.Now()
	pk, sk, err := bgn.NewKeyGen(KEYBITS, big.NewInt(MSGSPACE), POLYBASE, FPSCALEBASE, FPPREC, DET)
	if err != nil {
		panic(err)
	}
	bgn.ComputeDecryptionPreprocessing(pk, sk)

	//encrypting a bid
	c1_a := pk.Encrypt(big.NewInt(int64(bid)))
	//encrypting r_1(A)
	c2_a := pk.Encrypt(big.NewInt(int64(r1a)))
	//encrypting r_2(A)
	c3_a := pk.Encrypt(big.NewInt(int64(r2a)))

	//encrypting r_1(B)
	c2_b := (pk.Encrypt(big.NewInt(int64(r1b))))
	//encrypting r_2(B)
	c3_b := (pk.Encrypt(big.NewInt(int64(r2b))))

	//Add r_1(A) and r_2(A)
	c4_a := pk.Add(c2_a, c2_b)

	//multiply r_1(A) and r_2(A) with encrypted bid
	c5_a := pk.Mult(c1_a, c4_a)

	//Add r_1(B) and r_2(B)
	c4_b := pk.Mult(c3_a, c3_b)

	//multiply r_1(B) and r_2(B) with overall encrypted bid
	c5_b := pk.Add(c5_a, c4_b)

	//final user encrypted bid
	//bid := sk.DecryptFailSafe(c5_b, pk)

	d := sk.DecryptFailSafe(c1_a, pk)
	c := sk.DecryptFailSafe(c4_a, pk)
	e := sk.DecryptFailSafe(c5_a, pk)
	f := sk.DecryptFailSafe(c4_b, pk)
	g := sk.DecryptFailSafe(c5_b, pk)

	fmt.Println("working", d.String())
	//fmt.Println("working",bid.String())
	fmt.Println("working", c.String())
	fmt.Println("working", e.String())
	fmt.Println("working", f.String())
	fmt.Println("working", g.String())

	elapsed := time.Since(start)
	log.Printf("time %s", elapsed)
	return *g

}

//auction operation
func auction_winner(dslice []big.Int) big.Int {
	//fmt.Println("Domain bid:", dslice)
	//to find the min bid value among domains
	//var min big.Int
	min := dslice[0]
	size := len(dslice)
	// var x *big.Int
	// var y *big.Int

	for i := 0; i < size; i++ {
		fmt.Println("%d ", dslice[i])

		//x= &dslice[i]
		//y= *x.Cmp(&min)

		// adr_mul= & domain_slice[i][j]
		// adr_mul.Mul(adr_mul, &(num[i]))
		// dslice[i] > min
		//if (dslice[i] < min) {
		//
		//    min = dslice[i]
		//}

	}
	fmt.Println("maximum bid: %d", min)
	return min
}

func readFile(filePath string) (numbers []int) {
	fd, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", filePath, err))
	}
	var line int
	for {

		_, err := fmt.Fscanf(fd, "%d\n", &line)

		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				return
			}
			panic(fmt.Sprintf("Scan Failed %s: %v", filePath, err))

		}
		numbers = append(numbers, line)
	}
	return
}

func optimal_domain(dom int, para int) []big.Int {

	//generate random integers
	rand.Seed(time.Now().UnixNano())
	//cost bid of resources
	// c_min := 1
	//c_max := 10000000
	//bandwidth of resources(in mbps)
	p_min := 0
	p_max := 1000
	//generate random numbers for both parties
	r1aa := rand.Int()
	r2aa := rand.Int()
	r1bb := rand.Int()
	r2bb := rand.Int()

	// no. of domain
	//var n, m int
	//var bw_weight int
	// var latency_weight int
	// var ploss_weight int

	// fmt.Print("enter no. of domains=")
	//  fmt.Scanln(&m)
	//  fmt.Print("enter no. of parameters=")
	// fmt.Scanln(&n)
	//weightage given to various parameter for resources
	//   fmt.Print("assign weight to the parameters:")
	/*   fmt.Scanln(&bw_weight)
	     fmt.Scanln(&latency_weight)
	     fmt.Scanln(&ploss_weight)  */
	fmt.Print("\n")
	fmt.Println("******Assigning bids to the domains******")
	fmt.Print("\n")
	domain_slice := [][]big.Int{}
	//  var bid_array[5] int

	//assign bid to the domain
	for i := 0; i < para; i++ {
		p_row := make([]big.Int, dom)
		for j := 0; j < dom; j++ {
			// domain_slice[i][j]= rand.Intn(p_max - p_min + 1) + p_min
			bids := rand.Intn(p_max-p_min+1) + p_min
			p_row[j] = auction_bid(bids, r1aa, r2aa, r1bb, r2bb)
		}
		//fmt.Println("Printing current Row", p_row)
		domain_slice = append(domain_slice, p_row)

	}
	fmt.Println("----Printing current 2d array----")
	fmt.Println(domain_slice)

	fmt.Print("\n")
	fmt.Println("******Adding weight to the bids of the domains******")
	fmt.Println("===weights w.r.t to various parameters===")
	num := make([]big.Int, para)

	for i := 0; i < 100; i++ {
		//num[rand.Intn(para)]++
	}
	fmt.Println(num)

	var adr *big.Int
	var adr_mul *big.Int
	//  fmt.Println("Printing current 2d array", domain_slice[0][0])
	//domain_slice1 := [][]int{}
	for i := 0; i < para; i++ {
		for j := 0; j < dom; j++ {
			adr_mul = &domain_slice[i][j]
			adr_mul.Mul(adr_mul, &(num[i]))
			//domain_slice[i][j]= domain_slice[i][j] * num[i]
			//p_row[j] = p_row[j]* num[j]

		}

	}
	fmt.Print("\n")
	fmt.Println("----Printing current 2d weighted array----")
	fmt.Println(domain_slice)
	fmt.Print("\n")

	fmt.Println("----selection of optimal domain----")
	fmt.Print("\n")
	// domain_sum:=0
	domain_sum := make([]big.Int, dom)

	for i := 0; i < para; i++ {
		//sum :=0
		for j := 0; j < dom; j++ {
			adr = &domain_sum[j]
			adr.Add(adr, &(domain_slice[i][j]))
			//domain_sum[j] = domain_sum[j] + domain_slice[i][j]
			// sum = sum + domain_slice[i][j]
			// fmt.Println("The Sum of Each Column Item in a Matrix  = ", domain_sum)
		}

	}
	fmt.Println(domain_sum)

	fmt.Print("\n")
	//MDO performs auction operation

	return domain_sum

}

func main() {
	var m, n int

	//r1aa :=rand.Int()
	//r2aa :=rand.Int()
	//r1bb :=rand.Int()
	//r2bb :=rand.Int()

	printWelcome()

	//keyBits := 512 // length of q1 and q2
	//messageSpace := big.NewInt(1021)
	//polyBase := 3 // base for the ciphertext polynomial
	//fpScaleBase := 3
	//fpPrecision := 0.0001

	//runSimpleCheck(keyBits, polyBase)
	//runPolyArithmeticCheck(keyBits, messageSpace, polyBase, fpScaleBase, fpPrecision)

	numbers := readFile("foo.in.txt")
	//  fmt.Println(numbers[0])
	m = numbers[0]
	n = numbers[1]
	fmt.Println("no. of domains are:", m)
	fmt.Println("no. of parameters are:", n)
	opt_domain_array := optimal_domain(m, n)

	min_bid := auction_winner(opt_domain_array)
	fmt.Printf("Maximum bid: %d", min_bid)

}

//user bids

// fmt.Println(domain_slice)

//  min_bid :=auction(domain_slice)
// fmt.Printf("Minimum bid: %d", min_bid)

/*
func auctionencoding() {

      // p_min := 0
    //p_max := 10000
      bid :=90
       pk, sk, err := bgn.NewKeyGen(KEYBITS, big.NewInt(MSGSPACE), POLYBASE, FPSCALEBASE, FPPREC, DET)
	if err != nil {
		panic(err)
	}
	bgn.ComputeDecryptionPreprocessing(pk, sk)

	//encrypting a bid
	c1_a := pk.Encrypt(big.NewInt(int64(bid)))
	//encrypting r_1(A)
	c2_a := pk.Encrypt(big.NewInt(40))
	//encrypting r_2(A)
	c3_a := pk.Encrypt(big.NewInt(2))

	//Add r_1(A) and r_2(A)
	c4_a :=pk.Add(c2_a, c3_a)

	//multiply r_1(A) and r_2(A) with encrypted bid
	c5_a :=pk.Mult(c1_a, c4_a)

	//encrypting r_1(B)
	c2_b := (pk.Encrypt(big.NewInt(5)))
	//encrypting r_2(B)
	c3_b := (pk.Encrypt(big.NewInt(300)))

	//Add r_1(B) and r_2(B)
	c4_b :=pk.Mult(c2_b, c3_b)

	//multiply r_1(B) and r_2(B) with overall encrypted bid
	c5_b :=pk.Add(c5_a, c4_b)

	//final user encrypted bid
	//bid := sk.DecryptFailSafe(c5_b, pk)


	d := sk.DecryptFailSafe(c1_a, pk)
	c := sk.DecryptFailSafe(c4_a, pk)
	e := sk.DecryptFailSafe(c5_a, pk)
	f := sk.DecryptFailSafe(c4_b, pk)
	g := sk.DecryptFailSafe(c5_b, pk)

	fmt.Println("working",d.String())
	//fmt.Println("working",bid.String())
	fmt.Println("working",c.String())
	fmt.Println("working",e.String())
	fmt.Println("working",f.String())
	fmt.Println("working",g.String())

}
func runPolyArithmeticCheck(keyBits int, messageSpace *big.Int, polyBase int, fpScaleBase int, fpPrecision float64) {

	pk, sk, _ := bgn.NewKeyGen(keyBits, messageSpace, polyBase, fpScaleBase, fpPrecision, true)
	bgn.ComputeDecryptionPreprocessing(pk, sk)

	m1 := pk.NewPolyPlaintext(big.NewFloat(0.0111))
	m2 := pk.NewPolyPlaintext(big.NewFloat(9.1))
	m3 := pk.NewPolyPlaintext(big.NewFloat(2.75))
	m4 := pk.NewPolyPlaintext(big.NewFloat(2.99))

	c1 := pk.EncryptPoly(m1)
	c2 := pk.EncryptPoly(m2)
	c3 := pk.EncryptPoly(m3)
	c4 := pk.EncryptPoly(m4)
	c6 := pk.NegPoly(c4)

	print("\n----------RUNNING ARITHMETIC TEST----------\n\n")

	fmt.Printf("c1 = E(%s)\n", sk.DecryptPoly(c1, pk).String())
	fmt.Printf("c2 = E(%s)\n", sk.DecryptPoly(c2, pk).String())
	fmt.Printf("c3 = E(%s)\n", sk.DecryptPoly(c3, pk).String())
	fmt.Printf("c4 = E(%s)\n", sk.DecryptPoly(c4, pk).String())
	fmt.Println()

	r1 := pk.AddPoly(c1, c4)
	fmt.Printf("[Add] E(%s) ⊞ E(%s) = E(%s)\n\n", m1, m4, sk.DecryptPoly(r1, pk).String())

	const1 := big.NewFloat(10.0)
	r2 := pk.MultConstPoly(c2, const1)
	fmt.Printf("[MultConst] E(%s) ⊠ %f = E(%s)\n\n", m2, const1, sk.DecryptPoly(r2, pk).String())

	r3 := pk.MultPoly(c3, c4)
	dr3 := sk.DecryptPoly(r3, pk)
	fmt.Printf("[Mult] E(%s) ⊠ E(%s) = E(%s)\n\n", m3, m4, sk.DecryptPoly(r3, pk).String())

	const2 := big.NewFloat(0.5)
	r4 := pk.MultConstPoly(r3, const2)
	dr4 := sk.DecryptPoly(r4, pk)
	fmt.Printf("[MultConst] E(%s) ⊠ %f = E(%s)\n\n", dr3.String(), const2, dr4.String())

	r5 := pk.AddPoly(r3, r3)
	fmt.Printf("[Add] E(%s) ⊞ E(%s) = E(%s)\n\n", dr3.String(), dr3.String(), sk.DecryptPoly(r5, pk).String())

	r6 := pk.AddPoly(c1, c6)
	fmt.Printf("[Add] E(%s) ⊞ Neg(E(%s)) = E(%s)\n\n", m1, m4, sk.DecryptPoly(r6, pk).String())

	fmt.Println("\n----------DONE----------")

}

func runSimpleCheck(keyBits int, polyBase int) {

	pk, sk, _ := bgn.NewKeyGen(keyBits, big.NewInt(1021), polyBase, 3, 2, true)
	bgn.ComputeDecryptionPreprocessing(pk, sk)

	zero := pk.Encrypt(big.NewInt(0))
	one := pk.Encrypt(big.NewInt(1))
	negone := pk.Encrypt(big.NewInt(-1.0))

	fmt.Print("\n---------RUNNING BASIC CHECK----------\n\n")
	fmt.Println("0 + 0 = " + sk.DecryptFailSafe(pk.Add(zero, zero), pk).String())
	fmt.Println("0 + 1 = " + sk.DecryptFailSafe(pk.Add(zero, one), pk).String())
	fmt.Println("1 + 1 = " + sk.DecryptFailSafe(pk.Add(one, one), pk).String())
	fmt.Println("1 + 0 = " + sk.DecryptFailSafe(pk.Add(one, zero), pk).String())

	fmt.Println("0 * 0 = " + sk.DecryptFailSafe(pk.Mult(zero, zero), pk).String())
	fmt.Println("0 * 1 = " + sk.DecryptFailSafe(pk.Mult(zero, one), pk).String())
	fmt.Println("1 * 0 = " + sk.DecryptFailSafe(pk.Mult(one, zero), pk).String())
	fmt.Println("1 * 1 = " + sk.DecryptFailSafe(pk.Mult(one, one), pk).String())

	fmt.Println("0 - 0 = " + sk.DecryptFailSafe(pk.Add(zero, pk.Neg(zero)), pk).String())
	fmt.Println("0 - 1 = " + sk.DecryptFailSafe(pk.Add(zero, pk.Neg(one)), pk).String())
	fmt.Println("0 + (-1) = " + sk.DecryptFailSafe(pk.Add(zero, negone), pk).String())
	fmt.Println("1 - 1 = " + sk.DecryptFailSafe(pk.Add(one, pk.Neg(one)), pk).String())
	fmt.Println("1 - 0 = " + sk.DecryptFailSafe(pk.Add(one, pk.Neg(zero)), pk).String())

	fmt.Println("0 * (-0) = " + sk.DecryptFailSafe(pk.Mult(zero, pk.Neg(zero)), pk).String())
	fmt.Println("0 * (-1) = " + sk.DecryptFailSafe(pk.Mult(zero, pk.Neg(one)), pk).String())
	fmt.Println("1 * (-0) = " + sk.DecryptFailSafe(pk.Mult(one, pk.Neg(zero)), pk).String())
	fmt.Println("1 * (-1) = " + sk.DecryptFailSafe(pk.Mult(one, pk.Neg(one)), pk).String())
	fmt.Println("(-1) * (-1) = " + sk.DecryptFailSafe(pk.Mult(pk.Neg(one), pk.Neg(one)), pk).String())
	fmt.Println("\n---------DONE----------")

}
*/
