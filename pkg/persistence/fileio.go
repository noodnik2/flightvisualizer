package persistence

import "os"

type Saver interface {
    Save(string, []byte) error
}

type Loader interface {
    Load(string) ([]byte, error)
}

type FileSaverWriter func(filePath string, contents []byte) error

type FileLoaderReader func(filePath string) ([]byte, error)

type FileSaver struct {
    Writer FileSaverWriter
}

type FileLoader struct {
    Reader FileLoaderReader
}

func (rs *FileSaver) Save(fnRef string, contents []byte) error {
    return rs.save(os.WriteFile, fnRef, contents)
}

func (fl *FileLoader) Load(fnRef string) (contents []byte, err error) {
    return fl.load(os.ReadFile, fnRef)
}

type underWriter func(name string, data []byte, perm os.FileMode) error

type underReader func(name string) ([]byte, error)

func (rs *FileSaver) save(uw underWriter, fnRef string, contents []byte) error {
    if rs.Writer != nil {
        return rs.Writer(fnRef, contents)
    }
    return uw(fnRef, contents, 0644)
}

func (fl *FileLoader) load(ur underReader, fnRef string) (contents []byte, err error) {
    if fl.Reader != nil {
        return fl.Reader(fnRef)
    }
    return ur(fnRef)
}
