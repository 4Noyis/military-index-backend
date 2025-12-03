package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8080/health"
	numRequests := 10

	fmt.Println("========================================")
	fmt.Println("Rate Limiter Burst Test")
	fmt.Println("Sending", numRequests, "simultaneous requests")
	fmt.Println("Expected limit: 3 requests/minute")
	fmt.Println("========================================\n")

	var wg sync.WaitGroup
	results := make(chan string, numRequests)

	// Send all requests at once
	startTime := time.Now()
	for i := 1; i <= numRequests; i++ {
		wg.Add(1)
		go func(reqNum int) {
			defer wg.Done()

			resp, err := http.Get(url)
			if err != nil {
				results <- fmt.Sprintf("Request %2d: ERROR - %v", reqNum, err)
				return
			}
			defer resp.Body.Close()
			io.ReadAll(resp.Body) // Consume body

			if resp.StatusCode == 200 {
				results <- fmt.Sprintf("[OK]  Request %2d: SUCCESS (200)", reqNum)
			} else if resp.StatusCode == 429 {
				results <- fmt.Sprintf("[LIMIT] Request %2d: RATE LIMITED (429)", reqNum)
			} else {
				results <- fmt.Sprintf("[?] Request %2d: HTTP %d", reqNum, resp.StatusCode)
			}
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	close(results)
	elapsed := time.Since(startTime)

	// Print results
	success := 0
	limited := 0
	for result := range results {
		fmt.Println(result)
		if result[1] == 'O' {
			success++
		} else if result[1] == 'L' {
			limited++
		}
	}

	fmt.Printf("\nCompleted in: %v\n", elapsed)
	fmt.Printf("Results: %d successful, %d rate-limited\n", success, limited)

	if limited > 0 {
		fmt.Println("\n✅ Rate limiter is WORKING!")
	} else {
		fmt.Println("\n⚠️  All requests succeeded - rate limiter may not be triggering")
	}
}
