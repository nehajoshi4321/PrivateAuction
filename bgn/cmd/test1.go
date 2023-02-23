package main
import (
        "fmt"
        "github.com/sachaservan/bgn"
        "math/big"
        "math/rand"
        "time"
        "log"
        "os"
        "bufio"
        "strconv"
        "strings"
        "runtime"
)


const POLYBASE = 3
const FPSCALEBASE = 3
const FPPREC = 0.0001
const DET = true // deterministic ops

var (
	KEYBITS     int
	MSGSPACE    int64
	NUM_BIDDERS int
	MAX_RAND    int64
	MAX_BID     int64
)


type Bidder struct {
        identity int
        bid,rA,rB int64
        pubK *bgn.PublicKey
        secK *bgn.SecretKey
        eBid, eRA, eRB *bgn.Ciphertext
}
//create pairwise keys.
func createPairwiseKey() (*bgn.PublicKey, *bgn.SecretKey, error) {
        pk, sk, err := bgn.NewKeyGen(KEYBITS, big.NewInt(MSGSPACE), POLYBASE, FPSCALEBASE, FPPREC, DET)

        if err != nil {
                panic(err)
        }
        return pk, sk, err
}

// Generates the random values (rA, rB), encrypted bids (eBid), and encrypted random values (eRA, eRB).
// TODO: Just send on bidder by referece - Do not need idx
func genEncodingParameter(bidders []Bidder, idx int) {
    // Generating random values
    bidders[idx].rA = rand.Int63n(MAX_RAND)
    bidders[idx].rB = rand.Int63n(MAX_RAND)

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
func auction(bidders []Bidder) int{
    // Assuming 1st bidder as winner
    start := time.Now()
    var winner int = 0
    var partyB int

    for i := 1; i < NUM_BIDDERS; i = i + 1 {

        partyB = i
        var winnerCrossEncBid, partyBCrossEncBid *bgn.Ciphertext
        log.Println("Comparing " , winner, " and " , partyB, "raw values:" , bidders[winner].bid, ",", bidders[partyB].bid)
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
        if(decodedWinner.Cmp(decodedPartyB) == -1){
            winner = partyB
        }
    }
     elapsed := time.Since(start)
    log.Printf("Total time taken by the auction process: %s", elapsed)
    return winner

}

// Generating plaintext bids and public-private key pair for each bidder.
func initBidders(bidders []Bidder, bidValues []int64){

    var _err error
    for i := 0; i < NUM_BIDDERS; i = i + 1 {
        bidders[i].identity = i + 1
        bidders[i].bid = bidValues[i]
        fmt.Print("    ",  bidders[i].identity ," : ", bidders[i].bid )
        bidders[i].pubK,bidders[i].secK, _err = createPairwiseKey()
        genEncodingParameter(bidders, i)
        if _err != nil {
                panic(_err)
        }
    }
    fmt.Println()
}

func initBiddersRand(bidders []Bidder){
    start := time.Now()
    var _err error
    for i := 0; i < NUM_BIDDERS; i = i + 1 {
        bidders[i].identity = i + 1
        bidders[i].bid = rand.Int63n(MAX_BID)
        fmt.Print("    ",  bidders[i].identity ," : ", bidders[i].bid )
        bidders[i].pubK,bidders[i].secK, _err = createPairwiseKey()
        genEncodingParameter(bidders, i)
        if _err != nil {
                panic(_err)
        }
    }
    fmt.Println()
     elapsed := time.Since(start)
            log.Printf("Time for calculating the initBiddersRand: %s", elapsed)
}

func readConfig(filename string) error {
	// Open the text file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Split the line into a variable name and value
		line := scanner.Text()
		parts := strings.Split(line, " = ")
		if len(parts) != 2 {
			fmt.Println("Invalid line:", line)
			continue
		}
		name := parts[0]
		valueStr := parts[1]

		// Convert the value string to an integer
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
		   fmt.Println("Error converting value:", valueStr)
		   continue
		}

		// Assign the value to the appropriate variable
		switch name {
		case "KEYBITS":
			KEYBITS = int(value)
		case "MSGSPACE":
			MSGSPACE = value
		case "NUM_BIDDERS":
			NUM_BIDDERS = int(value)
		case "MAX_RAND":
			MAX_RAND = value
		case "MAX_BID":
			MAX_BID = value
		default:
			fmt.Println("Unknown variable:", name)
		}
	}

	return nil
}


func main() {

        n := 10 // number of times to run the code
        var totalTime time.Duration = 0
        memory := uint64(0)
           for i := 0; i < n; i++ {
               start := time.Now()
               // code to be executed n times
               // Read the variables from the file
                       	if err := readConfig("Input.txt"); err != nil {
                       		fmt.Println("Error reading config:", err)
                       		return
                       	}

                       file, e := os.OpenFile("output.txt", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)

                       if e !=nil{
                        log.Print(e.Error() + "\r\n")

                       }

                       log.SetOutput(file)

                       //log.Println("\n MSGSPACE, MAX_RAND, MAX_BID", MSGSPACE, ", ", MAX_RAND, ", ", MAX_BID)
                       // Print the values of the variables
                       log.Println("\niteration:", i,  "KEYBITS:", KEYBITS ,"\t MSGSPACE:", MSGSPACE,"\t NUM_BIDDERS:", NUM_BIDDERS,
                                                   "\t MAX_RAND:", MAX_RAND,"\t MAX_BID:", MAX_BID)
                       //log.Println("MSGSPACE:", MSGSPACE)

                       // Initializing the seed using current time for random number generation.
                       rand.Seed(time.Now().UnixNano())
                       bidders := make([]Bidder, NUM_BIDDERS)
                       // Initializing the bidders. Generating plaintext bids and public-private key pair.
                       // Static bid values for reproducing error
                               // values :=  []int64{825616, 54460, 857406, 129782, 181565, 552263, 258629}
                       // initBidders(bidders, values)
                               initBiddersRand(bidders)

                       // Performing auction
                       var winnerIdx = auction(bidders)
                       log.Println("Winner is  bidder: ", bidders[winnerIdx].identity, " with bid: " , bidders[winnerIdx].bid)
                       elapsed := time.Since(start)
                       log.Printf("Time for calculating the winner: %s", elapsed)

               var m runtime.MemStats
               runtime.ReadMemStats(&m)
               memory +=m.Alloc

               time.Sleep(time.Millisecond * 100) // example code that takes 100 milliseconds to execute
               end := time.Now()
               totalTime += end.Sub(start)


           }
           averageTime := totalTime / time.Duration(n)
           fmt.Println("Average execution time:", averageTime)
           avgMemory := float64(memory)/ float64(n)
           fmt.Printf("amount of memory currently allocated by the code = %.2f MB\n", avgMemory/1024/1024)

        }