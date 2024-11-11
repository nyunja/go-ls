package lsfunctions

import "strings"

func getFileType(entry Entry) string {
	mod := entry.Mode

	if strings.ContainsAny(mod, "x") {
		return "exec"
	}
	name := strings.Trim(entry.Name, "'")
	tokens := strings.Split(strings.ToLower(name), ".")
	ext := tokens[len(tokens)-1]
	switch "." + ext {
	case ".log":
		return "text"
	case ".pdf":
		return "pdf"
	case "rar", ".zip", ".tar", ".gz", ".7z", "deb":
		return "archive"
	case ".mp3", ".wav", ".flac":
		return "audio"
	case ".webp",".jpg", ".jpeg", ".png", ".gif":
		return "image"
	case ".crdownload":
		return "crd"
	default:
		return "other"
	}
}
