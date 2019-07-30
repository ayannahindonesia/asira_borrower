package migration

import (
	"asira/asira"
	"asira/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func Seed() {
	seeder := asira.App.DB.Begin()
	defer seeder.Commit()

	if asira.App.ENV == "development" {
		// seed borrowers
		borrowers := []models.Borrower{
			models.Borrower{
				Fullname:             "Full Name A",
				Gender:               "M",
				IdCardNumber:         "9876123451234567789",
				TaxIDnumber:          "0987654321234567890",
				Email:                "emaila@domain.com",
				Birthday:             time.Now(),
				Birthplace:           "a birthplace",
				LastEducation:        "a last edu",
				MotherName:           "a mom",
				Phone:                "081234567890",
				MarriedStatus:        "single",
				SpouseName:           "a spouse",
				SpouseBirthday:       time.Now(),
				SpouseLastEducation:  "master",
				Dependants:           0,
				Address:              "a street address",
				Province:             "a province",
				City:                 "a city",
				NeighbourAssociation: "a rt",
				Hamlets:              "a rw",
				HomePhoneNumber:      "021837163",
				Subdistrict:          "a camat",
				UrbanVillage:         "a lurah",
				HomeOwnership:        "privately owned",
				LivedFor:             5,
				Occupation:           "accupation",
				EmployerName:         "amployer",
				EmployerAddress:      "amployer address",
				Department:           "a department",
				BeenWorkingFor:       2,
				DirectSuperior:       "a boss",
				EmployerNumber:       "02188776655",
				MonthlyIncome:        5000000,
				OtherIncome:          2000000,
				RelatedPersonName:    "a big sis",
				RelatedPhoneNumber:   "08987654321",
				OTPverified:          true,
				BankAccountNumber:    "520384716",
				Password:             "password",
			},
			models.Borrower{
				Fullname:             "Full Name B",
				Gender:               "F",
				IdCardNumber:         "9876123451234567781",
				TaxIDnumber:          "0987654321234567891",
				Email:                "emailb@domain.com",
				Birthday:             time.Now(),
				Birthplace:           "b birthplace",
				LastEducation:        "b last edu",
				MotherName:           "b mom",
				Phone:                "081234567891",
				MarriedStatus:        "single",
				SpouseName:           "b spouse",
				SpouseBirthday:       time.Now(),
				SpouseLastEducation:  "master",
				Dependants:           0,
				Address:              "b street address",
				Province:             "b province",
				City:                 "b city",
				NeighbourAssociation: "b rt",
				Hamlets:              "b rw",
				HomePhoneNumber:      "021837163",
				Subdistrict:          "b camat",
				UrbanVillage:         "b lurah",
				HomeOwnership:        "privately owned",
				LivedFor:             5,
				Occupation:           "bccupation",
				EmployerName:         "bmployer",
				EmployerAddress:      "bmployer address",
				Department:           "b department",
				BeenWorkingFor:       2,
				DirectSuperior:       "b boss",
				EmployerNumber:       "02188776655",
				MonthlyIncome:        5000000,
				OtherIncome:          2000000,
				RelatedPersonName:    "b big sis",
				RelatedPhoneNumber:   "08987654321",
				RelatedAddress:       "big sis address",
				OTPverified:          false,
				Password:             "password",
			},
		}
		for _, borrower := range borrowers {
			borrower.Create()
		}
		images := []models.Image{
			models.Image{
				Image_string: "iVBORw0KGgoAAAANSUhEUgAAAP0AAACnCAYAAADaIFptAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAHiSURBVHhe7dMBAQAACMMg+5e+QQYduAEp0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPSQsj3nqr5USY34xwAAAABJRU5ErkJggg==",
			},
			models.Image{
				Image_string: "iVBORw0KGgoAAAANSUhEUgAAAP0AAACnCAYAAADaIFptAAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAHiSURBVHhe7dMBAQAACMMg+5e+QQYduAEp0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPQQIz3ESA8x0kOM9BAjPcRIDzHSQ4z0ECM9xEgPMdJDjPSQsj3nqr5USY34xwAAAABJRU5ErkJggg==",
			},
		}
		for _, image := range images {
			image.Create()
		}
		// seed loans
		loans := []models.Loan{
			models.Loan{
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				LoanAmount:       1000000,
				Installment:      6,
				LoanIntention:    "a loan 1 intention",
				IntentionDetails: "a loan 1 intention details",
			},
			models.Loan{
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Status:           "accepted",
				LoanAmount:       500000,
				Installment:      2,
				LoanIntention:    "a loan 2 intention",
				IntentionDetails: "a loan 2 intention details",
			},
			models.Loan{
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Status:           "rejected",
				LoanAmount:       2000000,
				Installment:      8,
				LoanIntention:    "a loan 3 intention",
				IntentionDetails: "a loan 3 intention details",
			},
		}
		for _, loan := range loans {
			loan.Create()
		}
	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"banks",
				"bank_products",
				"images",
				"borrowers",
				"loans",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = asira.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
