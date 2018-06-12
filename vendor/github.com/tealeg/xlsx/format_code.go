package xlsx

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Do not edit these attributes once this struct is created. This struct should only be created by
// parseFullNumberFormatString() from a number format string. If the format for a cell needs to change, change
// the number format string and getNumberFormat() will invalidate the old struct and re-parse the string.
type parsedNumberFormat struct {
	numFmt                        string
	isTimeFormat                  bool
	negativeFormatExpectsPositive bool
	positiveFormat                *formatOptions
	negativeFormat                *formatOptions
	zeroFormat                    *formatOptions
	textFormat                    *formatOptions
	parseEncounteredError         *error
}

type formatOptions struct {
	isTimeFormat        bool
	showPercent         bool
	fullFormatString    string
	reducedFormatString string
	prefix              string
	suffix              string
}

// FormatValue returns a value, and possibly an error condition
// from a Cell.  If it is possible to apply a format to the cell
// value, it will do so, if not then an error will be returned, along
// with the raw value of the Cell.
//
// This is the documentation of the "General" Format in the Office Open XML spec:
//
// Numbers
// The application shall attempt to display the full number up to 11 digits (inc. decimal point). If the number is too
// large*, the application shall attempt to show exponential format. If the number has too many significant digits, the
// display shall be truncated. The optimal method of display is based on the available cell width. If the number cannot
// be displayed using any of these formats in the available width, the application shall show "#" across the width of
// the cell.
//
// Conditions for switching to exponential format:
// 1. The cell value shall have at least five digits for xE-xx
// 2. If the exponent is bigger than the size allowed, a floating point number cannot fit, so try exponential notation.
// 3. Similarly, for negative exponents, check if there is space for even one (non-zero) digit in floating point format**.
// 4. Finally, if there isn't room for all of the significant digits in floating point format (for a negative exponent),
// exponential format shall display more digits if the exponent is less than -3. (The 3 is because E-xx takes 4
// characters, and the leading 0 in floating point takes only 1 character. Thus, for an exponent less than -3, there is
// more than 3 additional leading 0's, more than enough to compensate for the size of the E-xx.)
//
// Floating point rule:
// For general formatting in cells, max overall length for cell display is 11, not including negative sign, but includes
// leading zeros and decimal separator.***
//
// Added Notes:
// * "If the number is too large" can also mean "if the number has more than 11 digits", so greater than or equal to
// 1e11 and less than 1e-9.
// ** Means that you should switch to scientific if there would be 9 zeros after the decimal (the decimal and first zero
// count against the 11 character limit), so less than 1e9.
// *** The way this is written, you can get numbers that are more than 11 characters because the golang Float fmt
// does not support adjusting the precision while not padding with zeros, while also not switching to scientific
// notation too early.
func (fullFormat *parsedNumberFormat) FormatValue(cell *Cell) (string, error) {
	switch cell.cellType {
	case CellTypeError:
		// The error type is what XLSX uses in error cases such as when formulas are invalid.
		// There will be text in the cell's value that can be shown, something ugly like #NAME? or #######
		return cell.Value, nil
	case CellTypeBool:
		if cell.Value == "0" {
			return "FALSE", nil
		} else if cell.Value == "1" {
			return "TRUE", nil
		} else {
			return cell.Value, errors.New("invalid value in bool cell")
		}
	case CellTypeString:
		fallthrough
	case CellTypeInline:
		fallthrough
	case CellTypeStringFormula:
		textFormat := cell.parsedNumFmt.textFormat
		// This switch statement is only for String formats
		switch textFormat.reducedFormatString {
		case builtInNumFmt[builtInNumFmtIndex_GENERAL]: // General is literally "general"
			return cell.Value, nil
		case builtInNumFmt[builtInNumFmtIndex_STRING]: // String is "@"
			return textFormat.prefix + cell.Value + textFormat.suffix, nil
		case "":
			// If cell is not "General" and there is not an "@" symbol in the format, then the cell's value is not
			// used when determining what to display. It would be completely legal to have a format of "Error"
			// for strings, and all values that are not numbers would show up as "Error". In that case, this code would
			// have a prefix of "Error" and a reduced format string of "" (empty string).
			return textFormat.prefix + textFormat.suffix, nil
		default:
			return cell.Value, errors.New("invalid or unsupported format, unsupported string format")
		}
	case CellTypeDate:
		// These are dates that are stored in date format instead of being stored as numbers with a format to turn them
		// into a date string.
		return cell.Value, nil
	case CellTypeNumeric:
		return fullFormat.formatNumericCell(cell)
	default:
		return cell.Value, errors.New("unknown cell type")
	}
}

