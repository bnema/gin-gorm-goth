// WARNING: This snippet must be used in a new project. It will not work in the current project.

package main

import (
	"fmt"
	"net/http"
	"time"
)

var counter int

func main() {
	counter = 1000
	// Create new empty array to store the response times
	var responseTimes []int64

	// Print starting message
	fmt.Println("Please wait while", counter, "requests are sent...")
	// Send 1000 GET requests to the /benchmark route
	for i := 0; i < counter; i++ {
		// Create a new HTTP GET request
		req, err := http.NewRequest("GET", "http://localhost:3000/benchmark", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		// Send the request and measure the response time
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()
		elapsed := time.Since(start)

		// Add the response time to the array
		responseTimes = append(responseTimes, elapsed.Milliseconds())
	}

	// Calculate the average response time
	var sum int64
	for _, v := range responseTimes {
		sum += v
	}
	average := sum / int64(len(responseTimes))

	// Print the average response time
	fmt.Println("Average response time:", average, "ms")

	// Print the requests per second
	fmt.Println("Requests per second:", 1000/average)
}
