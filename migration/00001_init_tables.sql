-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE "bank_types" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255),
    "description" text,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "banks" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255),
    "image" text,
    "type" varchar(255),
    "address" text,
    "province" varchar(255),
    "city" varchar(255),
    "pic" varchar(255),
    "phone" varchar(255),
    "services" int ARRAY,
    "products" int ARRAY,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "services" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255),
    "image" text,
    "status" varchar(255),
    "description" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "agent_providers" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255),
    "pic" varchar(255),
    "phone" varchar(255) UNIQUE,
    "address" text,
    "status" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "agents" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255),
    "username" varchar(255) UNIQUE,
    "password" text,
    "image" text,
    "email" varchar(255) UNIQUE,
    "phone" varchar(255) UNIQUE,
    "category" varchar(255),
    "agent_provider" bigint,
    "banks" int ARRAY,
    "status" varchar(255),
    "fcm_token" varchar(255),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("agent_provider") REFERENCES agent_providers(id)
) WITH (OIDS = FALSE);

CREATE TABLE "products" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255),
    "status" varchar(255),
    "service_id" bigserial,
    "min_timespan" int,
    "max_timespan" int,
    "interest" int,
    "interest_type" varchar(255),
    "min_loan" int,
    "max_loan" int,
    "fees" jsonb DEFAULT '[]',
    "collaterals" varchar(255) ARRAY,
    "financing_sector" varchar(255) ARRAY,
    "assurance" varchar(255),
    "form" jsonb DEFAULT '[]',
    FOREIGN KEY ("service_id") REFERENCES services(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "borrowers" (
    "id" bigserial,    
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "suspended_time" timestamptz,
    "fullname" varchar(255) NOT NULL,
    "nickname" varchar(255),
    "gender" varchar(1),
    "image" text,
    "idcard_number" varchar(255),
    "idcard_image" text,
    "taxid_number" varchar(255),
    "taxid_image" text,
    "nationality" varchar(255),
    "email" varchar(255),
    "birthday" DATE,
    "birthplace" varchar(255),
    "last_education" varchar(255),
    "mother_name" varchar(255) ,
    "phone" varchar(255) NOT NULL,
    "marriage_status" varchar(255),
    "spouse_name" varchar(255),
    "spouse_birthday" DATE,
    "spouse_lasteducation" varchar(255),
    "dependants" int DEFAULT (0),
    "address" text,
    "province" varchar(255),
    "city" varchar(255),
    "neighbour_association" varchar(255),
    "hamlets" varchar(255),
    "home_phonenumber" varchar(255),
    "subdistrict" varchar(255),
    "urban_village" varchar(255),
    "home_ownership" varchar(255),
    "lived_for" int,
    "occupation" varchar(255),
    "employee_id" varchar(255),
    "employer_name" varchar(255),
    "employer_address" text,
    "department" varchar(255),
    "been_workingfor" int,
    "direct_superiorname" varchar(255),
    "employer_number" varchar(255),
    "monthly_income" int,
    "other_income" int,
    "other_incomesource" varchar(255),
    "field_of_work" varchar(255) ,
    "related_personname" varchar(255),
    "related_relation" varchar(255) ,
    "related_phonenumber" varchar(255) ,
    "related_homenumber" varchar(255),
    "related_address" text,
    "bank" bigserial,
    "bank_accountnumber" varchar(255),
    "otp_verified" BOOLEAN,
    "agent_referral"  bigint,
    "status" varchar(255) DEFAULT 'active',
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "loan_purposes" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "name" varchar(255),
    "status" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "loans" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "borrower" bigserial,
    "product" bigserial,
    "status" varchar(255) DEFAULT  ('processing'),
    "loan_amount" FLOAT NOT NULL,
    "installment" int NOT NULL,
    "installment_details" int ARRAY,
    "fees" jsonb DEFAULT '[]',
    "interest" FLOAT NOT NULL,
    "total_loan" FLOAT NOT NULL,
    "disburse_amount" FLOAT NOT NULL,
    "due_date" timestamptz,
    "layaway_plan" FLOAT NOT NULL,
    "loan_intention" varchar(255) NOT NULL,
    "intention_details" text NOT NULL,
    "borrower_info" jsonb DEFAULT '[]',
    "otp_verified" BOOLEAN,
    "disburse_date" timestamptz,
    "disburse_date_changed" BOOLEAN,
    "disburse_status" varchar(255) DEFAULT ('processing'),
    "approval_date" timestamptz,
    "reject_reason" text,
    "form_info" jsonb DEFAULT '[]',
    FOREIGN KEY ("borrower") REFERENCES borrowers(id),
    FOREIGN KEY ("product") REFERENCES products(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "installments" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "period" int,
    "loan_payment" FLOAT,
    "interest_payment" FLOAT,
    "paid_date" timestamptz,
    "paid_status" BOOLEAN,
    "paid_amount" FLOAT,
    "underpayment" FLOAT,
    "penalty" FLOAT,
    "due_date" timestamptz,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "uuid_reset_passwords" (
    "id" bigserial,
    "uuid" varchar(255) NOT NULL,
    "borrower" bigserial,
    "expired" timestamptz,
    "used"  BOOLEAN,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    FOREIGN KEY ("borrower") REFERENCES borrowers(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "clients" (
    "id" bigserial,
    "name" varchar(255) NOT NULL,
    "role" varchar(255) NOT NULL,
    "secret" varchar(255) NOT NULL,
    "key" varchar(255) NOT NULL,    
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "notifications" (
    "id" bigserial,
    "client_id" bigserial,
    "recipient_id" varchar(255),
    "title" varchar(255) NOT NULL,
    "message_body" text,
    "firebase_token" varchar(255),
    "topic" varchar(125),
    "response" varchar(255),
    "send_time" timestamptz DEFAULT CURRENT_TIMESTAMP,    
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "users" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "borrower" bigserial NOT NULL,
    "password" varchar(255) NOT NULL,
    "fcm_token" varchar(255),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("borrower") REFERENCES borrowers(id)
) WITH (OIDS = FALSE);

CREATE TABLE "faqs" (
    "id" bigserial,
    "created_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "title" varchar(255),
    "description" text,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS "products" CASCADE;
DROP TABLE IF EXISTS "services" CASCADE;
DROP TABLE IF EXISTS "banks" CASCADE;
DROP TABLE IF EXISTS "bank_types" CASCADE;
DROP TABLE IF EXISTS "images" CASCADE;
DROP TABLE IF EXISTS "borrowers" CASCADE;
DROP TABLE IF EXISTS "loan_purposes" CASCADE;
DROP TABLE IF EXISTS "loans" CASCADE;
DROP TABLE IF EXISTS "installments" CASCADE;
DROP TABLE IF EXISTS "uuid_reset_passwords" CASCADE;
DROP TABLE IF EXISTS "clients" CASCADE;
DROP TABLE IF EXISTS "agents" CASCADE;
DROP TABLE IF EXISTS "agent_providers" CASCADE;
DROP TABLE IF EXISTS "notifications" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "faqs" CASCADE;
