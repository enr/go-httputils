package httputils

import (
	"fmt"
	"github.com/enr/go-files/files"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	// server is a test HTTP server used to provide mock API responses.
	server                   *httptest.Server
	localFilePath            string
	localFileOriginalContent string
	remoteFilePath           string
	remoteFileContent        string
	downloadUrl              string
)

// setup sets up a test HTTP server along with variables representing paths and contents.
func setup() {
	localFilePath = "localFilePath.txt"
	localFileOriginalContent = "localFileOriginalContent"
	remoteFilePath = "/remoteFilePath"
	remoteFileContent = "remoteFileContent"

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, remoteFileContent)
	}))
	downloadUrl = server.URL + remoteFilePath
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

var validUrls = []string{
	"http://www.google.com",
	"https://localhost:8080",
}

var invalidUrls = []string{
	"",
	"   ",
	":",
	"noturl",
	"http",
	"www.google.com",
}

func TestIsValidUrl(t *testing.T) {
	for _, url := range invalidUrls {
		valid := IsValidUrl(url)
		if valid {
			t.Errorf(`Expected invalid "%s"`, url)
		}
	}
	for _, url := range validUrls {
		valid := IsValidUrl(url)
		if !valid {
			t.Errorf(`Expected valid "%s"`, url)
		}
	}
}

func TestDownload_invalidUrl(t *testing.T) {
	for _, url := range invalidUrls {
		deleteFile(localFilePath, t)
		err := download(url, localFilePath)
		if err == nil {
			t.Errorf(`Expected error downloading invalid url "%s" to "%s"`, url, localFilePath)
		}
		if files.Exists(localFilePath) {
			t.Errorf(`File %s created after a failing download`, localFilePath)
		}
	}
	defer deleteFile(localFilePath, t)
}

// if destination file exists, download shouldn't be done
func TestDownloadIfNotExists_fileExists(t *testing.T) {
	setup()
	defer teardown()

	createFileWithContent(localFilePath, localFileOriginalContent, t)
	defer deleteFile(localFilePath, t)

	err := DownloadIfNotExists(downloadUrl, localFilePath)
	if err != nil {
		t.Errorf(`Error downloading "%s" to "%s"`, downloadUrl, localFilePath)
	}

	localFileActualContent := readFile(localFilePath, t)
	if localFileActualContent != localFileOriginalContent {
		t.Errorf(`Local file content, expected [%s] but was [%s]`, localFileOriginalContent, localFileActualContent)
	}
}

func TestDownloadIfNotExists_fileDoesntExist(t *testing.T) {
	setup()
	defer teardown()

	deleteFile(localFilePath, t)
	defer deleteFile(localFilePath, t)

	err := DownloadIfNotExists(downloadUrl, localFilePath)
	if err != nil {
		t.Errorf(`Error downloading "%s" to "%s"`, downloadUrl, localFilePath)
	}

	localFileActualContent := readFile(localFilePath, t)
	if localFileActualContent != remoteFileContent {
		t.Errorf(`Local file content, expected [%s] but was [%s]`, remoteFileContent, localFileActualContent)
	}
}

// if destination file exists, it should be overwritten
func TestDownloadOverwriting_fileExists(t *testing.T) {
	setup()
	defer teardown()

	createFileWithContent(localFilePath, localFileOriginalContent, t)
	defer deleteFile(localFilePath, t)

	err := DownloadOverwriting(downloadUrl, localFilePath)
	if err != nil {
		t.Errorf(`Error downloading "%s" to "%s"`, downloadUrl, localFilePath)
	}

	localFileActualContent := readFile(localFilePath, t)
	if localFileActualContent != remoteFileContent {
		t.Errorf(`Local file content, expected [%s] but was [%s]`, remoteFileContent, localFileActualContent)
	}
}

// destination file should be written even if the server responds with error.
func TestDownloadOverwriting_serverError(t *testing.T) {
	errorMessage := "Bad Request"
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, errorMessage, 400)
	}))
	defer server.Close()
	u := server.URL + remoteFilePath

	createFileWithContent(localFilePath, localFileOriginalContent, t)
	defer deleteFile(localFilePath, t)

	err := DownloadOverwriting(u, localFilePath)
	if err != nil {
		t.Errorf(`Error downloading "%s" to "%s"`, u, localFilePath)
	}

	localFileActualContent := strings.TrimSpace(readFile(localFilePath, t))
	if localFileActualContent != errorMessage {
		t.Errorf(`Local file content, expected [%s] but was [%s]`, errorMessage, localFileActualContent)
	}
}

// if destination file exists, it should be saved in backupPath
func TestDownloadPreservingOld_fileExists(t *testing.T) {
	setup()
	defer teardown()

	backupFilePath := "TestDownloadPreservingOld_fileExists.txt.OLD"

	deleteFile(backupFilePath, t)
	defer deleteFile(backupFilePath, t)
	createFileWithContent(localFilePath, localFileOriginalContent, t)
	defer deleteFile(localFilePath, t)

	err := DownloadPreservingOld(downloadUrl, localFilePath, backupFilePath)
	if err != nil {
		t.Errorf(`Error downloading "%s" to "%s"`, downloadUrl, localFilePath)
	}

	if !files.Exists(backupFilePath) {
		t.Errorf(`Backup file not created: %s`, backupFilePath)
	}

	localFileActualContent := readFile(localFilePath, t)
	if localFileActualContent != remoteFileContent {
		t.Errorf(`Local file content, expected [%s] but was [%s]`, remoteFileContent, localFileActualContent)
	}

	backupFileActualContent := readFile(backupFilePath, t)
	if backupFileActualContent != localFileOriginalContent {
		t.Errorf(`Backup file content, expected [%s] but was [%s]`, localFileOriginalContent, backupFileActualContent)
	}
}

type args struct {
	src         string
	dst         string
	destination string
}

var existsData = []args{
	{"http://example/afile.zip", "myfile.zip", "myfile.zip"},
	{"http://example/afile.zip", "./", "afile.zip"},
	{"http://example/afile.zip", ".", "afile.zip"},
	{"http://example/afile.zip", "   ", "afile.zip"},
}

func TestBuildDestinationUrl(t *testing.T) {
	for _, data := range existsData {
		destination := buildDestinationPath(data.dst, data.src)
		if destination != data.destination {
			t.Errorf(`destination expected "%s" but was "%s"`, data.destination, destination)
		}
	}
}

func createFileWithContent(path, content string, t *testing.T) {
	deleteFile(path, t)
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Errorf(`Error writing file "%s"`, path)
	}
}

func readFile(path string, t *testing.T) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf(`Error reading file "%s"`, path)
	}
	return string(b)
}

func deleteFile(path string, t *testing.T) {
	if files.Exists(path) {
		err := os.Remove(path)
		if err != nil {
			t.Error("error deleting path", path)
		}
	}
}
