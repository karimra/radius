package dictionarygen

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"layeh.com/radius/dictionary"
)

func TestTestData(t *testing.T) {
	tbl := []struct {
		Name          string
		InitParser    func(*dictionary.Parser)
		InitGenerator func(*Generator)
		Err           string
	}{
		{
			Name: "identical-attributes",
			InitParser: func(p *dictionary.Parser) {
				p.IgnoreIdenticalAttributes = true
			},
		},
		{
			Name: "identifier-collision",
			Err:  "conflicting identifier between First_Name (200) and First-Name (201)",
		},
	}

	for _, tt := range tbl {
		t.Run(tt.Name, func(t *testing.T) {
			parser := &dictionary.Parser{
				Opener: &dictionary.FileSystemOpener{},
			}
			if tt.InitParser != nil {
				tt.InitParser(parser)
			}

			dictFile := filepath.Join("testdata", tt.Name+".dictionary")
			dict, err := parser.ParseFile(dictFile)
			if err != nil {
				t.Fatalf("could not parse file: %s", err)
			}

			generator := &Generator{
				Package: "main",
			}
			if tt.InitGenerator != nil {
				tt.InitGenerator(generator)
			}

			generatedCode, err := generator.Generate(dict)
			if err != nil {
				if tt.Err != "" {
					if !strings.Contains(err.Error(), tt.Err) {
						t.Fatalf("got generate error %v; expected %v", err, tt.Err)
					}
					return
				}
				t.Fatalf("could not generate dictionary code: %s", err)
			}

			generatedFile := filepath.Join("testdata", tt.Name+".generated")
			if err := ioutil.WriteFile(generatedFile, generatedCode, 0644); err != nil {
				t.Fatalf("could not write generated file: %s", err)
			}

			expectedFile := filepath.Join("testdata", tt.Name+".expected")
			expectedCode, err := ioutil.ReadFile(expectedFile)
			if err != nil {
				t.Fatalf("could not read expected output: %s", err)
			}

			if !bytes.Equal(generatedCode, expectedCode) {
				t.Fatal("generated code does not equal expected")
			}

			os.Remove(generatedFile)
		})
	}
}
