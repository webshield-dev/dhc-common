package verification_test

import (
	"github.com/stretchr/testify/require"
	"github.com/webshield-dev/dhc-common/vaccinemd"
	"github.com/webshield-dev/dhc-common/verification"
	"testing"
	"time"
)

func Test_CardStateNotVerified(t *testing.T) {
	//checks that card signature can be verified and issuer is trusted

	type testCase struct {
		name                          string
		cardStructureResults          *verification.CardStructureVerificationResults
		issuerResults *verification.IssuerVerificationResults
		expectedState                 verification.CardVerificationState
		expectedCardStructureVerified bool
	}

	testCases := []testCase{
		{
			name:                          "all verifications met so should be valid",
			expectedState:                 verification.CardVerificationStateValid,
			expectedCardStructureVerified: true,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       true,
				SignatureValid:   true,
				Expired:          false,
			},
			issuerResults: &verification.IssuerVerificationResults{
				Trusted: true,
			},
		},
		{
			name:                          "if issuer not trusted should not be verified",
			expectedState:                 verification.CardVerificationStateNotVerified,
			expectedCardStructureVerified: true,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       true,
				SignatureValid:   true,
				Expired:          false,
			},
			issuerResults: &verification.IssuerVerificationResults{
				Trusted: false,
			},
		},
		{
			name:                          "should be corrupt if signature not valid",
			expectedState:                 verification.CardVerificationStateCorrupt,
			expectedCardStructureVerified: false,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       true,
				SignatureValid:   false,
				Expired:          false,
			},
			issuerResults: &verification.IssuerVerificationResults{
				Trusted: true,
			},
		},
		{
			name:                          "should not verify if signature not checked",
			expectedState:                 verification.CardVerificationStateNotVerified,
			expectedCardStructureVerified: false,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: false,
				FetchedKey:       false,
				SignatureValid:   false,
				Expired:          false,
			},
			issuerResults: &verification.IssuerVerificationResults{
				Trusted: true,
			},
		},
		{
			name:                          "should not verify if did not fetch key",
			expectedState:                 verification.CardVerificationStateNotVerified,
			expectedCardStructureVerified: false,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       false,
				SignatureValid:   false,
				Expired:          false,
			},
			issuerResults: &verification.IssuerVerificationResults{
				Trusted: true,
			},
		},
		{
			name:                          "should not verify if card expired",
			expectedState:                 verification.CardVerificationStateNotVerified,
			expectedCardStructureVerified: false,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       true,
				SignatureValid:   true,
				Expired:          true,
			},
			issuerResults: &verification.IssuerVerificationResults{
				Trusted: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			processor := verification.NewProcessor()

			//set so immunization ok, as want to just test if verified attributes alter result
			setImmunizationResultsOK(t, processor)

			if tc.issuerResults != nil {
				if tc.issuerResults.Trusted {
					processor.SetIssuerTrusted()
				}
			}

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
			require.Equal(t, tc.expectedCardStructureVerified, results.CardStructure.AllChecksPassed)

			if tc.issuerResults != nil {
				if tc.issuerResults.Trusted {
					require.True(t, processor.IssuerVerified(), "issuer should be verified")
					require.True(t, results.Issuer.AllChecksPassed, "issuer all checks should have passed")
				} else {
					require.False(t, processor.IssuerVerified(), "issuer verified not expected")
					require.False(t, results.Issuer.AllChecksPassed, "issuer all checks should not have passed")

				}
			}

			require.True(t, processor.ImmunizationCriteriaMet(), "imm should be be met as set to ok")

		})
	}
}

