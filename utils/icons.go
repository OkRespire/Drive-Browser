package utils

import (
	"path/filepath"
	"strings"
)

// File type icons with Unicode values
const (
	// Folders
	IconFolder     = "📁" // U+1F4C1
	IconFolderOpen = "📂" // U+1F4C2

	// Documents
	IconDocument    = "📄" // U+1F4C4
	IconTextFile    = "📝" // U+1F4DD
	IconPDF         = "📕" // U+1F4D5
	IconWord        = "📘" // U+1F4D8
	IconPowerPoint  = "📙" // U+1F4D9
	IconExcel       = "📊" // U+1F4CA
	IconSpreadsheet = "📈" // U+1F4C8

	// Media
	IconImage = "🖼️" // U+1F5BC + U+FE0F
	IconVideo = "🎬"  // U+1F3AC
	IconAudio = "🎵"  // U+1F3B5
	IconMusic = "🎶"  // U+1F3B6

	// Code files
	IconCode       = "💻"  // U+1F4BB
	IconHTML       = "🌐"  // U+1F310
	IconCSS        = "🎨"  // U+1F3A8
	IconJavaScript = "⚡"  // U+26A1
	IconPython     = "🐍"  // U+1F40D
	IconJava       = "☕"  // U+2615
	IconCPlusPlus  = "⚙️" // U+2699 + U+FE0F
	IconGo         = "🐹"  // U+1F439 (gopher)
	IconRust       = "🦀"  // U+1F980

	// Archives
	IconArchive = "📦"  // U+1F4E6
	IconZip     = "🗜️" // U+1F5DC + U+FE0F

	// Other
	IconDatabase = "🗃️" // U+1F5C3 + U+FE0F
	IconFont     = "🔤"  // U+1F524
	IconConfig   = "⚙️" // U+2699 + U+FE0F
	IconLog      = "📋"  // U+1F4CB
	IconMarkdown = "📝"  // U+1F4DD
	IconJSON     = "📋"  // U+1F4CB
	IconXML      = "📄"  // U+1F4C4
	IconDefault  = "📄"  // U+1F4C4
)

// MIME type to icon mapping
var mimeTypeIcons = map[string]string{
	// Google Workspace
	"application/vnd.google-apps.folder":       IconFolder,
	"application/vnd.google-apps.document":     IconWord,
	"application/vnd.google-apps.spreadsheet":  IconExcel,
	"application/vnd.google-apps.presentation": IconPowerPoint,
	"application/vnd.google-apps.drawing":      IconImage,
	"application/vnd.google-apps.form":         IconDocument,
	"application/vnd.google-apps.site":         IconHTML,

	// Microsoft Office
	"application/msword": IconWord,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": IconWord,
	"application/vnd.ms-excel": IconExcel,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         IconExcel,
	"application/vnd.ms-powerpoint":                                             IconPowerPoint,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": IconPowerPoint,

	// Documents
	"application/pdf": IconPDF,
	"text/plain":      IconTextFile,
	"text/markdown":   IconMarkdown,
	"application/rtf": IconDocument,

	// Images
	"image/jpeg":    IconImage,
	"image/jpg":     IconImage,
	"image/png":     IconImage,
	"image/gif":     IconImage,
	"image/bmp":     IconImage,
	"image/svg+xml": IconImage,
	"image/webp":    IconImage,
	"image/tiff":    IconImage,

	// Audio
	"audio/mpeg":     IconMusic,
	"audio/mp3":      IconMusic,
	"audio/wav":      IconAudio,
	"audio/ogg":      IconAudio,
	"audio/flac":     IconAudio,
	"audio/aac":      IconMusic,
	"audio/x-ms-wma": IconAudio,

	// Video
	"video/mp4":       IconVideo,
	"video/avi":       IconVideo,
	"video/quicktime": IconVideo,
	"video/x-msvideo": IconVideo,
	"video/x-ms-wmv":  IconVideo,
	"video/webm":      IconVideo,
	"video/ogg":       IconVideo,
	"video/3gpp":      IconVideo,

	// Archives
	"application/zip":              IconZip,
	"application/x-rar-compressed": IconArchive,
	"application/x-7z-compressed":  IconArchive,
	"application/x-tar":            IconArchive,
	"application/gzip":             IconArchive,
	"application/x-bzip2":          IconArchive,

	// Code files
	"text/html":              IconHTML,
	"text/css":               IconCSS,
	"application/javascript": IconJavaScript,
	"text/javascript":        IconJavaScript,
	"application/json":       IconJSON,
	"application/xml":        IconXML,
	"text/xml":               IconXML,

	// Other
	"application/x-sqlite3": IconDatabase,
	"application/sql":       IconDatabase,
}

