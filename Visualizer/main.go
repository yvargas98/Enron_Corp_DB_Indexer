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

	// "strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func createSearchRequest(stream string, value string, from int, size int) []byte {
	searchRequest := struct {
		Query struct {
			SQL            string `json:"sql"`
			StartTime      int64  `json:"start_time"`
			EndTime        int64  `json:"end_time"`
			From           int    `json:"from"`
			Size           int    `json:"size"`
			TrackTotalHits bool   `json:"track_total_hits"`
			SQLMode        string `json:"sql_mode"`
		} `json:"query"`
	}{
		Query: struct {
			SQL            string `json:"sql"`
			StartTime      int64  `json:"start_time"`
			EndTime        int64  `json:"end_time"`
			From           int    `json:"from"`
			Size           int    `json:"size"`
			TrackTotalHits bool   `json:"track_total_hits"`
			SQLMode        string `json:"sql_mode"`
		}{
			SQL:            fmt.Sprintf("SELECT * FROM %s WHERE match_all('%s') ORDER BY id", stream, value),
			StartTime:      1703900002074496,
			EndTime:        time.Now().UnixMicro(),
			From:           from,
			Size:           size,
			TrackTotalHits: true,
			SQLMode:        "full",
		},
	}

	searchRequestJSON, err := json.Marshal(searchRequest)
	if err != nil {
		panic(err)
	}

	return searchRequestJSON
}

// Realiza la búsqueda
func search(stream string, value string, from int, size int) ([]byte, error) {
	searchRequestJSON := createSearchRequest(stream, value, from, size)

	client := &http.Client{}
	request, err := http.NewRequest("POST", "http://localhost:5080/api/default/_search", bytes.NewReader(searchRequestJSON))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("yvargas.vargasgodoy@gmail.com", "@Va221998")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error when searching: %d", response.StatusCode)
	}

	searchResponseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return searchResponseBytes, nil
}

func sendSearchResponse(w http.ResponseWriter, searchResponseBytes []byte) {
	var searchResponse map[string]interface{}
	if err := json.Unmarshal(searchResponseBytes, &searchResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hitsBytes, err := json.Marshal(searchResponse["hits"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		size, err := strconv.Atoi(r.URL.Query().Get("size"))

		searchResponseBytes, err := search(stream, value, from, size)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		sendSearchResponse(w, searchResponseBytes)
	})
	apiV1 := chi.NewRouter()
	router.Mount("/api/default", apiV1)
	fmt.Printf("Server is running at port %v\n", port)
	http.ListenAndServe(":"+port, router)
}

func static(r *chi.Mux) {
	root := "./View/dist"
	fs := http.FileServer(http.Dir(root))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
