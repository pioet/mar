/*
Copyright © 2025 2025 Pioet <pioet@aliyun.com>
*/
package models

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/pioet/mar/internal/utils"
)

// URI type
type URIType int

const (
	TypeUnknown URIType = iota
	TypeWebLink
	TypeLocalFile
)

type Uri struct {
	UriType   URIType
	UriInput  string
	UriOutput string
	Title     string
}

func NewUri(uriInput string) (*Uri, error) {
	u := &Uri{
		UriInput: uriInput,
	}
	// case: uri -> web link
	if u.isWebLink() {
		u.UriType = TypeWebLink
		u.UriOutput = uriInput
		webTitle, err := utils.GetWebTitle(uriInput)
		if err != nil {
			u.Title = "Untitled"
		} else {
			u.Title = webTitle
		}
	} else if u.isLocalFile() {
		// case: uri -> local file
		u.UriType = TypeLocalFile
		absPath, err := filepath.Abs(uriInput)
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path: %w", err)
		}
		u.UriOutput = absPath
		u.Title = filepath.Base(absPath) // file name as title
	} else {
		return u, fmt.Errorf("unrecognized URI type: %s", uriInput)
	}
	return u, nil
}

func (u *Uri) isWebLink() bool {
	parsedURL, err := url.Parse(u.UriInput)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

func (u *Uri) isLocalFile() bool {
	// check existence of the file
	info, err := os.Stat(u.UriInput)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false // it is folder, not a file
	}
	return true
}

func (u *Uri) GetTitle() (string, error) {
	return u.Title, nil
}
