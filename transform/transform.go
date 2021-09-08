package transform

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
)

type kind string

const (
	city          kind = "MUNICCSV"
	cnae               = "CNAECSV"
	company            = "EMPRECSV"
	country            = "PAISCSV"
	facility           = "ESTABELE"
	motive             = "MOTICSV"
	nature             = "NATJUCSV"
	partner            = "SOCIOCSV"
	qualification      = "QUALSCSV"
	simple             = "SIMPLES"
)

const CompressionAlgorithms = "xz, gz"

type dataset struct {
	kind        kind
	dir         string
	compression string
	done        bool
	fileHandler *os.File
	ioWriter    io.WriteCloser
	csvWriter   *csv.Writer
}

func newDataset(k kind, d, c string) *dataset {
	a := dataset{k, d, c, false, nil, nil, nil}
	return &a
}

func (a *dataset) Close() {
	if a.csvWriter != nil {
		a.csvWriter.Flush()
	}

	if a.fileHandler != nil {
		a.fileHandler.Close()
	}

	if a.ioWriter != nil {
		a.ioWriter.Close()
	}

	a.done = true
}

func (a *dataset) files() ([]string, error) {
	var o []string

	ls, err := os.ReadDir(a.dir)
	if err != nil {
		return []string{}, err
	}

	for _, f := range ls {
		if !f.IsDir() {
			o = append(o, filepath.Join(a.dir, f.Name()))
		}
	}
	return filesFor(a, o), nil
}

func (a *dataset) Writer(i io.WriteCloser) (io.WriteCloser, error) {
	switch a.compression {
	case "xz":
		return xz.NewWriter(i)
	case "gz":
		return gzip.NewWriter(i), nil
	default:
		return i, nil
	}
}

func validateCompressionAlgorithm(c string) error {
	if c == "" {
		return nil
	}

	for _, o := range strings.Split(CompressionAlgorithms, ",") {
		if c == strings.TrimSpace(o) {
			return nil
		}
	}

	return fmt.Errorf(
		"Unknown compression algorithm %s, options are: %s",
		c,
		CompressionAlgorithms,
	)
}

// Transform unzips the downloaded files and merge them into CSV files.
func Transform(dir string, compression string) error {
	if err := validateCompressionAlgorithm(compression); err != nil {
		return err
	}

	var as []*dataset
	for _, k := range []kind{
		city,
		cnae,
		company,
		country,
		facility,
		motive,
		nature,
		partner,
		qualification,
		simple,
	} {
		as = append(as, newDataset(k, dir, compression))
	}

	c := make(chan error)
	for _, a := range as {
		go a.writeCsv(c)
	}

	q := make(chan struct{})
	go startSpinners(q, as)

	for range as {
		err := <-c
		if err != nil {
			return err
		}
	}

	q <- struct{}{} // ask the spinner to wrap up
	<-q             // wait for the spinner to wrap up
	close(q)

	return nil
}
