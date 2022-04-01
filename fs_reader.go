package archive

import (
	"fmt"
	"os"
	pathapi "path"
	"sort"
)

func CopyPath(w Writer, path string) error {
	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	return copyPathFI(w, path, fi)
}

func copyPathFI(w Writer, path string, fi os.FileInfo) error {
	mode := fi.Mode()
	switch typ := mode & os.ModeType; typ {
	default:
		return fmt.Errorf("nix/archive: unsupported file type %s", typ)
	case os.ModeSymlink:
		target, err := os.Readlink(path)
		if err != nil {
			return err
		}
		return w.Symlink(target)
	case 0:
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		executable := mode&0100 != 0
		bytes := uint64(fi.Size())
		reader := &contentReader{f, bytes}
		return w.File(executable, reader)
	case os.ModeDir:
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		files, err := f.Readdir(0)
		if err != nil {
			return err
		}
		directoryWriter, err := w.Directory()
		if err != nil {
			return err
		}
		fileInfos(files).Sort()
		for _, fi := range files {
			name := fi.Name()
			w := directoryWriter.Entry(name)
			path := pathapi.Join(path, name)
			if err := copyPathFI(w, path, fi); err != nil {
				return err
			}
		}
		return directoryWriter.Close()
	}
}

type fileInfos []os.FileInfo

func (f fileInfos) Len() int           { return len(f) }
func (f fileInfos) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f fileInfos) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f fileInfos) Sort()              { sort.Sort(f) }
