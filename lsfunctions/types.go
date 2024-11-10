package lsfunctions

import "os"

// FileInfo struct to store file information from readDir function
type FileInfo struct {
	Name       string
	Info       os.FileInfo
	LinkTarget string
	Rdev       uint64
}

type Entry struct {
	Name, Mode, User, Owner, Group, Type,
	LinkTarget, LinkCount, Size, Minor, Time string
	IsDirectory bool
}

type TotalBlocks int64

var ShowTotals bool

type Widths struct {
	sizeCol, ownerCol, groupCol, linkCol, timeCol, modCol, minorCol int
}

