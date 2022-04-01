package archive

import (
	"io"

	"github.com/mutable/wire"
)

func WriteDump(w io.Writer) Writer {
	if err := wire.WriteString(w, archiveVersion1); err != nil {
		return &errWriter{err}
	}
	return &writer{writer: w}
}

type writer struct {
	writer io.Writer
	child  bool
}

func (w *writer) Symlink(target string) error {
	attrs := []string{"(", "type", "symlink", "target", target, ")"}
	for _, s := range attrs {
		if err := wire.WriteString(w.writer, s); err != nil {
			return err
		}
	}
	return w.close()
}

func (w *writer) File(executable bool, contents ContentReader) error {
	attrs := []string{"(", "type", "regular"}
	if executable {
		attrs = append(attrs, "executable", "")
	}
	attrs = append(attrs, "contents")
	for _, s := range attrs {
		if err := wire.WriteString(w.writer, s); err != nil {
			return err
		}
	}
	n := contents.Bytes()
	if err := wire.WriteUint64(w.writer, n); err != nil {
		return err
	}
	if _, err := io.Copy(w.writer, contents); err != nil {
		return err
	}
	if err := wire.WritePadding(w.writer, n); err != nil {
		return err
	}
	if err := wire.WriteString(w.writer, ")"); err != nil {
		return err
	}
	return w.close()
}

func (w *writer) Directory() (DirectoryWriter, error) {
	for _, s := range []string{"(", "type", "directory"} {
		if err := wire.WriteString(w.writer, s); err != nil {
			return nil, err
		}
	}
	return (*directoryWriter)(w), nil
}

type directoryWriter writer

func (w *directoryWriter) Entry(name string) Writer {
	for _, s := range []string{"entry", "(", "name", name, "node"} {
		if err := wire.WriteString(w.writer, s); err != nil {
			return &errWriter{err}
		}
	}
	return &writer{writer: w.writer, child: true}
}

func (w *directoryWriter) Close() error {
	if err := wire.WriteString(w.writer, ")"); err != nil {
		return err
	}
	return (*writer)(w).close()
}

func (w *writer) close() error {
	if w.child {
		if err := wire.WriteString(w.writer, ")"); err != nil {
			return err
		}
	}
	return nil
}

type errWriter struct{ error }

func (e *errWriter) Directory() (DirectoryWriter, error) {
	return nil, e.error
}
func (e *errWriter) File(bool, ContentReader) error {
	return e.error
}
func (e *errWriter) Symlink(string) error {
	return e.error
}
