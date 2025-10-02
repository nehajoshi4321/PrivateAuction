package main

import (
	"fmt"
	"math/rand"
	"time"
	//"miracl/core/BN254"
	"io"
	"os"
	//  "reflect"
)

//auction operation
func auction(dslice []int) int {
	//fmt.Println("Domain bid:", dslice)
	//to find the min bid value among domains
	min := dslice[0]
	size := len(dslice)
	// var i int
	for i := 0; i < size; i++ {

		if dslice[i] > min {

			min = dslice[i]
		}

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

func optimal_domain(dom int, para int) []int {

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
	domain_slice := [][]int{}
	//  var bid_array[5] int

	//assign bid to the domain
	for i := 0; i < para; i++ {
		p_row := make([]int, dom)
		for j := 0; j < dom; j++ {
			// domain_slice[i][j]= rand.Intn(p_max - p_min + 1) + p_min
			Bid := rand.Intn(p_max-p_min+1) + p_min
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
	num := make([]int, para)

	for i := 0; i < 100; i++ {
		num[rand.Intn(para)]++
	}
	fmt.Println(num)
	//  fmt.Println("Printing current 2d array", domain_slice[0][0])
	//domain_slice1 := [][]int{}
	for i := 0; i < para; i++ {
		for j := 0; j < dom; j++ {
			domain_slice[i][j] = domain_slice[i][j] * num[i]
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
	domain_sum := make([]int, dom)

	for i := 0; i < para; i++ {
		//sum :=0
		for j := 0; j < dom; j++ {
			domain_sum[j] = domain_sum[j] + domain_slice[i][j]
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

	r1aa := rand.Int()
	r2aa := rand.Int()
	r1bb := rand.Int()
	r2bb := rand.Int()

	numbers := readFile("foo.in.txt")
	//  fmt.Println(numbers[0])
	m = numbers[0]
	n = numbers[1]
	fmt.Println("no. of domains are:", m)
	fmt.Println("no. of parameters are:", n)
	opt_domain_array := optimal_domain(m, n)
	min_bid := auction(opt_domain_array)
	fmt.Printf("Maximum bid: %d", min_bid)

}

//user bids

// fmt.Println(domain_slice)

//  min_bid :=auction(domain_slice)
// fmt.Printf("Minimum bid: %d", min_bid)
