package parser

import (
	"encoding/binary"
	"fmt"
	"io"
)

func WithMessage(msg string) message {
	return func(io.ReaderAt, int64) string {
		return msg
	}
}

func WithMessageEmpty() message {
	return WithMessage("")
}

func WithMessagePattern(pattern string) message {
	return func(r io.ReaderAt, pos int64) string {
		var (
			data []any
		)
		for i := 0; i < len(pattern); i++ {
			c := pattern[i]
			if c != '%' {
				continue
			}
			var b []byte
			// it indicates a pattern
			i++
			c = pattern[i]
			// any additional flags, width, precision or modifiers?
			switch c {
			case '#', '+', '-', ' ', '0':
				i++
				c = pattern[i]
			}
			switch c {
			case 's':
				// string, we hit a null
				b2 := make([]byte, 1)
				for {
					n, err := r.ReadAt(b2, pos)
					if err != nil || n != len(b2) || b2[0] == 0 {
						break
					}
					b = append(b, b2[0])
					pos++
				}
				data = append(data, string(b))
			case '%':
				// '%%' is a literal '%'
			case 'd', 'i':
				// signed decimal integer
				b = make([]byte, 4)
				n, err := r.ReadAt(b, pos)
				if err != nil || n != len(b) {
					break
				}
				data = append(data, int32(binary.LittleEndian.Uint32(b)))
			case 'u', 'x', 'X':
				// signed decimal integer
				b = make([]byte, 4)
				n, err := r.ReadAt(b, pos)
				if err != nil || n != len(b) {
					break
				}
				data = append(data, binary.LittleEndian.Uint32(b))
			}
		}
		return fmt.Sprintf(pattern, data...)
	}
}
