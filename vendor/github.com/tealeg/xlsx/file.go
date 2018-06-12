package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

// File is a high level structure providing a slice of Sheet structs
// to the user.
type File struct {
	worksheets     map[string]*zip.File
	referenceTable *RefTable
	Date1904       bool
	styles         *xlsxStyleSheet
	Sheets         []*Sheet
	Sheet          map[string]*Sheet
	theme          *theme
	DefinedNames   []*xlsxDefinedName
}

const NoRowLimit int = -1

// Create a new File
func NewFile() *File {
	return &File{
		Sheet:        make(map[string]*Sheet),
		Sheets:       make([]*Sheet, 0),
		DefinedNames: make([]*xlsxDefinedName, 0),
	}
}

// OpenFile() take the name of an XLSX file and returns a populated
// xlsx.File struct for it.
func OpenFile(fileName string) (file *File, err error) {
	return OpenFileWithRowLimit(fileName, NoRowLimit)
}

// OpenFileWithRowLimit() will open the file, but will only read the specified number of rows.
// If you save this file, it will be truncated to the number of rows specified.
func OpenFileWithRowLimit(fileName string, rowLimit int) (file *File, err error) {
	var z *zip.ReadCloser
	z, err = zip.OpenReader(fileName)
	if err != nil {
		return nil, err
	}
	return ReadZipWithRowLimit(z, rowLimit)
}

// OpenBinary() take bytes of an XLSX file and returns a populated
// xlsx.File struct for it.
func OpenBinary(bs []byte) (*File, error) {
	return OpenBinaryWithRowLimit(bs, NoRowLimit)
}

// OpenBinaryWithRowLimit() take bytes of an XLSX file and returns a populated
// xlsx.File struct for it.
func OpenBinaryWithRowLimit(bs []byte, rowLimit int) (*File, error) {
	r := bytes.NewReader(bs)
	return OpenReaderAtWithRowLimit(r, int64(r.Len()), rowLimit)
}

// OpenReaderAt() take io.ReaderAt of an XLSX file and returns a populated
// xlsx.File struct for it.
func OpenReaderAt(r io.ReaderAt, size int64) (*File, error) {
	return OpenReaderAtWithRowLimit(r, size, NoRowLimit)
}

// OpenReaderAtWithRowLimit() take io.ReaderAt of an XLSX file and returns a populated
// xlsx.File struct for it.
func OpenReaderAtWithRowLimit(r io.ReaderAt, size int64, rowLimit int) (*File, error) {
	file, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	return ReadZipReaderWithRowLimit(file, rowLimit)
}

// A convenient wrapper around File.ToSlice, FileToSlice will
// return the raw data contained in an Excel XLSX file as three
// dimensional slice.  The first index represents the sheet number,
// the second the row number, and the third the cell number.
//
// For example:
//
//    var mySlice [][][]string
//    var value string
//    mySlice = xlsx.FileToSlice("myXLSX.xlsx")
//    value = mySlice[0][0][0]
//
// Here, value would be set to the raw value of the cell A1 in the
// first sheet in the XLSX file.
func FileToSlice(path string) ([][][]string, error) {
	f, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	return f.ToSlice()
}

// FileToSliceUnmerged is a wrapper around File.ToSliceUnmerged.
// It returns the raw data contained in an Excel XLSX file as three
// dimensional slice. Merged cells will be unmerged. Covered cells become the
// values of theirs origins.
func FileToSliceUnmerged(path string) ([][][]string, error) {
	f, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	return f.ToSliceUnmerged()
}

// Save the File to an xlsx file at the provided path.
func (f *File) Save(path string) (err error) {
	target, err := os.Create(path)
	if err != nil {
		return err
	}
	err = f.Write(target)
	if err != nil {
		return err
	}
	return target.Close()
}

// Write the File to io.Writer as xlsx
func (f *File) Write(writer io.Writer) (err error) {
	parts, err := f.MarshallParts()
	if err != nil {
		return
	}
	zipWriter := zip.NewWriter(writer)
	for partName, part := range parts {
		w, err := zipWriter.Create(partName)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(part))
		if err != nil {
			return err
		}
	}
	return zipWriter.Close()
}

