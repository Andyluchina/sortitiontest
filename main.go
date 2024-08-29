package main

import (
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/coniks-sys/coniks-go/crypto/vrf"
)

func main() {

	test_iterations := 100000

	results := make([]int64, test_iterations)

	for i := 0; i < test_iterations; i++ {

		start := time.Now()
		sk, err := vrf.GenerateKey(nil)

		if err != nil {
			fmt.Errorf("Failed to generate key")
		}
		_, succ := sk.Public()

		if !succ {
			fmt.Errorf("Failed to generate public key")
		}

		// URL to request
		url := "https://beacon.nist.gov/beacon/2.0/pulse/last"

		// Make the GET request
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unexpected status code: %d", resp.StatusCode)
		}

		// Decode into a generic map
		var result map[string]interface{}

		// Decode the JSON response
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Fatalf("Failed to decode JSON response: %v", err)
		}

		result = result["pulse"].(map[string]interface{})
		// Print the JSON response (dynamic fields)

		hexString := result["localRandomValue"].(string)

		// fmt.Println(hexString)

		byteArray, err := hex.DecodeString(hexString)
		if err != nil {
			log.Fatalf("Failed to decode hex string: %v", err)
		}

		// Print the byte array
		// fmt.Printf("Byte Array: %v\n", byteArray)

		// fmt.Println("Compute VRF of public randomness beacon")
		_, _ = sk.Prove(byteArray)
		// fmt.Println("VRF: ", v)
		// fmt.Println("Proof: ", p)
		elapsed := time.Since(start)

		results[i] = elapsed.Milliseconds()

		if i%100 == 0 {
			fmt.Println("Iteration: ", i)
		}
	}

	fmt.Println("Results: ", results)
	// Open the file for writing (create if not exists, overwrite if exists)
	file, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Convert the slice of int64 to a slice of strings (CSV format expects strings)
	stringValues := make([]string, len(results))
	for i, num := range results {
		stringValues[i] = strconv.FormatInt(num, 10) // Convert int64 to string
	}

	// Write the slice of strings as a single row to the CSV
	err = writer.Write(stringValues)
	if err != nil {
		log.Fatalf("Failed to write to CSV: %v", err)
	}

	log.Println("CSV written successfully")

}
