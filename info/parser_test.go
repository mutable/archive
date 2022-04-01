package info

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNarInfo(t *testing.T) {
	narinfoContents := `
StorePath: /nix/store/qbwimh3fpbnz0g875i1jn6vihvg1xiik-kubectl-1.18.8
URL: nar/1awci4i0ddg32q255kbyzxi55wfhm9964wpbsz2fg8rbag46h7p8.nar.xz
Compression: xz
FileHash: sha256:1awci4i0ddg32q255kbyzxi55wfhm9964wpbsz2fg8rbag46h7p8
FileSize: 8780328
NarHash: sha256:1q4553kc9i03mlav8qpkvigxs9k2xnqcv3kq9cd8wrxbrrwicrsx
NarSize: 44292552
References: gr6rjrscqn0ldhrkq937l5b0qvx5adi1-tzdata-2019c ifnmhjrvk3f0hbz3f25s3izlb9yk8x0f-iana-etc-20200729 r2wvgnr54vmwnjvzyqdixv8xbn362jgh-mailcap-2.1.48
Sig: cache.nixos.org-1:rH4wxlNRbTbViQon40C15og5zlcFEphwoF26IQGHi2QCwVYyaLj6LOag+MeWcZ65SWzy6PnOlXjriLNcxE0hAQ==
Sig: nix-cache-mutable-1:/UAhsnPTuNzbw9iRsPhh7M+yDB1uPPLSE6IMIn8D0nAC++tAOgKkwBoCqVs0Aqz+cmmepPFPR4e3JPBUKuJWAw==
`[1:]
	r := strings.NewReader(narinfoContents)

	ni, err := ParseNarInfo(r)

	if assert.NoError(t, err) {
		assert.Equal(t, "/nix/store/qbwimh3fpbnz0g875i1jn6vihvg1xiik-kubectl-1.18.8", ni.StorePath)
		assert.Equal(t, "nar/1awci4i0ddg32q255kbyzxi55wfhm9964wpbsz2fg8rbag46h7p8.nar.xz", ni.URL)
		assert.Equal(t, "xz", ni.Compression)
		assert.Equal(t, "sha256:1awci4i0ddg32q255kbyzxi55wfhm9964wpbsz2fg8rbag46h7p8", ni.FileHash)
		assert.Equal(t, uint64(8780328), ni.FileSize)
		assert.Equal(t, "sha256:1q4553kc9i03mlav8qpkvigxs9k2xnqcv3kq9cd8wrxbrrwicrsx", ni.NarHash)
		assert.Equal(t, uint64(44292552), ni.NarSize)
		assert.Equal(t, []string{"gr6rjrscqn0ldhrkq937l5b0qvx5adi1-tzdata-2019c", "ifnmhjrvk3f0hbz3f25s3izlb9yk8x0f-iana-etc-20200729", "r2wvgnr54vmwnjvzyqdixv8xbn362jgh-mailcap-2.1.48"}, ni.References)
		assert.Equal(t, []string{"cache.nixos.org-1:rH4wxlNRbTbViQon40C15og5zlcFEphwoF26IQGHi2QCwVYyaLj6LOag+MeWcZ65SWzy6PnOlXjriLNcxE0hAQ==", "nix-cache-mutable-1:/UAhsnPTuNzbw9iRsPhh7M+yDB1uPPLSE6IMIn8D0nAC++tAOgKkwBoCqVs0Aqz+cmmepPFPR4e3JPBUKuJWAw=="}, ni.Sig)
	}

	// close the circle, by testing NarInfo's String() method to return the same as our initial narinfo
	assert.Equal(t, ni.String(), narinfoContents)
}
