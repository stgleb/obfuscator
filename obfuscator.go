package obfuscator

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

const LowWord = 2
const HiWord = 5

type Script struct {
	Buf []byte
}

// TODO(stgleb): Use better pattern matching algorithm
// Creates a dictionary of the given word length and counts occurrences.
func dict(buf []byte, wlen int) map[string]int {

	dict := map[string]int{}

	l := len(buf)

	for i := 0; i+wlen < l; i++ {
		s := string(buf[i : i+wlen])
		if _, ok := dict[s]; ok == false {
			dict[s] = 0
		} else {
			dict[s] = dict[s] + 1
		}
	}

	return dict
}

// Returns the first unused byte, we need it to store a chunk.
func getFirstUnusedByte(used map[byte]bool) (byte, error) {
	var c byte
	for i := byte(126); i > byte(0); i-- {
		switch i {
		case
			// Non-printable bytes.
			byte(10),
			byte(13),
			byte(34),
			byte(39),
			byte(92):
		default:
			c = i
			if !used[c] {
				used[c] = true
				return c, nil
			}
		}
	}
	return byte(0), errors.New("No more unused bytes\n")
}

// Compresses a buffer into a chunk and chunk keys.
func Compress(buf []byte) ([]byte, []byte, error) {
	used := make(map[byte]bool)
	keys := []byte{}
	mv := -1

	for {
		t, err := getFirstUnusedByte(used)

		if err != nil {
			break
		}

		mk := ""

		for i := LowWord; i < HiWord; i++ {
			d := dict(buf, i)

			for k, v := range d {
				if v > 1 {
					// Calculating length.
					tlen := len(buf) - v*len(k) + v
					slen := tlen + len(k) + len(keys) + 2

					if slen < mv || mv < 0 {
						if tlen < len(buf) {
							mk = k
							mv = slen
						}
					}
				}
			}
		}

		if mk != "" {
			buf = bytes.Replace(buf, []byte(mk), []byte{t}, -1)
			buf = append(buf, t)
			buf = append(buf, []byte(mk)...)
			keys = append([]byte{t}, keys...)
		} else {
			break
		}
	}

	return buf, keys, nil
}

// Uses buf and key to create the javascript decompressor.
func Pack(buf []byte, keys []byte) string {
	return fmt.Sprintf("for(s='%s',i=0;j='%s'[i++];)with(s.split(j))s=join(pop());eval(s)", string(buf), string(keys))
}

func Minify(data []byte) ([]byte, error) {
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)

	b, err := m.Bytes("text/javascript", data)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Obfuscate(data []byte) (string, error) {
	buf, err := Minify(data)
	t, k, err := Compress(buf)

	if err != nil {
		return "", err
	}

	return Pack(t, k), nil
}
