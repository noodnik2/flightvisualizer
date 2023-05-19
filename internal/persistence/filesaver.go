package persistence

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type FileSaver struct {
	OutputDir      string
	Timestamp      *time.Time
	FilenameSuffix string
}

func (rs *FileSaver) SaveEndpointResponse(endpoint string, contents []byte) (string, error) {
	return rs.SaveNewFile(rs.fnFragmentFromEndpoint(endpoint), contents)
}

func (rs *FileSaver) SaveNewFile(fnFragment string, contents []byte) (saveFilename string, err error) {
	saveFilename = rs.fnFromFragment(fnFragment)
	err = rs.saveNewFile(saveFilename, contents)
	return
}

const maximumFilenameLength = 64

func (rs *FileSaver) fnFragmentFromEndpoint(endpoint string) string {

	encodedEndpoint := url.QueryEscape(endpoint)
	compile, _ := regexp.Compile(`%[A-Fa-f0-9]{2}`)
	target := "_"
	urlFnFragment := compile.ReplaceAllString(encodedEndpoint, target)
	if len(urlFnFragment) > 1 {
		noRepeatTargetChars := urlFnFragment[:1]
		for i := 1; i < len(urlFnFragment); i++ {
			if !strings.ContainsRune(target, rune(urlFnFragment[i])) || urlFnFragment[i] != urlFnFragment[i-1] {
				noRepeatTargetChars += string(urlFnFragment[i])
			}
		}
		urlFnFragment = noRepeatTargetChars
	}
	return trimFragment(urlFnFragment)
}

const FnPrefixTimestampFormat = "20060102150405Z"

func (rs *FileSaver) fnFromFragment(fnFragment string) string {
	var filenamePrefix string
	if rs.Timestamp != nil {
		timeStamp := rs.Timestamp.Format(FnPrefixTimestampFormat)
		filenamePrefix = fmt.Sprintf("%s-", timeStamp[2:])
	}
	fullLength := filenamePrefix + fnFragment

	resultFilename := fullLength
	maxLengthWithoutSuffix := maximumFilenameLength - len(rs.FilenameSuffix)
	if len(resultFilename) > maxLengthWithoutSuffix {
		resultFilename = trimFragment(resultFilename[0:maxLengthWithoutSuffix])
		log.Printf("INFO: truncated filename(%s) to length %d\n", fullLength, maximumFilenameLength)
	}

	withSuffix := resultFilename
	if rs.FilenameSuffix != "" && !strings.HasSuffix(strings.ToLower(withSuffix), rs.FilenameSuffix) {
		withSuffix += rs.FilenameSuffix
	}
	return withSuffix
}

func (rs *FileSaver) saveNewFile(saveFilename string, data []byte) error {
	saveFile, createErr := rs.createFile(saveFilename)
	if createErr != nil {
		return createErr
	}
	_, writeErr := saveFile.Write(data)
	return rs.closeFile(saveFile, writeErr)
}

func (rs *FileSaver) closeFile(saveFile *os.File, existingError error) error {
	if closeErr := saveFile.Close(); closeErr != nil {
		if existingError == nil {
			return closeErr
		}
		// an existing error is more important than a close error
	}
	return existingError
}

func (rs *FileSaver) createFile(saveFilename string) (*os.File, error) {
	filePath := filepath.Join(rs.OutputDir, saveFilename)
	log.Printf("INFO: creating file: %s\n", filePath)
	saveFile, openErr := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0664)
	if openErr != nil {
		return nil, openErr
	}
	return saveFile, nil
}

func trimFragment(urlFnFragment string) string {
	return strings.Trim(urlFnFragment, "+_")
}
