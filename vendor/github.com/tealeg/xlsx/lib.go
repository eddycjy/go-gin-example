package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
)

const (
	sheetEnding = `</sheetData></worksheet>`
)

// XLSXReaderError is the standard error type for otherwise undefined
// errors in the XSLX reading process.
type XLSXReaderError struct {
	Err string
}

// Error returns a string value from an XLSXReaderError struct in order
// that it might comply with the builtin.error interface.
func (e *XLSXReaderError) Error() string {
	return e.Err
}

// getRangeFromString is an internal helper function that converts
// XLSX internal range syntax to a pair of integers.  For example,
// the range string "1:3" yield the upper and lower integers 1 and 3.
func getRangeFromString(rangeString string) (lower int, upper int, error error) {
	var parts []string
	parts = strings.SplitN(rangeString, ":", 2)
	if parts[0] == "" {
		error = errors.New(fmt.Sprintf("Invalid range '%s'\n", rangeString))
	}
	if parts[1] == "" {
		error = errors.New(fmt.Sprintf("Invalid range '%s'\n", rangeString))
	}
	lower, error = strconv.Atoi(parts[0])
	if error != nil {
		error = errors.New(fmt.Sprintf("Invalid range (not integer in lower bound) %s\n", rangeString))
	}
	upper, error = strconv.Atoi(parts[1])
	if error != nil {
		error = errors.New(fmt.Sprintf("Invalid range (not integer in upper bound) %s\n", rangeString))
	}
	return lower, upper, error
}

// ColLettersToIndex is used to convert a character based column
// reference to a zero based numeric column identifier.
func ColLettersToIndex(letters string) int {
	sum, mul, n := 0, 1, 0
	for i := len(letters) - 1; i >= 0; i, mul, n = i-1, mul*26, 1 {
		c := letters[i]
		switch {
		case 'A' <= c && c <= 'Z':
			n += int(c - 'A')
		case 'a' <= c && c <= 'z':
			n += int(c - 'a')
		}
		sum += n * mul
	}
	return sum
}

// Get the largestDenominator that is a multiple of a basedDenominator
// and fits at least once into a given numerator.
func getLargestDenominator(numerator, multiple, baseDenominator, power int) (int, int) {
	if numerator/multiple == 0 {
		return 1, power
	}
	next, nextPower := getLargestDenominator(
		numerator, multiple*baseDenominator, baseDenominator, power+1)
	if next > multiple {
		return next, nextPower
	}
	return multiple, power
}

// Convers a list of numbers representing a column into a alphabetic
// representation, as used in the spreadsheet.
func formatColumnName(colId []int) string {
	lastPart := len(colId) - 1

	result := ""
	for n, part := range colId {
		if n == lastPart {
			// The least significant number is in the
			// range 0-25, all other numbers are 1-26,
			// hence we use a differente offset for the
			// last part.
			result += string(part + 65)
		} else {
			// Don't output leading 0s, as there is no
			// representation of 0 in this format.
			if part > 0 {
				result += string(part + 64)
			}
		}
	}
	return result
}

func smooshBase26Slice(b26 []int) []int {
	// Smoosh values together, eliminating 0s from all but the
	// least significant part.
	lastButOnePart := len(b26) - 2
	for i := lastButOnePart; i > 0; i-- {
		part := b26[i]
		if part == 0 {
			greaterPart := b26[i-1]
			if greaterPart > 0 {
				b26[i-1] = greaterPart - 1
				b26[i] = 26
			}
		}
	}
	return b26
}

func intToBase26(x int) (parts []int) {
	// Excel column codes are pure evil - in essence they're just
	// base26, but they don't represent the number 0.
	b26Denominator, _ := getLargestDenominator(x, 1, 26, 0)

	// This loop terminates because integer division of 1 / 26
	// returns 0.
	for d := b26Denominator; d > 0; d = d / 26 {
		value := x / d
		remainder := x % d
		parts = append(parts, value)
		x = remainder
	}
	return parts
}

// ColIndexToLetters is used to convert a zero based, numeric column
// indentifier into a character code.
func ColIndexToLetters(colRef int) string {
	parts := intToBase26(colRef)
	return formatColumnName(smooshBase26Slice(parts))
}

