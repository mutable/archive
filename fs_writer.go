package archive

import (
	"io"
	"os"
	"path"
)

type fsWriter struct {
	path string
}

func NewFSWriter(path string) Writer {
	return &fsWriter{path}
}

func (w *fsWriter) Symlink(target string) error {
	return os.Symlink(target, w.path)
}

func (w *fsWriter) File(executable bool, contents ContentReader) error {
	mode := os.FileMode(0644)
	if executable {
		mode |= 0111
	}

	f, err := os.OpenFile(w.path, os.O_RDWR|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, contents); err != nil {
		return err
	}
	return f.Close()
}

func (w *fsWriter) Directory() (DirectoryWriter, error) {
	if err := os.Mkdir(w.path, 0755); err != nil {
		return nil, err
	}

	return &fsDirectoryWriter{w.path}, nil
}

type fsDirectoryWriter struct {
	path string
}

func (w *fsDirectoryWriter) Entry(name string) Writer {
	return NewFSWriter(path.Join(w.path, name))
}

func (w *fsDirectoryWriter) Close() error {
	return nil
}
