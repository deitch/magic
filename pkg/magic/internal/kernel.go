package internal

import (
	"encoding/binary"

	parser "github.com/deitch/magic/pkg/magic/parser"
)

func init() {
	AllTests = append(AllTests, kernelTests...)
}

var kernelTests = []parser.MagicTest{
	{Test: parser.StringTest(parser.WithOffset(514), "HdrS"), Message: "Linux kernel", Children: []parser.MagicTest{
		// the original for this is "leshort" - we are just treating this as a short, and would handle endianness when parsing
		// the magic file
		{Test: parser.Uint16Test(parser.WithOffset(510), 0xAA55, parser.Equal, binary.LittleEndian), Message: "x86 boot executable", Children: []parser.MagicTest{
			{Test: parser.Uint16Test(parser.WithOffset(518), 0x1ff, parser.GreaterThan, binary.LittleEndian), Children: []parser.MagicTest{
				{Test: parser.Int8Test(parser.WithOffset(529), 0, parser.Equal), Message: "zImage"},
				{Test: parser.Int8Test(parser.WithOffset(529), 1, parser.Equal), Message: "bzImage"},
				{Test: parser.Uint32Test(parser.WithOffset(526), 0, parser.GreaterThan, binary.LittleEndian), Children: []parser.MagicTest{
					{Test: parser.Int8Test(parser.WithChainedOffsetReaders(parser.WithIndirectOffsetShort(526, binary.LittleEndian), parser.WithOffset(0x200)), 0, parser.GreaterThan), Message: "version %s"},
				}},
			}},
			{Test: parser.Uint16Test(parser.WithOffset(498), 1, parser.Equal, binary.LittleEndian), Message: "RO-rootFS"},
			{Test: parser.Uint16Test(parser.WithOffset(498), 0, parser.Equal, binary.LittleEndian), Message: "RW-rootFS"},
			{Test: parser.Uint16Test(parser.WithOffset(508), 0, parser.GreaterThan, binary.LittleEndian), Message: "root_dev %#X"},
			{Test: parser.Uint16Test(parser.WithOffset(502), 0, parser.GreaterThan, binary.LittleEndian), Message: "swap_dev %#X"},
			{Test: parser.Uint16Test(parser.WithOffset(504), 0, parser.GreaterThan, binary.LittleEndian), Message: "RAMdisksize %u KB"},
			{Test: parser.Uint16Test(parser.WithOffset(506), 0xffff, parser.Equal, binary.LittleEndian), Message: "Normal VGA"},
			{Test: parser.Uint16Test(parser.WithOffset(506), 0xfffe, parser.Equal, binary.LittleEndian), Message: "Extended VGA"},
			{Test: parser.Uint16Test(parser.WithOffset(506), 0xfffd, parser.Equal, binary.LittleEndian), Message: "Prompt for Videomode"},
			{Test: parser.Uint16Test(parser.WithOffset(506), 0x0, parser.GreaterThan, binary.LittleEndian), Message: "Video mode %d"},
		},
		},
	},
	}}