// letterOnlyMapF is used in conjunction with strings.Map to return
// only the characters A-Z and a-z in a string
func letterOnlyMapF(rune rune) rune {
	switch {
	case 'A' <= rune && rune <= 'Z':
		return rune
	case 'a' <= rune && rune <= 'z':
		return rune - 32
	}
	return -1
}

// intOnlyMapF is used in conjunction with strings.Map to return only
// the numeric portions of a string.
func intOnlyMapF(rune rune) rune {
	if rune >= 48 && rune < 58 {
		return rune
	}
	return -1
}

// GetCoordsFromCellIDString returns the zero based cartesian
// coordinates from a cell name in Excel format, e.g. the cellIDString
// "A1" returns 0, 0 and the "B3" return 1, 2.
func GetCoordsFromCellIDString(cellIDString string) (x, y int, error error) {
	var letterPart string = strings.Map(letterOnlyMapF, cellIDString)
	y, error = strconv.Atoi(strings.Map(intOnlyMapF, cellIDString))
	if error != nil {
		return x, y, error
	}
	y -= 1 // Zero based
	x = ColLettersToIndex(letterPart)
	return x, y, error
}

// GetCellIDStringFromCoords returns the Excel format cell name that
// represents a pair of zero based cartesian coordinates.
func GetCellIDStringFromCoords(x, y int) string {
	letterPart := ColIndexToLetters(x)
	numericPart := y + 1
	return fmt.Sprintf("%s%d", letterPart, numericPart)
}

// getMaxMinFromDimensionRef return the zero based cartesian maximum
// and minimum coordinates from the dimension reference embedded in a
// XLSX worksheet.  For example, the dimension reference "A1:B2"
// returns "0,0", "1,1".
func getMaxMinFromDimensionRef(ref string) (minx, miny, maxx, maxy int, err error) {
	var parts []string
	parts = strings.Split(ref, ":")
	minx, miny, err = GetCoordsFromCellIDString(parts[0])
	if err != nil {
		return -1, -1, -1, -1, err
	}
	maxx, maxy, err = GetCoordsFromCellIDString(parts[1])
	if err != nil {
		return -1, -1, -1, -1, err
	}
	return
}

// calculateMaxMinFromWorkSheet works out the dimensions of a spreadsheet
// that doesn't have a DimensionRef set.  The only case currently
// known where this is true is with XLSX exported from Google Docs.
// This is also true for XLSX files created through the streaming APIs.
func calculateMaxMinFromWorksheet(worksheet *xlsxWorksheet) (minx, miny, maxx, maxy int, err error) {
	// Note, this method could be very slow for large spreadsheets.
	var x, y int
	var maxVal int
	maxVal = int(^uint(0) >> 1)
	minx = maxVal
	miny = maxVal
	maxy = 0
	maxx = 0
	for _, row := range worksheet.SheetData.Row {
		for _, cell := range row.C {
			x, y, err = GetCoordsFromCellIDString(cell.R)
			if err != nil {
				return -1, -1, -1, -1, err
			}
			if x < minx {
				minx = x
			}
			if x > maxx {
				maxx = x
			}
			if y < miny {
				miny = y
			}
			if y > maxy {
				maxy = y
			}
		}
	}
	if minx == maxVal {
		minx = 0
	}
	if miny == maxVal {
		miny = 0
	}
	return
}

// makeRowFromSpan will, when given a span expressed as a string,
// return an empty Row large enough to encompass that span and
// populate it with empty cells.  All rows start from cell 1 -
// regardless of the lower bound of the span.
func makeRowFromSpan(spans string, sheet *Sheet) *Row {
	var error error
	var upper int
	var row *Row
	var cell *Cell

	row = new(Row)
	row.Sheet = sheet
	_, upper, error = getRangeFromString(spans)
	if error != nil {
		panic(error)
	}
	error = nil
	row.Cells = make([]*Cell, upper)
	for i := 0; i < upper; i++ {
		cell = new(Cell)
		cell.Value = ""
		row.Cells[i] = cell
	}
	return row
}

