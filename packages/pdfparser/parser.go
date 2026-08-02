package pdfparser

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

// skipPatterns are lines to ignore when parsing.
var skipPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^T\.C\.\s*SAĞLIK`),
	regexp.MustCompile(`(?i)^Sağlık\s+Bilgi`),
	regexp.MustCompile(`(?i)^Date$`),
	regexp.MustCompile(`(?i)^Test\s+Name$`),
	regexp.MustCompile(`(?i)^Result$`),
	regexp.MustCompile(`(?i)^Unit$`),
	regexp.MustCompile(`(?i)^Referenc$`),
	regexp.MustCompile(`(?i)^e\s+Value$`),
	regexp.MustCompile(`(?i)^enabiz\.gov\.tr$`),
	regexp.MustCompile(`^\d\s+850\s+240`),
	regexp.MustCompile(`(?i)^Page\s+\d+`),
	// Patient info header labels and their values (each on its own line).
	regexp.MustCompile(`(?i)^Name\s*/?\s*Surname\s*:`),
	regexp.MustCompile(`(?i)^Gender\s*:`),
	regexp.MustCompile(`(?i)^Date\s*:`),
	regexp.MustCompile(`(?i)^Birthday\s*:`),
	regexp.MustCompile(`(?i)^Healthcare\s+Facility\s*:`),
}

// Compiled patterns.
var (
	reDateLine = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	reTimeLine = regexp.MustCompile(`^\d{2}:\d{2}$`)
)

// knownUnits is used to distinguish a unit line from a reference value line.
var knownUnits = map[string]bool{
	"K/uL": true, "M/uL": true, "U/L": true,
	"g/dL": true, "mg/dL": true, "mg/dl": true,
	"µg/dl": true, "ug/L": true, "ng/mL": true, "ng/dL": true,
	"pg/mL": true, "pg": true, "fl": true, "fL": true,
	"%": true, "µIU/ml": true, "uIU/ml": true, "uIU/mL": true,
	"µU/ml": true, "mEq/L": true, "10 3/uL": true, "10 6/uL": true,
	"mL/min/": true,
}

// ParseBloodTestPDF extracts structured blood test data from PDF bytes.
func ParseBloodTestPDF(pdfData []byte) (*BloodTestReport, error) {
	text, err := extractText(pdfData)
	if err != nil {
		return nil, fmt.Errorf("pdfparser: failed to extract text: %w", err)
	}

	return parseText(text), nil
}

// extractText uses ledongthuc/pdf to pull all text content from a PDF.
func extractText(pdfData []byte) (string, error) {
	r, err := pdf.NewReader(bytes.NewReader(pdfData), int64(len(pdfData)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var sb strings.Builder
	for i := 1; i <= r.NumPage(); i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", i, err)
		}
		sb.WriteString(text)
	}

	return sb.String(), nil
}

// parseText is the deterministic parser that processes line-by-line output
// from ledongthuc/pdf. The library outputs each table cell on its own line,
// so a test row appears as 3 or 4 consecutive lines:
//
//	TestName
//	Result
//	Unit          (may be absent for unit-less tests)
//	ReferenceValue
func parseText(text string) *BloodTestReport {
	rawLines := strings.Split(text, "\n")

	// Clean: trim whitespace, drop skip-pattern lines and blanks.
	// Header labels (e.g. "Name/ Surname:") are followed by their value
	// on the next line — both must be excluded from the test-parsing stream.
	var lines []string
	skipNext := false
	for _, l := range rawLines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if skipNext {
			skipNext = false
			continue
		}
		if shouldSkip(l) {
			// Header label lines are followed by a value line to skip.
			if isHeaderLabel(l) {
				skipNext = true
			}
			continue
		}
		lines = append(lines, l)
	}

	report := &BloodTestReport{}
	report.PatientInfo = extractPatientInfo(rawLines)
	report.TestGroups = extractTestGroups(lines)

	total := 0
	for _, g := range report.TestGroups {
		total += len(g.Tests)
	}
	report.TotalTestCount = total

	return report
}

// shouldSkip returns true for header/footer/table-header lines.
func shouldSkip(line string) bool {
	for _, pat := range skipPatterns {
		if pat.MatchString(line) {
			return true
		}
	}
	return false
}

// headerLabelPatterns identify lines that are patient info labels.
// The line immediately after each label is the label's value and
// should also be skipped.
var headerLabelPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^Name\s*/?\s*Surname\s*:`),
	regexp.MustCompile(`(?i)^Gender\s*:`),
	regexp.MustCompile(`(?i)^Date\s*:`),
	regexp.MustCompile(`(?i)^Birthday\s*:`),
	regexp.MustCompile(`(?i)^Healthcare\s+Facility\s*:`),
}

// isHeaderLabel returns true if a line is a patient-info label whose
// following line (the value) should also be removed.
func isHeaderLabel(line string) bool {
	for _, pat := range headerLabelPatterns {
		if pat.MatchString(line) {
			return true
		}
	}
	return false
}

// extractPatientInfo scans raw lines for header fields.
func extractPatientInfo(lines []string) PatientInfo {
	pi := PatientInfo{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Name/ Surname:") || strings.HasPrefix(line, "Name/Surname:") {
			// Value is on the next line.
			if i+1 < len(lines) {
				pi.Name = strings.TrimSpace(lines[i+1])
			}
		}
		if strings.HasPrefix(line, "Gender:") {
			if i+1 < len(lines) {
				pi.Gender = strings.TrimSpace(lines[i+1])
			}
		}
		if strings.HasPrefix(line, "Date:") {
			if i+1 < len(lines) {
				pi.Date = strings.TrimSpace(lines[i+1])
			}
		}
		if strings.HasPrefix(line, "Birthday:") {
			if i+1 < len(lines) {
				pi.Birthday = strings.TrimSpace(lines[i+1])
			}
		}
		if strings.HasPrefix(line, "Healthcare Facility:") {
			if i+1 < len(lines) {
				pi.HealthcareFacility = strings.TrimSpace(lines[i+1])
			}
		}
	}
	return pi
}

