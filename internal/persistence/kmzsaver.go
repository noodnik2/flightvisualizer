package persistence

import (
	"os"

	gokml "github.com/twpayne/go-kml/v3"
)

type KmzSaver struct {
	FileSaver
	Assets map[string]any
}

func (rs *KmzSaver) SaveNewKmz(fnFragment string, contents []byte) (saveFilename string, err error) {
	var file *os.File
	file, err = rs.createFile(rs.fnFromFragment(fnFragment))
	if err != nil {
		return
	}
	files := make(map[string]any)
	for assetKey, assetValue := range rs.Assets {
		files[assetKey] = assetValue
	}
	files["doc.kml"] = contents

	err = gokml.WriteKMZ(file, files)
	err = rs.closeFile(file, err)
	return
}