// makeRowFromRaw returns the Row representation of the xlsxRow.
func makeRowFromRaw(rawrow xlsxRow, sheet *Sheet) *Row {
	var upper int
	var row *Row
	var cell *Cell

	row = new(Row)
	row.Sheet = sheet
	upper = -1

	for _, rawcell := range rawrow.C {
		if rawcell.R != "" {
			x, _, error := GetCoordsFromCellIDString(rawcell.R)
			if error != nil {
				panic(fmt.Sprintf("Invalid Cell Coord, %s\n", rawcell.R))
			}
			if x > upper {
				upper = x
			}
			continue
		}
		upper++
	}
	upper++

	row.OutlineLevel = rawrow.OutlineLevel

	row.Cells = make([]*Cell, upper)
	for i := 0; i < upper; i++ {
		cell = new(Cell)
		cell.Value = ""
		row.Cells[i] = cell
	}
	return row
}

func makeEmptyRow(sheet *Sheet) *Row {
	row := new(Row)
	row.Cells = make([]*Cell, 0)
	row.Sheet = sheet
	return row
}

type sharedFormula struct {
	x, y    int
	formula string
}

func formulaForCell(rawcell xlsxC, sharedFormulas map[int]sharedFormula) string {
	var res string

	f := rawcell.F
	if f == nil {
		return ""
	}
	if f.T == "shared" {
		x, y, err := GetCoordsFromCellIDString(rawcell.R)
		if err != nil {
			res = f.Content
		} else {
			if f.Ref != "" {
				res = f.Content
				sharedFormulas[f.Si] = sharedFormula{x, y, res}
			} else {
				sharedFormula := sharedFormulas[f.Si]
				dx := x - sharedFormula.x
				dy := y - sharedFormula.y
				orig := []byte(sharedFormula.formula)
				var start, end int
				var stringLiteral bool
				for end = 0; end < len(orig); end++ {
					c := orig[end]

					if c == '"' {
						stringLiteral = !stringLiteral
					}

					if stringLiteral {
						continue // Skip characters in quotes
					}

					if c >= 'A' && c <= 'Z' || c == '$' {
						res += string(orig[start:end])
						start = end
						end++
						foundNum := false
						for ; end < len(orig); end++ {
							idc := orig[end]
							if idc >= '0' && idc <= '9' || idc == '$' {
								foundNum = true
							} else if idc >= 'A' && idc <= 'Z' {
								if foundNum {
									break
								}
							} else {
								break
							}
						}
						if foundNum {
							cellID := string(orig[start:end])
							res += shiftCell(cellID, dx, dy)
							start = end
						}
					}
				}
				if start < len(orig) {
					res += string(orig[start:])
				}
			}
		}
	} else {
		res = f.Content
	}
	return strings.Trim(res, " \t\n\r")
}

// shiftCell returns the cell shifted according to dx and dy taking into consideration of absolute
// references with dollar sign ($)
func shiftCell(cellID string, dx, dy int) string {
	fx, fy, _ := GetCoordsFromCellIDString(cellID)

	// Is fixed column?
	fixedCol := strings.Index(cellID, "$") == 0

	// Is fixed row?
	fixedRow := strings.LastIndex(cellID, "$") > 0

	if !fixedCol {
		// Shift column
		fx += dx
	}

	if !fixedRow {
		// Shift row
		fy += dy
	}

	// New shifted cell
	shiftedCellID := GetCellIDStringFromCoords(fx, fy)

	if !fixedCol && !fixedRow {
		return shiftedCellID
	}

	// There are absolute references, need to put the $ back into the formula.
	letterPart := strings.Map(letterOnlyMapF, shiftedCellID)
	numberPart := strings.Map(intOnlyMapF, shiftedCellID)

	result := ""

	if fixedCol {
		result += "$"
	}

	result += letterPart

	if fixedRow {
		result += "$"
	}

	result += numberPart

	return result
}

