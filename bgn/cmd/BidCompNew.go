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
        eBid, erA, erB *bgn.Ciphertext
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

// A party calls this function to add its own randomization to the received encrypted bid of the other party.
func addRandom(eBid *bgn.Ciphertext, erA *bgn.Ciphertext, erB *bgn.Ciphertext, pubK *bgn.PublicKey, rA int, rB int) *bgn.Ciphertext {

        // Encrypting the random values with the received public key
        selfErA := pubK.Encrypt(big.NewInt(int64(rA)))
        selfErB := pubK.Encrypt(big.NewInt(int64(rB)))

        // Adding encrypted random values to the received eBid
        tempErA := pubK.Add(erA, selfErA)
        tempErB := pubK.Mult(erB, selfErB)

        tempEbid := pubK.Mult(eBid, tempErA)
        eBid = pubK.Add(tempEbid, tempErB)

        return eBid
}
 


// Generates the random values (rA, rB), encrypted bids (eBid), and encrypted random values (erA, erB).
func genPar(bidders []Bidder, idx int){
    // Generating random values
    bidders[idx].rA = rand.Intn(MAX_RAND)
    bidders[idx].rB = rand.Intn(MAX_RAND)

    // Encrypting plaintext bid and random values
    bidders[idx].eBid = bidders[idx].pubK.Encrypt(big.NewInt(int64(bidders[idx].bid)))
    bidders[idx].erA = bidders[idx].pubK.Encrypt(big.NewInt(int64(bidders[idx].rA)))
    bidders[idx].erB = bidders[idx].pubK.Encrypt(big.NewInt(int64(bidders[idx].rB)))

}

// Performs private auction on encrypted bids and return the winner.
func auction(bidders []Bidder) int{
    // Assuming 1st bidder as winner
    var winner int = 0
    var partyB int

    for i := 1; i < NUM_BIDDERS; i = i + 1 {
        partyB = i
        // Generating random values and encrypted bids for winner
        genPar(bidders,winner)
        // Generating random values and encrypted bids for partyB
        genPar(bidders,partyB)

        // In real world, below two function calls are to be executed locally by the respective parties.

        // partyB calling the addRandom function to add randomization on the encrypted bids of winner party.
        bidders[winner].eBid = addRandom(bidders[winner].eBid, bidders[winner].erA, bidders[winner].erB, bidders[winner].pubK, bidders[partyB].rA, bidders[partyB].rB)

        // winner party calling the addRandom function to add randomization on the encrypted bids of partyB
        bidders[partyB].eBid = addRandom(bidders[partyB].eBid, bidders[partyB].erA, bidders[partyB].erB, bidders[partyB].pubK, bidders[winner].rA, bidders[winner].rB)

        // winner party decrypting the encrypted bid to compute encoded bid
        bgn.ComputeDecryptionPreprocessing(bidders[winner].pubK, bidders[winner].secK)
        encodedWinner := bidders[winner].secK.DecryptFailSafe(bidders[winner].eBid, bidders[winner].pubK)

        // partyB decrypting the encrypted bid to compute encoded bid
        bgn.ComputeDecryptionPreprocessing(bidders[partyB].pubK, bidders[partyB].secK)
        encodedPartyB := bidders[partyB].secK.DecryptFailSafe(bidders[partyB].eBid, bidders[partyB].pubK)

        if(encodedWinner.Cmp(encodedPartyB) == -1){
            winner = partyB
        }
    }

    return bidders[winner].identity
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
        var winner = auction(bidders)
        fmt.Println("Winner is  bidder: ", winner)
}
