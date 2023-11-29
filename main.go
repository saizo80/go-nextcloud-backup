package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

type NextcloudOptions struct {
	Url      string
	User     string
	Password string
}

func main() {
	localFilePath := "./Vault-Cryptomator.zip"
	webdavFilePath := "/Vault-Cryptomator.zip"

	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()

	options := NextcloudOptions{
		Url:      os.Getenv("NEXTCLOUD_URL"),
		User:     os.Getenv("NEXTCLOUD_USER"),
		Password: os.Getenv("NEXTCLOUD_PASSWORD"),
	}

	// if bigger than 512MB, use chunked upload
	if fileSize > 512*1024*1024 {
		chunkedUpload(localFilePath, fileSize, webdavFilePath, options)
	}
}

func chunkedUpload(localFilePath string, fileSize int64, webdavFilePath string, options NextcloudOptions) {
	chunkUploadUrl := strings.Replace(options.Url, "/remote.php/dav/files/", "/remote.php/dav/uploads/", 1)
	destination := fmt.Sprintf("Destination %s%s", options.Url, webdavFilePath)
	client := &http.Client{}

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	uuidStr := uuid.String()

	err = makeCol(client, uuidStr, chunkUploadUrl, destination, options)
	if err != nil {
		panic(err)
	}
	fmt.Println("MKCOL request successful")
}

func makeCol(
	client *http.Client,
	uuidStr string,
	chunkUploadUrl string,
	destination string,
	options NextcloudOptions,
) error {
	// Create a new request
	req, err := http.NewRequest("MKCOL", chunkUploadUrl+"/"+uuidStr, nil)
	if err != nil {
		return err
	}
	// Set basic authentication
	req.SetBasicAuth(options.User, options.Password)
	// Set the Destination header
	req.Header.Set("Destination", destination)
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		return (fmt.Errorf("MKCOL request failed with status code: %d", resp.StatusCode))
	}
	return nil
}