// fillCellData attempts to extract a valid value, usable in
// CSV form from the raw cell value.  Note - this is not actually
// general enough - we should support retaining tabs and newlines.
func fillCellData(rawCell xlsxC, refTable *RefTable, sharedFormulas map[int]sharedFormula, cell *Cell) {
	val := strings.Trim(rawCell.V, " \t\n\r")
	cell.formula = formulaForCell(rawCell, sharedFormulas)
	switch rawCell.T {
	case "s": // Shared String
		cell.cellType = CellTypeString
		if val != "" {
			ref, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			cell.Value = refTable.ResolveSharedString(ref)
		}
	case "inlineStr":
		cell.cellType = CellTypeInline
		fillCellDataFromInlineString(rawCell, cell)
	case "b": // Boolean
		cell.Value = val
		cell.cellType = CellTypeBool
	case "e": // Error
		cell.Value = val
		cell.cellType = CellTypeError
	case "str":
		// String Formula (special type for cells with formulas that return a string value)
		// Unlike the other string cell types, the string is stored directly in the value.
		cell.Value = val
		cell.cellType = CellTypeStringFormula
	case "d": // Date: Cell contains a date in the ISO 8601 format.
		cell.Value = val
		cell.cellType = CellTypeDate
	case "": // Numeric is the default
		fallthrough
	case "n": // Numeric
		cell.Value = val
		cell.cellType = CellTypeNumeric
	default:
		panic(errors.New("invalid cell type"))
	}
}

// fillCellDataFromInlineString attempts to get inline string data and put it into a Cell.
func fillCellDataFromInlineString(rawcell xlsxC, cell *Cell) {
	cell.Value = ""
	if rawcell.Is != nil {
		if rawcell.Is.T != "" {
			cell.Value = strings.Trim(rawcell.Is.T, " \t\n\r")
		} else {
			for _, r := range rawcell.Is.R {
				cell.Value += r.T
			}
		}
	}
}

// readRowsFromSheet is an internal helper function that extracts the
// rows from a XSLXWorksheet, populates them with Cells and resolves
// the value references from the reference table and stores them in
// the rows and columns.
func readRowsFromSheet(Worksheet *xlsxWorksheet, file *File, sheet *Sheet, rowLimit int) ([]*Row, []*Col, int, int) {
	var rows []*Row
	var cols []*Col
	var row *Row
	var minCol, maxCol, maxRow, colCount, rowCount int
	var reftable *RefTable
	var err error
	var insertRowIndex, insertColIndex int
	sharedFormulas := map[int]sharedFormula{}

	if len(Worksheet.SheetData.Row) == 0 {
		return nil, nil, 0, 0
	}
	reftable = file.referenceTable
	if len(Worksheet.Dimension.Ref) > 0 && len(strings.Split(Worksheet.Dimension.Ref, ":")) == 2 && rowLimit == NoRowLimit {
		minCol, _, maxCol, maxRow, err = getMaxMinFromDimensionRef(Worksheet.Dimension.Ref)
	} else {
		minCol, _, maxCol, maxRow, err = calculateMaxMinFromWorksheet(Worksheet)
	}
	if err != nil {
		panic(err.Error())
	}

	rowCount = maxRow + 1
	colCount = maxCol + 1
	rows = make([]*Row, rowCount)
	cols = make([]*Col, colCount)
	for i := range cols {
		cols[i] = &Col{
			Hidden: false,
		}
	}

	if Worksheet.Cols != nil {
		// Columns can apply to a range, for convenience we expand the
		// ranges out into individual column definitions.
		for _, rawcol := range Worksheet.Cols.Col {
			// Note, below, that sometimes column definitions can
			// exist outside the defined dimensions of the
			// spreadsheet - we deliberately exclude these
			// columns.
			for i := rawcol.Min; i <= rawcol.Max && i <= colCount; i++ {
				col := &Col{
					Min:          rawcol.Min,
					Max:          rawcol.Max,
					Hidden:       rawcol.Hidden,
					Width:        rawcol.Width,
					OutlineLevel: rawcol.OutlineLevel}
				cols[i-1] = col
				if file.styles != nil {
					col.style = file.styles.getStyle(rawcol.Style)
					col.numFmt, col.parsedNumFmt = file.styles.getNumberFormat(rawcol.Style)
				}
			}
		}
	}

	numRows := len(rows)
	for rowIndex := 0; rowIndex < len(Worksheet.SheetData.Row); rowIndex++ {
		rawrow := Worksheet.SheetData.Row[rowIndex]
		// Some spreadsheets will omit blank rows from the
		// stored data
		for rawrow.R > (insertRowIndex + 1) {
			// Put an empty Row into the array
			if insertRowIndex < numRows {
				rows[insertRowIndex] = makeEmptyRow(sheet)
			}
			insertRowIndex++
		}
		// range is not empty and only one range exist
		if len(rawrow.Spans) != 0 && strings.Count(rawrow.Spans, ":") == 1 {
			row = makeRowFromSpan(rawrow.Spans, sheet)
		} else {
			row = makeRowFromRaw(rawrow, sheet)
		}

		row.Hidden = rawrow.Hidden
		height, err := strconv.ParseFloat(rawrow.Ht, 64)
		if err == nil {
			row.Height = height
		}
		row.isCustom = rawrow.CustomHeight
		row.OutlineLevel = rawrow.OutlineLevel

		insertColIndex = minCol
		for _, rawcell := range rawrow.C {
			h, v, err := Worksheet.MergeCells.getExtent(rawcell.R)
			if err != nil {
				panic(err.Error())
			}
			x, _, _ := GetCoordsFromCellIDString(rawcell.R)

			// K1000000: Prevent panic when the range specified in the spreadsheet
			//           view exceeds the actual number of columns in the dataset.

			// Some spreadsheets will omit blank cells
			// from the data.
			for x > insertColIndex {
				// Put an empty Cell into the array
				if insertColIndex < len(row.Cells) {
					row.Cells[insertColIndex] = new(Cell)
				}
				insertColIndex++
			}
			cellX := insertColIndex

			if cellX < len(row.Cells) {
				cell := row.Cells[cellX]
				cell.HMerge = h
				cell.VMerge = v
				fillCellData(rawcell, reftable, sharedFormulas, cell)
				if file.styles != nil {
					cell.style = file.styles.getStyle(rawcell.S)
					cell.NumFmt, cell.parsedNumFmt = file.styles.getNumberFormat(rawcell.S)
				}
				cell.date1904 = file.Date1904
				// Cell is considered hidden if the row or the column of this cell is hidden
				cell.Hidden = rawrow.Hidden || (len(cols) > cellX && cols[cellX].Hidden)
				insertColIndex++
			}
		}
		if len(rows) > insertRowIndex {
			rows[insertRowIndex] = row
		}
		insertRowIndex++
	}

	// insert trailing empty rows for the rest of the file
	for ; insertRowIndex < rowCount; insertRowIndex++ {
		rows[insertRowIndex] = makeEmptyRow(sheet)
	}
	return rows, cols, colCount, rowCount
}

