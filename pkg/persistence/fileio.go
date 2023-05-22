package persistence

import "os"

type Saver interface {
    Save(string, []byte) error
}

type Loader interface {
    Load(string) ([]byte, error)
}

type FileSaver struct {
    Writer func(filePath string, contents []byte) error
}

type FileLoader struct {
    Reader func(filePath string) ([]byte, error)
}

func (rs *FileSaver) Save(fnRef string, contents []byte) error {
    if rs.Writer != nil {
        return rs.Writer(fnRef, contents)
    }
    return os.WriteFile(fnRef, contents, 0644)
}

func (fl *FileLoader) Load(fnRef string) (contents []byte, err error) {
    if fl.Reader != nil {
        return fl.Reader(fnRef)
    }
    return os.ReadFile(fnRef)
}
