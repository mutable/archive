package info

import (
	"bytes"
	"fmt"
)

// NarInfo contains a parsed narinfo file
type NarInfo struct {
	StorePath   string
	URL         string // relative path
	Compression string
	FileHash    string
	FileSize    uint64
	NarHash     string
	NarSize     uint64
	References  []string
	Deriver     string
	Sig         []string
	CA          string
}

func (n *NarInfo) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "StorePath: %s\n", n.StorePath)
	fmt.Fprintf(&buf, "URL: %s\n", n.URL)
	fmt.Fprintf(&buf, "Compression: %s\n", n.Compression)
	fmt.Fprintf(&buf, "FileHash: %s\n", n.FileHash)
	fmt.Fprintf(&buf, "FileSize: %d\n", n.FileSize)
	fmt.Fprintf(&buf, "NarHash: %s\n", n.NarHash)
	fmt.Fprintf(&buf, "NarSize: %d\n", n.NarSize)

	buf.WriteString("References:")
	if len(n.References) == 0 {
		buf.WriteByte(' ')
	} else {
		for _, r := range n.References {
			buf.WriteByte(' ')
			buf.WriteString(r)
		}
	}
	buf.WriteByte('\n')

	if n.Deriver != "" {
		fmt.Fprintf(&buf, "Deriver: %s\n", n.Deriver)
	}

	for _, s := range n.Sig {
		fmt.Fprintf(&buf, "Sig: %s\n", s)
	}

	if n.CA != "" {
		fmt.Fprintf(&buf, "CA: %s\n", n.CA)
	}

	return buf.String()
}
