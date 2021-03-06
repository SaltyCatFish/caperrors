package product

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	file "github.com/SaltyCatFish/caperrors/pkg/file"
)

type FileResult struct {
	fileBase        string
	filePath        string
	logFilePath     string
	expectedID      string
	expectedMessage string
}

func (r FileResult) paths() (results []string) {
	results = append(results, filepath.Join(r.filePath, r.fileBase))
	results = append(results, r.logFilePath)
	return
}

var fileResults = []FileResult{
	{
		filePath:        "../../testdata/files/dummy_test.xml",
		logFilePath:     "../../testdata/files/CapHandler.log",
		expectedID:      "wea_003_imminent_threat_test12020624141331",
		expectedMessage: "",
	},
	{
		filePath:        "../../testdata/files/09.Jul.2020.15.50.24.3219ceafa85-ef99-41c9-a0ae-c1ef620afb42.xml",
		logFilePath:     "../../testdata/files/CapHandler.log",
		expectedID:      "Chan_Amber_test1202069114831",
		expectedMessage: "[INFO ] - Unable to deserialize <alert> block #1 in /var/lib/caphandler/input/ipaws/09.Jul.2020.15.50.24.3219ceafa85-ef99-41c9-a0ae-c1ef620afb42.xml into an Alert object (gov.noaa.nws.cap.exception.DecoderException: <alert> block #1 in /var/lib/caphandler/input/ipaws/09.Jul.2020.15.50.24.3219ceafa85-ef99-41c9-a0ae-c1ef620afb42.xml has no <Signature><KeyInfo><X509Data><X509SubjectName> element)",
	},
	{
		filePath:        "../../testdata/files/ESFPUB.09164211.443799.txt",
		logFilePath:     "../../testdata/files/CapHandler.log",
		expectedID:      "47a9525a06cce867a5659855ef2600fe51d4da78",
		expectedMessage: "[INFO ] - Unable to decode the text product in /var/lib/caphandler/input/awips/wmo/ESFPUB.09164211.443799.txt (gov.noaa.nws.cap.exception.DecoderException: Probabilistic Hydrologic Outlook products aren't supported)",
	},
}

func TestMain(m *testing.M) {
	setup("../../testdata/files")
	code := m.Run()
	teardown("../../testdata/files")
	os.Exit(code)
}

func TestIDAndMessage(t *testing.T) {
	for _, test := range fileResults {
		f, err := os.Open(test.filePath)
		defer f.Close()
		if err != nil {
			t.Fatalf("Could not open %v", test.filePath)
		}
		fileinfo, err := f.Stat()
		if err != nil {
			t.Fatalf("Could not stat %v", test.filePath)
		}
		x := NewFile(file.NewFile(test.filePath, fileinfo))
		if id, _ := x.ID(test.logFilePath); id != test.expectedID {
			t.Fatalf("%v does not equal %v", id, test.expectedID)
		}
		if message, _ := x.ErrorMessage(test.logFilePath); !strings.Contains(message, test.expectedMessage) {
			t.Fatalf("%v does not equal %v", message, test.expectedMessage)
		}
	}
}

func setup(path string) {
	results := []string{}
	filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln(err)
			return nil
		}
		if !file.IsDir() && filepath.Ext(path) == ".gz" {
			results = append(results, path)
		}
		return nil
	})
	gunzip(results)
}

func gunzip(paths []string) {
	for _, path := range paths {
		reader, err := ioutil.ReadFile(path)
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		gr, err := gzip.NewReader(bytes.NewBuffer(reader))
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		defer gr.Close()
		data, err := ioutil.ReadAll(gr)
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
		err = ioutil.WriteFile(strings.TrimSuffix(path, filepath.Ext(path)), data, 777)
		if err != nil {
			panic(fmt.Errorf(err.Error()))
		}
	}
	return
}

func teardown(path string) {
	filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln(err)
			return nil
		}
		if !file.IsDir() && filepath.Ext(path) != ".gz" {
			os.RemoveAll(path)
		}
		return nil
	})
}
