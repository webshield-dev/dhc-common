package verification_test

import (
    "github.com/stretchr/testify/require"
    "github.com/webshield-dev/dhc-common/verification"
    "testing"
)

func Test_VerifyCardStructure(t *testing.T) {

	type testCase struct {
		name string
        cardStructureResults *verification.CardStructureVerificationResults
        expectedState verification.CardVerificationState
        expectedCardStructureVerified bool
	}

	testCases := []testCase{
		{
			name: "all verifications passed",
            expectedState: verification.CardVerificationStatePartlyVerified,
            expectedCardStructureVerified: true,
            cardStructureResults: &verification.CardStructureVerificationResults{
                    SignatureChecked: true,
                    FetchedKey: true,
                    SignatureValid: true,
                    Expired: false,
                },
		},
        {
            name: "should be corrupt if signature not valid",
            expectedState: verification.CardVerificationStateCorrupt,
            expectedCardStructureVerified: false,
            cardStructureResults: &verification.CardStructureVerificationResults{
                SignatureChecked: true,
                FetchedKey: true,
                SignatureValid: false,
                Expired: false,
            },
        },
        {
            name: "should not verify if signature not checked",
            expectedState: verification.CardVerificationStatePartlyVerified,
            expectedCardStructureVerified: false,
            cardStructureResults: &verification.CardStructureVerificationResults{
                SignatureChecked: false,
                FetchedKey: false,
                SignatureValid: false,
                Expired: false,
            },
        },
        {
            name: "should not verify if did not fetch key",
            expectedState: verification.CardVerificationStatePartlyVerified,
            expectedCardStructureVerified: false,
            cardStructureResults: &verification.CardStructureVerificationResults{
                SignatureChecked: true,
                FetchedKey: false,
                SignatureValid: false,
                Expired: false,
            },
        },
        {
            name: "should not verify if card expired",
            expectedState: verification.CardVerificationStatePartlyVerified,
            expectedCardStructureVerified: false,
            cardStructureResults: &verification.CardStructureVerificationResults{
                SignatureChecked: true,
                FetchedKey: true,
                SignatureValid: true,
                Expired: true,
            },
        },

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

            processor := verification.NewProcessor()

            cs := tc.cardStructureResults
            if cs.SignatureChecked {
                processor.SetSignatureChecked()
            }
            if cs.FetchedKey {
                processor.SetFetchedKey()
            }
            if cs.SignatureValid {
                processor.SetSignatureValid()
            }
            if cs.Expired {
                processor.SetExpired()
            }


            results := processor.GetVerificationResults()
            require.Equal(t, tc.expectedState, results.State)
            require.Equal(t, tc.expectedCardStructureVerified, processor.CardStructureVerified(), "card structure verified not expected")
            require.False(t, processor.IssuerVerified(), "issuer verified not expected")
            require.False(t, processor.ImmunizationCriteriaMet(), "imm met not expected")


		})
	}
}


func Test_VerifyIssuer(t *testing.T) {

    type testCase struct {
        name string
        issuerResults *verification.IssuerVerificationResults
        expectedState verification.CardVerificationState
        expectedIssuerVerified bool
    }

    testCases := []testCase{
        {
            name: "trusted issuer",
            expectedState: verification.CardVerificationStatePartlyVerified,
            expectedIssuerVerified: true,
            issuerResults: &verification.IssuerVerificationResults{
                Trusted: true,
            },
        },
        {
            name: "untrusted issuer",
            expectedState: verification.CardVerificationStatePartlyVerified,
            expectedIssuerVerified: false,
            issuerResults: &verification.IssuerVerificationResults{
                Trusted: false,
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {

            processor := verification.NewProcessor()


            is := tc.issuerResults
            if is.Trusted {
                processor.SetIssuerTrusted()
            }


            results := processor.GetVerificationResults()
            require.Equal(t, tc.expectedState, results.State)
            require.False(t, processor.CardStructureVerified(), "card structure verified not expected")
            require.Equal(t, tc.expectedIssuerVerified, processor.IssuerVerified(), "issuer verified not expected")
            require.False(t, processor.ImmunizationCriteriaMet(), "imm met not expected")
        })
    }
}


func Test_VerifyImmunization(t *testing.T) {

    type testCase struct {
        name string
        doses []verification.Dose
        expectedState verification.CardVerificationState
        expectedImmunizationCriteria bool
    }

    testCases := []testCase{
        {
            name: "check verified card",
            expectedState: verification.CardVerificationStatePartlyVerified,
            doses: nil,
            expectedImmunizationCriteria: false,
        },

    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {

            processor := verification.NewProcessor()


            results := processor.GetVerificationResults()
            require.Equal(t, tc.expectedState, results.State)
            require.False(t, processor.CardStructureVerified(), "card structure verified not expected")
            require.False(t, processor.IssuerVerified(), "issuer verified not expected")
            require.Equal(t, tc.expectedImmunizationCriteria, processor.ImmunizationCriteriaMet(), "imm met not expected")


        })
    }
}