func (fullFormat *parsedNumberFormat) formatNumericCell(cell *Cell) (string, error) {
	rawValue := strings.TrimSpace(cell.Value)
	// If there wasn't a value in the cell, it shouldn't have been marked as Numeric.
	// It's better to support this case though.
	if rawValue == "" {
		return "", nil
	}

	if fullFormat.isTimeFormat {
		return fullFormat.parseTime(rawValue, cell.date1904)
	}
	var numberFormat *formatOptions
	floatVal, floatErr := strconv.ParseFloat(rawValue, 64)
	if floatErr != nil {
		return rawValue, floatErr
	}
	// Choose the correct format. There can be different formats for positive, negative, and zero numbers.
	// Excel only uses the zero format if the value is literally zero, even if the number is so small that it shows
	// up as "0" when the positive format is used.
	if floatVal > 0 {
		numberFormat = fullFormat.positiveFormat
	} else if floatVal < 0 {
		// If format string specified a different format for negative numbers, then the number should be made positive
		// before getting formatted. The format string itself will contain formatting that denotes a negative number and
		// this formatting will end up in the prefix or suffix. Commonly if there is a negative format specified, the
		// number will get surrounded by parenthesis instead of showing it with a minus sign.
		if fullFormat.negativeFormatExpectsPositive {
			floatVal = math.Abs(floatVal)
		}
		numberFormat = fullFormat.negativeFormat
	} else {
		numberFormat = fullFormat.zeroFormat
	}

	// When showPercent is true, multiply the number by 100.
	// The percent sign will be in the prefix or suffix already, so it does not need to be added in this function.
	// The number format itself will be the same as any other number format once the value is multiplied by 100.
	if numberFormat.showPercent {
		floatVal = 100 * floatVal
	}

	// Only the most common format strings are supported here.
	// Eventually this switch needs to be replaced with a more general solution.
	// Some of these "supported" formats should have thousand separators, but don't get them since Go fmt
	// doesn't have a way to request thousands separators.
	// The only things that should be supported here are in the array formattingCharacters,
	// everything else has been stripped out before and will be placed in the prefix or suffix.
	// The formatting characters can have non-formatting characters mixed in with them and those should be maintained.
	// However, at this time we fail to parse those formatting codes and they get replaced with "General"
	var formattedNum string
	switch numberFormat.reducedFormatString {
	case builtInNumFmt[builtInNumFmtIndex_GENERAL]: // General is literally "general"
		// prefix, showPercent, and suffix cannot apply to the general format
		// The logic for showing numbers when the format is "general" is much more complicated than the rest of these.
		generalFormatted, err := generalNumericScientific(cell.Value, true)
		if err != nil {
			return rawValue, nil
		}
		return generalFormatted, nil
	case builtInNumFmt[builtInNumFmtIndex_STRING]: // String is "@"
		formattedNum = cell.Value
	case builtInNumFmt[builtInNumFmtIndex_INT], "#,##0": // Int is "0"
		// Previously this case would cast to int and print with %d, but that will not round the value correctly.
		formattedNum = fmt.Sprintf("%.0f", floatVal)
	case "0.0", "#,##0.0":
		formattedNum = fmt.Sprintf("%.1f", floatVal)
	case builtInNumFmt[builtInNumFmtIndex_FLOAT], "#,##0.00": // Float is "0.00"
		formattedNum = fmt.Sprintf("%.2f", floatVal)
	case "0.000", "#,##0.000":
		formattedNum = fmt.Sprintf("%.3f", floatVal)
	case "0.0000", "#,##0.0000":
		formattedNum = fmt.Sprintf("%.4f", floatVal)
	case "0.00e+00", "##0.0e+0":
		formattedNum = fmt.Sprintf("%e", floatVal)
	case "":
		// Do nothing.
	default:
		return rawValue, nil
	}
	return numberFormat.prefix + formattedNum + numberFormat.suffix, nil
}

