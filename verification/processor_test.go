package verification_test

import (
	"github.com/stretchr/testify/require"
	"github.com/webshield-dev/dhc-common/vaccinemd"
	"github.com/webshield-dev/dhc-common/verification"
	"testing"
    "time"
)

func Test_VerifyCardStructure(t *testing.T) {

	type testCase struct {
		name                          string
		cardStructureResults          *verification.CardStructureVerificationResults
		expectedState                 verification.CardVerificationState
		expectedCardStructureVerified bool
	}

	testCases := []testCase{
		{
			name:                          "all verifications passed",
			expectedState:                 verification.CardVerificationStatePartlyVerified,
			expectedCardStructureVerified: true,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       true,
				SignatureValid:   true,
				Expired:          false,
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
		},
		{
			name:                          "should not verify if signature not checked",
			expectedState:                 verification.CardVerificationStatePartlyVerified,
			expectedCardStructureVerified: false,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: false,
				FetchedKey:       false,
				SignatureValid:   false,
				Expired:          false,
			},
		},
		{
			name:                          "should not verify if did not fetch key",
			expectedState:                 verification.CardVerificationStatePartlyVerified,
			expectedCardStructureVerified: false,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       false,
				SignatureValid:   false,
				Expired:          false,
			},
		},
		{
			name:                          "should not verify if card expired",
			expectedState:                 verification.CardVerificationStatePartlyVerified,
			expectedCardStructureVerified: false,
			cardStructureResults: &verification.CardStructureVerificationResults{
				SignatureChecked: true,
				FetchedKey:       true,
				SignatureValid:   true,
				Expired:          true,
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
		name                   string
		issuerResults          *verification.IssuerVerificationResults
		expectedState          verification.CardVerificationState
		expectedIssuerVerified bool
	}

	testCases := []testCase{
		{
			name:                   "trusted issuer",
			expectedState:          verification.CardVerificationStatePartlyVerified,
			expectedIssuerVerified: true,
			issuerResults: &verification.IssuerVerificationResults{
				Trusted: true,
			},
		},
		{
			name:                   "untrusted issuer",
			expectedState:          verification.CardVerificationStatePartlyVerified,
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
		name                            string
		region                          verification.Region
		doses                           []*verification.Dose
		expectedState                   verification.CardVerificationState
		expectedMetImmunizationCriteria bool
	}

	testCases := []testCase{
		{
			name:          "all criteria met two doses",
			expectedState: verification.CardVerificationStatePartlyVerified,
			region:        verification.RegionUSA,
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
            expectedState: verification.CardVerificationStatePartlyVerified,
            region:        verification.RegionUSA,
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
			expectedState: verification.CardVerificationStatePartlyVerified,
			region:        verification.RegionUSA,
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
			expectedState:                   verification.CardVerificationStatePartlyVerified,
			region:                          verification.RegionUSA,
			doses:                           nil,
			expectedMetImmunizationCriteria: false,
		},
		{
			name:          "criteria not met as need two does and pass one",
			expectedState: verification.CardVerificationStatePartlyVerified,
			region:        verification.RegionUSA,
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
			expectedState: verification.CardVerificationStatePartlyVerified,
			region:        verification.RegionUSA,
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
            expectedState: verification.CardVerificationStatePartlyVerified,
            region:        verification.RegionUSA,
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

			immVerifed, err := processor.VerifyImmunization(tc.region, tc.doses)
			require.NoError(t, err)
			require.Equal(t, tc.expectedMetImmunizationCriteria, immVerifed)

			//check results
			results := processor.GetVerificationResults()
			require.Equal(t, tc.expectedState, results.State)
			require.False(t, processor.CardStructureVerified(), "card structure verified not expected")
			require.False(t, processor.IssuerVerified(), "issuer verified not expected")
			require.Equal(t, tc.expectedMetImmunizationCriteria, processor.ImmunizationCriteriaMet(), "imm met not expected")

		})
	}
}