type indexedSheet struct {
	Index int
	Sheet *Sheet
	Error error
}

func readSheetViews(xSheetViews xlsxSheetViews) []SheetView {
	if xSheetViews.SheetView == nil || len(xSheetViews.SheetView) == 0 {
		return nil
	}
	sheetViews := []SheetView{}
	for _, xSheetView := range xSheetViews.SheetView {
		sheetView := SheetView{}
		if xSheetView.Pane != nil {
			xlsxPane := xSheetView.Pane
			pane := &Pane{}
			pane.XSplit = xlsxPane.XSplit
			pane.YSplit = xlsxPane.YSplit
			pane.TopLeftCell = xlsxPane.TopLeftCell
			pane.ActivePane = xlsxPane.ActivePane
			pane.State = xlsxPane.State
			sheetView.Pane = pane
		}
		sheetViews = append(sheetViews, sheetView)
	}
	return sheetViews
}

// readSheetFromFile is the logic of converting a xlsxSheet struct
// into a Sheet struct.  This work can be done in parallel and so
// readSheetsFromZipFile will spawn an instance of this function per
// sheet and get the results back on the provided channel.
func readSheetFromFile(sc chan *indexedSheet, index int, rsheet xlsxSheet, fi *File, sheetXMLMap map[string]string, rowLimit int) (errRes error) {
	result := &indexedSheet{Index: index, Sheet: nil, Error: nil}
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) {
			case error:
				result.Error = e.(error)
				errRes = e.(error)
			default:
				result.Error = errors.New("unexpected error")
			}
			// The only thing here, is if one close the channel. but its not the case
			sc <- result
		}
	}()

	worksheet, err := getWorksheetFromSheet(rsheet, fi.worksheets, sheetXMLMap, rowLimit)
	if err != nil {
		result.Error = err
		sc <- result
		return err
	}
	sheet := new(Sheet)
	sheet.File = fi
	sheet.Rows, sheet.Cols, sheet.MaxCol, sheet.MaxRow = readRowsFromSheet(worksheet, fi, sheet, rowLimit)
	sheet.Hidden = rsheet.State == sheetStateHidden || rsheet.State == sheetStateVeryHidden
	sheet.SheetViews = readSheetViews(worksheet.SheetViews)

	sheet.SheetFormat.DefaultColWidth = worksheet.SheetFormatPr.DefaultColWidth
	sheet.SheetFormat.DefaultRowHeight = worksheet.SheetFormatPr.DefaultRowHeight
	sheet.SheetFormat.OutlineLevelCol = worksheet.SheetFormatPr.OutlineLevelCol
	sheet.SheetFormat.OutlineLevelRow = worksheet.SheetFormatPr.OutlineLevelRow

	result.Sheet = sheet
	sc <- result
	return nil
}