func generalNumericScientific(value string, allowScientific bool) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value, err
	}
	if allowScientific {
		absF := math.Abs(f)
		// When using General format, numbers that are less than 1e-9 (0.000000001) and greater than or equal to
		// 1e11 (100,000,000,000) should be shown in scientific notation.
		// Numbers less than the number after zero, are assumed to be zero.
		if (absF >= math.SmallestNonzeroFloat64 && absF < minNonScientificNumber) || absF >= maxNonScientificNumber {
			return strconv.FormatFloat(f, 'E', -1, 64), nil
		}
	}
	// This format (fmt="f", prec=-1) will prevent padding with zeros and will never switch to scientific notation.
	// However, it will show more than 11 characters for very precise numbers, and this cannot be changed.
	// You could also use fmt="g", prec=11, which doesn't pad with zeros and allows the correct precision,
	// but it will use scientific notation on numbers less than 1e-4. That value is hardcoded in Go and cannot be
	// configured or disabled.
	return strconv.FormatFloat(f, 'f', -1, 64), nil
}

// Format strings are a little strange to compare because empty string needs to be taken as general, and general needs
// to be compared case insensitively.
func compareFormatString(fmt1, fmt2 string) bool {
	if fmt1 == fmt2 {
		return true
	}
	if fmt1 == "" || strings.EqualFold(fmt1, "general") {
		fmt1 = "general"
	}
	if fmt2 == "" || strings.EqualFold(fmt2, "general") {
		fmt2 = "general"
	}
	return fmt1 == fmt2
}

func parseFullNumberFormatString(numFmt string) *parsedNumberFormat {
	parsedNumFmt := &parsedNumberFormat{
		numFmt: numFmt,
	}
	if isTimeFormat(numFmt) {
		// Time formats cannot have multiple groups separated by semicolons, there is only one format.
		// Strings are unaffected by the time format.
		parsedNumFmt.isTimeFormat = true
		parsedNumFmt.textFormat, _ = parseNumberFormatSection("general")
		return parsedNumFmt
	}

	var fmtOptions []*formatOptions
	formats, err := splitFormatOnSemicolon(numFmt)
	if err == nil {
		for _, formatSection := range formats {
			parsedFormat, err := parseNumberFormatSection(formatSection)
			if err != nil {
				// If an invalid number section is found, fall back to general
				parsedFormat = fallbackErrorFormat
				parsedNumFmt.parseEncounteredError = &err
			}
			fmtOptions = append(fmtOptions, parsedFormat)
		}
	} else {
		fmtOptions = append(fmtOptions, fallbackErrorFormat)
		parsedNumFmt.parseEncounteredError = &err
	}
	if len(fmtOptions) > 4 {
		fmtOptions = []*formatOptions{fallbackErrorFormat}
		err = errors.New("invalid number format, too many format sections")
		parsedNumFmt.parseEncounteredError = &err
	}

	if len(fmtOptions) == 1 {
		// If there is only one option, it is used for all
		parsedNumFmt.positiveFormat = fmtOptions[0]
		parsedNumFmt.negativeFormat = fmtOptions[0]
		parsedNumFmt.zeroFormat = fmtOptions[0]
		if strings.Contains(fmtOptions[0].fullFormatString, "@") {
			parsedNumFmt.textFormat = fmtOptions[0]
		} else {
			parsedNumFmt.textFormat, _ = parseNumberFormatSection("general")
		}
	} else if len(fmtOptions) == 2 {
		// If there are two formats, the first is used for positive and zeros, the second gets used as a negative format,
		// and strings are not formatted.
		// When negative numbers now have their own format, they should become positive before having the format applied.
		// The format will contain a negative sign if it is desired, but they may be colored red or wrapped in
		// parenthesis instead.
		parsedNumFmt.negativeFormatExpectsPositive = true
		parsedNumFmt.positiveFormat = fmtOptions[0]
		parsedNumFmt.negativeFormat = fmtOptions[1]
		parsedNumFmt.zeroFormat = fmtOptions[0]
		parsedNumFmt.textFormat, _ = parseNumberFormatSection("general")
	} else if len(fmtOptions) == 3 {
		// If there are three formats, the first is used for positive, the second gets used as a negative format,
		// the third is for negative, and strings are not formatted.
		parsedNumFmt.negativeFormatExpectsPositive = true
		parsedNumFmt.positiveFormat = fmtOptions[0]
		parsedNumFmt.negativeFormat = fmtOptions[1]
		parsedNumFmt.zeroFormat = fmtOptions[2]
		parsedNumFmt.textFormat, _ = parseNumberFormatSection("general")
	} else {
		// With four options, the first is positive, the second is negative, the third is zero, and the fourth is strings
		// Negative numbers should be still become positive before having the negative formatting applied.
		parsedNumFmt.negativeFormatExpectsPositive = true
		parsedNumFmt.positiveFormat = fmtOptions[0]
		parsedNumFmt.negativeFormat = fmtOptions[1]
		parsedNumFmt.zeroFormat = fmtOptions[2]
		parsedNumFmt.textFormat = fmtOptions[3]
	}
	return parsedNumFmt
}

