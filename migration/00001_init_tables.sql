-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE "banks" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "type" varchar(255),
    "address" text,
    "province" varchar(255),
    "city" varchar(255),
    "services" jsonb DEFAULT '[]',
    "products" jsonb DEFAULT '[]',
    "pic" varchar(255),
    "phone" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "bank_services" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "image_id" bigint,
    "status" varchar(255),
    FOREIGN KEY ("image_id") REFERENCES images(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);
COMMENT ON COLUMN "bank_services"."status" IS '0 = inactive, 1 = active';

CREATE TABLE "service_products" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "min_timespan" int,
    "max_timespan" int,
    "interest" int,
    "min_loan" int,
    "max_loan" int,
    "fees" jsonb DEFAULT '[]',
    "asn_fee" varchar(255),
    "service" bigint,
    "collaterals" jsonb DEFAULT '[]',
    "financing_sector" jsonb DEFAULT '[]',
    "assurance" varchar(255),
    "status" varchar(255),
    FOREIGN KEY ("service") REFERENCES bank_services(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "images" (
    "id" bigserial,
    "image_string" text,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "borrowers" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "suspended_time" timestamptz,
    "fullname" varchar(255) NOT NULL,
    "gender" varchar(1) NOT NULL,
    "idcard_number" varchar(255) NOT NULL UNIQUE,
    "idcard_image" bigint,
    "taxid_number" varchar(255),
    "taxid_image" bigint,
    "email" varchar(255) NOT NULL UNIQUE,
    "birthday" DATE NOT NULL,
    "birthplace" varchar(255) NOT NULL,
    "last_education" varchar(255) NOT NULL,
    "mother_name" varchar(255) NOT NULL,
    "phone" varchar(255) NOT NULL UNIQUE,
    "marriage_status" varchar(255) NOT NULL,
    "spouse_name" varchar(255),
    "spouse_birthday" DATE,
    "spouse_lasteducation" varchar(255),
    "dependants" int DEFAULT (0),
    "address" text NOT NULL,
    "province" varchar(255) NOT NULL,
    "city" varchar(255) NOT NULL,
    "neighbour_association" varchar(255) NOT NULL,
    "hamlets" varchar(255) NOT NULL,
    "home_phonenumber" varchar(255) NOT NULL,
    "subdistrict" varchar(255) NOT NULL,
    "urban_village" varchar(255) NOT NULL,
    "home_ownership" varchar(255) NOT NULL,
    "lived_for" int NOT NULL,
    "occupation" varchar(255) NOT NULL,
    "employee_id" varchar(255),
    "employer_name" varchar(255) NOT NULL,
    "employer_address" text NOT NULL,
    "department" varchar(255) NOT NULL,
    "been_workingfor" int NOT NULL,
    "direct_superiorname" varchar(255),
    "employer_number" varchar(255) NOT NULL,
    "monthly_income" int NOT NULL,
    "other_income" int,
    "other_incomesource" varchar(255),
    "field_of_work" varchar(255) NOT NULL,
    "related_personname" varchar(255) NOT NULL,
    "related_relation" varchar(255) NOT NULL,
    "related_phonenumber" varchar(255) NOT NULL,
    "related_homenumber" varchar(255),
    "related_address" text,
    "bank" bigint,
    "bank_accountnumber" varchar(255),
    "otp_verified" BOOLEAN,
    "password" varchar(255) NOT NULL,
    FOREIGN KEY ("idcard_image") REFERENCES images(id),
    FOREIGN KEY ("taxid_image") REFERENCES images(id),
    FOREIGN KEY ("bank") REFERENCES banks(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "loans" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "owner" bigint,
    "status" varchar(255) DEFAULT  ('processing'),
    "loan_amount" FLOAT NOT NULL,
    "installment" int NOT NULL,
    "fees" jsonb DEFAULT '[]',
    "interest" FLOAT NOT NULL,
    "total_loan" FLOAT NOT NULL,
    "due_date" timestamptz,
    "layaway_plan" FLOAT NOT NULL,
    "loan_intention" varchar(255) NOT NULL,
    "intention_details" text NOT NULL,
    "borrower_info" jsonb DEFAULT '[]',
    "otp_verified" BOOLEAN,
    FOREIGN KEY ("owner") REFERENCES borrowers(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);


CREATE TABLE "uuid_reset_passwords" (
    "id" bigserial,
    "uuid" varchar(255) NOT NULL,
    "borrower" bigint,
    "expired" timestamptz,
    "used"  BOOLEAN,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("borrower") REFERENCES borrowers(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS "banks" CASCADE;
DROP TABLE IF EXISTS "bank_services" CASCADE;
DROP TABLE IF EXISTS "service_products" CASCADE;
DROP TABLE IF EXISTS "images" CASCADE;
DROP TABLE IF EXISTS "borrowers" CASCADE;
DROP TABLE IF EXISTS "loans" CASCADE;
DROP TABLE IF EXISTS "uuid_reset_passwords" CASCADE;