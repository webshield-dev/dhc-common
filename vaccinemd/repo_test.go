package vaccinemd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webshield-dev/dhc-common/vaccinemd"
)

func Test_FindVaccine(t *testing.T) {

	type testCase struct {
		name   string
		code   string
		system string
		found  bool
		expectedManufacturerName string
	}

	testCases := []testCase{
		{
			name:   "should find a known code",
			system: "http://hl7.org/fhir/sid/cvx",
			code:   "208",
			found:  true,
			expectedManufacturerName: "Pfizer",
		},
		{
			name:   "should not find a unknown code",
			system: "bogus",
			code:   "208",
			found:  false,
		},
	}

	repo := vaccinemd.MakeRepo()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			vmd := repo.FindCovidVaccine(tc.system, tc.code)
			if tc.found {
				require.NotNil(t, vmd)
			} else {
				require.Nil(t, vmd)
			}

		})
	}

}
