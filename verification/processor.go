package verification

//Processor can be created by a verifier to manage the verification state and calculate cards verification state
//not designed to be thread safe. Create one per card verification
type Processor interface {

	//GetResults returns the current result state
	GetResults() *CardVerificationResults

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

	//ImmunizationCriteriaMet true if all the immunization criteria have been met
	ImmunizationCriteriaMet() bool

	//SetTrustedVaccineType set as trusted vaccine
	SetTrustedVaccineType()

	//SetMetDosedRequiredCriteria set met doses required
	SetMetDosedRequiredCriteria()

	//SetMetDaysBetweenDoesCriteria set met days between
	SetMetDaysBetweenDoesCriteria()

	//SetMetDaysSinceLastDoseCriteria set met days since last dose
	SetMetDaysSinceLastDoseCriteria()
}

func NewProcessor() Processor {

	return &v1Processor{
		results: &CardVerificationResults{
			State:         CardVerificationStateUnknown,
			CardStructure: &CardStructureVerificationResults{},
			Issuer:        &IssuerVerificationResults{},
			Immunization:  &ImmunizationVerificationResults{},
		},
	}

}

type v1Processor struct {
	results *CardVerificationResults
}

func (e *v1Processor) GetResults() *CardVerificationResults {
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

	if 	e.results.CardStructure.SignatureChecked &&
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
	if e.results.Immunization.TrustedVaccineType &&
		e.results.Immunization.MetDosedRequiredCriteria &&
		e.results.Immunization.MetDaysBetweenDoesCriteria &&
		e.results.Immunization.MetDaysSinceLastDoseCriteria {
		return true
	}

	return false
}


func (e *v1Processor) SetTrustedVaccineType() {
	e.results.Immunization.TrustedVaccineType = true
}

func (e *v1Processor) SetMetDosedRequiredCriteria() {
	e.results.Immunization.MetDosedRequiredCriteria = true
}

func (e *v1Processor) SetMetDaysBetweenDoesCriteria() {
	e.results.Immunization.MetDaysBetweenDoesCriteria = true
}

func (e *v1Processor) SetMetDaysSinceLastDoseCriteria() {
	e.results.Immunization.MetDaysSinceLastDoseCriteria = true
}


