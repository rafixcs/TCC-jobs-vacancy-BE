CREATE TABLE "companies"(
    "id" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "creation_date" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "description" TEXT NOT NULL,
    "location" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "companies" ADD PRIMARY KEY("id");
CREATE TABLE "user_applies"(
    "id" VARCHAR(255) NOT NULL,
    "job_vacancy_id" VARCHAR(255) NOT NULL,
    "user_id" VARCHAR(255) NOT NULL,
    "cover_letter" TEXT NULL,
    "email" VARCHAR(255) NOT NULL,
    "full_name" TEXT NOT NULL,
    "phone" VARCHAR(255) NULL,
    "url_resume" TEXT NOT NULL
);
ALTER TABLE
    "user_applies" ADD PRIMARY KEY("id");
CREATE TABLE "user_logins"(
    "id" VARCHAR(255) NOT NULL,
    "user_id" VARCHAR(255) NOT NULL,
    "login_date" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "logout_date" TIMESTAMP(0) WITHOUT TIME ZONE NULL
);
ALTER TABLE
    "user_logins" ADD PRIMARY KEY("id");
CREATE TABLE "job_vacancies"(
    "id" VARCHAR(255) NOT NULL,
    "user_id" VARCHAR(255) NOT NULL,
    "company_id" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "creation_date" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "location" VARCHAR(255) NOT NULL,
    "salary" VARCHAR(255) NOT NULL,
    "requirements" TEXT NOT NULL,
    "responsibilities" TEXT NOT NULL,
    "job_type" VARCHAR(255) NOT NULL,
    "experience_level" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "job_vacancies" ADD PRIMARY KEY("id");
CREATE TABLE "users"(
    "id" VARCHAR(255) NOT NULL,
    "role_id" INTEGER NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "phone" VARCHAR(255) NULL
);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "user_roles"(
    "id" INTEGER NOT NULL,
    "type" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "user_roles" ADD PRIMARY KEY("id");
CREATE TABLE "company_users"(
    "company_id" VARCHAR(255) NOT NULL,
    "user_id" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "users" ADD CONSTRAINT "users_role_id_foreign" FOREIGN KEY("role_id") REFERENCES "user_roles"("id");
ALTER TABLE
    "job_vacancies" ADD CONSTRAINT "job_vacancies_company_id_foreign" FOREIGN KEY("company_id") REFERENCES "companies"("id");
ALTER TABLE
    "job_vacancies" ADD CONSTRAINT "job_vacancies_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "user_applies" ADD CONSTRAINT "user_applies_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "user_logins" ADD CONSTRAINT "user_logins_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "company_users" ADD CONSTRAINT "company_users_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "company_users" ADD CONSTRAINT "company_users_company_id_foreign" FOREIGN KEY("company_id") REFERENCES "companies"("id");
ALTER TABLE
    "user_applies" ADD CONSTRAINT "user_applies_job_vacancy_id_foreign" FOREIGN KEY("job_vacancy_id") REFERENCES "job_vacancies"("id");