// Authors: Ryan Hollis (ryanh@)

// The purpose of StreamFileBuilder and StreamFile is to allow streamed writing of XLSX files.
// Directions:
// 1. Create a StreamFileBuilder with NewStreamFileBuilder() or NewStreamFileBuilderForPath().
// 2. Add the sheets and their first row of data by calling AddSheet().
// 3. Call Build() to get a StreamFile. Once built, all functions on the builder will return an error.
// 4. Write to the StreamFile with Write(). Writes begin on the first sheet. New rows are always written and flushed
// to the io. All rows written to the same sheet must have the same number of cells as the header provided when the sheet
// was created or an error will be returned.
// 5. Call NextSheet() to proceed to the next sheet. Once NextSheet() is called, the previous sheet can not be edited.
// 6. Call Close() to finish.

// Future work suggestions:
// Currently the only supported cell type is string, since the main reason this library was written was to prevent
// strings from being interpreted as numbers. It would be nice to have support for numbers and money so that the exported
// files could better take advantage of XLSX's features.
// All text is written with the same text style. Support for additional text styles could be added to highlight certain
// data in the file.
// The current default style uses fonts that are not on Macs by default so opening the XLSX files in Numbers causes a
// pop up that says there are missing fonts. The font could be changed to something that is usually found on Mac and PC.

package xlsx

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type StreamFileBuilder struct {
	built              bool
	xlsxFile           *File
	zipWriter          *zip.Writer
	cellTypeToStyleIds map[CellType]int
	maxStyleId         int
	styleIds           [][]int
}

const (
	sheetFilePathPrefix = "xl/worksheets/sheet"
	sheetFilePathSuffix = ".xml"
	endSheetDataTag     = "</sheetData>"
	dimensionTag        = `<dimension ref="%s"></dimension>`
	// This is the index of the max style that this library will insert into XLSX sheets by default.
	// This allows us to predict what the style id of styles that we add will be.
	// TestXlsxStyleBehavior tests that this behavior continues to be what we expect.
	initMaxStyleId = 1
)

var BuiltStreamFileBuilderError = errors.New("StreamFileBuilder has already been built, functions may no longer be used")

// NewStreamFileBuilder creates an StreamFileBuilder that will write to the the provided io.writer
func NewStreamFileBuilder(writer io.Writer) *StreamFileBuilder {
	return &StreamFileBuilder{
		zipWriter:          zip.NewWriter(writer),
		xlsxFile:           NewFile(),
		cellTypeToStyleIds: make(map[CellType]int),
		maxStyleId:         initMaxStyleId,
	}
}

// NewStreamFileBuilderForPath takes the name of an XLSX file and returns a builder for it.
// The file will be created if it does not exist, or truncated if it does.
func NewStreamFileBuilderForPath(path string) (*StreamFileBuilder, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return NewStreamFileBuilder(file), nil
}

// AddSheet will add sheets with the given name with the provided headers. The headers cannot be edited later, and all
// rows written to the sheet must contain the same number of cells as the header. Sheet names must be unique, or an
// error will be thrown.
func (sb *StreamFileBuilder) AddSheet(name string, headers []string, cellTypes []*CellType) error {
	if sb.built {
		return BuiltStreamFileBuilderError
	}
	if len(cellTypes) > len(headers) {
		return errors.New("cellTypes is longer than headers")
	}
	sheet, err := sb.xlsxFile.AddSheet(name)
	if err != nil {
		// Set built on error so that all subsequent calls to the builder will also fail.
		sb.built = true
		return err
	}
	sb.styleIds = append(sb.styleIds, []int{})
	row := sheet.AddRow()
	if count := row.WriteSlice(&headers, -1); count != len(headers) {
		// Set built on error so that all subsequent calls to the builder will also fail.
		sb.built = true
		return errors.New("failed to write headers")
	}
	for i, cellType := range cellTypes {
		var cellStyleIndex int
		var ok bool
		if cellType != nil {
			// The cell type is one of the attributes of a Style.
			// Since it is the only attribute of Style that we use, we can assume that cell types
			// map one to one with Styles and their Style ID.
			// If a new cell type is used, a new style gets created with an increased id, if an existing cell type is
			// used, the pre-existing style will also be used.
			cellStyleIndex, ok = sb.cellTypeToStyleIds[*cellType]
			if !ok {
				sb.maxStyleId++
				cellStyleIndex = sb.maxStyleId
				sb.cellTypeToStyleIds[*cellType] = sb.maxStyleId
			}
			sheet.Cols[i].SetType(*cellType)
		}
		sb.styleIds[len(sb.styleIds)-1] = append(sb.styleIds[len(sb.styleIds)-1], cellStyleIndex)
	}
	return nil
}

