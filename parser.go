package archive

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/mutable/wire"
)

type Writer interface {
	File(executable bool, contents ContentReader) error
	Directory() (DirectoryWriter, error)
	Symlink(target string) error
}

type DirectoryWriter interface {
	Entry(name string) Writer
	Close() error
}

const archiveVersion1 = "nix-archive-1"

const (
	// for small tokens,
	// we use this to limit how large an invalid token we'll read
	tokenLenMax = 32
	// maximum length for a single path element
	// NAME_MAX is 255 on Linux
	nameLenMax = 255
	// maximum length for a relative path
	// PATH_MAX is 4096 on Linux, but that includes a null byte
	pathLenMax = 4096 - 1
)

func ParseDump(w Writer, r io.Reader) error {
	version, err := wire.ReadString(r, tokenLenMax)
	if err != nil {
		return err
	}
	if version != archiveVersion1 {
		return errors.New("nix/archive: input doesn't look like a Nix archive")
	}
	return parse(w, r)
}

func parse(w Writer, r io.Reader) error {
	if s, err := wire.ReadString(r, tokenLenMax); err != nil {
		return err
	} else if s != "(" {
		return errors.New("nix/archive: expected open tag")
	}

	const (
		tpUnknown = iota
		tpRegular
		tpDirectory
		tpSymlink
	)
	var (
		typ             = tpUnknown
		executable      = false
		directoryWriter DirectoryWriter
	)

	for {
		s, err := wire.ReadString(r, tokenLenMax)
		if err != nil {
			return err
		}
		switch s {
		default:
			return fmt.Errorf("nix/archive: unknown field: %q", s)
		case ")":
			if directoryWriter != nil {
				return directoryWriter.Close()
			}
			return nil
		case "type":
			if typ != tpUnknown {
				return fmt.Errorf("nix/archive: multiple type fields")
			}
			t, err := wire.ReadString(r, tokenLenMax)
			if err != nil {
				return err
			}
			switch t {
			default:
				return fmt.Errorf("nix/archive: unknown type: %q", t)
			case "regular":
				typ = tpRegular
			case "directory":
				typ = tpDirectory
			case "symlink":
				typ = tpSymlink
			}
		case "executable":
			if typ != tpRegular {
				return fmt.Errorf("nix/archive: unexpected field: %q", s)
			}
			if _, err := wire.ReadString(r, tokenLenMax); err != nil {
				return err
			}
			executable = true
		case "contents":
			if typ != tpRegular {
				return fmt.Errorf("nix/archive: unexpected field: %q", s)
			}
			contents, err := parseContents(r)
			if err != nil {
				return err
			}
			n := contents.n
			if err := w.File(executable, contents); err != nil {
				return err
			}
			// ensure we've consumed all file contents before continuing
			if _, err := io.Copy(ioutil.Discard, contents); err != nil {
				return err
			}
			if err := wire.ReadPadding(r, n); err != nil {
				return err
			}
		case "target":
			if typ != tpSymlink {
				return fmt.Errorf("nix/archive: unexpected field: %q", s)
			}
			target, err := wire.ReadString(r, pathLenMax)
			if err != nil {
				return err
			}
			if err := w.Symlink(target); err != nil {
				return err
			}
		case "entry":
			if directoryWriter == nil {
				if directoryWriter, err = w.Directory(); err != nil {
					return err
				}
			}
			if err := parseDirectory(directoryWriter, r); err != nil {
				return err
			}
		}
	}
}

func parseContents(r io.Reader) (*contentReader, error) {
	n, err := wire.ReadUint64(r)
	if err != nil {
		return nil, err
	}
	return &contentReader{r, n}, nil
}

func parseDirectory(w DirectoryWriter, r io.Reader) error {
	if s, err := wire.ReadString(r, tokenLenMax); err != nil {
		return err
	} else if s != "(" {
		return errors.New("nix/archive: expected open tag")
	}

	var entryWriter Writer
	for {
		s, err := wire.ReadString(r, tokenLenMax)
		if err != nil {
			return err
		}
		switch s {
		default:
			return fmt.Errorf("nix/archive: unknown field: %q", s)
		case ")":
			return nil
		case "name":
			name, err := wire.ReadString(r, nameLenMax)
			if err != nil {
				return err
			}
			entryWriter = w.Entry(name)
		case "node":
			if err := parse(entryWriter, r); err != nil {
				return err
			}
		}
	}
}
