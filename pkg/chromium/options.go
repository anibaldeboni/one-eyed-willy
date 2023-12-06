package chromium

import "time"

// Options are the available expectedOptions for converting HTML document to PDF.
type Options struct {
	// FailOnConsoleExceptions sets if the conversion should fail if there are
	// exceptions in the Chromium console.
	// Optional.
	FailOnConsoleExceptions bool

	// WaitDelay is the duration to wait when loading an HTML document before
	// converting it to PDF.
	// Optional.
	WaitDelay time.Duration

	// WaitWindowStatus is the window.status value to wait for before
	// converting an HTML document to PDF.
	// Optional.
	WaitWindowStatus string

	// WaitForExpression is the custom JavaScript expression to wait before
	// converting an HTML document to PDF until it returns true
	// Optional.
	WaitForExpression string

	// ExtraHttpHeaders are the HTTP headers to send by Chromium while loading
	// the HTML document.
	// Optional.
	ExtraHttpHeaders map[string]string

	// EmulatedMediaType is the media type to emulate, either "screen" or
	// "print".
	// Optional.
	EmulatedMediaType string

	// Landscape sets the paper orientation.
	// Optional.
	Landscape bool

	// PrintBackground prints the background graphics.
	// Optional.
	PrintBackground bool

	// OmitBackground hides default white background and allows generating PDFs
	// with transparency.
	// Optional.
	OmitBackground bool

	// Scale is the scale of the page rendering.
	// Optional.
	Scale float64

	// PaperWidth is the paper width, in inches.
	// Optional.
	PaperWidth float64

	// PaperHeight is the paper height, in inches.
	// Optional.
	PaperHeight float64

	// MarginTop is the top margin, in inches.
	// Optional.
	MarginTop float64

	// MarginBottom is the bottom margin, in inches.
	// Optional.
	MarginBottom float64

	// MarginLeft is the left margin, in inches.
	// Optional.
	MarginLeft float64

	// MarginRight is the right margin, in inches.
	// Optional.
	MarginRight float64

	// Page ranges to print, e.g., '1-5, 8, 11-13'. Empty means all pages.
	// Optional.
	PageRanges string

	// HeaderTemplate is the HTML template of the header. It should be valid
	// HTML  markup with following classes used to inject printing values into
	// them:
	// - date: formatted print date
	// - title: document title
	// - url: document location
	// - pageNumber: current page number
	// - totalPages: total pages in the document
	// For example, <span class=title></span> would generate span containing
	// the title.
	// Optional.
	HeaderTemplate string

	// FooterTemplate is the HTML template of the footer. It should use the
	// same format as the HeaderTemplate.
	// Optional.
	FooterTemplate string

	// PreferCssPageSize defines whether to prefer page size as defined by CSS.
	// If false, the content will be scaled to fit the paper size.
	// Optional.
	PreferCssPageSize bool
}

const (
	A4_WIDTH_IN_INCHES  = 8.3
	A4_HEIGHT_IN_INCHES = 11.7
)

// DefaultOptions returns the default values for Options.
func DefaultOptions() Options {
	return Options{
		FailOnConsoleExceptions: false,
		WaitDelay:               0,
		WaitWindowStatus:        "",
		WaitForExpression:       "",
		ExtraHttpHeaders:        nil,
		EmulatedMediaType:       "",
		Landscape:               false,
		PrintBackground:         false,
		OmitBackground:          false,
		Scale:                   1.0,
		PaperWidth:              A4_WIDTH_IN_INCHES,
		PaperHeight:             A4_HEIGHT_IN_INCHES,
		MarginTop:               0.39,
		MarginBottom:            0.39,
		MarginLeft:              0.39,
		MarginRight:             0.39,
		PageRanges:              "",
		HeaderTemplate:          "<html><head></head><body></body></html>",
		FooterTemplate:          "<html><head></head><body></body></html>",
		PreferCssPageSize:       false,
	}
}
