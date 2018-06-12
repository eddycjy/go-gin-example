package xlsx

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	maxNonScientificNumber = 1e11
	minNonScientificNumber = 1e-9
)

// CellType is an int type for storing metadata about the data type in the cell.
type CellType int

// These are the cell types from the ST_CellType spec
const (
	CellTypeString CellType = iota
	// CellTypeStringFormula is a specific format for formulas that return string values. Formulas that return numbers
	// and booleans are stored as those types.
	CellTypeStringFormula
	CellTypeNumeric
	CellTypeBool
	// CellTypeInline is not respected on save, all inline string cells will be saved as SharedStrings
	// when saving to an XLSX file. This the same behavior as that found in Excel.
	CellTypeInline
	CellTypeError
	// d (Date): Cell contains a date in the ISO 8601 format.
	// That is the only mention of this format in the XLSX spec.
	// Date seems to be unused by the current version of Excel, it stores dates as Numeric cells with a date format string.
	// For now these cells will have their value output directly. It is unclear if the value is supposed to be parsed
	// into a number and then formatted using the formatting or not.
	CellTypeDate
)

func (ct CellType) Ptr() *CellType {
	return &ct
}

// Cell is a high level structure intended to provide user access to
// the contents of Cell within an xlsx.Row.
type Cell struct {
	Row          *Row
	Value        string
	formula      string
	style        *Style
	NumFmt       string
	parsedNumFmt *parsedNumberFormat
	date1904     bool
	Hidden       bool
	HMerge       int
	VMerge       int
	cellType     CellType
}

// CellInterface defines the public API of the Cell.
type CellInterface interface {
	String() string
	FormattedValue() string
}

// NewCell creates a cell and adds it to a row.
func NewCell(r *Row) *Cell {
	return &Cell{Row: r, NumFmt: "general"}
}

// Merge with other cells, horizontally and/or vertically.
func (c *Cell) Merge(hcells, vcells int) {
	c.HMerge = hcells
	c.VMerge = vcells
}

// Type returns the CellType of a cell. See CellType constants for more details.
func (c *Cell) Type() CellType {
	return c.cellType
}

// SetString sets the value of a cell to a string.
func (c *Cell) SetString(s string) {
	c.Value = s
	c.formula = ""
	c.cellType = CellTypeString
}

// String returns the value of a Cell as a string.  If you'd like to
// see errors returned from formatting then please use
// Cell.FormattedValue() instead.
func (c *Cell) String() string {
	// To preserve the String() interface we'll throw away errors.
	// Not that using FormattedValue is therefore strongly
	// preferred.
	value, _ := c.FormattedValue()
	return value
}

// SetFloat sets the value of a cell to a float.
func (c *Cell) SetFloat(n float64) {
	c.SetValue(n)
}

//GetTime returns the value of a Cell as a time.Time
func (c *Cell) GetTime(date1904 bool) (t time.Time, err error) {
	f, err := c.Float()
	if err != nil {
		return t, err
	}
	return TimeFromExcelTime(f, date1904), nil
}

/*
	The following are samples of format samples.

	* "0.00e+00"
	* "0", "#,##0"
	* "0.00", "#,##0.00", "@"
	* "#,##0 ;(#,##0)", "#,##0 ;[red](#,##0)"
	* "#,##0.00;(#,##0.00)", "#,##0.00;[red](#,##0.00)"
	* "0%", "0.00%"
	* "0.00e+00", "##0.0e+0"
*/

// SetFloatWithFormat sets the value of a cell to a float and applies
// formatting to the cell.
func (c *Cell) SetFloatWithFormat(n float64, format string) {
	c.SetValue(n)
	c.NumFmt = format
	c.formula = ""
}

var timeLocationUTC, _ = time.LoadLocation("UTC")

func TimeToUTCTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), timeLocationUTC)
}

func TimeToExcelTime(t time.Time) float64 {
	return float64(t.UnixNano())/8.64e13 + 25569.0
}

// DateTimeOptions are additional options for exporting times
type DateTimeOptions struct {
	// Location allows calculating times in other timezones/locations
	Location *time.Location
	// ExcelTimeFormat is the string you want excel to use to format the datetime
	ExcelTimeFormat string
}

var (
	DefaultDateFormat     = builtInNumFmt[14]
	DefaultDateTimeFormat = builtInNumFmt[22]

	DefaultDateOptions = DateTimeOptions{
		Location:        timeLocationUTC,
		ExcelTimeFormat: DefaultDateFormat,
	}

	DefaultDateTimeOptions = DateTimeOptions{
		Location:        timeLocationUTC,
		ExcelTimeFormat: DefaultDateTimeFormat,
	}
)

// SetDate sets the value of a cell to a float.
func (c *Cell) SetDate(t time.Time) {
	c.SetDateWithOptions(t, DefaultDateOptions)
}

func (c *Cell) SetDateTime(t time.Time) {
	c.SetDateWithOptions(t, DefaultDateTimeOptions)
}

// SetDateWithOptions allows for more granular control when exporting dates and times
func (c *Cell) SetDateWithOptions(t time.Time, options DateTimeOptions) {
	_, offset := t.In(options.Location).Zone()
	t = time.Unix(t.Unix()+int64(offset), 0)
	c.SetDateTimeWithFormat(TimeToExcelTime(t.In(timeLocationUTC)), options.ExcelTimeFormat)
}

func (c *Cell) SetDateTimeWithFormat(n float64, format string) {
	c.Value = strconv.FormatFloat(n, 'f', -1, 64)
	c.NumFmt = format
	c.formula = ""
	c.cellType = CellTypeNumeric
}

// Float returns the value of cell as a number.
func (c *Cell) Float() (float64, error) {
	f, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		return math.NaN(), err
	}
	return f, nil
}

