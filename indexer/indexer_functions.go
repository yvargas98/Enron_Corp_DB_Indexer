package indexer

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type PropertyDetail struct {
	Type  string `json:"type"`
	Index bool   `json:"index"`
	Store bool   `json:"store"`
}

type Mapping struct {
	Properties map[string]PropertyDetail `json:"properties"`
}

type IndexerData struct {
	Name         string  `json:"name"`
	StorageType  string  `json:"storage_type"`
	ShardNum     int     `json:"shard_num"`
	MappingField Mapping `json:"mappings"`
}

func processFile(path string) (ECEmail, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ECEmail{}, err
	}
	mimeData := extractMIMEFields(content)
	return ECEmail{
		From:                      mimeData["From"],
		To:                        mimeData["To"],
		Subject:                   mimeData["Subject"],
		Content:                   mimeData["Text"],
		Message_ID:                mimeData["Message-ID"],
		Date:                      mimeData["Date"],
		Content_Type:              mimeData["Content-Type"],
		Mime_version:              mimeData["Mime-Version"],
		Content_Transfer_Encoding: mimeData["Content-Transfer-Encoding"],
		X_from:                    mimeData["X-From"],
		X_to:                      mimeData["X-To"],
		X_cc:                      mimeData["X-cc"],
		X_bcc:                     mimeData["X-bcc"],
		X_folder:                  mimeData["X-Folder"],
		X_origin:                  mimeData["X-Origin"],
		X_filename:                mimeData["X-Filename"],
	}, nil
}

func CreateIndexerFromJsonFile(filepath string) (IndexerData, error) {
	var indexerData IndexerData

	file, err := os.Open(filepath)
	if err != nil {
		return indexerData, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&indexerData)
	if err != nil {
		return indexerData, err
	}

	return indexerData, nil
}

func CreateIndexOnZincSearch(indexerData IndexerData) error {
	jsonData, err := json.Marshal(indexerData)
	if err != nil {
		log.Fatal(err)
	}

	ZincSearchUrl := os.Getenv("SEARCH_SERVER_URL")
	ZSusername := os.Getenv("SEARCH_SERVER_USERNAME")
	ZSpassword := os.Getenv("SEARCH_SERVER_PASSWORD")

	req, err := http.NewRequest("POST", ZincSearchUrl+"/enron_corp/_json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(ZSusername, ZSpassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to create indexer, status code: %d", resp.StatusCode)
	}

	return nil
}

func extractMIMEFields(content []byte) map[string]string {
	headers := make(map[string]string)
	var headerBuffer bytes.Buffer
	var bodyBuffer bytes.Buffer

	inHeaders := true

	for _, line := range bytes.Split(content, []byte{'\n'}) {
		if len(line) == 0 {
			inHeaders = false
			continue
		}

		if inHeaders {
			headerBuffer.Write(line)
			headerBuffer.WriteByte('\n')
		} else {
			bodyBuffer.Write(line)
			bodyBuffer.WriteByte('\n')
		}
	}

	mime, err := parseHeaders(headerBuffer.Bytes())
	if err != nil {
		log.Println("Error parsing MIME headers:", err)
		return headers
	}

	headers["From"] = mime.Get("From")
	headers["To"] = mime.Get("To")
	headers["Subject"] = mime.Get("Subject")
	headers["Message-ID"] = mime.Get("Message-ID")
	headers["Date"] = mime.Get("Date")
	headers["Content-Type"] = mime.Get("Content-Type")
	headers["Mime-Version"] = mime.Get("Mime-Version")
	headers["Content-Transfer-Encoding"] = mime.Get("Content-Transfer-Encoding")
	headers["X-From"] = mime.Get("X-From")
	headers["X-To"] = mime.Get("X-To")
	headers["X-cc"] = mime.Get("X-cc")
	headers["X-bcc"] = mime.Get("X-bcc")
	headers["X-Folder"] = mime.Get("X-Folder")
	headers["X-Origin"] = mime.Get("X-Origin")
	headers["X-Filename"] = mime.Get("X-Filename")
	headers["Text"] = bodyBuffer.String()

	return headers
}

func parseHeaders(data []byte) (http.Header, error) {
	header := make(http.Header)
	var key, value string

	for _, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			// This is a continuation of the previous header
			if key != "" {
				header.Add(key, line)
			}
		} else {
			// This is a new header
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key = strings.TrimSpace(parts[0])
				value = strings.TrimSpace(parts[1])
				header.Set(key, value)
			}
		}
	}
	return header, nil
}
