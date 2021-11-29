package verification

import (
	"github.com/webshield-dev/dhc-common/vaccinemd"
)

// CardVerificationState the card's verification state, see below
type CardVerificationState string

const (

	//CardVerificationStateUnknown no verifications have been performed
	CardVerificationStateUnknown CardVerificationState = "unknown"

	//CardVerificationStateValid card structure is valid, issuer is trusted, and
	//immunization criteria has been met
	CardVerificationStateValid CardVerificationState = "valid"

	//CardVerificationStateUnVerified the card signature could not be verified
	CardVerificationStateUnVerified CardVerificationState = "unverified"

	//CardVerificationStateExpired the card has expired
	CardVerificationStateExpired CardVerificationState = "expired"

	//CardVerificationStateIssuerUnknown card issuer not on a white list
	CardVerificationStateIssuerUnknown CardVerificationState = "issuer_unknown"

	//CardVerificationStateSafetyCriteriaNotMet the card signature is valid and from a trusted issuer
	//but the immunization criteria were not met. Card expired is also part of this
	CardVerificationStateSafetyCriteriaNotMet CardVerificationState = "safety_criteria_not_met"

	//CardVerificationStateCorrupt the digital signature is invalid
	CardVerificationStateCorrupt CardVerificationState = "corrupt"

)

//CardVerificationResults all verifications for card
type CardVerificationResults struct {
	//State the rolled up state
	State CardVerificationState `json:"state"`

	//CardStructure the card structure verifications results
	CardStructure *CardStructureVerificationResults `json:"card_structure,omitempty"`

	Issuer *IssuerVerificationResults `json:"issuer,omitempty"`

	Immunization *ImmunizationVerificationResults `json:"immunization,omitempty"`
}

//CardStructureVerificationResults the card structure verifications results
type CardStructureVerificationResults struct {

	//AllChecksPassed all required checks passed
	AllChecksPassed bool `json:"all_checks_passed"`

	//SignatureChecked for some reason the verifier can choose to not verify the card signature
	SignatureChecked bool `json:"signature_not_checked"`

	//FetchedKey true if successful fetched the verification key
	FetchedKey bool `json:"fetched_key"`

	//SignatureValid true if checked the signature an is valid
	SignatureValid bool `json:"signature_valid"`

	//Expired true if an exp date and has passed
	Expired bool `json:"expired"`
}

//IssuerVerificationResults issuer verification results
type IssuerVerificationResults struct {

	//AllChecksPassed all required checks passed
	AllChecksPassed bool `json:"all_checks_passed"`

	//Trusted issue is on a trusted whitelist
	Trusted bool `json:"trusted"`
}

// ImmunizationVerificationResults immunization verification results
type ImmunizationVerificationResults struct {

	//AllChecksPassed all required checks passed
	AllChecksPassed bool `json:"all_checks_passed"`

	//UnKnownVaccineType the vaccine is on a regional whitelist
	UnKnownVaccineType bool `json:"unknown_vaccine_type"`

	//TrustedVaccineType the vaccine is on a regional whitelist
	TrustedVaccineType bool `json:"trusted_vaccine_type"`

	MetDosesRequiredCriteria bool `json:"met_doses_required_criteria"`

	MetDaysBetweenDoesCriteria bool `json:"met_days_between_does_criteria"`

	MetDaysSinceLastDoseCriteria bool `json:"met_days_since_last_dose_criteria"`
}

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
}


//Code https://www.hl7.org/fhir/datatypes.html#code
type Code string

const (
	//CodeCompleted action taken
	CodeCompleted Code = "completed"
)
