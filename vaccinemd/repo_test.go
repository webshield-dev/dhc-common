package vaccinemd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webshield-dev/dhc-common/vaccinemd"
)

func Test_FindVaccine(t *testing.T) {

	type testCase struct {
		name                     string
		code                     string
		system                   string
		found                    bool
		expectedManufacturerName string
	}

	testCases := []testCase{
		{
			name:                     "should find a known code",
			system:                   "http://hl7.org/fhir/sid/cvx",
			code:                     "208",
			found:                    true,
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

	require.Equal(t, 4, len(repo.CovidVaccines()), "should return all vaccines")

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			vmd := repo.FindCovidVaccine(tc.system, tc.code)
			if tc.found {
				require.NotNil(t, vmd)

				vmd2 := repo.FindCovidVaccineByID(vmd.ID)
				require.Equal(t, vmd, vmd2)

			} else {
				require.Nil(t, vmd)
			}

		})
	}

}

func Test_FindTrustedVaccine(t *testing.T) {

	type testCase struct {
		name                string
		region              vaccinemd.Region
		expectedResultCount int
	}

	testCases := []testCase{
		{
			name:                "should find trusted for USA",
			region:              vaccinemd.RegionUSA,
			expectedResultCount: 3,
		},
	}

	repo := vaccinemd.MakeRepo()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			vmd := repo.FindTrustedVaccinesForRegion(tc.region)
			require.Equal(t, tc.expectedResultCount, len(vmd))
		})
	}

}
