package verification

// CardVerificationState the card's verification state, see below
type CardVerificationState string

const (

	//CardVerificationStateUnknown no verifications have been performed
	CardVerificationStateUnknown CardVerificationState = "unknown"

	//CardVerificationStateVerified the card has been verified and all checks have passed
	CardVerificationStateVerified CardVerificationState = "verified"

	//CardVerificationStateInvalid one or more verifications have failed
	CardVerificationStateInvalid CardVerificationState = "invalid"

	//CardVerificationStateCorrupt the digital signature and is invalid
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

	//SignatureNotChecked for some reason the verifier choose to not verify the card signature
	SignatureNotChecked bool `json:"signature_not_checked"`

	//FetchedKey true if successful fetched the verification key
	FetchedKey bool `json:"fetched_key"`

	//SignatureValid true if checked the signature an is valid
	SignatureValid bool `json:"signature_valid"`

	//Expired true if an exp date and has passed
	Expired bool `json:"expired"`
}

//IssuerVerificationResults issuer verification results
type IssuerVerificationResults struct{

	//Trusted issue is on a trusted whitelist
	Trusted bool `json:"trusted"`

}

// ImmunizationVerificationResults immunization verification results
type ImmunizationVerificationResults struct{}
