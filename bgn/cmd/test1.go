package main

import (
	"fmt"
	"github.com/sachaservan/bgn"
	"log"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
	"reflect"
	"math"

)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"

	POLYBASE = 3

	FPSCALEBASE = 3

	FPPREC = 0.0001

	DET = true // deterministic ops

)

var (
	keyBitLength int

	msgSpace int64

	numBidders int

	randPercent int64

	maxBid int64
)

type Bidder struct {
	identity int

	bid, rA, rB int64

	pubK *bgn.PublicKey

	secK *bgn.SecretKey

	eBid, eRA, eRB *bgn.Ciphertext
}

//create pairwise keys.

func createPairwiseKey() (*bgn.PublicKey, *bgn.SecretKey, error) {

	pk, sk, err := bgn.NewKeyGen(keyBitLength, big.NewInt(msgSpace), POLYBASE, FPSCALEBASE, FPPREC, DET)

	if err != nil {

		panic(err)

	}

	return pk, sk, err

}

// Generates the random values (rA, rB), encrypted bids (eBid), and encrypted random values (eRA, eRB).

// TODO: Just send on bidder by referece - Do not need idx

func genEncodingParameter(bidders []Bidder, idx int) {

	// Generating random values

	bidders[idx].rA = rand.Int63n(randPercent)

	bidders[idx].rB = rand.Int63n(randPercent)

	// Encrypting plaintext bid and random values

	bidders[idx].eBid = bidders[idx].pubK.Encrypt(big.NewInt(bidders[idx].bid))

	bidders[idx].eRA = bidders[idx].pubK.Encrypt(big.NewInt(bidders[idx].rA))

	bidders[idx].eRB = bidders[idx].pubK.Encrypt(big.NewInt(bidders[idx].rB))

}

// A party calls this function to add its own randomization to the

//  encrypted bid received from the foreign(other) party.

func addRandomOnEncRec(fEncBid *bgn.Ciphertext, fRA *bgn.Ciphertext, fRB *bgn.Ciphertext,

	fPubK *bgn.PublicKey, sRA int64, sRB int64) *bgn.Ciphertext {

	// Encrypting the random values with the received public key

	selfERA := fPubK.Encrypt(big.NewInt(sRA))

	selfERB := fPubK.Encrypt(big.NewInt(sRB))

	// Adding encrypted random values to the received eBid

	// derefernce - as the library updates the value

	tempFRA := *fRA

	tempFRB := *fRB

	tempERA := fPubK.Add(&tempFRA, selfERA)

	tempERB := fPubK.Mult(&tempFRB, selfERB)

	tempEbid := fPubK.Mult(fEncBid, tempERA)

	fEncBid = fPubK.Add(tempEbid, tempERB)

	return fEncBid

}

// Performs private auction on encrypted bids and return the winner.

// Simple Bubble approach - O(n)

func auction(bidders []Bidder) int {

	// Assuming 1st bidder as winner

	start := time.Now()

	log.Println("Numbidders: ", numBidders)

	var winner int = 0

	var partyB int

	for i := 1; i < numBidders; i = i + 1 {

		partyB = i

		var winnerCrossEncBid, partyBCrossEncBid *bgn.Ciphertext

		log.Println("Comparing ", winner, " and ", partyB, "raw values:", bidders[winner].bid, ",", bidders[partyB].bid)

		// genEncodingParameter(bidders, winner)

		// genEncodingParameter(bidders, partyB)

		// In real world (TBD TODO - object oriented implementation)

		//  below two blocks calls are to be executed locally by the respective parties.

		// T1 partyB to add randomization on the encrypted bids of winner party.

		winnerCrossEncBid = addRandomOnEncRec(bidders[winner].eBid, bidders[winner].eRA, bidders[winner].eRB, bidders[winner].pubK, bidders[partyB].rA, bidders[partyB].rB)

		// T1 old winner party decrypting the encrypted bid to compute encoded bid

		bgn.ComputeDecryptionPreprocessing(bidders[winner].pubK, bidders[winner].secK)

		decodedWinner := bidders[winner].secK.DecryptFailSafe(winnerCrossEncBid, bidders[winner].pubK)

		// T2 old winner to add randomization on the encrypted bids of partyB

		partyBCrossEncBid = addRandomOnEncRec(bidders[partyB].eBid, bidders[partyB].eRA, bidders[partyB].eRB, bidders[partyB].pubK, bidders[winner].rA, bidders[winner].rB)

		// T2 partyB decrypting the encrypted bid to compute encoded bid

		bgn.ComputeDecryptionPreprocessing(bidders[partyB].pubK, bidders[partyB].secK)

		decodedPartyB := bidders[partyB].secK.DecryptFailSafe(partyBCrossEncBid, bidders[partyB].pubK)

		// Following code may be run on either side - after cross sharing encodedPartyB and encodedWinner

		// TODO - no communication yet

		log.Println("decodedWinner: ",decodedWinner,"  :::: decodedPartyB: ",decodedPartyB)

		if decodedWinner.Cmp(decodedPartyB) == -1 {

			winner = partyB

		}

	}

	elapsed := time.Since(start)

	log.Printf("Total time taken by the auction process: %s", elapsed)

	return winner

}

