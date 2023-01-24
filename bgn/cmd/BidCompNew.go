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
const NUM_BIDDERS = 4
const FPSCALEBASE = 3
const FPPREC = 0.0001
const MAX_RAND = 1000
const MAX_BID = 1000
const DET = true // deterministic ops
 
type Bidder struct {
        bid,rA,rB int
        pubK *bgn.PublicKey
        secK *bgn.SecretKey
}
 
func createPairwiseKey() (*bgn.PublicKey, *bgn.SecretKey, error) {
        start := time.Now()
        pk, sk, err := bgn.NewKeyGen(KEYBITS, big.NewInt(MSGSPACE), POLYBASE, FPSCALEBASE, FPPREC, DET)
        elapsed := time.Since(start)
        log.Printf("Time for pairwise setup %s", elapsed)
        if err != nil {
                panic(err)
        }
        return pk, sk
}
 
func encryptPubKey(pK *bgn.PublicKey, bid int, rA int, rB int) (*bgn.Ciphertext, *bgn.Ciphertext, *bgn.Ciphertext) {
        cBid := pK.Encrypt(big.NewInt(int64(bid)))
        cA := pK.Encrypt(big.NewInt(int64(rA)))
        cB := pK.Encrypt(big.NewInt(int64(rB)))
        return cBid, cA, cB
 
}
 
func encryptMergeRemote(sEncBid *bgn.Ciphertext, sEncA *bgn.Ciphertext, sEncB *bgn.Ciphertext, pk *bgn.PublicKey, r1b int, r2b int) *bgn.Ciphertext {
 
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
 
func encBid(selfIdx int, otherIdx int, Bidder[] bidders) big.Int {
 
 
        rand.Seed(time.Now().UnixNano())
    self = bidders[selfIdx]
    other = bidders[othrIdx]
 
        fmt.Println("===plain bid value")
        fmt.Println(bidders[self].bid)
        fmt.Println("\n===encrypted bid generation")
 
        a, b, c := encryptPubKey(self.pubK, self.bid, self.rA, other.rA)
        encr := bid_random(a, b, c, self.pubK, r1B, r2B)
 
        fmt.Println("===encrypted bid value")
        fmt.Println(encr)
 
        decr := bid_comp(encr, pk_i, sk_i)
 
        fmt.Println(decr)
 
        return *decr
}
 
func encodedBidPair(bid_i int, bid_j int) (big.Int, big.Int) {
 
        r1a := rand.Intn(MAX_RAND)
        r2a := rand.Intn(MAX_RAND)
        r1b := rand.Intn(MAX_RAND)
        r2b := rand.Intn(MAX_RAND)
 
        fmt.Println(r1a, r2a, r1b, r2b)
 
        user_i = encBidGenerate(bid_i, r1a, r2a, r1b, r2b)
        user_j = encBidGenerate(bid_j, r1b, r2b, r1a, r2a)
        return user_i, user_j
}
 
func encryptedCompare(Bidder[] bidders, int numParties) int {
 
        bid_i := rand.Intn(MAX_BID)
        bid_j := rand.Intn(MAX_BID)
 
        domain := make([]big.Int, NUM_BIDDERS)
        for i := 0; i < n; i = i + 2 {
                domain[i], domain[i+1] = bid_final_score(user_i, user_j)
        }
 
        fmt.Println("finalscore:", domain)
        var k int
        k = 0
        for j := 0; j < n; j++ {
 
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
 
func main() {
 
        rand.Seed(time.Now().UnixNano())
        start := time.Now()
        bids := make([]Bidder, NUM_BIDDERS)
        minBid = MAX_BID
        minIndex = MAX_BID
 
        fmt.Print(" BIDS are ")
        for i := 0; i < NUM_BIDDERS; i = i + 1 {
                bids[i].bid= rand.Intn(MAX_BID)
                bids[i].r1= rand.Intn(MAX_RAND)
                bids[i].r2= rand.Intn(MAX_RAND)
                if minBid > bids[i].bid{
                    minBid = bids[i].bid
                    minIndex = i
                }
                (bids[i].pubK,bids[i].secK) = createPairwiseKey()
                fmt.Print(" : ",  rawBids[i])
        }
        fmt.Println(" with MIN = ", minBid)
        //winning bid calculation
        winBid = encryptedCompare(bids, NUM_BIDDERS)
        fmt.Println("Winning bid in Main:", winBid)
}
