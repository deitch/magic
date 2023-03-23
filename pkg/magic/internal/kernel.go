package internal

import (
	parser "github.com/deitch/magic/pkg/magic/parser"
)

func init() {
	AllTests = append(AllTests, kernelTests...)
}

var kernelTests = []parser.MagicTest{
	{Test: parser.StringTest(parser.WithOffset(514), "HdrS"), Message: parser.WithMessage("Linux kernel"), Children: []parser.MagicTest{
		// the original for this is "leshort" - we are just treating this as a short, and would handle endianness when parsing
		// the magic file
		{Test: parser.ShortTestLittleEndian(parser.WithOffset(510), 0xAA55, parser.Equal), Message: parser.WithMessage("x86 boot executable"), Children: []parser.MagicTest{
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(518), 0x1ff, parser.GreaterThan), Message: parser.WithMessageEmpty(), Children: []parser.MagicTest{
				{Test: parser.ByteTest(parser.WithOffset(529), 0, parser.Equal), Message: parser.WithMessage("zImage")},
				{Test: parser.ByteTest(parser.WithOffset(529), 1, parser.Equal), Message: parser.WithMessage("bzImage")},
				{Test: parser.LongTestLittleEndian(parser.WithOffset(526), 0, parser.GreaterThan), Message: parser.WithMessageEmpty(), Children: []parser.MagicTest{
					{Test: parser.ByteTest(parser.WithChainedOffsetReaders(parser.WithIndirectOffsetShortLittleEndian(526), parser.WithOffset(0x200)), 0, parser.GreaterThan), Message: parser.WithMessagePattern("version %s")},
				}},
			}},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(498), 1, parser.Equal), Message: parser.WithMessage("RO-rootFS")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(498), 0, parser.Equal), Message: parser.WithMessage("RW-rootFS")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(508), 0, parser.GreaterThan), Message: parser.WithMessagePattern("root_dev %#X")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(502), 0, parser.GreaterThan), Message: parser.WithMessagePattern("swap_dev %#X")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(504), 0, parser.GreaterThan), Message: parser.WithMessagePattern("RAMdisksize %u KB")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(506), 0xffff, parser.Equal), Message: parser.WithMessage("Normal VGA")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(506), 0xfffe, parser.Equal), Message: parser.WithMessage("Extended VGA")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(506), 0xfffd, parser.Equal), Message: parser.WithMessage("Prompt for Videomode")},
			{Test: parser.ShortTestLittleEndian(parser.WithOffset(506), 0x0, parser.GreaterThan), Message: parser.WithMessagePattern("Video mode %d")},
		},
		},
	},
	}}
