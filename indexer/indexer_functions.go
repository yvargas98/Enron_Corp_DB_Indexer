package indexer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

func GetFolders(folderName string) []string {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	var folders []string

	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}
	return folders
}

func GetFiles(folderName string) []string {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	var fileNames []string
	for _, file := range files {
		if file.IsDir() == false && file.Name() != ".DS_Store" {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func ProcessFile(path string, id int) (ECEmail, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ECEmail{}, err
	}
	mimeData := extractMIMEFields(content)
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

func extractMIMEFields(content []byte) map[string]string {
	headers := make(map[string]string) //Mapa clave - valor
	var headerBuffer bytes.Buffer
	var bodyBuffer bytes.Buffer

	inHeaders := true

	for _, line := range bytes.Split(content, []byte{'\n'}) { //divide el archivo del email en lineas
		if len(line) == 0 {
			inHeaders = false
			continue
		}

		if inHeaders {
			headerBuffer.Write(line) //lineas de encabezados del email
			headerBuffer.WriteByte('\n')
		} else {
			bodyBuffer.Write(line) //mensaje del email
			bodyBuffer.WriteByte('\n')
		}
	}

	header, err := parseHeaders(headerBuffer.Bytes())
	if err != nil {
		log.Println("Error parsing MIME headers:", err)
		return headers
	}

	headers["From"] = header.Get("From")
	headers["To"] = header.Get("To")
	headers["Subject"] = header.Get("Subject")
	headers["Message-ID"] = header.Get("Message-ID")
	headers["Date"] = header.Get("Date")
	headers["Content-Type"] = header.Get("Content-Type")
	headers["Mime-Version"] = header.Get("Mime-Version")
	headers["Content-Transfer-Encoding"] = header.Get("Content-Transfer-Encoding")
	headers["X-From"] = header.Get("X-From")
	headers["X-To"] = header.Get("X-To")
	headers["X-cc"] = header.Get("X-cc")
	headers["X-bcc"] = header.Get("X-bcc")
	headers["X-Folder"] = header.Get("X-Folder")
	headers["X-Origin"] = header.Get("X-Origin")
	headers["X-Filename"] = header.Get("X-Filename")
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
			if key != "" { //linea de continuación de un header
				header.Add(key, line)
			}
		} else {
			parts := strings.SplitN(line, ":", 2) //header nuevo
			if len(parts) == 2 {
				key = strings.TrimSpace(parts[0]) //para remover espacio en blanco si existe
				value = strings.TrimSpace(parts[1])
				header.Set(key, value)
			}
		}
	}
	return header, nil
}

func getRequiredEnvVar(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("%s environment variable is not set", name)
	}
	return value, nil
}

func PostDataToOpenObserve(data []ECEmail) error {
	jsonData, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %s", err)
	}

	ZincSearchUrl, err := getRequiredEnvVar("SEARCH_SERVER_URL")
	if err != nil {
		return err
	}
	ZSusername, err := getRequiredEnvVar("SEARCH_SERVER_USERNAME")
	if err != nil {
		return err
	}
	ZSpassword, err := getRequiredEnvVar("SEARCH_SERVER_PASSWORD")
	if err != nil {
		return err
	}
	indexName, err := getRequiredEnvVar("INDEX_NAME")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, ZincSearchUrl+"/"+indexName+"/_json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error reading request.", err)
	}

	req.SetBasicAuth(ZSusername, ZSpassword)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error making request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