// Build begins streaming the XLSX file to the io, by writing all the XLSX metadata. It creates a StreamFile struct
// that can be used to write the rows to the sheets.
func (sb *StreamFileBuilder) Build() (*StreamFile, error) {
	if sb.built {
		return nil, BuiltStreamFileBuilderError
	}
	sb.built = true
	parts, err := sb.xlsxFile.MarshallParts()
	if err != nil {
		return nil, err
	}
	es := &StreamFile{
		zipWriter:      sb.zipWriter,
		xlsxFile:       sb.xlsxFile,
		sheetXmlPrefix: make([]string, len(sb.xlsxFile.Sheets)),
		sheetXmlSuffix: make([]string, len(sb.xlsxFile.Sheets)),
		styleIds:       sb.styleIds,
	}
	for path, data := range parts {
		// If the part is a sheet, don't write it yet. We only want to write the XLSX metadata files, since at this
		// point the sheets are still empty. The sheet files will be written later as their rows come in.
		if strings.HasPrefix(path, sheetFilePathPrefix) {
			if err := sb.processEmptySheetXML(es, path, data); err != nil {
				return nil, err
			}
			continue
		}
		metadataFile, err := sb.zipWriter.Create(path)
		if err != nil {
			return nil, err
		}
		_, err = metadataFile.Write([]byte(data))
		if err != nil {
			return nil, err
		}
	}

	if err := es.NextSheet(); err != nil {
		return nil, err
	}
	return es, nil
}

// processEmptySheetXML will take in the path and XML data of an empty sheet, and will save the beginning and end of the
// XML file so that these can be written at the right time.
func (sb *StreamFileBuilder) processEmptySheetXML(sf *StreamFile, path, data string) error {
	// Get the sheet index from the path
	sheetIndex, err := getSheetIndex(sf, path)
	if err != nil {
		return err
	}

	// Remove the Dimension tag. Since more rows are going to be written to the sheet, it will be wrong.
	// It is valid to for a sheet to be missing a Dimension tag, but it is not valid for it to be wrong.
	data, err = removeDimensionTag(data, sf.xlsxFile.Sheets[sheetIndex])
	if err != nil {
		return err
	}

	// Split the sheet at the end of its SheetData tag so that more rows can be added inside.
	prefix, suffix, err := splitSheetIntoPrefixAndSuffix(data)
	if err != nil {
		return err
	}
	sf.sheetXmlPrefix[sheetIndex] = prefix
	sf.sheetXmlSuffix[sheetIndex] = suffix
	return nil
}

// getSheetIndex parses the path to the XLSX sheet data and returns the index
// The files that store the data for each sheet must have the format:
// xl/worksheets/sheet123.xml
// where 123 is the index of the sheet. This file path format is part of the XLSX file standard.
func getSheetIndex(sf *StreamFile, path string) (int, error) {
	indexString := path[len(sheetFilePathPrefix) : len(path)-len(sheetFilePathSuffix)]
	sheetXLSXIndex, err := strconv.Atoi(indexString)
	if err != nil {
		return -1, errors.New("Unexpected sheet file name from xlsx package")
	}
	if sheetXLSXIndex < 1 || len(sf.sheetXmlPrefix) < sheetXLSXIndex ||
		len(sf.sheetXmlSuffix) < sheetXLSXIndex || len(sf.xlsxFile.Sheets) < sheetXLSXIndex {
		return -1, errors.New("Unexpected sheet index")
	}
	sheetArrayIndex := sheetXLSXIndex - 1
	return sheetArrayIndex, nil
}

// removeDimensionTag will return the passed in XLSX Spreadsheet XML with the dimension tag removed.
// data is the XML data for the sheet
// sheet is the Sheet struct that the XML was created from.
// Can return an error if the XML's dimension tag does not match was is expected based on the provided Sheet
func removeDimensionTag(data string, sheet *Sheet) (string, error) {
	x := len(sheet.Cols) - 1
	y := len(sheet.Rows) - 1
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	var dimensionRef string
	if x == 0 && y == 0 {
		dimensionRef = "A1"
	} else {
		endCoordinate := GetCellIDStringFromCoords(x, y)
		dimensionRef = "A1:" + endCoordinate
	}
	dataParts := strings.Split(data, fmt.Sprintf(dimensionTag, dimensionRef))
	if len(dataParts) != 2 {
		return "", errors.New("unexpected Sheet XML: dimension tag not found")
	}
	return dataParts[0] + dataParts[1], nil
}

// splitSheetIntoPrefixAndSuffix will split the provided XML sheet into a prefix and a suffix so that
// more spreadsheet rows can be inserted in between.
func splitSheetIntoPrefixAndSuffix(data string) (string, string, error) {
	// Split the sheet at the end of its SheetData tag so that more rows can be added inside.
	sheetParts := strings.Split(data, endSheetDataTag)
	if len(sheetParts) != 2 {
		return "", "", errors.New("unexpected Sheet XML: SheetData close tag not found")
	}
	return sheetParts[0], sheetParts[1], nil
}
