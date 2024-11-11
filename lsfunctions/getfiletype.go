package lsfunctions

import "strings"

func getFileType(entry Entry) string {
	mod := entry.Mode

	if strings.ContainsAny(mod, "x") {
		return "exec"
	}
	name := entry.Name
	tokens := strings.Split(strings.ToLower(name), ".")
	ext := tokens[len(tokens)-1]
	switch ext {
	case ".txt", ".md", ".log":
		return "text"
	case ".pdf":
		return "pdf"
	case ".doc", ".docx":
		return "word"
	case ".xls", ".xlsx":
		return "excel"
	case ".ppt", ".pptx":
		return "powerpoint"
	case ".zip", ".tar", ".gz", ".7z", "deb":
		return "archive"
	case ".mp3", ".wav", ".flac":
		return "audio"
	case ".mp4", ".avi", ".mkv":
		return "video"
	case ".jpg", ".jpeg", ".png", ".gif":
		return "image"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	default:
		return "other"
	}
}