// splitFormatOnSemicolon will split the format string into the format sections
// This logic to split the different formats on semicolon is fully correct, and will skip all literal semicolons,
// and will catch all breaking semicolons.
func splitFormatOnSemicolon(format string) ([]string, error) {
	var formats []string
	prevIndex := 0
	for i := 0; i < len(format); i++ {
		if format[i] == ';' {
			formats = append(formats, format[prevIndex:i])
			prevIndex = i + 1
		} else if format[i] == '\\' {
			i++
		} else if format[i] == '"' {
			endQuoteIndex := strings.Index(format[i+1:], "\"")
			if endQuoteIndex == -1 {
				// This is an invalid format string, fall back to general
				return nil, errors.New("invalid format string, unmatched double quote")
			}
			i += endQuoteIndex + 1
		}
	}
	return append(formats, format[prevIndex:]), nil
}

var fallbackErrorFormat = &formatOptions{
	fullFormatString:    "general",
	reducedFormatString: "general",
}

// parseNumberFormatSection takes in individual format and parses out most of the options.
// Some options are parsed, removed from the string, and set as settings on formatOptions.
// There remainder of the format string is put in the reducedFormatString attribute, and supported values for these
// are handled in a switch in the Cell.FormattedValue() function.
// Ideally more and more of the format string would be parsed out here into settings until there is no remainder string
// at all.
// Features that this supports:
// - Time formats are detected, and marked in the options. Time format strings are handled when doing the formatting.
//   The logic to detect time formats is currently not correct, and can catch formats that are not time formats as well
//   as miss formats that are time formats.
// - Color formats are detected and removed.
// - Currency annotations are handled properly.
// - Literal strings wrapped in quotes are handled and put into prefix or suffix.
// - Numbers that should be percent are detected and marked in the options.
// - Conditionals are detected and removed, but they are not obeyed. The conditional groups will be used just like the
//   positive;negative;zero;string format groups. Here is an example of a conditional format: "[Red][<=100];[Blue][>100]"
// Decoding the actual number formatting portion is out of scope, that is placed into reducedFormatString and is used
// when formatting the string. The string there will be reduced to only the things in the formattingCharacters array.
// Everything not in that array has been parsed out and put into formatOptions.
func parseNumberFormatSection(fullFormat string) (*formatOptions, error) {
	reducedFormat := strings.TrimSpace(fullFormat)

	// general is the only format that does not use the normal format symbols notations
	if compareFormatString(reducedFormat, "general") {
		return &formatOptions{
			fullFormatString:    "general",
			reducedFormatString: "general",
		}, nil
	}

	prefix, reducedFormat, showPercent1, err := parseLiterals(reducedFormat)
	if err != nil {
		return nil, err
	}

	reducedFormat, suffixFormat := splitFormatAndSuffixFormat(reducedFormat)

	suffix, remaining, showPercent2, err := parseLiterals(suffixFormat)
	if err != nil {
		return nil, err
	}
	if len(remaining) > 0 {
		// This paradigm of codes consisting of literals, number formats, then more literals is not always correct, they can
		// actually be intertwined. Though 99% of the time number formats will not do this.
		// Excel uses this format string for Social Security Numbers: 000\-00\-0000
		// and this for US phone numbers: [<=9999999]###\-####;\(###\)\ ###\-####
		return nil, errors.New("invalid or unsupported format string")
	}

	return &formatOptions{
		fullFormatString:    fullFormat,
		isTimeFormat:        false,
		reducedFormatString: reducedFormat,
		prefix:              prefix,
		suffix:              suffix,
		showPercent:         showPercent1 || showPercent2,
	}, nil
}

