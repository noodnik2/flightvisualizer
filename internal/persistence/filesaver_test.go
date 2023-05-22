package persistence

//
//import (
//    "testing"
//
//    "github.com/stretchr/testify/require"
//)
//
//func TestHrefSaver_Save(t *testing.T) {
//    testCases := []struct {
//        name           string
//        href           string
//        expectedFnPath string
//    }{
//        {
//            name:           "complete with 2 query pairs",
//            href:           "http://one.two:8901/myUri/subform/one?option=very_long&also=this",
//            expectedFnPath: "http_one.two_8901_myUri_subform_one_option_very_long_also_this",
//        },
//        {
//            name:           "localhost case with port no query trailing space",
//            href:           "http://localhost:2324 ",
//            expectedFnPath: "http_localhost_2324",
//        },
//        {
//            name:           "uri path only with query",
//            href:           "/_test?arg=one ",
//            expectedFnPath: "test_arg_one",
//        },
//        {
//            name:           "localhost case with single query",
//            href:           "http://localhost/test?arg=one",
//            expectedFnPath: "http_localhost_test_arg_one",
//        },
//        {
//            name:           "quotes and spaces in query",
//            href:           `https://foreign.service.org:2324/one/two?spice=' nice '&that="one  two"`,
//            expectedFnPath: "https_foreign.service.org_2324_one_two_spice_+nice+_that_one++two",
//        },
//    }
//
//    for _, tc := range testCases {
//        t.Run(tc.name, func(t *testing.T) {
//            requirer := require.New(t)
//            hs := HrefSaver{FileSaver{Writer: func(filePath string, contents []byte) error {
//                requirer.Equal(tc.expectedFnPath, filePath)
//                return nil
//            }}}
//            save, err := hs.Save(tc.href, []byte{})
//            requirer.NoError(err)
//            requirer.Equal(tc.expectedFnPath, save)
//        })
//    }
//}
