package archive

import "io"

type ContentReader interface {
	// Embed a regular io.Reader
	io.Reader
	// Returns remaining byte count
	Bytes() uint64
}

type contentReader struct {
	r io.Reader
	n uint64
}

func (r *contentReader) Read(p []byte) (n int, err error) {
	if r.n == 0 {
		return 0, io.EOF
	}
	if uint64(len(p)) > r.n {
		p = p[:r.n]
	}
	n, err = r.r.Read(p)
	r.n -= uint64(n)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

func (r *contentReader) Bytes() uint64 {
	return r.n
}
