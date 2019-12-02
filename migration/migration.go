package migration

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

func Seed() {
	seeder := asira.App.DB.Begin()
	defer seeder.Commit()

	if asira.App.ENV == "development" {
		// seed internals
		client := []models.Client_config{
			models.Client_config{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Client_config{
				Name:   "android",
				Key:    "androkey",
				Role:   "client",
				Secret: "androsecret",
			},
		}
		for _, clients := range client {
			clients.Create()
		}

		// seed images
		file, _ := os.Open("migration/image_dummy.txt")
		defer file.Close()
		b64image, _ := ioutil.ReadAll(file)
		images := []models.Image{
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
		}
		for _, image := range images {
			image.Create()
		}

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
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
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
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
			},
		}
		for _, borrower := range borrowers {
			borrower.Create()
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
				LoanIntention:    "Pendidikan",
				IntentionDetails: "a loan 1 intention details",
				Product:          1,
			},
			models.Loan{
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Status:           "approved",
				LoanAmount:       500000,
				Installment:      2,
				LoanIntention:    "Rumah Tangga",
				IntentionDetails: "a loan 2 intention details",
				Product:          1,
				OTPverified:      true,
			},
			models.Loan{
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Status:           "rejected",
				LoanAmount:       2000000,
				Installment:      8,
				LoanIntention:    "Kesehatan",
				IntentionDetails: "a loan 3 intention details",
				Product:          1,
				OTPverified:      true,
			},
		}
		for _, loan := range loans {
			loan.Create()
		}

		//agent migration
		agents := []models.Agent{
			models.Agent{
				Name:     "Agent K",
				Username: "agentK",
				Password: "password",
				Email:    "agentk@mib.com",
				Phone:    "081234567890",
				Category: "agent",
				AgentProvider: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Banks:  pq.Int64Array{1, 2},
				Status: "active",
			},
			models.Agent{
				Name:     "Agent J",
				Username: "agentJ",
				Password: "password",
				Email:    "agentj@mib.com",
				Phone:    "081234567891",
				Category: "account_executive",
				Banks:    pq.Int64Array{1},
				Status:   "active",
			},
		}
		for _, agent := range agents {
			agent.Create()
		}

		//seed uuid
		uuid := models.Uuid_Reset_Password{
			UUID: "f4f71eae-2cc9-4289-94e4-2421df67d4d7",
			Borrower: sql.NullInt64{
				Int64: int64(1),
				Valid: true,
			},
		}
		uuid.Create()

		roles := models.InternalRoles{
			Name:        "Admin",
			System:      "Core",
			Description: "Role Admin",
			Status:      true,
		}
		roles.Create()
	}
}