// formattingCharacters will be left in the reducedNumberFormat
// It is important that these be looked for in order so that the slash cases are handled correctly.
// / (slash) is a fraction format if preceded by 0, #, or ?, otherwise it is not a formatting character
// E- E+ e- e+ are scientific notation, but E, e, -, + are not formatting characters independently
// \ (back slash) makes the next character a literal (not formatting)
// " Anything in double quotes is not a formatting character
// _ (underscore) skips the width of the next character, so the next character cannot be formatting
var formattingCharacters = []string{"0/", "#/", "?/", "E-", "E+", "e-", "e+", "0", "#", "?", ".", ",", "@", "*"}

// The following are also time format characters, but since this is only used for detecting, not decoding, they are
// redundant here: ee, gg, ggg, rr, ss, mm, hh, yyyy, dd, ddd, dddd, mm, mmm, mmmm, mmmmm, ss.0000, ss.000, ss.00, ss.0
// The .00 type format is very tricky, because it only counts if it comes after ss or s or [ss] or [s]
// .00 is actually a valid number format by itself.
var timeFormatCharacters = []string{"m", "d", "yy", "h", "m", "AM/PM", "A/P", "am/pm", "a/p", "r", "g", "e", "b1", "b2", "[hh]", "[h]", "[mm]", "[m]",
	"s.0000", "s.000", "s.00", "s.0", "s", "[ss].0000", "[ss].000", "[ss].00", "[ss].0", "[ss]", "[s].0000", "[s].000", "[s].00", "[s].0", "[s]"}

