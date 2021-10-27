package vaccinemd

//
// SEE https://www.cdc.gov/vaccines/programs/iis/COVID-19-related-codes.html
//



//MakeRepo fixme in the future pass in the metadata config
func MakeRepo() Repo {

	vaccineMD := createCovidVaccineMetadata()

	code2CodingMap := make(map[string]*CovidVaccineMetadata)

	for _, vmd := range vaccineMD {

		for _, code := range vmd.Codes {
			//code is unique within system, start with code as more unique
			key := string(code.Code) + "_" + code.System
			code2CodingMap[key] = vmd
		}
	}

	return &v1Repo{vaccineMD: vaccineMD, code2CodingMap: code2CodingMap}
}

//Repo provides methods to find out vaccine info
type Repo interface {

	//FindCovidVaccine return vaccine metadata if the passed in coding is known CovidVaccine
	FindCovidVaccine(system string, code string) *CovidVaccineMetadata


	//FindTrustedVaccinesForRegion find for the specified region
	FindTrustedVaccinesForRegion(region Region) []*CovidVaccineMetadata
}


type v1Repo struct {
	vaccineMD []*CovidVaccineMetadata
	code2CodingMap map[string]*CovidVaccineMetadata
}

func (vmi *v1Repo) FindCovidVaccine(system string, code string) *CovidVaccineMetadata {

	key := code + "_" + system
	return vmi.code2CodingMap[key]

}


func (vmi *v1Repo) FindTrustedVaccinesForRegion(region Region) []*CovidVaccineMetadata {

	result := make([]*CovidVaccineMetadata, 0)
	for _, md := range vmi.vaccineMD {
		switch region {
		case RegionUSA: {
			if md.CVXStatus == CVSStatusActive {
				result = append(result, md)
			}
		}
		//fixme how to handle other regions
		}
	}

	return result

}