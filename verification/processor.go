package verification

//Processor can be created by a verifier to manage the verification state and calculate cards verification state
//not designed to be thread safe. Create one per card verification
type Processor interface {

	//GetResults returns the current result state
	GetResults() *CardVerificationResults

	//
	// Card struct setters
	//

	//SetSignatureNotChecked do not check signature for some reason
	SetSignatureNotChecked()

	//SetFetchedKey record key fetched
	SetFetchedKey()

	//SetSignatureValid record signature valid
	SetSignatureValid()

	//SetExpired record expired
	SetExpired()

	//
	// Issuer verifications
	//

	//SetIssuerTrusted issuer is on a trusted whitelist
	SetIssuerTrusted()
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

	//
	// Check if card is corrupted
	//
	if !e.results.CardStructure.SignatureNotChecked && //the signature is being checked
		e.results.CardStructure.FetchedKey &&
		!e.results.CardStructure.SignatureValid {
		e.results.State = CardVerificationStateCorrupt
		return
	}

	//
	//see if any of the card structure verifications failed if so the card is invalid
	//
	if e.results.CardStructure.SignatureNotChecked ||
		!e.results.CardStructure.FetchedKey ||
		!e.results.CardStructure.SignatureValid ||
		e.results.CardStructure.Expired {
		e.results.State = CardVerificationStateInvalid
		return
	}

	//
	// see if any of the issuer verifications have failed, if so the card is invalid
	//
	if !e.results.Issuer.Trusted {
		e.results.State = CardVerificationStateInvalid
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

func (e *v1Processor) SetSignatureNotChecked() {
	e.results.CardStructure.SignatureNotChecked = true
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

func (e *v1Processor) SetIssuerTrusted() {
	e.results.Issuer.Trusted = true
}
