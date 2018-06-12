package xlsx

// Default column width in excel
const ColWidth = 9.5

type Col struct {
	Min          int
	Max          int
	Hidden       bool
	Width        float64
	Collapsed    bool
	OutlineLevel uint8
	numFmt       string
	parsedNumFmt *parsedNumberFormat
	style        *Style
}

// SetType will set the format string of a column based on the type that you want to set it to.
// This function does not really make a lot of sense.
func (c *Col) SetType(cellType CellType) {
	switch cellType {
	case CellTypeString:
		c.numFmt = builtInNumFmt[builtInNumFmtIndex_STRING]
	case CellTypeNumeric:
		c.numFmt = builtInNumFmt[builtInNumFmtIndex_INT]
	case CellTypeBool:
		c.numFmt = builtInNumFmt[builtInNumFmtIndex_GENERAL] //TEMP
	case CellTypeInline:
		c.numFmt = builtInNumFmt[builtInNumFmtIndex_STRING]
	case CellTypeError:
		c.numFmt = builtInNumFmt[builtInNumFmtIndex_GENERAL] //TEMP
	case CellTypeDate:
		// Cells that are stored as dates are not properly supported in this library.
		// They should instead be stored as a Numeric with a date format.
		c.numFmt = builtInNumFmt[builtInNumFmtIndex_GENERAL]
	case CellTypeStringFormula:
		c.numFmt = builtInNumFmt[builtInNumFmtIndex_STRING]
	}
}

// GetStyle returns the Style associated with a Col
func (c *Col) GetStyle() *Style {
	return c.style
}

// SetStyle sets the style of a Col
func (c *Col) SetStyle(style *Style) {
	c.style = style
}
