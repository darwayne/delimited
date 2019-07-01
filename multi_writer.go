package delimited

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// MultiWriter writes data to multiple delimited files based on a given set of column names
// all new files are written to a temp directory and can be merged into a single file via Merge
type MultiWriter struct {
	opts          []WriterOptsFn
	groupCols     []int
	tempDir       string
	currentWriter *Writer
	currentFile   *os.File
	currentGroup  string
}

// NewMultiWriter creates a new multi-writer
func NewMultiWriter(groupCols []int, opts ...WriterOptsFn) (*MultiWriter, error) {
	tempDir, err := ioutil.TempDir("", "multi-writer-files")
	if err != nil {
		return nil, err
	}

	return &MultiWriter{
		groupCols: groupCols,
		opts:      opts,
		tempDir:   tempDir,
	}, nil
}

func (w *MultiWriter) Write(record []string) error {
	if err := w.prepNextWriter(record); err != nil {
		return err
	}

	return w.currentWriter.Write(record)
}

func (w *MultiWriter) prepNextWriter(record []string) error {
	filename := w.getFileName(record)
	if filename != w.currentGroup || w.currentFile == nil || w.currentWriter == nil {
		if w.currentWriter != nil {
			w.currentWriter.Flush()
		}
		if w.currentFile != nil {
			if err := w.currentFile.Close(); err != nil {
				return err
			}
		}

		w.currentGroup = filename
		fileLoc := path.Join(w.tempDir, filename)
		var err error
		w.currentFile, err = os.OpenFile(fileLoc, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			return err
		}

		w.currentWriter = NewWriter(w.currentFile, w.opts...)
	}
	return nil
}

func (w MultiWriter) getFileName(record []string) string {
	var str strings.Builder

	for _, col := range w.groupCols {
		str.WriteString(record[col])
	}
	//str.WriteString(".txt")

	return str.String()
}

func (w *MultiWriter) Flush() {
	if w.currentWriter == nil {
		return
	}

	w.currentWriter.Flush()
}

func (w *MultiWriter) TempDir() string {
	return w.tempDir
}

// Merge writes data for each group into the writer
func (w *MultiWriter) Merge(writer io.Writer) error {
	// read all files in tempDir and write contents to writer
	files, err := ioutil.ReadDir(w.tempDir)
	if err != nil {
		return err
	}
	w.Flush()

	for _, file := range files {
		if err := w.mergeFileContents(file.Name(), writer); err != nil {
			return err
		}
	}

	return nil
}

func (w *MultiWriter) mergeFileContents(filename string, writer io.Writer) error {
	filePath := path.Join(w.tempDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("an error occurred closing:", filename, err)
		}
		if err := os.Remove(filePath); err != nil {
			fmt.Println("an error occurred removing:", filename, err)
		}
	}()

	if _, err := io.Copy(writer, file); err != nil {
		return err
	}

	return nil
}

func (w *MultiWriter) Close() error {
	w.Flush()
	if w.currentFile != nil {
		w.currentFile.Close()
	}

	return os.RemoveAll(w.tempDir)
}