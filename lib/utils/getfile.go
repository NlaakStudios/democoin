package utils

import (
	"io"
	"net/http"
	"os"
)

type SyncFiles struct {
}

/* Example Usage:
func main() {

	fileUrl := "https://golangcode.com/images/avatar.jpg"

	err := DownloadFile("avatar.jpg", fileUrl)
	if err != nil {
		panic(err)
	}

}

*/
// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func (s *SyncFiles) downloadFile(filepath string, url string) error {

	// Create the file. Full path with filename
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

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