// Add a new Sheet, with the provided name, to a File
func (f *File) AddSheet(sheetName string) (*Sheet, error) {
	if _, exists := f.Sheet[sheetName]; exists {
		return nil, fmt.Errorf("duplicate sheet name '%s'.", sheetName)
	}
	if utf8.RuneCountInString(sheetName) >= 31 {
		return nil, fmt.Errorf("sheet name must be less than 31 characters long.  It is currently '%d' characters long", utf8.RuneCountInString(sheetName))
	}
	sheet := &Sheet{
		Name:     sheetName,
		File:     f,
		Selected: len(f.Sheets) == 0,
	}
	f.Sheet[sheetName] = sheet
	f.Sheets = append(f.Sheets, sheet)
	return sheet, nil
}

// Appends an existing Sheet, with the provided name, to a File
func (f *File) AppendSheet(sheet Sheet, sheetName string) (*Sheet, error) {
	if _, exists := f.Sheet[sheetName]; exists {
		return nil, fmt.Errorf("duplicate sheet name '%s'.", sheetName)
	}
	sheet.Name = sheetName
	sheet.File = f
	sheet.Selected = len(f.Sheets) == 0
	f.Sheet[sheetName] = &sheet
	f.Sheets = append(f.Sheets, &sheet)
	return &sheet, nil
}

func (f *File) makeWorkbook() xlsxWorkbook {
	return xlsxWorkbook{
		FileVersion: xlsxFileVersion{AppName: "Go XLSX"},
		WorkbookPr:  xlsxWorkbookPr{ShowObjects: "all"},
		BookViews: xlsxBookViews{
			WorkBookView: []xlsxWorkBookView{
				{
					ShowHorizontalScroll: true,
					ShowSheetTabs:        true,
					ShowVerticalScroll:   true,
					TabRatio:             204,
					WindowHeight:         8192,
					WindowWidth:          16384,
					XWindow:              "0",
					YWindow:              "0",
				},
			},
		},
		Sheets: xlsxSheets{Sheet: make([]xlsxSheet, len(f.Sheets))},
		CalcPr: xlsxCalcPr{
			IterateCount: 100,
			RefMode:      "A1",
			Iterate:      false,
			IterateDelta: 0.001,
		},
	}
}

// Some tools that read XLSX files have very strict requirements about
// the structure of the input XML.  In particular both Numbers on the Mac
// and SAS dislike inline XML namespace declarations, or namespace
// prefixes that don't match the ones that Excel itself uses.  This is a
// problem because the Go XML library doesn't multiple namespace
// declarations in a single element of a document.  This function is a
// horrible hack to fix that after the XML marshalling is completed.
func replaceRelationshipsNameSpace(workbookMarshal string) string {
	newWorkbook := strings.Replace(workbookMarshal, `xmlns:relationships="http://schemas.openxmlformats.org/officeDocument/2006/relationships" relationships:id`, `r:id`, -1)
	// Dirty hack to fix issues #63 and #91; encoding/xml currently
	// "doesn't allow for additional namespaces to be defined in the
	// root element of the document," as described by @tealeg in the
	// comments for #63.
	oldXmlns := `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">`
	newXmlns := `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`
	return strings.Replace(newWorkbook, oldXmlns, newXmlns, 1)
}

