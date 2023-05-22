package persistence

//
//import (
//	"fmt"
//	"log"
//	"net/url"
//	"os"
//	"path/filepath"
//	"regexp"
//	"strings"
//	"time"
//)
//
//type FileContext struct {
//	Folder     string
//	Suffix     string
//	MaxBaseLen int
//	Reader     func(filePath string) ([]byte, error)
//	Writer     func(filePath string, contents []byte) error
//}
//
//func (rs *FileContext) Reconcile(filePath string) (FileContext, string) {
//	newFileContext := *rs
//	fullFilePath := JoinPathIfNeeded(rs.Folder, filePath)
//	newFileContext.Folder = filepath.Dir(fullFilePath)
//	baseFn := filepath.Base(fullFilePath)
//	return newFileContext, JoinPathIfNeeded(newFileContext.Folder, rs.baseFromFnRef(baseFn))
//}
//
//func (rs *FileContext) SaveFromHref(href string, contents []byte) (string, error) {
//	return rs.SaveFromFnRef(fnFragmentFromHref(href), contents)
//}
//
//func (rs *FileContext) SaveFromFnRef(fnRef string, contents []byte) (filePath string, err error) {
//	_, filePath = rs.Reconcile(fnRef)
//	if rs.Writer != nil {
//		err = rs.Writer(filePath, contents)
//		return
//	}
//	err = os.WriteFile(filePath, contents, 0644)
//	return
//}
//
//func (rs *FileContext) LoadFromFnRef(fnRef string) (filePath string, contents []byte, err error) {
//	_, filePath = rs.Reconcile(fnRef)
//	if rs.Reader != nil {
//		contents, err = rs.Reader(filePath)
//		return
//	}
//	contents, err = os.ReadFile(filePath)
//	return
//}
//
//// GetTsFromTo returns a string representation of a time range using fnPrefixTimestampFormat
//// to format the "from" time, and a subsequence of that for the "to" time, with leading common
//// prefix removed.  Example:
////
//// { 2023010203040506Z, 2023010203050506Z } => "23010203040506-50506Z" ('5' differs with '4' in tsBase)
//func GetTsFromTo(from, to time.Time) string {
//	fromFmt := from.Format(fnPrefixTimestampFormat)[2:]
//	toFmt := to.Format(fnPrefixTimestampFormat)[2:]
//
//	i := 0
//	for i < len(fromFmt) && i < len(toFmt) && fromFmt[i] == toFmt[i] {
//		i++
//	}
//	return fmt.Sprintf("%s-%s", fromFmt, toFmt[i:])
//}
//
//// JoinPathIfNeeded TODO privatize this - shouldn't be needed outside "persistence" package
//func JoinPathIfNeeded(leftPart, rightPart string) string {
//	//if filepath.IsAbs(rightPart) {
//	//	return rightPart
//	//}
//	if filepath.Dir(rightPart) != "." {
//		return rightPart
//	}
//	return filepath.Join(leftPart, rightPart)
//}
//
//// fnFragmentFromHref generates a "filename fragment" from an endpoint address through
//// URL escaping, then converting those URL escapes into underscores
//func fnFragmentFromHref(endpoint string) string {
//
//	// convert an endpoint reference into a valid filename fragment
//	encodedEndpoint := url.QueryEscape(endpoint)
//	compile, _ := regexp.Compile(`%[A-Fa-f0-9]{2}`)
//	target := "_"
//	urlFnFragment := compile.ReplaceAllString(encodedEndpoint, target)
//
//	// enforce a length limit on the filename fragment
//	if len(urlFnFragment) > 1 {
//		noRepeatTargetChars := urlFnFragment[:1]
//		for i := 1; i < len(urlFnFragment); i++ {
//			if !strings.ContainsRune(target, rune(urlFnFragment[i])) || urlFnFragment[i] != urlFnFragment[i-1] {
//				noRepeatTargetChars += string(urlFnFragment[i])
//			}
//		}
//		urlFnFragment = noRepeatTargetChars
//	}
//
//	// return the trimmed filename fragment
//	return trimFragment(urlFnFragment)
//}
//
//// fnPrefixTimestampFormat TODO privatize this - I believe it shouldn't be needed outside of this package
//const fnPrefixTimestampFormat = "20060102150405Z"
//
//// baseFromFnRef returns the "base filename" (i.e., just the last part of
//// the path) given the input "fnRef"
//func (rs *FileContext) baseFromFnRef(fnRef string) string {
//	resultFilename := fnRef
//	maxBaseLenToUse := rs.MaxBaseLen
//	if maxBaseLenToUse == 0 {
//		maxBaseLenToUse = 64
//	}
//	if len(resultFilename) > maxBaseLenToUse {
//		resultFilename = trimFragment(resultFilename[:maxBaseLenToUse])
//		log.Printf("INFO: truncated filename reference(%s) to length %d\n", fnRef, maxBaseLenToUse)
//	}
//
//	withSuffix := resultFilename
//	if rs.Suffix != "" && !strings.HasSuffix(strings.ToLower(withSuffix), strings.ToLower(rs.Suffix)) {
//		withSuffix += rs.Suffix
//	}
//	return withSuffix
//}
//
//func trimFragment(urlFnFragment string) string {
//	return strings.Trim(urlFnFragment, "+_")
//}
