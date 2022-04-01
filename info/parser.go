package info

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseNarInfo reads .narinfo file contents
// and returns a NarInfo struct with the parsed data
func ParseNarInfo(r io.Reader) (*NarInfo, error) {
	narInfo := &NarInfo{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var err error

		line := scanner.Text()
		parts := strings.Split(line, ": ")

		if len(parts) != 2 {
			return nil, fmt.Errorf("Unable to split line %s", line)
		}

		k := parts[0]
		v := parts[1]

		switch k {
		case "StorePath":
			narInfo.StorePath = v
		case "URL":
			narInfo.URL = v
		case "Compression":
			narInfo.Compression = v
		case "FileHash":
			narInfo.FileHash = v
		case "FileSize":
			narInfo.FileSize, err = strconv.ParseUint(v, 10, 64)
		case "NarHash":
			narInfo.NarHash = v
		case "NarSize":
			narInfo.NarSize, err = strconv.ParseUint(v, 10, 64)
		case "References":
			if v == "" {
				continue
			}
			narInfo.References = append(narInfo.References, strings.Split(v, " ")...)
		case "CA":
			narInfo.CA = v
		case "Deriver":
			narInfo.Deriver = v
		case "Sig":
			narInfo.Sig = append(narInfo.Sig, v)
		}

		if err != nil {
			return nil, fmt.Errorf("Unable to parse %s", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return narInfo, nil
}
