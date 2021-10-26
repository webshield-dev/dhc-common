package vaccinemd

//
// SEE https://www.cdc.gov/vaccines/programs/iis/COVID-19-related-codes.html
//

//CovidVaccineMetadata metadata about a covid vaccine
//from https://www.cdc.gov/vaccines/programs/iis/COVID-19-related-codes.html
type CovidVaccineMetadata struct {
	Codes []Coding

	//CVXStatus cvx status from the cdc table
	CVXStatus CVSStatus  `json:"cvs_status"`

	//Doses number of doses required
	Doses int `json:"doses"`

	//DaysSinceLastDoseCriteria took from common pass recommendations
	DaysSinceLastDoseCriteria int `json:"days_since_last_dose_criteria"`

	//DaysBetweenDoesCriteriaBegin begin of range took from common pass recommendations
	DaysBetweenDoesCriteriaBegin int `json:"days_between_does_criteria_begin"`

	//DaysBetweenDoesCriteriaEnd end of range took from common pass recommendations
	DaysBetweenDoesCriteriaEnd int `json:"days_between_does_criteria_end"`

	//SaleProprietaryName from cdc table
	SaleProprietaryName string `json:"sale_proprietary_name"`

	//ManufacturerName short name used to display in cards
	ManufacturerName string `json:"manufacturer_name"`
}

//CVSStatus if CDC states from table
type CVSStatus string

const (
	//CVSStatusActive active in US
	CVSStatusActive CVSStatus = "Active"

	//VSStatusNonUS active outside of US
	CVSStatusNonUS CVSStatus = "Non-US"
)



//Coding http://hl7.org/fhir/R4/datatypes.html#Coding
type Coding struct {

	//System http://hl7.org/fhir/R4/datatypes-definitions.html#Coding.system
	System string `json:"system,omitempty" yaml:"system,omitempty"`

	//Code http://hl7.org/fhir/R4/datatypes-definitions.html#Coding.code
	Code string `json:"code,omitempty" yaml:"code,omitempty"`
}
