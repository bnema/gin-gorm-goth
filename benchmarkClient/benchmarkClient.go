package main

import (
	"fmt"
	"net/http"
	"time"
)

var counter int

func main() {
	counter := 1000
	// Create a new HTTP client
	client := &http.Client{}

	// Create new empty array to store the response times
	var responseTimes []int64

	// Send the requests and measure the response time
	for i := 0; i < counter; i++ {
		// Create a new GET request
		req, err := http.NewRequest("GET", "http://localhost:3000/benchmark/test", nil)
		if err != nil {
			fmt.Println("Erreur lors de la création de la requête:", err)
			return
		}

		// Send the request and measure the response time
		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erreur lors de l'envoi de la requête:", err)
			return
		}
		defer resp.Body.Close()
		elapsed := time.Since(start)

		// Add the response time to the array
		fmt.Println("Temps de réponse:", elapsed.Milliseconds(), "ms")

		// Add the response time to the array
		responseTimes = append(responseTimes, elapsed.Milliseconds())
	}

	// Calculate the average response time
	var sum int64
	for _, v := range responseTimes {
		sum += v
	}
	average := sum / int64(len(responseTimes))

	// Show the average response time
	fmt.Println("Moyenne des temps de réponse:", average, "ms")

	// Show the number of requests per second
	fmt.Println("Nombre de requêtes par seconde:", 1000/average)

}
