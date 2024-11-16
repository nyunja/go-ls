package lsfunctions

import (
	"os"
	"syscall"
)

func formatPermissionsWithACL(path string, mode os.FileMode) (string, error) {
	permissions := formatPermissions(mode)
	hasACL, err := hasExtendedAttributes(path)

	if err != nil {
		return permissions, nil
	}
	if hasACL {
		permissions += "+"
	}
	return permissions, nil
}

func hasExtendedAttributes(path string) (bool, error) {
	dest := []byte{}
	sz, err := syscall.Listxattr(path, dest)
	if err != nil {
		// Extended attributes not supported on this platform
		if err == syscall.ENOTSUP {
			return false, nil
		}
		return false, err
	}
	return sz > 0, nil
}

func formatPermissions(mode os.FileMode) string {
	var buf [10]byte

	// Permissions(rwx)
	rwx := []os.FileMode{0400, 0200, 0100, 0040, 0020, 0010, 0004, 0002, 0001}
	symbols := "rwxrwxrwx"
	for i, p := range rwx {
		if mode&p != 0 {
			buf[i+1] = symbols[i]
		} else {
			buf[i+1] = '-'
		}
	}
	if mode&os.ModeDir != 0 {
		buf[0] = 'd'
	} else if mode&os.ModeSymlink != 0 {
		buf[0] = 'l'
	} else if mode&os.ModeNamedPipe != 0 {
		buf[0] = 'p'
	} else if mode&os.ModeSocket != 0 {
		buf[0] = 's'
	} else if mode&os.ModeDevice != 0 {
		if mode&os.ModeCharDevice != 0 {
			buf[0] = 'c'
		} else {
			buf[0] = 'b'
		}
	} else {
		buf[0] = '-'
	}

	if mode&os.ModeSticky != 0 {
		if buf[9] == 'x' {
			buf[9] = 't'
		} else {
			buf[9] = 'T'
		}
	}
	if mode&os.ModeSetuid != 0 {
		if buf[3] == 'x' {
			buf[3] = 's'
		} else {
			buf[3] = 'S'
		}
	}
	if mode&os.ModeSetgid != 0 {
		if buf[6] == 'x' {
			buf[6] = 's'
		} else {
			buf[6] = 'S'
		}
	}
	return string(buf[:])
}
