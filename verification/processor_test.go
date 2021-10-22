package verification_test

import (
    "github.com/stretchr/testify/require"
    "github.com/webshield-dev/dhc-common/verification"
    "testing"
)

func Test_FindVaccine(t *testing.T) {

	type testCase struct {
		name string
        results *verification.CardVerificationResults
        expectedState verification.CardVerificationState
	}

	testCases := []testCase{
		{
			name: "check verified card",
            expectedState: verification.CardVerificationStateVerified,
            results: &verification.CardVerificationResults{
                CardStructure: &verification.CardStructureVerificationResults{
                    SignatureNotChecked: false,
                    FetchedKey: true,
                    SignatureValid: true,
                    Expired: false,
                },
                Issuer: &verification.IssuerVerificationResults{
                    Trusted: true,
                },
            },
		},
        {
            name: "check corrupt card",
            expectedState: verification.CardVerificationStateCorrupt,
            results: &verification.CardVerificationResults{
                CardStructure: &verification.CardStructureVerificationResults{
                    SignatureNotChecked: false,
                    FetchedKey: true,
                    SignatureValid: false,
                    Expired: false,
                },
                Issuer: &verification.IssuerVerificationResults{
                    Trusted: false,
                },
            },
        },
        {
            name: "check invalid card, could not fetch verification key",
            expectedState: verification.CardVerificationStateInvalid,
            results: &verification.CardVerificationResults{
                CardStructure: &verification.CardStructureVerificationResults{
                    SignatureNotChecked: false,
                    FetchedKey: false,
                    SignatureValid: false,
                    Expired: false,
                },
                Issuer: &verification.IssuerVerificationResults{
                    Trusted: true,
                },
            },
        },
        {
            name: "check invalid card, card has expired ",
            expectedState: verification.CardVerificationStateInvalid,
            results: &verification.CardVerificationResults{
                CardStructure: &verification.CardStructureVerificationResults{
                    SignatureNotChecked: false,
                    FetchedKey: true,
                    SignatureValid: true,
                    Expired: true,
                },
                Issuer: &verification.IssuerVerificationResults{
                    Trusted: true,
                },
            },
        },

        {
            name: "check invalid card, issuer is not trusted",
            expectedState: verification.CardVerificationStateInvalid,
            results: &verification.CardVerificationResults{
                CardStructure: &verification.CardStructureVerificationResults{
                    SignatureNotChecked: false,
                    FetchedKey: true,
                    SignatureValid: true,
                    Expired: false,
                },
                Issuer: &verification.IssuerVerificationResults{
                    Trusted: false,
                },
            },
        },

	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

            processor := verification.NewProcessor()

            cs := tc.results.CardStructure
            if cs.SignatureNotChecked {
                processor.SetSignatureNotChecked()
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

            is := tc.results.Issuer
            if is.Trusted {
                processor.SetIssuerTrusted()
            }

            results := processor.GetResults()
            require.Equal(t, tc.expectedState, results.State)

		})
	}
}
