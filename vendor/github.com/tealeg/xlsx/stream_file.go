package xlsx

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

type StreamFile struct {
	xlsxFile       *File
	sheetXmlPrefix []string
	sheetXmlSuffix []string
	zipWriter      *zip.Writer
	currentSheet   *streamSheet
	styleIds       [][]int
	err            error
}

type streamSheet struct {
	// sheetIndex is the XLSX sheet index, which starts at 1
	index int
	// The number of rows that have been written to the sheet so far
	rowCount int
	// The number of columns in the sheet
	columnCount int
	// The writer to write to this sheet's file in the XLSX Zip file
	writer   io.Writer
	styleIds []int
}

var (
	NoCurrentSheetError     = errors.New("no Current Sheet")
	WrongNumberOfRowsError  = errors.New("invalid number of cells passed to Write. All calls to Write on the same sheet must have the same number of cells")
	AlreadyOnLastSheetError = errors.New("NextSheet() called, but already on last sheet")
)

// Write will write a row of cells to the current sheet. Every call to Write on the same sheet must contain the
// same number of cells as the header provided when the sheet was created or an error will be returned. This function
// will always trigger a flush on success. Currently the only supported data type is string data.
func (sf *StreamFile) Write(cells []string) error {
	if sf.err != nil {
		return sf.err
	}
	err := sf.write(cells)
	if err != nil {
		sf.err = err
		return err
	}
	return sf.zipWriter.Flush()
}

func (sf *StreamFile) WriteAll(records [][]string) error {
	if sf.err != nil {
		return sf.err
	}
	for _, row := range records {
		err := sf.write(row)
		if err != nil {
			sf.err = err
			return err
		}
	}
	return sf.zipWriter.Flush()
}

func (sf *StreamFile) write(cells []string) error {
	if sf.currentSheet == nil {
		return NoCurrentSheetError
	}
	if len(cells) != sf.currentSheet.columnCount {
		return WrongNumberOfRowsError
	}
	sf.currentSheet.rowCount++
	if err := sf.currentSheet.write(`<row r="` + strconv.Itoa(sf.currentSheet.rowCount) + `">`); err != nil {
		return err
	}
	for colIndex, cellData := range cells {
		// documentation for the c.t (cell.Type) attribute:
		// b (Boolean): Cell containing a boolean.
		// d (Date): Cell contains a date in the ISO 8601 format.
		// e (Error): Cell containing an error.
		// inlineStr (Inline String): Cell containing an (inline) rich string, i.e., one not in the shared string table.
		// If this cell type is used, then the cell value is in the is element rather than the v element in the cell (c element).
		// n (Number): Cell containing a number.
		// s (Shared String): Cell containing a shared string.
		// str (String): Cell containing a formula string.
		cellCoordinate := GetCellIDStringFromCoords(colIndex, sf.currentSheet.rowCount-1)
		cellType := "inlineStr"
		cellOpen := `<c r="` + cellCoordinate + `" t="` + cellType + `"`
		// Add in the style id if the cell isn't using the default style
		if colIndex < len(sf.currentSheet.styleIds) && sf.currentSheet.styleIds[colIndex] != 0 {
			cellOpen += ` s="` + strconv.Itoa(sf.currentSheet.styleIds[colIndex]) + `"`
		}
		cellOpen += `><is><t>`
		cellClose := `</t></is></c>`

		if err := sf.currentSheet.write(cellOpen); err != nil {
			return err
		}
		if err := xml.EscapeText(sf.currentSheet.writer, []byte(cellData)); err != nil {
			return err
		}
		if err := sf.currentSheet.write(cellClose); err != nil {
			return err
		}
	}
	if err := sf.currentSheet.write(`</row>`); err != nil {
		return err
	}
	return sf.zipWriter.Flush()
}

// Error reports any error that has occurred during a previous Write or Flush.
func (sf *StreamFile) Error() error {
	return sf.err
}

func (sf *StreamFile) Flush() {
	if sf.err != nil {
		sf.err = sf.zipWriter.Flush()
	}
}

// NextSheet will switch to the next sheet. Sheets are selected in the same order they were added.
// Once you leave a sheet, you cannot return to it.
func (sf *StreamFile) NextSheet() error {
	if sf.err != nil {
		return sf.err
	}
	var sheetIndex int
	if sf.currentSheet != nil {
		if sf.currentSheet.index >= len(sf.xlsxFile.Sheets) {
			sf.err = AlreadyOnLastSheetError
			return AlreadyOnLastSheetError
		}
		if err := sf.writeSheetEnd(); err != nil {
			sf.currentSheet = nil
			sf.err = err
			return err
		}
		sheetIndex = sf.currentSheet.index
	}
	sheetIndex++
	sf.currentSheet = &streamSheet{
		index:       sheetIndex,
		columnCount: len(sf.xlsxFile.Sheets[sheetIndex-1].Cols),
		styleIds:    sf.styleIds[sheetIndex-1],
		rowCount:    1,
	}
	sheetPath := sheetFilePathPrefix + strconv.Itoa(sf.currentSheet.index) + sheetFilePathSuffix
	fileWriter, err := sf.zipWriter.Create(sheetPath)
	if err != nil {
		sf.err = err
		return err
	}
	sf.currentSheet.writer = fileWriter

	if err := sf.writeSheetStart(); err != nil {
		sf.err = err
		return err
	}
	return nil
}

// Close closes the Stream File.
// Any sheets that have not yet been written to will have an empty sheet created for them.
func (sf *StreamFile) Close() error {
	if sf.err != nil {
		return sf.err
	}
	// If there are sheets that have not been written yet, call NextSheet() which will add files to the zip for them.
	// XLSX readers may error if the sheets registered in the metadata are not present in the file.
	if sf.currentSheet != nil {
		for sf.currentSheet.index < len(sf.xlsxFile.Sheets) {
			if err := sf.NextSheet(); err != nil {
				sf.err = err
				return err
			}
		}
		// Write the end of the last sheet.
		if err := sf.writeSheetEnd(); err != nil {
			sf.err = err
			return err
		}
	}
	err := sf.zipWriter.Close()
	if err != nil {
		sf.err = err
	}
	return err
}

// writeSheetStart will write the start of the Sheet's XML
func (sf *StreamFile) writeSheetStart() error {
	if sf.currentSheet == nil {
		return NoCurrentSheetError
	}
	return sf.currentSheet.write(sf.sheetXmlPrefix[sf.currentSheet.index-1])
}

// writeSheetEnd will write the end of the Sheet's XML
func (sf *StreamFile) writeSheetEnd() error {
	if sf.currentSheet == nil {
		return NoCurrentSheetError
	}
	if err := sf.currentSheet.write(endSheetDataTag); err != nil {
		return err
	}
	return sf.currentSheet.write(sf.sheetXmlSuffix[sf.currentSheet.index-1])
}

func (ss *streamSheet) write(data string) error {
	_, err := ss.writer.Write([]byte(data))
	return err
}