// SetInt64 sets a cell's value to a 64-bit integer.
func (c *Cell) SetInt64(n int64) {
	c.SetValue(n)
}

// Int64 returns the value of cell as 64-bit integer.
func (c *Cell) Int64() (int64, error) {
	f, err := strconv.ParseInt(c.Value, 10, 64)
	if err != nil {
		return -1, err
	}
	return f, nil
}

// GeneralNumeric returns the value of the cell as a string. It is formatted very closely to the the XLSX spec for how
// to display values when the storage type is Number and the format type is General. It is not 100% identical to the
// spec but is as close as you can get using the built in Go formatting tools.
func (c *Cell) GeneralNumeric() (string, error) {
	return generalNumericScientific(c.Value, true)
}

// GeneralNumericWithoutScientific returns numbers that are always formatted as numbers, but it does not follow
// the rules for when XLSX should switch to scientific notation, since sometimes scientific notation is not desired,
// even if that is how the document is supposed to be formatted.
func (c *Cell) GeneralNumericWithoutScientific() (string, error) {
	return generalNumericScientific(c.Value, false)
}

// SetInt sets a cell's value to an integer.
func (c *Cell) SetInt(n int) {
	c.SetValue(n)
}

// SetInt sets a cell's value to an integer.
func (c *Cell) SetValue(n interface{}) {
	switch t := n.(type) {
	case time.Time:
		c.SetDateTime(t)
		return
	case int, int8, int16, int32, int64:
		c.setNumeric(fmt.Sprintf("%d", n))
	case float64:
		// When formatting floats, do not use fmt.Sprintf("%v", n), this will cause numbers below 1e-4 to be printed in
		// scientific notation. Scientific notation is not a valid way to store numbers in XML.
		// Also not not use fmt.Sprintf("%f", n), this will cause numbers to be stored as X.XXXXXX. Which means that
		// numbers will lose precision and numbers with fewer significant digits such as 0 will be stored as 0.000000
		// which causes tests to fail.
		c.setNumeric(strconv.FormatFloat(t, 'f', -1, 64))
	case float32:
		c.setNumeric(strconv.FormatFloat(float64(t), 'f', -1, 32))
	case string:
		c.SetString(t)
	case []byte:
		c.SetString(string(t))
	case nil:
		c.SetString("")
	default:
		c.SetString(fmt.Sprintf("%v", n))
	}
}

// setNumeric sets a cell's value to a number
func (c *Cell) setNumeric(s string) {
	c.Value = s
	c.NumFmt = builtInNumFmt[builtInNumFmtIndex_GENERAL]
	c.formula = ""
	c.cellType = CellTypeNumeric
}

// Int returns the value of cell as integer.
// Has max 53 bits of precision
// See: float64(int64(math.MaxInt))
func (c *Cell) Int() (int, error) {
	f, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		return -1, err
	}
	return int(f), nil
}

// SetBool sets a cell's value to a boolean.
func (c *Cell) SetBool(b bool) {
	if b {
		c.Value = "1"
	} else {
		c.Value = "0"
	}
	c.cellType = CellTypeBool
}

// Bool returns a boolean from a cell's value.
// TODO: Determine if the current return value is
// appropriate for types other than CellTypeBool.
func (c *Cell) Bool() bool {
	// If bool, just return the value.
	if c.cellType == CellTypeBool {
		return c.Value == "1"
	}
	// If numeric, base it on a non-zero.
	if c.cellType == CellTypeNumeric {
		return c.Value != "0"
	}
	// Return whether there's an empty string.
	return c.Value != ""
}

// SetFormula sets the format string for a cell.
func (c *Cell) SetFormula(formula string) {
	c.formula = formula
	c.cellType = CellTypeNumeric
}

func (c *Cell) SetStringFormula(formula string) {
	c.formula = formula
	c.cellType = CellTypeStringFormula
}

// Formula returns the formula string for the cell.
func (c *Cell) Formula() string {
	return c.formula
}

// GetStyle returns the Style associated with a Cell
func (c *Cell) GetStyle() *Style {
	if c.style == nil {
		c.style = NewStyle()
	}
	return c.style
}

// SetStyle sets the style of a cell.
func (c *Cell) SetStyle(style *Style) {
	c.style = style
}

// GetNumberFormat returns the number format string for a cell.
func (c *Cell) GetNumberFormat() string {
	return c.NumFmt
}

func (c *Cell) formatToFloat(format string) (string, error) {
	f, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		return c.Value, err
	}
	return fmt.Sprintf(format, f), nil
}

func (c *Cell) formatToInt(format string) (string, error) {
	f, err := strconv.ParseFloat(c.Value, 64)
	if err != nil {
		return c.Value, err
	}
	return fmt.Sprintf(format, int(f)), nil
}

// getNumberFormat will update the parsedNumFmt struct if it has become out of date, since a cell's NumFmt string is a
// public field that could be edited by clients.
func (c *Cell) getNumberFormat() *parsedNumberFormat {
	if c.parsedNumFmt == nil || c.parsedNumFmt.numFmt != c.NumFmt {
		c.parsedNumFmt = parseFullNumberFormatString(c.NumFmt)
	}
	return c.parsedNumFmt
}

// FormattedValue returns a value, and possibly an error condition
// from a Cell.  If it is possible to apply a format to the cell
// value, it will do so, if not then an error will be returned, along
// with the raw value of the Cell.
func (c *Cell) FormattedValue() (string, error) {
	fullFormat := c.getNumberFormat()
	returnVal, err := fullFormat.FormatValue(c)
	if fullFormat.parseEncounteredError != nil {
		return returnVal, *fullFormat.parseEncounteredError
	}
	return returnVal, err
}
