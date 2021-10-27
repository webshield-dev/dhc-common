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
			DisplayName:                  "Moderna",
			SaleProprietaryName:          "Moderna COVID-19 Vaccine",
			ManufacturerName:             "Moderna US, Inc",
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
			DisplayName:                  "Pfizer",
			SaleProprietaryName:          "Pfizer-BioNTech COVID-19 Vaccine",
			ManufacturerName:             "Pfizer-BioNTech",
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
			DisplayName:               "AstraZeneca",
			SaleProprietaryName:       "AstraZeneca COVID-19 Vaccine",
			ManufacturerName:          "AstraZeneca Pharmaceuticals LP",
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
			DisplayName:               "Johnson & Johnson Janssen",
			SaleProprietaryName:       "Janssen COVID-19 Vaccine",
			ManufacturerName:          "Janssen Products, LP",
		},
	}

	return result
}