// Extension to icon mapping (fallback when MIME type is not available)
var extensionIcons = map[string]string{
	// Documents
	".pdf":  IconPDF,
	".doc":  IconWord,
	".docx": IconWord,
	".xls":  IconExcel,
	".xlsx": IconExcel,
	".ppt":  IconPowerPoint,
	".pptx": IconPowerPoint,
	".txt":  IconTextFile,
	".md":   IconMarkdown,
	".rtf":  IconDocument,

	// Images
	".jpg":  IconImage,
	".jpeg": IconImage,
	".png":  IconImage,
	".gif":  IconImage,
	".bmp":  IconImage,
	".svg":  IconImage,
	".webp": IconImage,
	".tiff": IconImage,
	".ico":  IconImage,

	// Audio
	".mp3":  IconMusic,
	".wav":  IconAudio,
	".ogg":  IconAudio,
	".flac": IconAudio,
	".aac":  IconMusic,
	".wma":  IconAudio,
	".m4a":  IconMusic,

	// Video
	".mp4":  IconVideo,
	".avi":  IconVideo,
	".mov":  IconVideo,
	".wmv":  IconVideo,
	".webm": IconVideo,
	".mkv":  IconVideo,
	".flv":  IconVideo,
	".3gp":  IconVideo,

	// Archives
	".zip": IconZip,
	".rar": IconArchive,
	".7z":  IconArchive,
	".tar": IconArchive,
	".gz":  IconArchive,
	".bz2": IconArchive,

	// Code files
	".html": IconHTML,
	".htm":  IconHTML,
	".css":  IconCSS,
	".js":   IconJavaScript,
	".json": IconJSON,
	".xml":  IconXML,
	".py":   IconPython,
	".java": IconJava,
	".cpp":  IconCPlusPlus,
	".c":    IconCode,
	".h":    IconCode,
	".go":   IconGo,
	".rs":   IconRust,
	".php":  IconCode,
	".rb":   IconCode,
	".sh":   IconCode,
	".bat":  IconCode,
	".ps1":  IconCode,

	// Config files
	".yml":    IconConfig,
	".yaml":   IconConfig,
	".toml":   IconConfig,
	".ini":    IconConfig,
	".conf":   IconConfig,
	".config": IconConfig,

	// Logs
	".log": IconLog,

	// Fonts
	".ttf":   IconFont,
	".otf":   IconFont,
	".woff":  IconFont,
	".woff2": IconFont,

	// Database
	".db":     IconDatabase,
	".sqlite": IconDatabase,
	".sql":    IconDatabase,
}

// GetFileIcon returns the appropriate icon for a file based on MIME type or extension
func GetFileIcon(filename, mimeType string) string {
	// Check if it's a folder first
	if mimeType == "application/vnd.google-apps.folder" {
		return IconFolder
	}

	// Try MIME type first
	if icon, exists := mimeTypeIcons[mimeType]; exists {
		return icon
	}

	// Fallback to extension
	ext := strings.ToLower(filepath.Ext(filename))
	if icon, exists := extensionIcons[ext]; exists {
		return icon
	}

	// Default icon
	return IconDefault
}
