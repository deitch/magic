package parser

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Operator byte

const (
	Equal Operator = iota
	NotEqual
	GreaterThan
	LessThan
	GreaterThanOrEqual
	LessThanOrEqual
	Any
)

type UnifiedReader interface {
	io.Reader
	io.ReaderAt
}
type Tester func(UnifiedReader, string) (bool, string, error)
type offsetReader func(io.ReaderAt) (int64, error)

type MagicTest struct {
	Test     Tester
	Message  string
	Children []MagicTest
}

func StringTest(offsetFunc offsetReader, compare string) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, len(compare))
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(compare) {
			return false, "", nil
		}
		isMatch := string(b) == compare
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Uint16Test(offsetFunc offsetReader, compare uint16, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 2)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		actual := endian.Uint16(b)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Int16Test(offsetFunc offsetReader, compare int16, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 2)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		interim := endian.Uint16(b)
		actual := int16(interim)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Uint32Test(offsetFunc offsetReader, compare uint32, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 4)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		actual := endian.Uint32(b)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Int32Test(offsetFunc offsetReader, compare int32, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 4)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		interim := endian.Uint32(b)
		actual := int32(interim)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Uint64Test(offsetFunc offsetReader, compare uint64, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 8)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		actual := endian.Uint64(b)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Int64Test(offsetFunc offsetReader, compare int64, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 8)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		interim := endian.Uint64(b)
		actual := int64(interim)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}
func Float32Test(offsetFunc offsetReader, compare float32, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 8)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		var actual float32
		buf := bytes.NewReader(b)
		if err := binary.Read(buf, endian, &actual); err != nil {
			return false, "", err
		}
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}
func Float64Test(offsetFunc offsetReader, compare float64, comparator Operator, endian binary.ByteOrder) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 8)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		var actual float64
		buf := bytes.NewReader(b)
		if err := binary.Read(buf, endian, &actual); err != nil {
			return false, "", err
		}

		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Uint8Test(offsetFunc offsetReader, compare uint8, comparator Operator) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 1)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		isMatch, err := compareNumbers(uint8(b[0]), compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Int8Test(offsetFunc offsetReader, compare int8, comparator Operator) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 1)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		isMatch, err := compareNumbers(int8(b[0]), compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func GuidTest(offsetFunc offsetReader, compare string, comparator Operator) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 16)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		// convert to a GUID
		g, err := uuid.FromBytes(b)
		if err != nil {
			return false, "", err
		}
		var isMatch bool
		switch comparator {
		case Equal:
			isMatch = strings.ToUpper(g.String()) == compare
		case Any:
			isMatch = true
		case NotEqual:
			isMatch = strings.ToUpper(g.String()) != compare
		default:
			return false, "", fmt.Errorf("invalid comparator: %v", comparator)
		}

		if !isMatch {
			return false, "", nil
		}
		return isMatch, messageParser(r, offset, pattern, nil), nil
	}
}

func Date32Test(offsetFunc offsetReader, compare uint32, comparator Operator, endian binary.ByteOrder, local bool) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 4)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		actual := endian.Uint32(b)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		converter := func(b []byte) string {
			sec := int64(actual)
			t := time.Unix(sec, 0)
			if !local {
				t = t.UTC()
			}
			return t.Format(time.UnixDate)
		}
		return isMatch, messageParser(r, offset, pattern, converter), nil
	}
}

func Date64Test(offsetFunc offsetReader, compare uint64, comparator Operator, endian binary.ByteOrder, local bool) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		b := make([]byte, 8)
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		n, err := r.ReadAt(b, offset)
		if err != nil {
			return false, "", err
		}
		if n != len(b) {
			return false, "", nil
		}
		actual := endian.Uint64(b)
		var isMatch bool
		isMatch, err = compareNumbers(actual, compare, comparator)
		if err != nil {
			return false, "", err
		}
		if !isMatch {
			return false, "", nil
		}
		converter := func(b []byte) string {
			sec := int64(actual)
			t := time.Unix(sec, 0)
			if !local {
				t = t.UTC()
			}
			return t.Format(time.UnixDate)
		}
		return isMatch, messageParser(r, offset, pattern, converter), nil
	}
}

func DefaultTest(offsetFunc offsetReader, compare string) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		return true, messageParser(r, offset, pattern, nil), nil
	}
}

func RegexTest(offsetFunc offsetReader, compare string, extender string) Tester {
	return func(r UnifiedReader, pattern string) (bool, string, error) {
		offset, err := offsetFunc(r)
		if err != nil {
			return false, "", err
		}
		// was there a read extender?
		if extender != "" {
			var (
				count           int64
				countStr        string
				caseInsensitive bool
				startOfMatch    bool
				useLines        bool
			)
			for _, c := range extender {
				switch c {
				case 'c':
					caseInsensitive = true
				case 'l':
					useLines = true
				case 's':
					startOfMatch = true
				default:
					countStr += string(c)
				}
			}
			if len(countStr) != 0 {
				count, err = strconv.ParseInt(countStr, 10, 64)
				if err != nil {
					return false, "", err
				}
			}
			if caseInsensitive {
				compare = "(?i)" + compare
			}
			re, err := regexp.Compile(compare)
			if err != nil {
				return false, "", err
			}
			var runeReader io.RuneReader = bufio.NewReader(r)
			if count > 0 {
				var b []byte
				if !useLines {
					b = make([]byte, count)
					n, err := r.ReadAt(b, offset)
					if err != nil {
						return false, "", err
					}
					if n != len(b) {
						return false, "", fmt.Errorf("unsufficient bytes available to read at pos %d", offset)
					}
				} else {
					// it said to use lines, so we need to read until we hit carriage returns
					b = make([]byte, 0)
					scanner := bufio.NewScanner(r)
					var lines int
					for scanner.Scan() {
						b = append(b, scanner.Bytes()...)
						lines++
						if lines >= int(count) {
							break
						}
					}
				}
				runeReader = bytes.NewReader(b)
			}
			loc := re.FindReaderIndex(runeReader)
			if len(loc) < 2 || loc[0] < 0 || loc[1] < 0 {
				return false, "", nil
			}
			if startOfMatch {
				offset = offset + int64(loc[0])
			} else {
				offset = offset + int64(loc[1])
			}
		}
		return true, messageParser(r, offset, pattern, nil), nil
	}
}
