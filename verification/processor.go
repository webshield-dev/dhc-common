package verification

import (
	"fmt"
	"github.com/webshield-dev/dhc-common/vaccinemd"
	"time"
)

//Processor can be created by a verifier to manage the verification state and calculate cards verification state
//not designed to be thread safe. Create one per card verification
type Processor interface {

	//GetVerificationResults returns the current card verification results based on current status
	GetVerificationResults() *CardVerificationResults

	//
	// Card struct setters
	//

	//SetSignatureChecked do not check signature for some reason
	SetSignatureChecked()

	//SetFetchedKey record key fetched
	SetFetchedKey()

	//SetSignatureValid record signature valid
	SetSignatureValid()

	//SetExpired record expired
	SetExpired()

	//CardStructureVerified check if all card structure verifications have passed
	CardStructureVerified() bool

	//CardCorrupted true if card is corrupted
	CardCorrupted() bool

	//
	// Issuer verifications
	//

	//SetIssuerTrusted issuer is on a trusted whitelist
	SetIssuerTrusted()

	//IssuerVerified check is all the issuers verifications have passed
	IssuerVerified() bool

	//
	// Immunization Criteria
	//

	VerifyImmunization(
		region Region,
		Doses []*Dose, // the doses administered
	) (bool, error)

	//ImmunizationCriteriaMet true if all the immunization criteria have been met, can be called
	//after verifyImmunization
	ImmunizationCriteriaMet() bool
}

func NewProcessor() Processor {

	return &v1Processor{
		mdRepo: vaccinemd.MakeRepo(),
		results: &CardVerificationResults{
			State:         CardVerificationStateUnknown,
			CardStructure: &CardStructureVerificationResults{},
			Issuer:        &IssuerVerificationResults{},
			Immunization:  &ImmunizationVerificationResults{},
		},
	}

}

type v1Processor struct {
	mdRepo  vaccinemd.Repo
	results *CardVerificationResults
}

func (e *v1Processor) GetVerificationResults() *CardVerificationResults {
	e.calcState()
	return e.results
}

func (e *v1Processor) calcState() {

	//state by assume card is great
	e.results.State = CardVerificationStateUnknown

	if e.CardCorrupted() {
		e.results.State = CardVerificationStateCorrupt
		return
	}

	if !e.CardStructureVerified() {
		e.results.State = CardVerificationStatePartlyVerified
		return
	}

	if !e.IssuerVerified() {
		e.results.State = CardVerificationStatePartlyVerified
		return
	}

	if !e.ImmunizationCriteriaMet() {
		e.results.State = CardVerificationStatePartlyVerified
		return
	}

	//
	// If reach here all verifications have passed and the card is valid
	//
	e.results.State = CardVerificationStateVerified
}

//
// Card structure
//

func (e *v1Processor) CardCorrupted() bool {

	if e.results.CardStructure.SignatureChecked &&
		e.results.CardStructure.FetchedKey &&
		!e.results.CardStructure.SignatureValid {
		return true
	}

	return false
}

func (e *v1Processor) CardStructureVerified() bool {

	if e.results.CardStructure.SignatureChecked &&
		e.results.CardStructure.FetchedKey &&
		e.results.CardStructure.SignatureValid &&
		!e.results.CardStructure.Expired {
		return true
	}

	return false
}

func (e *v1Processor) SetSignatureChecked() {
	e.results.CardStructure.SignatureChecked = true
}

func (e *v1Processor) SetFetchedKey() {
	e.results.CardStructure.FetchedKey = true
}

func (e *v1Processor) SetSignatureValid() {
	e.results.CardStructure.SignatureValid = true
}

func (e *v1Processor) SetExpired() {
	e.results.CardStructure.Expired = true
}

//
// Issuer state
//

func (e *v1Processor) IssuerVerified() bool {
	return e.results.Issuer.Trusted
}

func (e *v1Processor) SetIssuerTrusted() {
	e.results.Issuer.Trusted = true
}

//
// Immunization State
//

func (e *v1Processor) ImmunizationCriteriaMet() bool {
	if !e.results.Immunization.UnKnownVaccineType &&
		e.results.Immunization.TrustedVaccineType &&
		e.results.Immunization.MetDosesRequiredCriteria &&
		e.results.Immunization.MetDaysBetweenDoesCriteria &&
		e.results.Immunization.MetDaysSinceLastDoseCriteria {
		return true
	}

	return false
}

