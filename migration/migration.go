package migration

import (
	"asira_borrower/asira"
	"asira_borrower/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

var (
	AgentPassword string = "$2y$12$lpU2qJ5S.q0tcK.bJaUYAedNM1U63bpMRIr0KT4YIaOrwNqXqo9tq"
)

func Seed() {
	seeder := asira.App.DB.Begin()
	defer seeder.Commit()

	if asira.App.ENV == "development" {
		// seed internals
		client := []models.Client{
			models.Client{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Client{
				Name:   "android",
				Key:    "androkey",
				Role:   "client",
				Secret: "androsecret",
			},
		}
		for _, clients := range client {
			clients.Create()
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

	}
}

func TestSeed() {
	seeder := asira.App.DB.Begin()
	defer seeder.Commit()

	if asira.App.ENV == "development" {
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
				Name:   "Pinjaman PNS",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Status: "active",
			},
			models.Service{
				Name:   "Pinjaman Pensiun",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Status: "active",
			},
			models.Service{
				Name:   "Pinjaman UMKN",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Status: "active",
			},
			models.Service{
				Name:   "Pinjaman Mikro",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Status: "inactive",
			},
			models.Service{
				Name:   "Pinjaman Lainnya",
				Image:  "https://images.unsplash.com/photo-1576039716094-066beef36943?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80",
				Status: "inactive",
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
				"fee_method":  "deduct_loan",
			},
			map[string]interface{}{
				"description": "Convenience Fee",
				"amount":      "2%",
				"fee_method":  "charge_loan",
			},
		})
		products := []models.Product{
			models.Product{
				Name:            "Product A",
				ServiceID:       1,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				InterestType:    "flat",
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product B",
				ServiceID:       2,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        8,
				InterestType:    "fixed",
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product C",
				ServiceID:       1,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        10,
				InterestType:    "onetimepay",
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product D",
				ServiceID:       2,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        12,
				InterestType:    "efektif_menurun",
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product E",
				ServiceID:       5,
				MinTimeSpan:     3,
				MaxTimeSpan:     12,
				Interest:        5,
				InterestType:    "flat",
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product Jeruk Kecil",
				ServiceID:       6,
				MinTimeSpan:     6,
				MaxTimeSpan:     12,
				Interest:        5,
				InterestType:    "efektif_menurun",
				MinLoan:         5000000,
				MaxLoan:         8000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product Jeruk Besar",
				ServiceID:       6,
				MinTimeSpan:     6,
				MaxTimeSpan:     24,
				Interest:        5,
				InterestType:    "efektif_menurun",
				MinLoan:         5000000,
				MaxLoan:         20000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product Pisang Kecil",
				ServiceID:       7,
				MinTimeSpan:     6,
				MaxTimeSpan:     24,
				Interest:        5,
				InterestType:    "fixed",
				MinLoan:         5000000,
				MaxLoan:         10000000,
				Fees:            postgres.Jsonb{feesMarshal},
				Collaterals:     []string{"Surat Tanah", "BPKB"},
				FinancingSector: []string{"Pendidikan"},
				Assurance:       "an Assurance",
				Status:          "active",
			},
			models.Product{
				Name:            "Product Pisang Raja",
				ServiceID:       7,
				MinTimeSpan:     6,
				MaxTimeSpan:     36,
				Interest:        7,
				InterestType:    "flat",
				MinLoan:         5000000,
				MaxLoan:         30000000,
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
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				AgentReferral: sql.NullInt64{
					Int64: 0,
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
				Bank: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
				AgentReferral: sql.NullInt64{
					Int64: 0,
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
				Borrower:         1,
				LoanAmount:       1000000,
				Installment:      6,
				LoanIntention:    "Pendidikan",
				IntentionDetails: "a loan 1 intention details",
				Product:          1,
				PaymentStatus:    "terbayar",
			},
			models.Loan{
				Borrower:         1,
				Status:           "approved",
				LoanAmount:       500000,
				Installment:      2,
				LoanIntention:    "Rumah Tangga",
				IntentionDetails: "a loan 2 intention details",
				Product:          1,
				OTPverified:      true,
				PaymentStatus:    "terbayar",
			},
			models.Loan{
				Borrower:         1,
				Status:           "rejected",
				LoanAmount:       2000000,
				Installment:      8,
				LoanIntention:    "Kesehatan",
				IntentionDetails: "a loan 3 intention details",
				Product:          1,
				OTPverified:      true,
				PaymentStatus:    "terbayar",
			},
		}
		for _, loan := range loans {
			loan.Create()
		}

		//agent provider migration
		agentProviders := []models.AgentProvider{
			models.AgentProvider{
				Name:    "Agent Provider A",
				PIC:     "PIC A",
				Phone:   "081234567890",
				Address: "address of provider a",
				Status:  "active",
			},
			models.AgentProvider{
				Name:    "Agent Provider B",
				PIC:     "PIC B",
				Phone:   "081234567891",
				Address: "address of provider b",
				Status:  "active",
			},
			models.AgentProvider{
				Name:    "Agent Provider C",
				PIC:     "PIC C",
				Phone:   "081234567892",
				Address: "address of provider c",
				Status:  "active",
			},
		}
		for _, agentProvider := range agentProviders {
			agentProvider.Create()
		}

		//agent migration
		agents := []models.Agent{
			models.Agent{
				Name:     "Agent K",
				Username: "agentK",
				Password: AgentPassword,
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
				Password: AgentPassword,
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
		agentBorrowers := []models.Borrower{
			models.Borrower{
				Fullname:             "Full Name AA",
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
				AgentReferral: sql.NullInt64{
					Int64: 1,
					Valid: true,
				},
			},
			models.Borrower{
				Fullname:             "Full Name AB",
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
				AgentReferral: sql.NullInt64{
					Int64: 2,
					Valid: true,
				},
			},
		}
		for _, agentBorrower := range agentBorrowers {
			agentBorrower.Create()
		}

		// seed loans for agent's borrowers (Borrower : 4 owned by Agent : 2 )
		loans = []models.Loan{
			models.Loan{
				Borrower:         4,
				LoanAmount:       1000000,
				Installment:      6,
				LoanIntention:    "Pendidikan",
				IntentionDetails: "a loan 1 intention details",
				Product:          1,
				PaymentStatus:    "terbayar",
			},
			models.Loan{
				Borrower:         4,
				Status:           "approved",
				LoanAmount:       500000,
				Installment:      2,
				LoanIntention:    "Rumah Tangga",
				IntentionDetails: "a loan 2 intention details",
				Product:          1,
				OTPverified:      true,
				PaymentStatus:    "terbayar",
			},
			models.Loan{
				Borrower:         4,
				Status:           "rejected",
				LoanAmount:       2000000,
				Installment:      8,
				LoanIntention:    "Kesehatan",
				IntentionDetails: "a loan 3 intention details",
				Product:          1,
				OTPverified:      true,
				PaymentStatus:    "terbayar",
			},
		}
		for _, loan := range loans {
			loan.Create()
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

		client := []models.Client{
			models.Client{
				Name:   "admin",
				Key:    "adminkey",
				Role:   "admin",
				Secret: "adminsecret",
			},
			models.Client{
				Name:   "android",
				Key:    "androkey",
				Role:   "client",
				Secret: "androsecret",
			},
		}
		for _, clients := range client {
			clients.Create()
		}

		notifications := []models.Notification{
			models.Notification{
				ClientID:      2,
				RecipientID:   "borrower-1",
				Title:         "Status Peminjaman Anda",
				MessageBody:   "Loan id 1 diterima oleh Bank A 2019-11-06",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      2,
				RecipientID:   "borrower-2",
				Title:         "Status Peminjaman Anda",
				MessageBody:   "Loan id 2 diterima oleh Bank B 2019-11-06",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      2,
				RecipientID:   "agent-1",
				Title:         "Promosi",
				MessageBody:   "Promosi Produk A...",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
			models.Notification{
				ClientID:      2,
				RecipientID:   "agent-2",
				Title:         "Info",
				MessageBody:   "Info, informasi...",
				FirebaseToken: "dtRiFA0U7tE:APA91bGURmGuJLLD0Y_PdhHvFj3haiq8jj8haAfaw94yH0j-Jj6TiGM4wzRHmoLUhihlXKLUF7iFAsAalXzrUVLEVbbtK57_wVYCu2jhameVoG13rc1HgIM6xanEMlApnuwSmURdcdz1",
				Topic:         "promosi",
				SendTime:      time.Now(),
				Response:      "projects/asira-app-33ed7/messages/0:1573131084523969%2056f65bf9fd7ecd",
			},
		}
		for _, notification := range notifications {
			notification.Create()
		}

		users := []models.User{
			models.User{
				Borrower: 1,
				Password: "password",
			},
			models.User{
				Borrower: 2,
				Password: "password",
			},
		}
		for _, user := range users {
			user.Create()
		}

		faqs := []models.FAQ{
			models.FAQ{
				Title: "How to register",
				Description: `
				<html>
				<head>
				</head>
				<body>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 1</strong>
				<p>Lorem ipsum...1</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 2</strong>
				<p>Lorem ipsum...2</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 3</strong>
				<p>Lorem ipsum...3</p>
				</div>
				</body>
				</html>`,
			},
			models.FAQ{
				Title: "How to applying loan",
				Description: `
				<html>
				<head>
				</head>
				<body>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 1</strong>
				<p>Lorem ipsum...1</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 2</strong>
				<p>Lorem ipsum...2</p>
				</div>
				<div class="panel" style="background-color:white; max-height:0; padding:20px 20px; transition:max-height 0.2s ease-out;margin-bottom:50px" >
				<strong class="accordion">Section 3</strong>
				<p>Lorem ipsum...3</p>
				</div>
				</body>
				</html>`,
			},
		}
		for _, faq := range faqs {
			faq.Create()
		}

	}
}

// truncate defined tables. []string{"all"} to truncate all tables.
func Truncate(tableList []string) (err error) {
	if len(tableList) > 0 {
		if tableList[0] == "all" {
			tableList = []string{
				"banks",
				"bank_types",
				"services",
				"products",
				"borrowers",
				"loan_purposes",
				"loans",
				"installments",
				"uuid_reset_passwords",
				"clients",
				"agents",
				"agent_providers",
				"notifications",
				"users",
				"faqs",
			}
		}

		tables := strings.Join(tableList, ", ")
		sqlQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)
		err = asira.App.DB.Exec(sqlQuery).Error
		return err
	}

	return fmt.Errorf("define tables that you want to truncate")
}
