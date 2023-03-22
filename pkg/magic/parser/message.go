package parser

import (
	"fmt"
	"io"
)

func WithMessage(msg string) message {
	return func(io.ReaderAt, int64) string {
		return msg
	}
}

func WithEmptyMessage() message {
	return WithMessage("")
}

func WithString(prefix, suffix string) message {
	return func(r io.ReaderAt, pos int64) string {
		// read until we hit a null
		b := make([]byte, 1)
		var s string
		for {
			n, err := r.ReadAt(b, pos)
			if err != nil {
				return fmt.Sprintf("%s%s%s", prefix, s, suffix)
			}
			if n != len(b) {
				return fmt.Sprintf("%s%s%s", prefix, s, suffix)
			}
			if b[0] == 0 {
				return fmt.Sprintf("%s%s%s", prefix, s, suffix)
			}
			s += string(b[0])
			pos++
		}
	}
}