// extractTestGroups walks lines building groups and test rows.
//
// The line stream (after filtering) looks like:
//
//	13.02.2026              ← date
//	09:35                   ← time
//	Tam Kan Sayımı (...)    ← group name (optional, non-numeric, non-date)
//	BASO                    ← test name
//	0,05                    ← result
//	K/uL                    ← unit   (may be missing)
//	0 - 0,1                 ← reference value
//	BASO %                  ← next test name
//	...
func extractTestGroups(lines []string) []BloodTestGroup {
	var groups []BloodTestGroup
	var cur *BloodTestGroup

	flush := func() {
		if cur != nil && len(cur.Tests) > 0 {
			groups = append(groups, *cur)
		}
	}

	i := 0
	for i < len(lines) {
		line := lines[i]

		// --- Date line: start of a new group ---
		if reDateLine.MatchString(line) {
			datePart := line
			timePart := ""
			groupName := ""

			// Peek for time.
			if i+1 < len(lines) && reTimeLine.MatchString(lines[i+1]) {
				timePart = lines[i+1]
				i++
				// Peek for group name (non-numeric, non-date line).
				if i+1 < len(lines) && !isNumericValue(lines[i+1]) && !reDateLine.MatchString(lines[i+1]) && !isKnownUnit(lines[i+1]) {
					// Check if next line is a test name followed by a numeric result,
					// or a standalone group name. A group name is typically longer
					// and descriptive.
					if i+2 < len(lines) && isNumericValue(lines[i+2]) {
						// lines[i+1] is a test name, not a group name. Don't consume.
					} else {
						groupName = lines[i+1]
						i++
					}
				}
			}

			flush()
			groupDate := datePart
			if timePart != "" {
				groupDate = datePart + " " + timePart
			}
			cur = &BloodTestGroup{
				GroupName: groupName,
				GroupDate: groupDate,
			}
			i++
			continue
		}

		// --- Try to parse a test result (3 or 4 lines) ---
		test, consumed := tryParseTestLines(lines, i)
		if consumed > 0 {
			if cur == nil {
				cur = &BloodTestGroup{GroupName: "", GroupDate: ""}
			}
			cur.Tests = append(cur.Tests, test)
			i += consumed
			continue
		}

		// Unrecognised line — skip.
		i++
	}

	flush()
	return groups
}

// tryParseTestLines tries to consume 3 or 4 lines starting at idx as a test
// result row. Returns the test and how many lines were consumed (0 if failed).
//
// Pattern A (4 lines): TestName, Result, Unit, ReferenceValue
// Pattern B (3 lines): TestName, Result, ReferenceValue (no unit)
func tryParseTestLines(lines []string, idx int) (BloodTestResult, int) {
	r := BloodTestResult{}

	if idx >= len(lines) {
		return r, 0
	}

	// Line 0: Test name — must not be numeric, a date, or a time.
	testName := lines[idx]
	if isNumericValue(testName) || reDateLine.MatchString(testName) || reTimeLine.MatchString(testName) {
		return r, 0
	}

	// Line 1: Result — must be numeric.
	if idx+1 >= len(lines) || !isNumericValue(lines[idx+1]) {
		return r, 0
	}
	result := lines[idx+1]

	// Line 2: could be Unit or ReferenceValue.
	if idx+2 >= len(lines) {
		// Only have testName + result.
		r.TestName = testName
		r.Result = result
		return r, 2
	}

	line2 := lines[idx+2]

	// Check if line2 is a known unit.
	if isKnownUnit(line2) {
		r.TestName = testName
		r.Result = result
		r.ResultUnit = line2

		// Line 3: Reference value.
		if idx+3 < len(lines) && isReferenceValue(lines[idx+3]) {
			r.ReferenceValue = lines[idx+3]
			return r, 4
		}
		return r, 3
	}

	// Line2 is not a unit — treat as reference value (3-line pattern).
	if isReferenceValue(line2) {
		r.TestName = testName
		r.Result = result
		r.ReferenceValue = line2
		return r, 3
	}

	// Line2 is something else — just take testName + result.
	r.TestName = testName
	r.Result = result
	return r, 2
}

// isKnownUnit checks if a line matches a known measurement unit.
func isKnownUnit(s string) bool {
	return knownUnits[s]
}

// isReferenceValue checks if a line looks like a reference range/value.
// Examples: "0 - 0,1", "30 -", "< 0,5", "> 40 Normal <40 Belirgin risk",
// "20-50 Optimum Seviye >150 Toksisite Riski"
func isReferenceValue(s string) bool {
	// Must not be a date or time.
	if reDateLine.MatchString(s) || reTimeLine.MatchString(s) {
		return false
	}
	// Starts with a digit, <, >, or = — likely a reference value.
	if len(s) > 0 {
		c := s[0]
		if (c >= '0' && c <= '9') || c == '<' || c == '>' || c == '=' {
			return true
		}
	}
	return false
}

// isNumericValue checks if a string is a number (possibly with comma decimal).
func isNumericValue(s string) bool {
	if s == "" {
		return false
	}
	start := 0
	if s[0] == '-' {
		start = 1
		if len(s) == 1 {
			return false
		}
	}
	hasDigit := false
	for _, c := range s[start:] {
		switch {
		case c >= '0' && c <= '9':
			hasDigit = true
		case c == ',' || c == '.':
			// decimal separator
		default:
			return false
		}
	}
	return hasDigit
}
