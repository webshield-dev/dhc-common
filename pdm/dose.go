package pdm

import (
    "github.com/webshield-dev/dhc-common/vaccinemd"
)

//Dose a vaccine dose, use this as opposed to a FHIR record for now as want to use across SHC and EU DGC
//so seems easier to have a very simple structure
type Dose struct {

    //Code vaccine code
    Coding vaccinemd.Coding

    //Status http://hl7.org/fhir/R4/immunization-definitions.html#Immunization.status
    Status Code `json:"status,omitempty"`

    //OccurrenceDateTime if a dateTime see http://hl7.org/fhir/r4/immunization-definitions.html#Immunization.occurrence_x_
    OccurrenceDateTime string `json:"occurrenceDateTime,omitempty"`

    //OccurrenceDateTime if a string see http://hl7.org/fhir/r4/immunization-definitions.html#Immunization.occurrence_x_
    OccurrenceString string `json:"occurrenceString,omitempty"`

    //LotNumber http://hl7.org/fhir/R4/immunization-definitions.html#Immunization.lotNumber
    LotNumber string `json:"lotNumber,omitempty"`

    //Site where the dose was administered
    Site string `json:"site,omitempty"`
}



//Code https://www.hl7.org/fhir/datatypes.html#code
type Code string

const (
    //CodeCompleted action taken
    CodeCompleted Code = "completed"
)
