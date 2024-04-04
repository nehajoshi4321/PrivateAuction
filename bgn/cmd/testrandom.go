package main

import (
	"flag"
	"fmt"
	"github.com/sachaservan/bgn"
	"log"
	"math/big"
	"math/rand"
	"time"
)

const POLYBASE = 3
const FPSCALEBASE = 3
const FPPREC = 0.0001
const DET = true

const KEYBITS = 512
const MSGSPACE = 10000000000000
const MAX_BID = 100000000000000
const ITERATIONS = 14

type UserContext struct {
	bid         int
	pubK        *bgn.PublicKey
	secK        *bgn.SecretKey
	eBid        *bgn.Ciphertext
	partDecoded *big.Int
}

func createPairwiseKey(keybits int, msgspace int) (*bgn.PublicKey, *bgn.SecretKey, error) {
	start := time.Now()
	pk, sk, err := bgn.NewKeyGen(keybits, big.NewInt(int64(msgspace)), POLYBASE, FPSCALEBASE, FPPREC, DET)
	elapsed := time.Since(start)
	log.Printf("Time for pairwise setup %s", elapsed)
	if err != nil {
		panic(err)
	}
	return pk, sk, err
}

func main() {
	var keybits int
	var msgspace int
	var maxBid int
	var iterations int

	flag.IntVar(&keybits, "keybits", KEYBITS, "Number of key bits")
    	flag.IntVar(&msgspace, "msgspace", MSGSPACE, "Message space for polynomial coefficients")
    	flag.IntVar(&maxBid, "maxbid", MAX_BID, "Maximum bid")
    	flag.IntVar(&iterations, "iterations", ITERATIONS, "Number of iterations")
    	flag.Parse()

	//fmt.Println("Entering main()\n")
	rand.Seed(time.Now().UnixNano())

	bidders := make([]UserContext, 2)
	for i := range bidders {
		bidders[i] = UserContext{}
	}

	pkZero, skZero, err := createPairwiseKey(keybits, msgspace)
	if err != nil {
		panic(err)
	}
	bidders[0].pubK = pkZero
	bidders[0].secK = skZero
	bgn.ComputeDecryptionPreprocessing(bidders[0].pubK, bidders[0].secK)

	pkOne, skOne, err := createPairwiseKey(keybits, msgspace)
	if err != nil {
		panic(err)
	}
	bidders[1].pubK = pkOne
	bidders[1].secK = skOne
	bgn.ComputeDecryptionPreprocessing(bidders[1].pubK, bidders[1].secK)

	for i := 1; i < iterations; i++ {
		bidders[0].bid = rand.Intn(maxBid)
		bidders[1].bid = rand.Intn(maxBid)
		bidders[0].eBid = bidders[0].pubK.Encrypt(big.NewInt(int64(bidders[0].bid)))
		bidders[1].eBid = bidders[1].pubK.Encrypt(big.NewInt(int64(bidders[1].bid)))

		bidders[0].partDecoded = bidders[0].secK.DecryptFailSafe(bidders[0].eBid, bidders[0].pubK)
		bidders[1].partDecoded = bidders[1].secK.DecryptFailSafe(bidders[1].eBid, bidders[1].pubK)

		if bidders[0].bid != int(bidders[0].partDecoded.Int64()) ||
			bidders[1].bid != int(bidders[1].partDecoded.Int64()) {
			fmt.Println("ERROR values are: ", bidders[0].bid, ", ", bidders[1].bid, ", ", bidders[0].partDecoded, ", ", bidders[1].partDecoded, "\n")
		} else if i%(iterations/10) == 0 {
			fmt.Println("values are: ", bidders[0].bid, ", ", bidders[1].bid, ", ", bidders[0].partDecoded, ", ", bidders[1].partDecoded, "\n")
		}
	}
}
