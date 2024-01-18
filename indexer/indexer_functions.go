package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type ECEmail struct {
	ID                        int    `json:"id"`
	Message_ID                string `json:"message_id"`
	Date                      string `json:"date"`
	From                      string `json:"from"`
	To                        string `json:"to"`
	Subject                   string `json:"subject"`
	Cc                        string `json:"cc"`
	Mime_version              string `json:"mime_version"`
	Content_Type              string `json:"content_type"`
	Content_Transfer_Encoding string `json:"content_transfer_encoding"`
	Bcc                       string `json:"bcc"`
	X_from                    string `json:"x_from"`
	X_to                      string `json:"x_to"`
	X_cc                      string `json:"x_cc"`
	X_bcc                     string `json:"x_bcc"`
	X_folder                  string `json:"x_folder"`
	X_origin                  string `json:"x_origin"`
	X_filename                string `json:"x_filename"`
	Content                   string `json:"content"`
}

type IndexerError struct {
	Message string
	Err     error
}

func (e *IndexerError) Error() string { //método de la estructura IndexerError
	return fmt.Sprintf("%s: %s", e.Message, e.Err)
}

func GetFolders(folderName string) ([]string, error) {
	files, err := os.ReadDir(folderName)
	if err != nil {
		return nil, &IndexerError{Message: "Error getting folder", Err: err}
	}

	var folders []string
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}

	return folders, nil
}

func GetFiles(folderName string) ([]string, error) {
	files, err := os.ReadDir(folderName)
	if err != nil {
		return nil, &IndexerError{Message: "Error reading file", Err: err}
	}

	var fileNames []string
	for _, file := range files {
		if file.IsDir() == false && file.Name() != ".DS_Store" {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}

func ProcessFile(path string, id int) (ECEmail, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ECEmail{}, &IndexerError{Message: "Error procesing file", Err: err}
	}
	mimeData, err := formatEmailContent(content)
	if err != nil {
		return ECEmail{}, err
	}
	return ECEmail{
		ID:                        id,
		Message_ID:                mimeData["Message-ID"],
		Date:                      mimeData["Date"],
		From:                      mimeData["From"],
		To:                        mimeData["To"],
		Subject:                   mimeData["Subject"],
		Cc:                        mimeData["Cc"],
		Mime_version:              mimeData["Mime-Version"],
		Content_Type:              mimeData["Content-Type"],
		Content_Transfer_Encoding: mimeData["Content-Transfer-Encoding"],
		Bcc:                       mimeData["Bcc"],
		X_from:                    mimeData["X-From"],
		X_to:                      mimeData["X-To"],
		X_cc:                      mimeData["X-cc"],
		X_bcc:                     mimeData["X-bcc"],
		X_folder:                  mimeData["X-Folder"],
		X_origin:                  mimeData["X-Origin"],
		X_filename:                mimeData["X-Filename"],
		Content:                   mimeData["Text"],
	}, nil
}

func formatEmailContent(data []byte) (map[string]string, error) {
	headers := make(map[string]string)
	var key, value string
	var bodyBuffer bytes.Buffer
	inHeaders := true

	for _, line := range bytes.Split(data, []byte{'\n'}) {
		if len(line) == 0 {
			inHeaders = false
			continue
		}

		parts := strings.SplitN(string(line), ":", 2) //header nuevo
		if len(parts) == 2 {
			key = strings.TrimSpace(parts[0]) //para remover espacio en blanco si existe
			value = strings.TrimSpace(parts[1])
			headers[key] = value
		}

		if !inHeaders {
			bodyBuffer.Write(line) //mensaje del email
			bodyBuffer.WriteByte('\n')
		}

		headers["Text"] = bodyBuffer.String()
	}
	return headers, nil
}

func getRequiredEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(&IndexerError{Message: "An error ocurred", Err: fmt.Errorf("%s environment variable is not set", name)})
	}
	return value
}

func PostDataToOpenObserve(data []ECEmail) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return &IndexerError{Message: "Error marshaling JSON ", Err: err}
	}

	openObserveUrl := getRequiredEnvVar("SEARCH_SERVER_URL")
	openObserveUsername := getRequiredEnvVar("SEARCH_SERVER_USERNAME")
	openObservePassword := getRequiredEnvVar("SEARCH_SERVER_PASSWORD")
	indexName := getRequiredEnvVar("INDEX_NAME")

	req, err := http.NewRequest(http.MethodPost, openObserveUrl+"/"+indexName+"/_json", bytes.NewBuffer(jsonData))
	if err != nil {
		return &IndexerError{Message: "Error creating request", Err: err}
	}

	req.SetBasicAuth(openObserveUsername, openObservePassword)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &IndexerError{Message: fmt.Sprintf("Error making request, status code %s and error %s ", resp.Status, err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &IndexerError{Message: "An error ocurred", Err: fmt.Errorf("Unexpected status code: %d", resp.StatusCode)}
	}

	return nil
}
