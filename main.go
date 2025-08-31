package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"io"
)

type Worker struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	JobPosition string  `json:"job_position"`
	Salaries    float64 `json:"salaries"`
	Payroll     string  `json:"payroll"`
	Currency    string  `json:"currency"`
}

// Helper function to create Worker
func NewWorker(id int, name, email, job, payroll string, salary float64) Worker {
	return Worker{
		ID:          id,
		Name:        name,
		Email:       email,
		JobPosition: job,
		Salaries:    salary,
		Payroll:     payroll,
		Currency:    "NZD",
	}
}

func inputWorkerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read raw body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		return
	}

	// Log raw JSON from Postman
	log.Printf("Raw JSON request: %s\n", string(body))

	// Decode JSON from Postman into a temporary variable, Linearize Format
	var inputData map[string]interface{}
	// err := json.NewDecoder(r.Body).Decode(&inputData)
	err = json.Unmarshal(body, &inputData)
	if err != nil {
		log.Printf("Invalid JSON input: %v\n", err)
		return
	}

	// Pretty-print JSON for logging, Not Lineraize Format
	// prettyJSON, err := json.MarshalIndent(inputData, "", "  ")
	// if err != nil {
	// 	log.Printf("Error formatting JSON: %v\n", err)
	// } else {
	// 	log.Printf("Formatted JSON request:\n%s\n", string(prettyJSON))
	// }

	// Create Worker from variable
	worker := NewWorker(
		int(inputData["id"].(float64)),
		inputData["name"].(string),
		inputData["email"].(string),
		inputData["job_position"].(string),
		inputData["payroll"].(string),
		inputData["salaries"].(float64),
	)

	// Log the Worker data
	log.Printf("Received Worker: %+v\n", worker)

	// Log success message
	log.Printf("Successfully sent salary to %s\n", worker.Name)

	// Send response to Postman
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": fmt.Sprintf("Worker %s processed", worker.Name),
	})
}

func main() {
	// Open System.log file (create if not exists, append mode)
	logFile, err := os.OpenFile("System.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// Set log output to file + keep terminal (multi-output)
	multi := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multi)

	http.HandleFunc("/worker/input", inputWorkerHandler)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