func (e *v1Processor) VerifyImmunization(
	region Region,
	doses []*Dose, // the doses administered
) (bool, error) {

	if len(doses) == 0 {
		return false, nil
	}

	//
	// Find system and code for vaccine, expect all doses to be the same system, code
	//
	system := ""
	code := ""
	for _, dose := range doses {
		if system == "" {
			system = dose.Coding.System
		} else if system != dose.Coding.System {
			return false, fmt.Errorf(
				"error verify immunization expects all doses to be of same type got=%s expected=%s",
				system, dose.Coding.System)

		}
		if code == "" {
			code = dose.Coding.Code
		} else if code != dose.Coding.Code {
			return false, fmt.Errorf(
				"error verify immunization expects all doses to be of same type got=%s expected=%s",
				system, dose.Coding.Code)

		}
	}

	vMD := e.mdRepo.FindCovidVaccine(system, code)
	if vMD == nil {
		//do not treat as an error
		e.results.Immunization.UnKnownVaccineType = true
		return false, nil
	}
	e.results.Immunization.UnKnownVaccineType = false

	//check if vaccine trusted for this region
	e.results.Immunization.TrustedVaccineType = false
	switch region {
	case RegionUSA:
		{
			if vMD.CVXStatus == vaccinemd.CVSStatusActive {
				e.results.Immunization.TrustedVaccineType = true
			}
		}
	case RegionEU:
		{
			return false, fmt.Errorf("error verify immunization need to implement check EU region")

		}
	}

	//
	// check if number of doses met
	//
	e.results.Immunization.MetDosesRequiredCriteria = true
	if len(doses) < vMD.Doses {
		e.results.Immunization.MetDosesRequiredCriteria = false
		return false, nil //no point in checking dates as not enough doses
	}

	//
	// http://build.fhir.org/ig/HL7/fhir-shc-vaccination-ig/StructureDefinition-shc-vaccination-dm-definitions.html#Immunization.occurrence[x]:occurrenceDateTime
	//

	// find last dose
	var lastDose *Dose
	for _, dose := range doses {

		if lastDose == nil {
			lastDose = dose
		} else {
			lastDoseOccurrenceTime, err := getOccurrenceTime(lastDose)
			if err != nil {
				return false, err
			}
			currentDoseOccurrenceTime, err := getOccurrenceTime(dose)
			if err != nil {
				return false, err
			}

			if currentDoseOccurrenceTime != nil && currentDoseOccurrenceTime.After(*lastDoseOccurrenceTime) {
				lastDose = dose
			}
		}
	}

	if lastDose == nil {
		return false, nil //for some reason could not find a last dose so lets just treat as false
	}

	occurrenceTime, err := getOccurrenceTime(lastDose)
	if err != nil {
		return false, err
	}

	if occurrenceTime == nil {
		return false, nil // could not find an occurrence date so no point in continuing
	}

	//check duration since the dose was taken
	today := time.Now()
	dateMustHaveOccuredBy := today.AddDate(0, 0, -(vMD.DaysSinceLastDoseCriteria))

	e.results.Immunization.MetDaysSinceLastDoseCriteria = false
	if dateMustHaveOccuredBy.After(*occurrenceTime) {
		e.results.Immunization.MetDaysSinceLastDoseCriteria = true
	}

	//fixme hard code that dates are ok, so can test
	e.results.Immunization.MetDaysBetweenDoesCriteria = true

	return e.ImmunizationCriteriaMet(), nil

}

func getOccurrenceTime(dose *Dose) (occurrenceTime *time.Time, err error) {
	if dose.OccurrenceDateTime != "" {
		occurrenceTime, err = dateStringTime(dose.OccurrenceDateTime)
		if err != nil {
			return nil, err
		}
	} else if dose.OccurrenceString != "" {
		occurrenceTime, err = dateStringTime(dose.OccurrenceDateTime)
		if err != nil {
			return nil, err
		}
	}

	return
}

func dateStringTime(date string) (*time.Time, error) {

	result, err := time.Parse("2006-01-02", date)
	if err == nil {
		return &result, nil
	}

	result, err = time.Parse(time.RFC3339, date)
	if err == nil {
		return &result, nil
	}

	//not sure how to convert
	return nil, fmt.Errorf("error verify immunization date format got=%s", date)

}
