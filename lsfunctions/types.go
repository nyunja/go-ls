package lsfunctions

import "os"

// FileDetails struct to store file information from readDir function
type FileDetails struct {
	Path       string
	Name       string
	Info       os.FileInfo
	LinkTarget string
	IsBrokenLink bool
	Rdev       uint64
	TargetInfo TargetInfo
}

type Entry struct {
	Name, Mode, User, Owner, Group, Type,
	LinkTarget, LinkCount, Size, Minor, Time, Path string
	IsDirectory, IsBrokenLink bool
	TargetInfo TargetInfo
}

type TargetInfo struct {
	Name, Mode string
	IsBrokenLink bool
}

type TotalBlocks int64

var ShowTotals bool

type Widths struct {
	sizeCol, ownerCol, groupCol, linkCol, timeCol, modCol, minorCol int
}
