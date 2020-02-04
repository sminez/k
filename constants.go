package k

const (
	// DefaultPort is the default port to server the HTTP handler for fzf previews
	DefaultPort              = 1992
	helpfilePath             = "HELPFILE_PATH"
	defaultHelpfileDirectory = ".helpfiles"

	// Markup constants for snippet files
	endOfSectionMarker = "--"
	commentMarker      = '>'
	titleMarker        = '#'
	codeMarker         = '$'
	urlMarker          = '%'
	tagMarker          = '?'

	// UI and interaction with fzf
	previewCmd      = "\"curl -s -X POST -d \"{}\" 127.0.0.1:%d/\""
	previewPosition = "--preview-window=up:70%"
	seperator       = " | " // Delimits sections in the lines passed to fzf
	yellow          = 33    // ANSI terminal escape sequence
	blue            = 34    // ANSI terminal escape sequence
	grey            = 37    // ANSI terminal escape sequence
)