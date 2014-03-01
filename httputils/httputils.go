package httputils

import (
	//"fmt"
	"github.com/enr/go-files/files"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func IsValidUrl(arg string) bool {
	if strings.TrimSpace(arg) == "" {
		return false
	}
	//fmt.Printf("ARG %s\n", arg)
	u, err := url.Parse(arg)
	//fmt.Printf("U %v\n", u)
	//fmt.Printf("ERR %v\n", err)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func DownloadIfNotExists(src, destination string) error {
	dst := buildDestinationPath(destination, src)
	if files.Exists(dst) {
		return nil
	}
	return download(src, dst)
}

func DownloadOverwriting(src, destination string) error {
	dst := buildDestinationPath(destination, src)
	return download(src, dst)
	/*
		destination := dst
		if files.Exists(dst) {
			tmpDir := filepath.Dir(dst)
			tmp, err := ioutil.TempFile(tmpDir, "")
			if err != nil {
				return err
			}
			destination = path.Join(tmpDir, tmp.Name())
		}
		err := download(src, destination)
		if err != nil {
			return err
		}
		err = os.Remove(dst)
		if err != nil {
			return err
		}
		err = os.Rename(destination, dst)
		if err != nil {
			return err
		}
		return nil
	*/
}

func DownloadPreservingOld(src, destination, backup string) error {
	dst := buildDestinationPath(destination, src)
	if files.Exists(dst) {
		err := os.Rename(dst, backup)
		if err != nil {
			return err
		}
	}
	return download(src, dst)
}

func download(src, destination string) error {
	res, err := http.Get(src)
	if err != nil {
		return err
	}
	responseBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destination, responseBody, 0644)
	if err != nil {
		return err
	}
	return nil
}

func buildDestinationPath(destination, source string) string {
	dst := strings.TrimSpace(destination)
	if dst == "" || files.IsDir(destination) {
		return path.Join(dst, filepath.Base(source))
	}
	return dst
}
