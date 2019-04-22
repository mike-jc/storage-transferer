package modelsDataStorage

import "strings"

type FileInfo struct {
	Name      string
	Extension string
	Path      string
	Size      int
}

func NameFromPath(path string) string {
	parts := strings.Split(path, "/")
	fullName := parts[len(parts)-1]
	parts = strings.Split(fullName, ".")
	return parts[0]
}

func ExtensionFromPath(path string) string {
	parts := strings.Split(path, "/")
	fullName := parts[len(parts)-1]
	parts = strings.Split(fullName, ".")
	return parts[1]
}

func ExtensionFromMime(mime string) string {
	switch mime {
	case "video/mp4":
		return "mp4"
	case "audio/mp3":
		return "mp3"
	default:
		return ""
	}
}
