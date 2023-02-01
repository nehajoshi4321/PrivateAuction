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

const MSGSPACE = 10000000 // message space for polynomial coefficients

const NUM_BIDDERS = 7

const FPSCALEBASE = 3

const FPPREC = 0.0001

const MAX_RAND = 1000

const MAX_BID = 1000

const DET = true // deterministic ops

 

type Bidder struct {

        bid,rA,rB,identity int

        pubK *bgn.PublicKey

        secK *bgn.SecretKey

        eBid, eRA, eRB *bgn.Ciphertext

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

 

// Generates the random values (rA, rB), encrypted bids (eBid), and encrypted random values (eRA, eRB).

// TODO: Just send on bidder by referece - Do not need idx

func genEncodingParameter(bidders []Bidder, idx int) {

    // Generating random values

    bidders[idx].rA = rand.Intn(MAX_RAND)

    bidders[idx].rB = rand.Intn(MAX_RAND)

 

    // Encrypting plaintext bid and random values

    bidders[idx].eBid = bidders[idx].pubK.Encrypt(big.NewInt(int64(bidders[idx].bid)))

    bidders[idx].eRA = bidders[idx].pubK.Encrypt(big.NewInt(int64(bidders[idx].rA)))

    bidders[idx].eRB = bidders[idx].pubK.Encrypt(big.NewInt(int64(bidders[idx].rB)))

}

 

// A party calls this function to add its own randomization to the

//  encrypted bid received from the foreign(other) party.

func addRandomOnEncRec(fEncBid *bgn.Ciphertext, fRA *bgn.Ciphertext, fRB *bgn.Ciphertext,

                                                                fPubK *bgn.PublicKey, sRA int, sRB int) *bgn.Ciphertext {

 

        // Encrypting the random values with the received public key

        selfERA := fPubK.Encrypt(big.NewInt(int64(sRA)))

        selfERB := fPubK.Encrypt(big.NewInt(int64(sRB)))

 

        // Adding encrypted random values to the received eBid

        tempERA := fPubK.Add(fRA, selfERA)

        tempERB := fPubK.Mult(fRB, selfERB)

 

        tempEbid := fPubK.Mult(fEncBid, tempERA)

        fEncBid = fPubK.Add(tempEbid, tempERB)

 

        return fEncBid

}

 

 

 

// Performs private auction on encrypted bids and return the winner.

// Simple Bubble approach - O(n)

func auction(bidders []Bidder) int{

    // Assuming 1st bidder as winner

    var winner int = 0

    var partyB int



    for i := 1; i < NUM_BIDDERS; i = i + 1 {

        partyB = i

 
        genEncodingParameter(bidders, winner)

        genEncodingParameter(bidders, partyB)

        // In real world (TBD TODO - object oriented implementation)

        //  below two blocks calls are to be executed locally by the respective parties.

 

        // T1 partyB to add randomization on the encrypted bids of winner party.

        bidders[winner].eBid = addRandomOnEncRec(bidders[winner].eBid, bidders[winner].eRA, bidders[winner].eRB, bidders[winner].pubK, bidders[partyB].rA, bidders[partyB].rB)

        // T1 old winner party decrypting the encrypted bid to compute encoded bid

        bgn.ComputeDecryptionPreprocessing(bidders[winner].pubK, bidders[winner].secK)

        encodedWinner := bidders[winner].secK.DecryptFailSafe(bidders[winner].eBid, bidders[winner].pubK)

 

        // T2 old winner to add randomization on the encrypted bids of partyB

        bidders[partyB].eBid = addRandomOnEncRec(bidders[partyB].eBid, bidders[partyB].eRA, bidders[partyB].eRB, bidders[partyB].pubK, bidders[winner].rA, bidders[winner].rB)

        // T2 partyB decrypting the encrypted bid to compute encoded bid

        bgn.ComputeDecryptionPreprocessing(bidders[partyB].pubK, bidders[partyB].secK)

        encodedPartyB := bidders[partyB].secK.DecryptFailSafe(bidders[partyB].eBid, bidders[partyB].pubK)

 

                // Following code may be run on eithr side - after cross sharing encodedPartyB and encodedWinner

                // TODO - no communication yet

        if(encodedWinner.Cmp(encodedPartyB) == -1){

            winner = partyB

        }

    }

 

    return winner

}

 

 

// Generating plaintext bids and public-private key pair for each bidder.

func initBidders(bidders []Bidder){

    fmt.Println("Entering initBidders()\n")

    var _err error

    for i := 0; i < NUM_BIDDERS; i = i + 1 {

        bidders[i].identity = i + 1

        bidders[i].bid = rand.Intn(MAX_BID)

        fmt.Println("bidders[",i,"] value is: ", bidders[i].bid, "\n" )

        bidders[i].pubK,bidders[i].secK, _err = createPairwiseKey()



        if _err != nil {

                panic(_err)

        }

    }

    fmt.Println("Exiting initBidders()\n")

}

 

 

func main() {

        fmt.Println("Entering main()\n")

        // Initializing the seed using current time for random number generation.

        rand.Seed(time.Now().UnixNano())

        bidders := make([]Bidder, NUM_BIDDERS)

 

        // Initializing the bidders. Generating plaintext bids and public-private key pair.

        initBidders(bidders)

 

        // Performing auction.

        var winnerIdx = auction(bidders)

        fmt.Println("Winner is  bidder: ", bidders[winnerIdx].identity, " with bid: " , bidders[winnerIdx].bid)

}


