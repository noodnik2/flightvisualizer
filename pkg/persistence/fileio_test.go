package persistence

import (
    "errors"
    "os"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestFileSaver_Save(t *testing.T) {
    testCases := []struct {
        name      string
        hasWriter bool
        hasError  bool
    }{
        {
            name: "no writer, no error",
        },
        {
            name:      "has writer, no error",
            hasWriter: true,
        },
        {
            name:     "no writer, error",
            hasError: true,
        },
        {
            name:      "has writer, error",
            hasWriter: true,
            hasError:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            requirer := require.New(t)

            var writer FileSaverWriter
            var actualWriterContentsSaved []byte
            var actualWriterFilenameSavedTo string
            if tc.hasWriter {
                writer = func(filePath string, contents []byte) error {
                    if tc.hasError {
                        return errors.New(tc.name)
                    }
                    actualWriterFilenameSavedTo = filePath
                    actualWriterContentsSaved = contents
                    return nil
                }
            }

            var uw underWriter
            var actualUwContentsSaved []byte
            var actualUwFilenameSavedTo string
            uw = func(filePath string, contents []byte, perm os.FileMode) error {
                if tc.hasError {
                    return errors.New(tc.name)
                }
                actualUwFilenameSavedTo = filePath
                actualUwContentsSaved = contents
                return nil
            }

            saver := &FileSaver{Writer: writer}
            const filenameSavedTo = "this is the filename"
            const contentsSaved = "this is what was saved"
            err := saver.save(uw, filenameSavedTo, []byte(contentsSaved))
            if tc.hasError {
                requirer.Error(err)
                requirer.Equal(tc.name, err.Error())
            } else {
                requirer.NoError(err)
                if tc.hasWriter {
                    requirer.Equal(contentsSaved, string(actualWriterContentsSaved))
                    requirer.Equal(filenameSavedTo, actualWriterFilenameSavedTo)
                } else {
                    requirer.Equal(contentsSaved, string(actualUwContentsSaved))
                    requirer.Equal(filenameSavedTo, actualUwFilenameSavedTo)
                }
            }
        })
    }
}

func TestFileSaver_Load(t *testing.T) {
    testCases := []struct {
        name      string
        hasReader bool
        hasError  bool
    }{
        {
            name: "no reader, no error",
        },
        {
            name:      "has reader, no error",
            hasReader: true,
        },
        {
            name:     "no reader, error",
            hasError: true,
        },
        {
            name:      "has reader, error",
            hasReader: true,
            hasError:  true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            requirer := require.New(t)

            const contentsLoaded = "this is what was loaded"

            var reader FileLoaderReader
            var actualReaderFilenameLoadedFrom string
            if tc.hasReader {
                reader = func(filePath string) ([]byte, error) {
                    if tc.hasError {
                        return nil, errors.New(tc.name)
                    }
                    actualReaderFilenameLoadedFrom = filePath
                    return []byte(contentsLoaded), nil
                }
            }

            var ur underReader
            var actualUrFilenameLoadedFrom string
            ur = func(filePath string) ([]byte, error) {
                if tc.hasError {
                    return nil, errors.New(tc.name)
                }
                actualUrFilenameLoadedFrom = filePath
                return []byte(contentsLoaded), nil
            }

            loader := &FileLoader{Reader: reader}
            var err error
            var actualContentsLoaded []byte
            const filenameLoadedFrom = "this is the load filename"
            actualContentsLoaded, err = loader.load(ur, filenameLoadedFrom)
            if tc.hasError {
                requirer.Error(err)
                requirer.Equal(tc.name, err.Error())
            } else {
                requirer.NoError(err)
                if tc.hasReader {
                    requirer.Equal(contentsLoaded, string(actualContentsLoaded))
                    requirer.Equal(filenameLoadedFrom, actualReaderFilenameLoadedFrom)
                } else {
                    requirer.Equal(contentsLoaded, string(actualContentsLoaded))
                    requirer.Equal(filenameLoadedFrom, actualUrFilenameLoadedFrom)
                }
            }
        })
    }
}
