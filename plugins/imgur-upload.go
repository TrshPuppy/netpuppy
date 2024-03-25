package plugins

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "bytes"
    "time"
)

type ImgurResponse struct {
    Data struct {
        Link string `json:"link"`
    } `json:"data"`
}

func init() {
    Register("imgur-upload", &ImgurUpload{})
}

type ImgurUpload struct{}

func (i *ImgurUpload) Execute(pluginDataChan chan<- string) {
    fmt.Println("[Imgur Image Upload]")
    fmt.Println("Enter file path:")
    var filePath string
    fmt.Scanln(&filePath)

    fileData, err := ioutil.ReadFile(filePath)
    if err != nil {
        pluginDataChan <- fmt.Sprintf("Error reading file: %v", err)
        return
    }

    base64Data := base64.StdEncoding.EncodeToString(fileData)

    client := &http.Client{Timeout: 30 * time.Second}
    apiUrl := "https://api.imgur.com/3/image"
    payload := map[string]string{"image": base64Data}
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        pluginDataChan <- fmt.Sprintf("Error preparing request data: %v", err)
        return
    }

    req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(payloadBytes))
    if err != nil {
        pluginDataChan <- fmt.Sprintf("Error creating HTTP request: %v", err)
        return
    }

    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        pluginDataChan <- fmt.Sprintf("Error sending request to Imgur: %v", err)
        return
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        pluginDataChan <- fmt.Sprintf("Error reading response body: %v", err)
        return
    }

    var imgurResponse ImgurResponse
    err = json.Unmarshal(responseBody, &imgurResponse)
    if err != nil {
        pluginDataChan <- fmt.Sprintf("Error parsing JSON response: %v", err)
        return
    }

    pluginDataChan <- fmt.Sprintf("Imgur Image URL: %s", imgurResponse.Data.Link)
}