func TestSeed() {
	seeder := asira.App.DB.Begin()
	defer seeder.Commit()

	if asira.App.ENV == "development" {
		// seed images
		file, _ := os.Open("migration/image_dummy.txt")
		defer file.Close()
		b64image, _ := ioutil.ReadAll(file)
		images := []models.Image{
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
			models.Image{
				Image_string: string(b64image),
			},
		}
		for _, image := range images {
			image.Create()
		}

		// seed bank types
		bankTypes := []models.BankType{
			models.BankType{
				Name:        "BPD",
				Description: "Layanan BPD",
			},
			models.BankType{
				Name:        "BPR",
				Description: "Layanan BPR",
			},
			models.BankType{
				Name:        "Koperasi",
				Description: "Layanan Koperasi",
			},
		}
		for _, bankType := range bankTypes {
			bankType.Create()
		}

		// seed banks
		banks := []models.Bank{
			models.Bank{
				Name:     "Bank A",
				Type:     1,
				Address:  "Bank A Address",
				Province: "Province A",
				City:     "City A",
				PIC:      "Bank A PIC",
				Phone:    "081234567890",
				Services: pq.Int64Array{1, 2},
				Products: pq.Int64Array{1, 2},
			},
			models.Bank{
				Name:     "Bank B",
				Type:     2,
				Address:  "Bank B Address",
				Province: "Province B",
				City:     "City B",
				PIC:      "Bank B PIC",
				Phone:    "081234567891",
				Services: pq.Int64Array{1, 2},
				Products: pq.Int64Array{1, 2},
			},
		}
		for _, bank := range banks {
			bank.Create()
		}

		// seed bank services
		services := []models.Service{
			models.Service{
				Name:    "Pinjaman PNS",
				ImageID: 1,
				Status:  "active",
			},
			models.Service{
				Name:    "Pinjaman Pensiun",
				ImageID: 1,
				Status:  "active",
			},
			models.Service{
				Name:    "Pinjaman UMKN",
				ImageID: 1,
				Status:  "active",
			},
			models.Service{
				Name:    "Pinjaman Mikro",
				ImageID: 1,
				Status:  "inactive",
			},
			models.Service{
				Name:    "Pinjaman Lainnya",
				ImageID: 1,
				Status:  "inactive",
			},
		}
		for _, service := range services {
			service.Create()
		}

		// seed service products
		feesMarshal, _ := json.Marshal([]interface{}{
			map[string]interface{}{
				"description": "Admin Fee",
				"amount":      "2500",
			},
			map[string]interface{}{
				"description": "Convenience Fee",
				"amount":      "2%",
			},
		})
		products := []models.Product{
			models.Product{
				ServiceID:       1,
				Name:            "Product A",
				MinTimeSpan:     1,
				MaxTimeSpan:     6,
				Interest:        5,
				MinLoan:         1000000,
				MaxLoan:         10000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				ServiceID:       1,
				Name:            "Product B",
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
		}
		for _, product := range products {
			product.Create()
		}

		// seed borrowers
		borrowers := []models.Borrower{
			models.Borrower{
				Fullname:             "Full Name A",
				Nickname:             "A",
				Gender:               "M",
				IdCardNumber:         "9876123451234567789",
				TaxIDnumber:          "0987654321234567890",
				Nationality:          "WNI",
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
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
			},
			models.Borrower{
				Fullname:             "Full Name B",
				Nickname:             "Be",
				Gender:               "F",
				IdCardNumber:         "9876123451234567781",
				TaxIDnumber:          "0987654321234567891",
				Nationality:          "WNI",
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
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
			},
		}
		for _, borrower := range borrowers {
			borrower.Create()
		}

		purposes := []models.LoanPurpose{
			models.LoanPurpose{
				Name:   "Pendidikan",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Rumah Tangga",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Kesehatan",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Berdagang",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Bertani",
				Status: "active",
			},
			models.LoanPurpose{
				Name:   "Berjudi",
				Status: "inactive",
			},
		}
		for _, purpose := range purposes {
			purpose.Create()
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
				LoanIntention:    "Pendidikan",
				IntentionDetails: "a loan 1 intention details",
				Product:          1,
			},
			models.Loan{
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Status:           "approved",
				LoanAmount:       500000,
				Installment:      2,
				LoanIntention:    "Rumah Tangga",
				IntentionDetails: "a loan 2 intention details",
				Product:          1,
				OTPverified:      true,
			},
			models.Loan{
				Owner: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Status:           "rejected",
				LoanAmount:       2000000,
				Installment:      8,
				LoanIntention:    "Kesehatan",
				IntentionDetails: "a loan 3 intention details",
				Product:          1,
				OTPverified:      true,
			},
		}
		for _, loan := range loans {
			loan.Create()
		}

		//agent migration
		agents := []models.Agent{
			models.Agent{
				Name:     "Agent K",
				Username: "agentK",
				Password: "password",
				Email:    "agentk@mib.com",
				Phone:    "081234567890",
				Category: "agent",
				AgentProvider: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				Banks:  pq.Int64Array{1, 2},
				Status: "active",
			},
			models.Agent{
				Name:     "Agent J",
				Username: "agentJ",
				Password: "password",
				Email:    "agentj@mib.com",
				Phone:    "081234567891",
				Category: "account_executive",
				Banks:    pq.Int64Array{1},
				Status:   "active",
			},
		}
		for _, agent := range agents {
			agent.Create()
		}

		// seed agent's borrowers
		agentBorrowers := []models.AgentBorrower{
			models.AgentBorrower{
				Fullname:             "Full Name AB",
				Gender:               "M",
				IdCardNumber:         "9876123451234566689",
				TaxIDnumber:          "0987654321234566690",
				Birthday:             time.Now(),
				Birthplace:           "a birthplace",
				LastEducation:        "a last edu",
				MotherName:           "a mom",
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
				BankAccountNumber:    "520384666",
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				AgentID: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
			},
			models.AgentBorrower{
				Fullname:             "Full Name BB",
				Gender:               "M",
				IdCardNumber:         "9666123451234566689",
				TaxIDnumber:          "0966654321234566690",
				Birthday:             time.Now(),
				Birthplace:           "a birthplace",
				LastEducation:        "a last edu",
				MotherName:           "a mom",
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
				BankAccountNumber:    "520384666777",
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				AgentID: sql.NullInt64{
					Int64: 2,
					Valid: true,
				},
			},
		}
		for _, agentBorrower := range agentBorrowers {
			agentBorrower.Create()
		}
		//seed uuid
		uuid := models.Uuid_Reset_Password{
			UUID: "f4f71eae-2cc9-4289-94e4-2421df67d4d7",
			Borrower: sql.NullInt64{
				Int64: int64(1),
				Valid: true,
			},
		}
		uuid.Create()

		//seed internal roles
		iroles := []models.InternalRoles{
			models.InternalRoles{
				Name:        "admin",
				Description: "ini admin",
				Status:      true,
				System:      "Core",
			},
		}
		for _, irole := range iroles {
			irole.Create()
		}

		client := []models.Client_config{
			models.Client_config{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Client_config{
				Name:   "android",
				Key:    "androkey",
				Role:   "client",
				Secret: "androsecret",
			},
		}
		for _, clients := range client {
			clients.Create()
		}

	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"bank_types",
				"banks",
				"services",
				"products",
				"images",
				"borrowers",
				"loan_purposes",
				"loans",
				"uuid_reset_passwords",
				"client_configs",
				"internal_roles",
				"agent_borrowers",
				"agents",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = asira.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
