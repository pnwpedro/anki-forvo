// package anki_forvo_plugin
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type FileDownloader struct {
	client *http.Client
}

func NewFileDownloader() *FileDownloader {
	dl := FileDownloader{
		client: http.DefaultClient,
	}
	return &dl
}

func (fd *FileDownloader) Download(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
