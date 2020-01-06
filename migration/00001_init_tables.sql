-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE "images" (
    "id" bigserial,
    "image_string" text,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "bank_types" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "description" text,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

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
    "pic" varchar(255),
    "phone" varchar(255),
    "adminfee_setup" varchar(255),
    "convfee_setup" varchar(255),
    "services" int ARRAY,
    "products" int ARRAY,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "services" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "image" text,
    "status" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "agents" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "username" varchar(255),
    "password" text,
    "image_id" bigint,
    "email" varchar(255),
    "phone" varchar(255),
    "category" varchar(255),
    "agent_provider" bigint,
    "banks" int ARRAY,
    "status" varchar(255),
    "fcm_token" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "products" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "status" varchar(255),
    "service_id" bigserial,
    "min_timespan" int,
    "max_timespan" int,
    "interest" int,
    "min_loan" int,
    "max_loan" int,
    "fees" jsonb DEFAULT '[]',
    "collaterals" varchar(255) ARRAY,
    "financing_sector" varchar(255) ARRAY,
    "assurance" varchar(255),
    FOREIGN KEY ("service_id") REFERENCES services(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "borrowers" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "suspended_time" timestamptz,
    "fullname" varchar(255) NOT NULL,
    "nickname" varchar(255),
    "gender" varchar(1) NOT NULL,
    "idcard_number" varchar(255) NOT NULL,
    "idcard_image" bigserial,
    "taxid_number" varchar(255),
    "taxid_image" bigserial,
    "nationality" varchar(255),
    "email" varchar(255)  ,
    "birthday" DATE NOT NULL,
    "birthplace" varchar(255) NOT NULL,
    "last_education" varchar(255) NOT NULL,
    "mother_name" varchar(255) NOT NULL,
    "phone" varchar(255) ,
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
    "been_workingfor" int,
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
    "bank" bigserial,
    "bank_accountnumber" varchar(255),
    "otp_verified" BOOLEAN,
    "agent_referral"  bigserial,
    "status" varchar(255) DEFAULT 'active',
    FOREIGN KEY ("idcard_image") REFERENCES images(id),
    FOREIGN KEY ("taxid_image") REFERENCES images(id),
    FOREIGN KEY ("bank") REFERENCES banks(id), 
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);



CREATE TABLE "loan_purposes" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "name" varchar(255),
    "status" varchar(255),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "loans" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "borrower" bigserial,
    "product" bigserial,
    "status" varchar(255) DEFAULT  ('processing'),
    "loan_amount" FLOAT NOT NULL,
    "installment" int NOT NULL,
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
    "reject_reason" text,
    FOREIGN KEY ("borrower") REFERENCES borrowers(id),
    FOREIGN KEY ("product") REFERENCES products(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "uuid_reset_passwords" (
    "id" bigserial,
    "uuid" varchar(255) NOT NULL,
    "borrower" bigserial,
    "expired" timestamptz,
    "used"  BOOLEAN,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("borrower") REFERENCES borrowers(id),
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "client_configs" (
    "id" bigserial,
    "name" varchar(255) NOT NULL,
    "role" varchar(255) NOT NULL,
    "secret" varchar(255) NOT NULL,
    "key" varchar(255) NOT NULL,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "internal_roles" (
    "id" bigserial,
    "name" varchar(255) NOT NULL,
    "system" varchar(255) NOT NULL,
    "description" text,
    "status" BOOLEAN,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
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
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
) WITH (OIDS = FALSE);

CREATE TABLE "users" (
    "id" bigserial,
    "created_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "updated_time" timestamptz DEFAULT CURRENT_TIMESTAMP,
    "deleted_time" timestamptz,
    "borrower" bigserial NOT NULL,
    "password" varchar(255) NOT NULL,
    "fcm_token" varchar(255),
    PRIMARY KEY ("id"),
    FOREIGN KEY ("borrower") REFERENCES borrowers(id)
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
DROP TABLE IF EXISTS "uuid_reset_passwords" CASCADE;
DROP TABLE IF EXISTS "client_configs" CASCADE;
DROP TABLE IF EXISTS "internal_roles" CASCADE;
DROP TABLE IF EXISTS "agents" CASCADE;
DROP TABLE IF EXISTS "notifications" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