func Test_SafetyCriteriaNotMet(t *testing.T) {

	type testCase struct {
		name                            string
		region                          vaccinemd.Region
		doses                           []*verification.Dose
		expectedState                   verification.CardVerificationState
		expectedMetImmunizationCriteria bool
	}

	testCases := []testCase{
		{
			name:          "all criteria met two doses",
			expectedState: verification.CardVerificationStateValid,
			region:        vaccinemd.RegionUSA,
			doses: []*verification.Dose{
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "207", //moderna
					},

					OccurrenceDateTime: "2021-03-16",
				},
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "207", //moderna
					},

					OccurrenceDateTime: "2021-04-06",
				},
			},
			expectedMetImmunizationCriteria: true,
		},
		{
			name:          "all criteria met one doses",
			expectedState: verification.CardVerificationStateValid,
			region:        vaccinemd.RegionUSA,
			doses: []*verification.Dose{
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "212", //janseen
					},

					OccurrenceDateTime: "2021-03-16",
				},
			},
			expectedMetImmunizationCriteria: true,
		},
		{
			name:          "all criteria met doses array order not based on time",
			expectedState: verification.CardVerificationStateValid,
			region:        vaccinemd.RegionUSA,
			doses: []*verification.Dose{
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "207", //moderna
					},

					OccurrenceDateTime: "2021-04-06",
				},
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "207", //moderna
					},

					OccurrenceDateTime: "2021-03-16",
				},
			},
			expectedMetImmunizationCriteria: true,
		},
		{
			name:                            "criteria not met as no doses passed",
			expectedState:                   verification.CardVerificationStateSafetyCriteriaNotMet,
			region:                          vaccinemd.RegionUSA,
			doses:                           nil,
			expectedMetImmunizationCriteria: false,
		},
		{
			name:          "criteria not met as need two does and pass one",
			expectedState: verification.CardVerificationStateSafetyCriteriaNotMet,
			region:        vaccinemd.RegionUSA,
			doses: []*verification.Dose{
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "207", //moderna
					},
				},
			},
			expectedMetImmunizationCriteria: false,
		},
		{
			name:          "criteria not met as passed two does but no occurrence date so cannot get dates",
			expectedState: verification.CardVerificationStateSafetyCriteriaNotMet,
			region:        vaccinemd.RegionUSA,
			doses: []*verification.Dose{
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "207", //moderna
					},
				},
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "207", //moderna
					},
				},
			},
			expectedMetImmunizationCriteria: false,
		},
		{
			name:          "criteria NOT met one dose ok but occurence data too soon",
			expectedState: verification.CardVerificationStateSafetyCriteriaNotMet,
			region:        vaccinemd.RegionUSA,
			doses: []*verification.Dose{
				{
					Coding: vaccinemd.Coding{
						System: vaccinemd.CVXSystem,
						Code:   "212", //janseen
					},

					OccurrenceDateTime: time.Now().Format("2006-01-02"),
				},
			},
			expectedMetImmunizationCriteria: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			processor := verification.NewProcessor()
			setCardStructureOK(processor)
			setIssuerResultsOK(processor)

			immVerifed, err := processor.VerifyImmunization(tc.region, tc.doses)
			require.NoError(t, err)
			require.Equal(t, tc.expectedMetImmunizationCriteria, immVerifed)

			//check results
			results := processor.GetVerificationResults()
			require.Equal(t, tc.expectedState, results.State)
			require.True(t, processor.CardStructureVerified(), "card structure verified not expected")
			require.True(t, processor.IssuerVerified(), "issuer verified not expected")
			require.Equal(t, tc.expectedMetImmunizationCriteria, processor.ImmunizationCriteriaMet(), "imm met not expected")
			require.Equal(t, tc.expectedMetImmunizationCriteria, results.Immunization.AllChecksPassed, "all checks passed not expected")

		})
	}
}

//Helpers

//setCardStructureOK set so ok
func setCardStructureOK(ps verification.Processor) {
	ps.SetSignatureChecked()
	ps.SetSignatureValid()
	ps.SetFetchedKey()
}

//setIssuerResultOK set so ok
func setIssuerResultsOK(ps verification.Processor) {
	ps.SetIssuerTrusted()
}

//setImmunizationResultsOK set so ok
func setImmunizationResultsOK(t *testing.T, ps verification.Processor)  {
	doses :=  []*verification.Dose{
		{
			Coding: vaccinemd.Coding{
				System: vaccinemd.CVXSystem,
				Code:   "207", //moderna
			},

			OccurrenceDateTime: "2021-03-16",
		},
		{
			Coding: vaccinemd.Coding{
				System: vaccinemd.CVXSystem,
				Code:   "207", //moderna
			},

			OccurrenceDateTime: "2021-04-06",
		},
	}

	immVerifed, err := ps.VerifyImmunization(vaccinemd.RegionUSA, doses)
	require.NoError(t, err)
	require.True(t, immVerifed, "should be verified")
}
