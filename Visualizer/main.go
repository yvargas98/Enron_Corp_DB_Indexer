package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Query struct {
	Term  string `json:"term"`
	Field string `json:"field"`
}

type SearchRequest struct {
	SearchType string   `json:"search_type"`
	Query      Query    `json:"query"`
	SortFields []string `json:"sort_fields"`
	From       int      `json:"from"`
	MaxResults int      `json:"max_results"`
	Source     []string `json:"_source"`
}

func createSearchRequest(value string, from int, size int) ([]byte, error) {
	if value == "" {
		return nil, fmt.Errorf("Value cannot be empty")
	}
	query := Query{
		Term:  value,
		Field: "_all",
	}
	searchRequest := SearchRequest{
		SearchType: "match",
		Query:      query,
		SortFields: []string{"id"},
		From:       from,
		MaxResults: size,
		Source:     []string{"from", "to", "date", "subject", "content"},
	}

	searchRequestJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling search request %v", err)
	}

	return searchRequestJSON, nil
}

func getRequiredEnvVars(vars []string) (map[string]string, error) {
	envVars := make(map[string]string)

	for _, varName := range vars {
		varValue := os.Getenv(varName)

		if varValue == "" {
			return nil, fmt.Errorf("Environment variable %s is not set", varName)
		}

		envVars[varName] = varValue
	}

	return envVars, nil
}

func search(stream string, value string, from int, size int) ([]byte, error) {
	searchRequestJSON, err := createSearchRequest(value, from, size)
	if err != nil {
		return nil, err
	}

	envVars, err := getRequiredEnvVars([]string{"SEARCH_SERVER_URL", "SEARCH_SERVER_USERNAME", "SEARCH_SERVER_PASSWORD"})
	if err != nil {
		return nil, err
	}
	url, username, password := envVars["SEARCH_SERVER_URL"], envVars["SEARCH_SERVER_USERNAME"], envVars["SEARCH_SERVER_PASSWORD"]

	client := &http.Client{}
	request, err := http.NewRequest("POST", url+stream+"/_search", bytes.NewReader(searchRequestJSON))
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(username, password)

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error making HTTP request %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error when searching: %d", response.StatusCode)
	}

	searchResponseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body %v", err)
	}

	return searchResponseBytes, nil
}

func sendSearchResponse(w http.ResponseWriter, searchResponseBytes []byte) {
	var searchResponse map[string]interface{}
	if err := json.Unmarshal(searchResponseBytes, &searchResponse); err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshaling search response %v", err), http.StatusInternalServerError)
		return
	}

	hitsBytes, err := json.Marshal(searchResponse["hits"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshaling hits %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(hitsBytes)
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.AllowContentType("application/json", "text/xml"))

	static(router)

	if len(os.Args) < 2 {
		fmt.Println("Port is missing.")
		return
	}
	port := os.Args[2]

	router.Post("/api/_search", func(w http.ResponseWriter, r *http.Request) {

		var searchParams struct {
			Stream string `json:"stream"`
			Value  string `json:"value"`
			From   int    `json:"from"`
			Size   int    `json:"size"`
		}

		err := json.NewDecoder(r.Body).Decode(&searchParams)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding request body: %v", err), http.StatusBadRequest)
			return
		}

		searchResponseBytes, err := search(searchParams.Stream, searchParams.Value, searchParams.From, searchParams.Size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sendSearchResponse(w, searchResponseBytes)
	})

	fmt.Printf("Mamuro is running in http://localhost:%v\n", port)
	http.ListenAndServe(":"+port, router)
}

func static(r *chi.Mux) {
	root := "./View/dist"
	fs := http.FileServer(http.Dir(root))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
