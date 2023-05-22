package persistence

//import (
//    "os"
//)
//
//type Saver interface {
//    Save(string, []byte) (string, error)
//}
//
//type Loader interface {
//    Load(string) (string, []byte, error)
//}
//
//type FileSaver struct {
//    Writer func(filePath string, contents []byte) error
//}
//
////type HrefSaver struct {
////	FileSaver
////}
//
//type FileLoader struct {
//    Reader func(filePath string) ([]byte, error)
//}
//
////type FileContext struct {
////	Reader func(filePath string) ([]byte, error)
////	Writer func(filePath string, contents []byte) error
////}
//
////func (hs *HrefSaver) Save(href string, contents []byte) (string, error) {
////	return hs.FileSaver.Save(fnFragmentFromHref(href), contents)
////}
//
//func (rs *FileSaver) Save(fnRef string, contents []byte) (filePath string, err error) {
//    filePath = fnRef
//    if rs.Writer != nil {
//        err = rs.Writer(filePath, contents)
//        return
//    }
//    err = os.WriteFile(filePath, contents, 0644)
//    return
//}
//
//func (fl *FileLoader) Load(fnRef string) (filePath string, contents []byte, err error) {
//    filePath = fnRef
//    if fl.Reader != nil {
//        contents, err = fl.Reader(filePath)
//        return
//    }
//    contents, err = os.ReadFile(filePath)
//    return
//}
//
////// GetTsFromTo returns a string representation of a time range using fnPrefixTimestampFormat
////// to format the "from" time, and a subsequence of that for the "to" time, with leading common
////// prefix removed.  Example:
//////
////// { 2023010203040506Z, 2023010203050506Z } => "23010203040506-50506Z" ('5' differs with '4' in tsBase)
////func GetTsFromTo(from, to time.Time) string {
////	fromFmt := from.Format(fnPrefixTimestampFormat)[2:]
////	toFmt := to.Format(fnPrefixTimestampFormat)[2:]
////
////	i := 0
////	for i < len(fromFmt) && i < len(toFmt) && fromFmt[i] == toFmt[i] {
////		i++
////	}
////	return fmt.Sprintf("%s-%s", fromFmt, toFmt[i:])
////}
//
////// JoinPathIfNeeded TODO privatize this - shouldn't be needed outside "persistence" package
////func JoinPathIfNeeded(leftPart, rightPart string) string {
////	if filepath.Dir(rightPart) != "." {
////		// if the right part has a folder associated with it, leave it alone
////		return rightPart
////	}
////	return filepath.Join(leftPart, rightPart)
////}
//
//// fnFragmentFromHref generates a "filename fragment" from an endpoint address through
//// URL escaping, then converting those URL escapes into underscores
////func fnFragmentFromHref(endpoint string) string {
////
////	// convert an endpoint reference into a valid filename fragment
////	encodedEndpoint := url.QueryEscape(endpoint)
////	compile, _ := regexp.Compile(`%[A-Fa-f0-9]{2}`)
////	target := "_"
////	urlFnFragment := compile.ReplaceAllString(encodedEndpoint, target)
////
////	// remove any repeated URL escape substitution markers
////	if len(urlFnFragment) > 1 {
////		noRepeatTargetChars := urlFnFragment[:1]
////		for i := 1; i < len(urlFnFragment); i++ {
////			if !strings.ContainsRune(target, rune(urlFnFragment[i])) || urlFnFragment[i] != urlFnFragment[i-1] {
////				noRepeatTargetChars += string(urlFnFragment[i])
////			}
////		}
////		urlFnFragment = noRepeatTargetChars
////	}
////
////	// return the trimmed filename fragment
////	return trimFragment(urlFnFragment)
////}
//
////
////const fnPrefixTimestampFormat = "20060102150405Z"
//
////func trimFragment(urlFnFragment string) string {
////    return strings.Trim(urlFnFragment, "+_")
////}