func splitFormatAndSuffixFormat(format string) (string, string) {
	var i int
	for ; i < len(format); i++ {
		curReducedFormat := format[i:]
		var found bool
		for _, special := range formattingCharacters {
			if strings.HasPrefix(curReducedFormat, special) {
				// Skip ahead if the special character was longer than length 1
				i += len(special) - 1
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	suffixFormat := format[i:]
	format = format[:i]
	return format, suffixFormat
}

func parseLiterals(format string) (string, string, bool, error) {
	var prefix string
	showPercent := false
	for i := 0; i < len(format); i++ {
		curReducedFormat := format[i:]
		switch curReducedFormat[0] {
		case '\\':
			// If there is a slash, skip the next character, and add it to the prefix
			if len(curReducedFormat) > 1 {
				i++
				prefix += curReducedFormat[1:2]
			}
		case '_':
			// If there is an underscore, skip the next character, but don't add it to the prefix
			if len(curReducedFormat) > 1 {
				i++
			}
		case '*':
			// Asterisks are used to repeat the next character to fill the full cell width.
			// There isn't really a cell size in this context, so this will be ignored.
		case '"':
			// If there is a quote skip to the next quote, and add the quoted characters to the prefix
			endQuoteIndex := strings.Index(curReducedFormat[1:], "\"")
			if endQuoteIndex == -1 {
				return "", "", false, errors.New("invalid formatting code, unmatched double quote")
			}
			prefix = prefix + curReducedFormat[1:endQuoteIndex+1]
			i += endQuoteIndex + 1
		case '%':
			showPercent = true
			prefix += "%"
		case '[':
			// Brackets can be currency annotations (e.g. [$$-409])
			// color formats (e.g. [color1] through [color56], as well as [red] etc.)
			// conditionals (e.g. [>100], the valid conditionals are =, >, <, >=, <=, <>)
			bracketIndex := strings.Index(curReducedFormat, "]")
			if bracketIndex == -1 {
				return "", "", false, errors.New("invalid formatting code, invalid brackets")
			}
			// Currencies in Excel are annotated with this format: [$<Currency String>-<Language Info>]
			// Currency String is something like $, ¥, €, or £
			// Language Info is three hexadecimal characters
			if len(curReducedFormat) > 2 && curReducedFormat[1] == '$' {
				dashIndex := strings.Index(curReducedFormat, "-")
				if dashIndex != -1 && dashIndex < bracketIndex {
					// Get the currency symbol, and skip to the end of the currency format
					prefix += curReducedFormat[2:dashIndex]
				} else {
					return "", "", false, errors.New("invalid formatting code, invalid currency annotation")
				}
			}
			i += bracketIndex
		case '$', '-', '+', '/', '(', ')', ':', '!', '^', '&', '\'', '~', '{', '}', '<', '>', '=', ' ':
			// These symbols are allowed to be used as literal without escaping
			prefix += curReducedFormat[0:1]
		default:
			for _, special := range formattingCharacters {
				if strings.HasPrefix(curReducedFormat, special) {
					// This means we found the start of the actual number formatting portion, and should return.
					return prefix, format[i:], showPercent, nil
				}
			}
			// Symbols that don't have meaning and aren't in the exempt literal characters and are not escaped.
			return "", "", false, errors.New("invalid formatting code: unsupported or unescaped characters")
		}
	}
	return prefix, "", showPercent, nil
}

// parseTime returns a string parsed using time.Time
func (fullFormat *parsedNumberFormat) parseTime(value string, date1904 bool) (string, error) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return value, err
	}
	val := TimeFromExcelTime(f, date1904)
	format := fullFormat.numFmt
	// Replace Excel placeholders with Go time placeholders.
	// For example, replace yyyy with 2006. These are in a specific order,
	// due to the fact that m is used in month, minute, and am/pm. It would
	// be easier to fix that with regular expressions, but if it's possible
	// to keep this simple it would be easier to maintain.
	// Full-length month and days (e.g. March, Tuesday) have letters in them that would be replaced
	// by other characters below (such as the 'h' in March, or the 'd' in Tuesday) below.
	// First we convert them to arbitrary characters unused in Excel Date formats, and then at the end,
	// turn them to what they should actually be.
	// Based off: http://www.ozgrid.com/Excel/CustomFormats.htm
	replacements := []struct{ xltime, gotime string }{
		{"yyyy", "2006"},
		{"yy", "06"},
		{"mmmm", "%%%%"},
		{"dddd", "&&&&"},
		{"dd", "02"},
		{"d", "2"},
		{"mmm", "Jan"},
		{"mmss", "0405"},
		{"ss", "05"},
		{"mm:", "04:"},
		{":mm", ":04"},
		{"mm", "01"},
		{"am/pm", "pm"},
		{"m/", "1/"},
		{"%%%%", "January"},
		{"&&&&", "Monday"},
	}
	// It is the presence of the "am/pm" indicator that determins
	// if this is a 12 hour or 24 hours time format, not the
	// number of 'h' characters.
	if is12HourTime(format) {
		format = strings.Replace(format, "hh", "03", 1)
		format = strings.Replace(format, "h", "3", 1)
	} else {
		format = strings.Replace(format, "hh", "15", 1)
		format = strings.Replace(format, "h", "15", 1)
	}
	for _, repl := range replacements {
		format = strings.Replace(format, repl.xltime, repl.gotime, 1)
	}
	// If the hour is optional, strip it out, along with the
	// possible dangling colon that would remain.
	if val.Hour() < 1 {
		format = strings.Replace(format, "]:", "]", 1)
		format = strings.Replace(format, "[03]", "", 1)
		format = strings.Replace(format, "[3]", "", 1)
		format = strings.Replace(format, "[15]", "", 1)
	} else {
		format = strings.Replace(format, "[3]", "3", 1)
		format = strings.Replace(format, "[15]", "15", 1)
	}
	return val.Format(format), nil
}

// isTimeFormat checks whether an Excel format string represents a time.Time.
// This function is now correct, but it can detect time format strings that cannot be correctly handled by parseTime()
func isTimeFormat(format string) bool {
	var foundTimeFormatCharacters bool
	for i := 0; i < len(format); i++ {
		curReducedFormat := format[i:]
		switch curReducedFormat[0] {
		case '\\', '_':
			// If there is a slash, skip the next character, and add it to the prefix
			// If there is an underscore, skip the next character, but don't add it to the prefix
			if len(curReducedFormat) > 1 {
				i++
			}
		case '*':
			// Asterisks are used to repeat the next character to fill the full cell width.
			// There isn't really a cell size in this context, so this will be ignored.
		case '"':
			// If there is a quote skip to the next quote, and add the quoted characters to the prefix
			endQuoteIndex := strings.Index(curReducedFormat[1:], "\"")
			if endQuoteIndex == -1 {
				// This is not any type of valid format.
				return false
			}
			i += endQuoteIndex + 1
		case '$', '-', '+', '/', '(', ')', ':', '!', '^', '&', '\'', '~', '{', '}', '<', '>', '=', ' ':
			// These symbols are allowed to be used as literal without escaping
		case ',':
			// This is not documented in the XLSX spec as far as I can tell, but Excel and Numbers will include
			// commas in number formats without escaping them, so this should be supported.
		default:
			foundInThisLoop := false
			for _, special := range timeFormatCharacters {
				if strings.HasPrefix(curReducedFormat, special) {
					foundTimeFormatCharacters = true
					foundInThisLoop = true
					i += len(special) - 1
					break
				}
			}
			if foundInThisLoop {
				continue
			}
			if curReducedFormat[0] == '[' {
				// For number formats, this code would happen above in a case '[': section.
				// However, for time formats it must happen after looking for occurrences in timeFormatCharacters
				// because there are a few time formats that can be wrapped in brackets.

				// Brackets can be currency annotations (e.g. [$$-409])
				// color formats (e.g. [color1] through [color56], as well as [red] etc.)
				// conditionals (e.g. [>100], the valid conditionals are =, >, <, >=, <=, <>)
				bracketIndex := strings.Index(curReducedFormat, "]")
				if bracketIndex == -1 {
					// This is not any type of valid format.
					return false
				}
				i += bracketIndex
				continue
			}
			// Symbols that don't have meaning, aren't in the exempt literal characters, and aren't escaped are invalid.
			// The string could still be a valid number format string.
			return false
		}
	}
	// If the string doesn't have any time formatting characters, it could technically be a time format, but it
	// would be a pretty weak time format. A valid time format with no time formatting symbols will also be a number
	// format with no number formatting symbols, which is essentially a constant string that does not depend on the
	// cell's value in anyway. The downstream logic will do the right thing in that case if this returns false.
	return foundTimeFormatCharacters
}

// is12HourTime checks whether an Excel time format string is a 12
// hours form.
func is12HourTime(format string) bool {
	return strings.Contains(format, "am/pm") || strings.Contains(format, "AM/PM") || strings.Contains(format, "a/p") || strings.Contains(format, "A/P")
}