// readSheetsFromZipFile is an internal helper function that loops
// over the Worksheets defined in the XSLXWorkbook and loads them into
// Sheet objects stored in the Sheets slice of a xlsx.File struct.
func readSheetsFromZipFile(f *zip.File, file *File, sheetXMLMap map[string]string, rowLimit int) (map[string]*Sheet, []*Sheet, error) {
	var workbook *xlsxWorkbook
	var err error
	var rc io.ReadCloser
	var decoder *xml.Decoder
	var sheetCount int
	workbook = new(xlsxWorkbook)
	rc, err = f.Open()
	if err != nil {
		return nil, nil, err
	}
	decoder = xml.NewDecoder(rc)
	err = decoder.Decode(workbook)
	if err != nil {
		return nil, nil, err
	}
	file.Date1904 = workbook.WorkbookPr.Date1904

	for entryNum := range workbook.DefinedNames.DefinedName {
		file.DefinedNames = append(file.DefinedNames, &workbook.DefinedNames.DefinedName[entryNum])
	}

	// Only try and read sheets that have corresponding files.
	// Notably this excludes chartsheets don't right now
	var workbookSheets []xlsxSheet
	for _, sheet := range workbook.Sheets.Sheet {
		if f := worksheetFileForSheet(sheet, file.worksheets, sheetXMLMap); f != nil {
			workbookSheets = append(workbookSheets, sheet)
		}
	}
	sheetCount = len(workbookSheets)
	sheetsByName := make(map[string]*Sheet, sheetCount)
	sheets := make([]*Sheet, sheetCount)
	sheetChan := make(chan *indexedSheet, sheetCount)

	go func() {
		defer close(sheetChan)
		err = nil
		for i, rawsheet := range workbookSheets {
			if err := readSheetFromFile(sheetChan, i, rawsheet, file, sheetXMLMap, rowLimit); err != nil {
				return
			}
		}
	}()

	for j := 0; j < sheetCount; j++ {
		sheet := <-sheetChan
		if sheet.Error != nil {
			return nil, nil, sheet.Error
		}
		sheetName := workbookSheets[sheet.Index].Name
		sheetsByName[sheetName] = sheet.Sheet
		sheet.Sheet.Name = sheetName
		sheets[sheet.Index] = sheet.Sheet
	}
	return sheetsByName, sheets, nil
}

// readSharedStringsFromZipFile() is an internal helper function to
// extract a reference table from the sharedStrings.xml file within
// the XLSX zip file.
func readSharedStringsFromZipFile(f *zip.File) (*RefTable, error) {
	var sst *xlsxSST
	var error error
	var rc io.ReadCloser
	var decoder *xml.Decoder
	var reftable *RefTable

	// In a file with no strings it's possible that
	// sharedStrings.xml doesn't exist.  In this case the value
	// passed as f will be nil.
	if f == nil {
		return nil, nil
	}
	rc, error = f.Open()
	if error != nil {
		return nil, error
	}
	sst = new(xlsxSST)
	decoder = xml.NewDecoder(rc)
	error = decoder.Decode(sst)
	if error != nil {
		return nil, error
	}
	reftable = MakeSharedStringRefTable(sst)
	return reftable, nil
}

