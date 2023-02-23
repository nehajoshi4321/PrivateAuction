package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	KEYBITS     int
	MSGSPACE    int64
	NUM_BIDDERS int
	MAX_RAND    int64
	MAX_BID     int64
)

func main() {
	// Read the variables from the file
	if err := readConfig("Input.txt"); err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	// Print the values of the variables
	fmt.Println("KEYBITS:", KEYBITS)
	fmt.Println("MSGSPACE:", MSGSPACE)
	fmt.Println("NUM_BIDDERS:", NUM_BIDDERS)
	fmt.Println("MAX_RAND:", MAX_RAND)
	fmt.Println("MAX_BID:", MAX_BID)
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

		// Convert the value string to the appropriate type
		//var value int64 //interface{}
		/* if strings.HasPrefix(valueStr, "0x") {
			// Hexadecimal number
			n, err := strconv.ParseUint(valueStr[2:], 16, 64)
			if err != nil {
				fmt.Println("Invalid value:", valueStr)
				continue
			}
			value = n
		} else if strings.HasPrefix(valueStr, "0") {
			// Octal number
			n, err := strconv.ParseUint(valueStr[1:], 8, 64)
			if err != nil {
				fmt.Println("Invalid value:", valueStr)
				continue
			}
			value = n
		} else if strings.HasPrefix(valueStr, "0b") {
			// Binary number
			n, err := strconv.ParseUint(valueStr[2:], 2, 64)
			if err != nil {
				fmt.Println("Invalid value:", valueStr)
				continue
			}
			value = n
		} else if strings.Contains(valueStr, ".") {
			// Floating-point number
			f, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				fmt.Println("Invalid value:", valueStr)
				continue
			}
			value = f
		} else { */
			// Integer
			value, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				fmt.Println("Error converting value:", valueStr)
				continue
			}
		
		//}

		// Assign the value to the appropriate variable
		switch name {
		case "KEYBITS":
			KEYBITS = value
		case "MSGSPACE":
			MSGSPACE = value
		case "NUM_BIDDERS":
			NUM_BIDDERS = value
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

