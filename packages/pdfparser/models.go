package pdfparser

// PatientInfo holds the header information extracted from a blood test PDF.
type PatientInfo struct {
	Name               string `json:"name"`
	Gender             string `json:"gender"`
	Date               string `json:"date"`
	Birthday           string `json:"birthday"`
	HealthcareFacility string `json:"healthcareFacility"`
}

// BloodTestResult represents a single test row.
// All fields are always present; missing values are stored as "".
type BloodTestResult struct {
	TestName       string `json:"testName"`
	Result         string `json:"result"`
	ResultUnit     string `json:"resultUnit"`
	ReferenceValue string `json:"referenceValue"`
}

// BloodTestGroup is a collection of tests that share the same date/time group.
type BloodTestGroup struct {
	GroupName string            `json:"groupName"`
	GroupDate string            `json:"groupDate"`
	Tests     []BloodTestResult `json:"tests"`
}

// BloodTestReport is the top-level parsed output of a blood test PDF.
type BloodTestReport struct {
	PatientInfo    PatientInfo      `json:"patientInfo"`
	TestGroups     []BloodTestGroup `json:"testGroups"`
	TotalTestCount int              `json:"totalTestCount"`
}