// Generating plaintext bids and public-private key pair for each bidder.

func initBidders(bidders []Bidder, bidValues []int64) {

	var _err error

	for i := 0; i < numBidders; i = i + 1 {

		bidders[i].identity = i + 1

		bidders[i].bid = bidValues[i]

		log.Print("    ", bidders[i].identity, " : ", bidders[i].bid)

		bidders[i].pubK, bidders[i].secK, _err = createPairwiseKey()

		genEncodingParameter(bidders, i)

		if _err != nil {

			panic(_err)

		}

	}

	log.Println()

}

func initBiddersRand(bidders []Bidder) {

	//start := time.Now()

	var _err error

	for i := 0; i < numBidders; i = i + 1 {

		bidders[i].identity = i + 1

		bidders[i].bid = rand.Int63n(maxBid)

		bidders[i].pubK, bidders[i].secK, _err = createPairwiseKey()

		genEncodingParameter(bidders, i)

		if _err != nil {

			panic(_err)

		}

	}

}


func main() {

	if len(os.Args) < 7 {

		fmt.Println("Wrong argument count ", os.Args[0], " <run count> <keyBitLength> <msgSpace> <numBidders> ",

			"<randPercent>  <maxBid> <seed#>")
		os.Exit(1)
		//return

	}

	n, er := strconv.Atoi(os.Args[1])

	keyBitLengthLocal, er := strconv.Atoi(os.Args[2])
	keyBitLength = keyBitLengthLocal

	msgSpaceLocal, er := strconv.Atoi(os.Args[3])
	msgSpace = int64(msgSpaceLocal)

	numBiddersLocal, er := strconv.Atoi(os.Args[4])
	numBidders = numBiddersLocal

	randPercentLocal, er := strconv.Atoi(os.Args[5])
	randPercent = int64(randPercentLocal)

	maxBidLocal, er := strconv.Atoi(os.Args[6])
	maxBid = int64(maxBidLocal)
	if er != nil {
		// handle error
	}

	seed := time.Now().UnixNano()
	fmt.Println(reflect.TypeOf(seed))

    var seedInt int
	if len(os.Args) == 8 {
		var err error
		seedInt, err = strconv.Atoi(os.Args[7])
		if err != nil {
			// handle error

		}
	}
	seedValue := int64(seedInt)
	now := time.Now()

	fileName := fmt.Sprintf("%s_%d_%d_%d_%d_%d_%d_%d", now.Format(DDMMYYYYhhmmss), n, keyBitLength, msgSpace,
		numBidders, randPercent, maxBid, seedValue)
	//var totalTime time.Duration = 0
	//  minMemory := uint64(999999999)
	//  maxMemory := uint64(0)
	var m runtime.MemStats
	rand.Seed(seedValue)
	file, e := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if e != nil {

		log.Print(e.Error() + "\r\n")

		return

	}

	log.SetOutput(file)

    	start := time.Now()


	for i := 0; i < n; i++ {

		log.Println("\nIteration:", i, "keyBitLength:", keyBitLength, "\t msgSpace:", msgSpace, "\t numBidders:", numBidders,

			"\t randPercent:", randPercent, "\t maxBid:", maxBid)
		bidders := make([]Bidder, numBidders)

		initBiddersRand(bidders)

		aucStart := time.Now()


		var winnerIdx = auction(bidders)



		runtime.ReadMemStats(&m)

		log.Println("Winner is  bidder: ", bidders[winnerIdx].identity, " with bid: ", bidders[winnerIdx].bid)
        elapsed := time.Since(aucStart)
		log.Printf("Time during Auction %s memory taken = %.2f MB\n", elapsed, float64(m.Alloc)/(1024*1024))

}
      duration := time.Since(start)
            cpuTime := duration - time.Duration(float64(duration.Nanoseconds())*0.1)
    	  //fmt.Printf("CPU utilization: %.2f%%\n", 100*float64(duration.Nanoseconds())/float64(time.Second))
            cpuUtilization := 100 * float64(duration.Nanoseconds()) / float64(time.Second)
            fmt.Printf("CPU utilization: %.2f%%\n", math.Min(100, cpuUtilization))
          //  fmt.Printf("CPU utilization: %.2f%%\n", cpuUtilization)
            fmt.Printf("CPU utilization time: %v\n", cpuTime)

   // elapse1 := time.Since(start)
   // fmt.Println("Average execution time:", elapse1)
	averageTime := time.Since(start) / time.Duration(n)

	fmt.Println("Average execution time:", averageTime)

}
