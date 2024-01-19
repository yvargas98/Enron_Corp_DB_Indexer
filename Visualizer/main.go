package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type SearchQuery struct {
	SQL            string `json:"sql"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	From           int    `json:"from"`
	Size           int    `json:"size"`
	TrackTotalHits bool   `json:"track_total_hits"`
	SQLMode        string `json:"sql_mode"`
}

type SearchRequest struct {
	Query SearchQuery `json:"query"`
}

func createSearchRequest(stream string, value string, from int, size int) ([]byte, error) {
	searchQuery := SearchQuery{
		SQL:            fmt.Sprintf("SELECT * FROM %s WHERE match_all_ignore_case('%s') ORDER BY id", stream, value),
		StartTime:      1703900002074496,
		EndTime:        time.Now().UnixMicro(),
		From:           from,
		Size:           size,
		TrackTotalHits: true,
		SQLMode:        "full",
	}

	request := SearchRequest{Query: searchQuery}

	return json.Marshal(request)
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
	searchRequestJSON, err := createSearchRequest(stream, value, from, size)
	if err != nil {
		return nil, err
	}

	envVars, err := getRequiredEnvVars([]string{"SEARCH_SERVER_URL", "SEARCH_SERVER_USERNAME", "SEARCH_SERVER_PASSWORD"})
	if err != nil {
		return nil, err
	}
	url, username, password := envVars["SEARCH_SERVER_URL"], envVars["SEARCH_SERVER_USERNAME"], envVars["SEARCH_SERVER_PASSWORD"]

	client := &http.Client{}
	request, err := http.NewRequest("POST", url+"/_search", bytes.NewReader(searchRequestJSON))
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

	router.Post("/api/default/_search", func(w http.ResponseWriter, r *http.Request) {

		stream := r.URL.Query().Get("stream")
		value := r.URL.Query().Get("value")
		from, err := strconv.Atoi(r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error converting 'from' to int %v", err), http.StatusBadRequest)
		}
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error converting 'size' to int %v", err), http.StatusBadRequest)
		}

		searchResponseBytes, err := search(stream, value, from, size)
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
