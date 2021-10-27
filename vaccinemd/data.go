package vaccinemd

// createCovidVaccineMetadata
func createCovidVaccineMetadata() []*CovidVaccineMetadata {

	result := []*CovidVaccineMetadata{
		{
			Codes: []Coding{
				{
					System: CVXSystem,
					Code:   "207",
				},
			},
			CVXStatus:                    CVSStatusActive,
			Doses:                        2,
			DaysSinceLastDoseCriteria:    14,
			DaysBetweenDoesCriteriaBegin: 24,
			DaysBetweenDoesCriteriaEnd:   92,
			SaleProprietaryName:          "Moderna COVID-19 Vaccine",
			ManufacturerName:             "Moderna US Inc",
		},
		{
			Codes: []Coding{
				{
					System: CVXSystem,
					Code:   "208",
				},
			},
			CVXStatus:                    CVSStatusActive,
			Doses:                        2,
			DaysSinceLastDoseCriteria:    14,
			DaysBetweenDoesCriteriaBegin: 17,
			DaysBetweenDoesCriteriaEnd:   92,
			SaleProprietaryName:          "Pfizer-BioNTech COVID-19 Vaccine",
			ManufacturerName:             "Pfizer, Inc",
		},
		{
			Codes: []Coding{
				{
					System: CVXSystem,
					Code:   "210",
				},
			},
			CVXStatus:                 CVSStatusNonUS,
			Doses:                     2,
			DaysSinceLastDoseCriteria: 14,
			SaleProprietaryName:       "AstraZeneca COVID-19 Vaccine",
			ManufacturerName:          "AstraZeneca",
		},
		{
			Codes: []Coding{
				{
					System: CVXSystem,
					Code:   "212",
				},
			},
			CVXStatus:                 CVSStatusActive,
			Doses:                     1,
			DaysSinceLastDoseCriteria: 14,
			SaleProprietaryName:       "Janssen COVID-19 Vaccine",
			ManufacturerName:          "Janssen",
		},
	}

	return result
}
