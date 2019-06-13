package migration

import (
	"kayacredit/kc"
	"kayacredit/models"
	"time"
)

func Seed() {
	seeder := kc.App.DB.Begin()

	if kc.App.ENV == "development" {
		borrowers := []models.Borrower{
			models.Borrower{
				Fullname:             "Bene Cucumbatch",
				Gender:               "M",
				IdCardNumber:         "94918376436273",
				TaxIDnumber:          "933393777294",
				Email:                "email@mail.com",
				Birthday:             time.Now(),
				Birthplace:           "herender",
				LastEducation:        "master",
				MotherName:           "mom",
				Phone:                "0812345654321",
				MarriedStatus:        "single",
				SpouseName:           "minnie mouse",
				SpouseBirthday:       time.Now(),
				SpouseLastEducation:  "master",
				Dependants:           0,
				Address:              "sserdda street",
				Province:             "vinpro",
				City:                 "citcit",
				NeighbourAssociation: "erte",
				Hamlets:              "hamham",
				HomePhoneNumber:      "021342124",
				Subdistrict:          "xyzabcd",
				UrbanVillage:         "abcde",
				HomeOwnership:        "privately owned",
				LivedFor:             5,
				Occupation:           "magician",
				EmployeeID:           "08",
				EmployerName:         "Stark Industries",
				EmployerAddress:      "Cambridge",
				Department:           "Magic",
				BeenWorkingFor:       2,
				DirectSuperior:       "Tony Stark",
				EmployerNumber:       "08123546126",
				MonthlyIncome:        5000000,
				OtherIncome:          2000000,
				OtherIncomeSource:    "freelance",
				FieldOfWork:          "building architechture",
				RelatedPersonName:    "spiderman",
				RelatedPhoneNumber:   "0812341234",
				RelatedHomePhone:     "456765434",
				Password:             "password",
			},
		}
		for _, borrower := range borrowers {
			seeder.Create(&borrower)
		}
	}

	seeder.Commit()
}