// Construct a map of file name to XML content representing the file
// in terms of the structure of an XLSX file.
func (f *File) MarshallParts() (map[string]string, error) {
	var parts map[string]string
	var refTable *RefTable = NewSharedStringRefTable()
	refTable.isWrite = true
	var workbookRels WorkBookRels = make(WorkBookRels)
	var err error
	var workbook xlsxWorkbook
	var types xlsxTypes = MakeDefaultContentTypes()

	marshal := func(thing interface{}) (string, error) {
		body, err := xml.Marshal(thing)
		if err != nil {
			return "", err
		}
		return xml.Header + string(body), nil
	}

	parts = make(map[string]string)
	workbook = f.makeWorkbook()
	sheetIndex := 1

	if f.styles == nil {
		f.styles = newXlsxStyleSheet(f.theme)
	}
	f.styles.reset()
	if len(f.Sheets) == 0 {
		err := errors.New("Workbook must contains atleast one worksheet")
		return nil, err
	}
	for _, sheet := range f.Sheets {
		xSheet := sheet.makeXLSXSheet(refTable, f.styles)
		rId := fmt.Sprintf("rId%d", sheetIndex)
		sheetId := strconv.Itoa(sheetIndex)
		sheetPath := fmt.Sprintf("worksheets/sheet%d.xml", sheetIndex)
		partName := "xl/" + sheetPath
		types.Overrides = append(
			types.Overrides,
			xlsxOverride{
				PartName:    "/" + partName,
				ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"})
		workbookRels[rId] = sheetPath
		workbook.Sheets.Sheet[sheetIndex-1] = xlsxSheet{
			Name:    sheet.Name,
			SheetId: sheetId,
			Id:      rId,
			State:   "visible"}
		parts[partName], err = marshal(xSheet)
		if err != nil {
			return parts, err
		}
		sheetIndex++
	}

	workbookMarshal, err := marshal(workbook)
	if err != nil {
		return parts, err
	}
	workbookMarshal = replaceRelationshipsNameSpace(workbookMarshal)
	parts["xl/workbook.xml"] = workbookMarshal
	if err != nil {
		return parts, err
	}

	parts["_rels/.rels"] = TEMPLATE__RELS_DOT_RELS
	parts["docProps/app.xml"] = TEMPLATE_DOCPROPS_APP
	// TODO - do this properly, modification and revision information
	parts["docProps/core.xml"] = TEMPLATE_DOCPROPS_CORE
	parts["xl/theme/theme1.xml"] = TEMPLATE_XL_THEME_THEME

	xSST := refTable.makeXLSXSST()
	parts["xl/sharedStrings.xml"], err = marshal(xSST)
	if err != nil {
		return parts, err
	}

	xWRel := workbookRels.MakeXLSXWorkbookRels()

	parts["xl/_rels/workbook.xml.rels"], err = marshal(xWRel)
	if err != nil {
		return parts, err
	}

	parts["[Content_Types].xml"], err = marshal(types)
	if err != nil {
		return parts, err
	}

	parts["xl/styles.xml"], err = f.styles.Marshal()
	if err != nil {
		return parts, err
	}

	return parts, nil
}

// Return the raw data contained in the File as three
// dimensional slice.  The first index represents the sheet number,
// the second the row number, and the third the cell number.
//
// For example:
//
//    var mySlice [][][]string
//    var value string
//    mySlice = xlsx.FileToSlice("myXLSX.xlsx")
//    value = mySlice[0][0][0]
//
// Here, value would be set to the raw value of the cell A1 in the
// first sheet in the XLSX file.
func (f *File) ToSlice() (output [][][]string, err error) {
	output = [][][]string{}
	for _, sheet := range f.Sheets {
		s := [][]string{}
		for _, row := range sheet.Rows {
			if row == nil {
				continue
			}
			r := []string{}
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					// Recover from strconv.NumError if the value is an empty string,
					// and insert an empty string in the output.
					if numErr, ok := err.(*strconv.NumError); ok && numErr.Num == "" {
						str = ""
					} else {
						return output, err
					}
				}
				r = append(r, str)
			}
			s = append(s, r)
		}
		output = append(output, s)
	}
	return output, nil
}

// ToSliceUnmerged returns the raw data contained in the File as three
// dimensional slice (s. method ToSlice).
// A covered cell become the value of its origin cell.
// Example: table where A1:A2 merged.
// | 01.01.2011 | Bread | 20 |
// |            | Fish  | 70 |
// This sheet will be converted to the slice:
// [  [01.01.2011 Bread 20]
// 		[01.01.2011 Fish 70] ]
func (f *File) ToSliceUnmerged() (output [][][]string, err error) {
	output, err = f.ToSlice()
	if err != nil {
		return nil, err
	}

	for s, sheet := range f.Sheets {
		for r, row := range sheet.Rows {
			for c, cell := range row.Cells {
				if cell.HMerge > 0 {
					for i := c + 1; i <= c+cell.HMerge; i++ {
						output[s][r][i] = output[s][r][c]
					}
				}

				if cell.VMerge > 0 {
					for i := r + 1; i <= r+cell.VMerge; i++ {
						output[s][i][c] = output[s][r][c]
					}
				}
			}
		}
	}

	return output, nil
}
