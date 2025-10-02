package main

import (
	"fmt"
	"github.com/sachaservan/bgn"
	"log"
	"math/big"
	"math/rand"
	"time"
)

const KEYBITS = 1024
const POLYBASE = 3
const MSGSPACE = 100000000
const FPSCALEBASE = 3
const FPPREC = 0.0001
const DET = true
const MAX_BID = 100000
const ITERATIONS = 10

type UserContext struct {
	bid         int
	pubK        *bgn.PublicKey
	secK        *bgn.SecretKey
	eBid        *bgn.Ciphertext
	partDecoded *big.Int
}

func createPairwiseKey() (*bgn.PublicKey, *bgn.SecretKey, error) {
	start := time.Now()
	pk, sk, err := bgn.NewKeyGen(KEYBITS, big.NewInt(MSGSPACE), POLYBASE, FPSCALEBASE, FPPREC, DET)
	elapsed := time.Since(start)
	log.Printf("Time for pairwise setup %s", elapsed)
	if err != nil {
		panic(err)
	}
	return pk, sk, err
}

func main() {
	fmt.Println("Entering main()\n")
	rand.Seed(time.Now().UnixNano())

	bidders := make([]UserContext, 2)
	for i := range bidders {
		bidders[i] = UserContext{}
	}

	pk_zero, sk_zero, _err := createPairwiseKey()

	if _err != nil {
		panic(_err)
	}
	bidders[0].pubK=pk_zero
    bidders[0].secK=sk_zero
	bgn.ComputeDecryptionPreprocessing(bidders[0].pubK, bidders[0].secK)

	bidders[1].pubK, bidders[1].secK, _err = createPairwiseKey()
	if _err != nil {
		panic(_err)
	}
	bgn.ComputeDecryptionPreprocessing(bidders[1].pubK, bidders[1].secK)

	for i := 1; i < ITERATIONS; i++ {
		bidders[0].bid = rand.Intn(MAX_BID)
		bidders[1].bid = rand.Intn(MAX_BID)
		bidders[0].eBid = bidders[0].pubK.Encrypt(big.NewInt(int64(bidders[0].bid)))
		bidders[1].eBid = bidders[1].pubK.Encrypt(big.NewInt(int64(bidders[1].bid)))

		bidders[0].partDecoded = bidders[0].secK.DecryptFailSafe(bidders[0].eBid, bidders[0].pubK)
		bidders[1].partDecoded = bidders[1].secK.DecryptFailSafe(bidders[1].eBid, bidders[1].pubK)

		if bidders[0].bid != int(bidders[0].partDecoded.Int64()) ||
			bidders[1].bid != int(bidders[1].partDecoded.Int64()) {
			fmt.Println("ERROR values are: ", bidders[0].bid, ", ", bidders[1].bid, ", ", bidders[0].partDecoded, ", ", bidders[1].partDecoded, "\n")
		} else if i%(ITERATIONS/10) == 0 {
			fmt.Println("values are: ", bidders[0].bid, ", ", bidders[1].bid, ", ", bidders[0].partDecoded, ", ", bidders[1].partDecoded, "\n")
		}
	}
}