// readStylesFromZipFile() is an internal helper function to
// extract a style table from the style.xml file within
// the XLSX zip file.
func readStylesFromZipFile(f *zip.File, theme *theme) (*xlsxStyleSheet, error) {
	var style *xlsxStyleSheet
	var error error
	var rc io.ReadCloser
	var decoder *xml.Decoder
	rc, error = f.Open()
	if error != nil {
		return nil, error
	}
	style = newXlsxStyleSheet(theme)
	decoder = xml.NewDecoder(rc)
	error = decoder.Decode(style)
	if error != nil {
		return nil, error
	}
	buildNumFmtRefTable(style)
	return style, nil
}

func buildNumFmtRefTable(style *xlsxStyleSheet) {
	for _, numFmt := range style.NumFmts.NumFmt {
		// We do this for the side effect of populating the NumFmtRefTable.
		style.addNumFmt(numFmt)
	}
}

func readThemeFromZipFile(f *zip.File) (*theme, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}

	var themeXml xlsxTheme
	err = xml.NewDecoder(rc).Decode(&themeXml)
	if err != nil {
		return nil, err
	}

	return newTheme(themeXml), nil
}

type WorkBookRels map[string]string

func (w *WorkBookRels) MakeXLSXWorkbookRels() xlsxWorkbookRels {
	relCount := len(*w)
	xWorkbookRels := xlsxWorkbookRels{}
	xWorkbookRels.Relationships = make([]xlsxWorkbookRelation, relCount+3)
	for k, v := range *w {
		index, err := strconv.Atoi(k[3:])
		if err != nil {
			panic(err.Error())
		}
		xWorkbookRels.Relationships[index-1] = xlsxWorkbookRelation{
			Id:     k,
			Target: v,
			Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet"}
	}

	relCount++
	sheetId := fmt.Sprintf("rId%d", relCount)
	xWorkbookRels.Relationships[relCount-1] = xlsxWorkbookRelation{
		Id:     sheetId,
		Target: "sharedStrings.xml",
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings"}

	relCount++
	sheetId = fmt.Sprintf("rId%d", relCount)
	xWorkbookRels.Relationships[relCount-1] = xlsxWorkbookRelation{
		Id:     sheetId,
		Target: "theme/theme1.xml",
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"}

	relCount++
	sheetId = fmt.Sprintf("rId%d", relCount)
	xWorkbookRels.Relationships[relCount-1] = xlsxWorkbookRelation{
		Id:     sheetId,
		Target: "styles.xml",
		Type:   "http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"}

	return xWorkbookRels
}

// readWorkbookRelationsFromZipFile is an internal helper function to
// extract a map of relationship ID strings to the name of the
// worksheet.xml file they refer to.  The resulting map can be used to
// reliably derefence the worksheets in the XLSX file.
func readWorkbookRelationsFromZipFile(workbookRels *zip.File) (WorkBookRels, error) {
	var sheetXMLMap WorkBookRels
	var wbRelationships *xlsxWorkbookRels
	var rc io.ReadCloser
	var decoder *xml.Decoder
	var err error

	rc, err = workbookRels.Open()
	if err != nil {
		return nil, err
	}
	decoder = xml.NewDecoder(rc)
	wbRelationships = new(xlsxWorkbookRels)
	err = decoder.Decode(wbRelationships)
	if err != nil {
		return nil, err
	}
	sheetXMLMap = make(WorkBookRels)
	for _, rel := range wbRelationships.Relationships {
		if strings.HasSuffix(rel.Target, ".xml") && rel.Type == "http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" {
			_, filename := path.Split(rel.Target)
			sheetXMLMap[rel.Id] = strings.Replace(filename, ".xml", "", 1)
		}
	}
	return sheetXMLMap, nil
}

// ReadZip() takes a pointer to a zip.ReadCloser and returns a
// xlsx.File struct populated with its contents.  In most cases
// ReadZip is not used directly, but is called internally by OpenFile.
func ReadZip(f *zip.ReadCloser) (*File, error) {
	return ReadZipWithRowLimit(f, NoRowLimit)
}

// ReadZipWithRowLimit() takes a pointer to a zip.ReadCloser and returns a
// xlsx.File struct populated with its contents.  In most cases
// ReadZip is not used directly, but is called internally by OpenFile.
func ReadZipWithRowLimit(f *zip.ReadCloser, rowLimit int) (*File, error) {
	defer f.Close()
	return ReadZipReaderWithRowLimit(&f.Reader, rowLimit)
}

// ReadZipReader() can be used to read an XLSX in memory without
// touching the filesystem.
func ReadZipReader(r *zip.Reader) (*File, error) {
	return ReadZipReaderWithRowLimit(r, NoRowLimit)
}

// ReadZipReaderWithRowLimit() can be used to read an XLSX in memory without
// touching the filesystem.
// rowLimit is the number of rows that should be read from the file. If rowLimit is -1, no limit is applied.
// You can specify this with the constant NoRowLimit.
func ReadZipReaderWithRowLimit(r *zip.Reader, rowLimit int) (*File, error) {
	var err error
	var file *File
	var reftable *RefTable
	var sharedStrings *zip.File
	var sheetXMLMap map[string]string
	var sheetsByName map[string]*Sheet
	var sheets []*Sheet
	var style *xlsxStyleSheet
	var styles *zip.File
	var themeFile *zip.File
	var v *zip.File
	var workbook *zip.File
	var workbookRels *zip.File
	var worksheets map[string]*zip.File

	file = NewFile()
	// file.numFmtRefTable = make(map[int]xlsxNumFmt, 1)
	worksheets = make(map[string]*zip.File, len(r.File))
	for _, v = range r.File {
		switch v.Name {
		case "xl/sharedStrings.xml":
			sharedStrings = v
		case "xl/workbook.xml":
			workbook = v
		case "xl/_rels/workbook.xml.rels":
			workbookRels = v
		case "xl/styles.xml":
			styles = v
		case "xl/theme/theme1.xml":
			themeFile = v
		default:
			if len(v.Name) > 17 {
				if v.Name[0:13] == "xl/worksheets" {
					worksheets[v.Name[14:len(v.Name)-4]] = v
				}
			}
		}
	}
	if workbookRels == nil {
		return nil, fmt.Errorf("xl/_rels/workbook.xml.rels not found in input xlsx.")
	}
	sheetXMLMap, err = readWorkbookRelationsFromZipFile(workbookRels)
	if err != nil {
		return nil, err
	}
	if len(worksheets) == 0 {
		return nil, fmt.Errorf("Input xlsx contains no worksheets.")
	}
	file.worksheets = worksheets
	reftable, err = readSharedStringsFromZipFile(sharedStrings)
	if err != nil {
		return nil, err
	}
	file.referenceTable = reftable
	if themeFile != nil {
		theme, err := readThemeFromZipFile(themeFile)
		if err != nil {
			return nil, err
		}

		file.theme = theme
	}
	if styles != nil {
		style, err = readStylesFromZipFile(styles, file.theme)
		if err != nil {
			return nil, err
		}

		file.styles = style
	}
	sheetsByName, sheets, err = readSheetsFromZipFile(workbook, file, sheetXMLMap, rowLimit)
	if err != nil {
		return nil, err
	}
	if sheets == nil {
		readerErr := new(XLSXReaderError)
		readerErr.Err = "No sheets found in XLSX File"
		return nil, readerErr
	}
	file.Sheet = sheetsByName
	file.Sheets = sheets
	return file, nil
}

// truncateSheetXML will take in a reader to an XML sheet file and will return a reader that will read an equivalent
// XML sheet file with only the number of rows specified. This greatly speeds up XML unmarshalling when only
// a few rows need to be read from a large sheet.
// When sheets are truncated, all formatting present after the sheetData tag will be lost, but all of this formatting
// is related to printing and visibility, and is out of scope for most purposes of this library.
func truncateSheetXML(r io.Reader, rowLimit int) (io.Reader, error) {
	var rowCount int
	var token xml.Token
	var readErr error

	output := new(bytes.Buffer)
	r = io.TeeReader(r, output)
	decoder := xml.NewDecoder(r)

	for {
		token, readErr = decoder.Token()
		if readErr == io.EOF {
			break
		} else if readErr != nil {
			return nil, readErr
		}
		end, ok := token.(xml.EndElement)
		if ok && end.Name.Local == "row" {
			rowCount++
			if rowCount >= rowLimit {
				break
			}
		}
	}

	offset := decoder.InputOffset()
	output.Truncate(int(offset))

	if readErr != io.EOF {
		_, err := output.Write([]byte(sheetEnding))
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}
