package data

import (
	"os"
	"path"
	"physicsGUI/pkg/io"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestImport(t *testing.T) {
	filePath := path.Join("..", "..", "testdata", "syntheticdataset.dat")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Error(err)
	}

	data, err := Parse(fileContent)
	if err != nil {
		t.Error(err)
	}

	spew.Dump(data)

	// TODO: add test handling instead of printing
}

func TestImportAlt(t *testing.T) {
	filePath := path.Join("..", "..", "testdata", "dataset.dat")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Error(err)
	}

	data, err := Parse(fileContent)
	if err != nil {
		t.Error(err)
	}

	spew.Dump(data)

	// TODO: add test handling instead of printing
}

func TestOldImport(t *testing.T) {
	filePath := path.Join("..", "..", "testdata", "syntheticdataset.dat")

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Error(err)
	}

	parser := io.OldParser()
	data, err := parser.tryParse(fileContent)
	if err != nil {
		t.Error(err)
	}

	spew.Dump(data)

	// TODO: add test handling instead of printing
}
