package persistence

import (
	"bytes"

	gokml "github.com/twpayne/go-kml/v3"

	"github.com/noodnik2/flightvisualizer/pkg/persistence"
)

type KmzSaver struct {
	persistence.Saver
	Assets map[string]any
}

func (rs *KmzSaver) Save(fnFragment string, contents []byte) error {
	files := make(map[string]any)
	for assetKey, assetValue := range rs.Assets {
		files[assetKey] = assetValue
	}
	files["doc.kml"] = contents

	memoryWriter := &bytes.Buffer{}
	if writeErr := gokml.WriteKMZ(memoryWriter, files); writeErr != nil {
		return writeErr
	}

	return rs.Saver.Save(fnFragment, memoryWriter.Bytes())
}
