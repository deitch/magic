package compiler

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"

	"github.com/deitch/magic/pkg/magic/parser"
)

const matchAny = "x"

var nativeEndian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}

func ParseSource(r io.Reader) ([]parser.MagicTest, error) {
	var tests []parser.MagicTest
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// strip leading whitespace
		line = strings.TrimLeft(line, " \t")
		// skip empty lines and comments
		if line == "" || line[0] == '#' {
			continue
		}
		// the line will have 3 or 4 parts
		fields := strings.Fields(line)
		if strings.HasPrefix(fields[0], ">") {
			// this is a child of the parent
		} else {
			// this is a new parent
			offset, err := strconv.ParseInt(fields[0], 0, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid line '%s': %w", line, err)
			}
			readType := fields[1]
			testValue := fields[2]
			var readExtender string
			var message string
			if len(fields) > 3 {
				message = strings.Join(fields[3:], " ")
			}
			var unsigned bool
			if readType[0] == 'u' {
				readType = readType[1:]
				unsigned = true
			}
			endian := nativeEndian
			switch {
			case strings.HasPrefix(readType, "le"):
				endian = binary.LittleEndian
			case strings.HasPrefix(readType, "be"):
				endian = binary.BigEndian
			case strings.HasPrefix(readType, "me"):
				endian = parser.Pdp11ByteOrder
			}
			comparator := parser.Equal
			switch {
			case len(testValue) == 0:
				return nil, fmt.Errorf("invalid test value '%s': %w", testValue, err)
			case testValue[0] == 'x':
				comparator = parser.Any
			case testValue[0] == '>':
				comparator = parser.GreaterThan
				testValue = testValue[1:]
				if testValue[0] == '=' {
					comparator = parser.GreaterThanOrEqual
					testValue = testValue[1:]
				}
			case testValue[0] == '<':
				comparator = parser.LessThan
				testValue = testValue[1:]
				if testValue[0] == '=' {
					comparator = parser.LessThanOrEqual
					testValue = testValue[1:]
				}
			case testValue[0] == '=':
				comparator = parser.Equal
				testValue = testValue[1:]
			}

			parts := strings.SplitN(readType, "/", 2)
			if len(parts) > 1 {
				readType = parts[0]
				readExtender = parts[1]
			}
			var test parser.Tester
			switch readType {
			case "byte":
				num, err := strconv.ParseInt(testValue, 0, 8)
				if err != nil {
					return nil, fmt.Errorf("invalid byte test value '%s': %w", testValue, err)
				}
				if unsigned {
					test = parser.Uint8Test(parser.WithOffset(offset), uint8(num), comparator)
				} else {
					test = parser.Int8Test(parser.WithOffset(offset), int8(num), comparator)
				}
			case "short", "leshort", "beshort", "meshort":
				num, err := strconv.ParseInt(testValue, 0, 16)
				if err != nil {
					return nil, fmt.Errorf("invalid short test value '%s': %w", testValue, err)
				}
				if unsigned {
					test = parser.Uint16Test(parser.WithOffset(offset), uint16(num), comparator, endian)
				} else {
					test = parser.Int16Test(parser.WithOffset(offset), int16(num), comparator, endian)
				}
			case "long", "lelong", "belong", "melong":
				num, err := strconv.ParseInt(testValue, 0, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid long test value '%s': %w", testValue, err)
				}
				if unsigned {
					test = parser.Uint32Test(parser.WithOffset(offset), uint32(num), comparator, endian)
				} else {
					test = parser.Int32Test(parser.WithOffset(offset), int32(num), comparator, endian)
				}
			case "quad", "lequad", "bequad", "mequad":
				num, err := strconv.ParseInt(testValue, 0, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid quad test value '%s': %w", testValue, err)
				}
				if unsigned {
					test = parser.Uint64Test(parser.WithOffset(offset), uint64(num), comparator, endian)
				} else {
					test = parser.Int64Test(parser.WithOffset(offset), int64(num), comparator, endian)
				}
			case "float", "lefloat", "befloat", "mefloat":
				num, err := strconv.ParseFloat(testValue, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid float test value '%s': %w", testValue, err)
				}
				test = parser.Float32Test(parser.WithOffset(offset), float32(num), comparator, endian)
			case "double", "ledouble", "bedouble", "medouble":
				num, err := strconv.ParseFloat(testValue, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid float test value '%s': %w", testValue, err)
				}
				test = parser.Float64Test(parser.WithOffset(offset), float64(num), comparator, endian)
			case "string":
				test = parser.StringTest(parser.WithOffset(offset), testValue)
			case "pstring":
				//TODO: not yet supported
				continue
			case "date", "bedate", "ledate", "medate":
				// 4-byte date interpreted as Unix time UTC
				num, err := strconv.ParseInt(testValue, 0, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid 4-byte date test value '%s': %w", testValue, err)
				}
				test = parser.Date32Test(parser.WithOffset(offset), uint32(num), comparator, endian, false)
			case "qdate", "beqdate", "leqdate", "meqdate":
				// 8-byte date interpreted as Unix time UTC
				num, err := strconv.ParseInt(testValue, 0, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid 8-byte date test value '%s': %w", testValue, err)
				}
				test = parser.Date64Test(parser.WithOffset(offset), uint64(num), comparator, endian, false)
			case "ldate", "beldate", "leldate", "meldate":
				// 4-byte date interpreted as Unix time local timezone
				num, err := strconv.ParseInt(testValue, 0, 32)
				if err != nil {
					return nil, fmt.Errorf("invalid date test value '%s': %w", testValue, err)
				}
				test = parser.Date32Test(parser.WithOffset(offset), uint32(num), comparator, endian, true)
			case "qldate", "beqldate", "leqldate", "meqldate":
				// 8-byte date interpreted as Unix time local timezone
				num, err := strconv.ParseInt(testValue, 0, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid date test value '%s': %w", testValue, err)
				}
				test = parser.Date64Test(parser.WithOffset(offset), uint64(num), comparator, endian, true)
			case "qwdate", "beqwdate", "leqwdate", "meqwdate":
				// windows-style date
				num, err := strconv.ParseInt(testValue, 0, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid date test value '%s': %w", testValue, err)
				}
				test = parser.Date64Test(parser.WithOffset(offset), uint64(num), comparator, endian, true)
			case "beid3", "leid3":
				//TODO: not yet supported
				// ID3 tags
			case "indirect":
			case "name":
			case "use":
			case "regex":
				// supports extended posix regex
				test = parser.RegexTest(parser.WithOffset(offset), testValue, readExtender)
			case "search":
			case "default":
				if testValue != matchAny {
					return nil, fmt.Errorf("invalid default test value '%s'", testValue)
				}
				test = parser.DefaultTest(parser.WithOffset(offset), message)
			case "clear":
				if testValue != matchAny {
					return nil, fmt.Errorf("invalid clear test value '%s'", testValue)
				}
			case "der":
			case "guid":
				test = parser.GuidTest(parser.WithOffset(offset), testValue, comparator)
			case "offset":
			default:
				return nil, fmt.Errorf("invalid read type '%s'", readType)
			}
			tests = append(tests, parser.MagicTest{
				Test:    test,
				Message: message,
			})
		}

	}
	return tests, nil
}
