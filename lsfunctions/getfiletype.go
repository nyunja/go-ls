package lsfunctions

import "strings"

func getFileType(entry Entry) (Entry, string) {
	mod := entry.Mode
	// Handle setuid files
	if entry.Mode[0] == 'u' {
		entry.Mode = swapU(entry.Mode)
		return entry, "setuid"
	}
	// Handle setuid files
	if entry.Mode[0] == 'g' {
		entry.Mode = swapG(entry.Mode)
		return entry, "setgid"
	}
	// Handle directories
	if entry.Mode[0] == 'd' {
		if strings.Contains(entry.Mode, "t") {
			entry.Mode = swapT(entry.Mode)
		}
		return entry, "dir"
	}
	// Handle device files
	if entry.Mode[0] == 'D' {
		entry.Mode = "b" + strings.TrimPrefix(entry.Mode, "D")
		return entry, "dev"
	}

	// Ensure mode is properly formatted
	if len(entry.Mode) != 10 {
		entry.Mode = "b" + entry.Mode
		if len(entry.Mode) != 10 {
			entry.Mode += "-"
		}
	}
	if strings.ContainsAny(mod, "x") {
		return entry, "exec"
	}
	name := strings.Trim(entry.Name, "'")
	tokens := strings.Split(strings.ToLower(name), ".")
	ext := "." + tokens[len(tokens)-1]
	switch ext {
	case ".py", ".js", ".go":
		return entry, "exec"
	case ".log":
		return entry, "text"
	case ".pdf":
		return entry, "pdf"
	case "rar", ".zip", ".tar", ".gz", ".7z", "deb":
		return entry, "archive"
	case ".mp3", ".wav", ".flac":
		return entry, "audio"
	case ".webp", ".jpg", ".jpeg", ".png", ".gif":
		return entry, "image"
	case ".crdownload":
		return entry, "crd"
	default:
		return entry, "other"
	}
}

func swapG(s string) string {
	s = strings.Replace(s, "g", "-", 1)
	res := []rune{}
	count := 0
	for _, ch := range s {
		if ch == 'x' {
			count++
		}
		if count == 2 {
			res = append(res, 's')
			count++
		} else {
			res = append(res, ch)
		}
	}
	return string(res)
}
