package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Struct that holds the DB connection
type Store struct {
	db DBInterface
}

// Concrete implementations that wrap sql.DB, sql.Rows, and sql.Row
type DBWrapper struct {
	db *sql.DB
}

// NewStore now takes DBInterface instead of *sql.DB
func NewStore(db DBInterface) *Store {
	return &Store{db: db}
}

// NewStoreFromSQLDB creates a Store from sql.DB by wrapping it
func NewStoreFromSQLDB(db *sql.DB) *Store {
	return &Store{db: &DBWrapper{db}}
}

// Closes the DB connection
func CloseDB() error {
	return DB.Close()
}

// Open a DB connection
func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./resume.db")
	if err != nil {
		return fmt.Errorf("database connection failed: %v", err)
	}

	if err := createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	if err := fulfillTables(); err != nil {
		return fmt.Errorf("failed to insert data into the tables: %v", err)
	}

	return nil
}

// Initiates the DB tables
func createTables() error {
	query := `CREATE TABLE IF NOT EXISTS profile (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		firstname TEXT NOT NULL,
        lastname TEXT NOT NULL,
		pronoun TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		location TEXT NOT NULL,
		postal_code INTEGER NOT NULL,
        headline TEXT NOT NULL,
		about TEXT NOT NULL,
        birthdate DATETIME NOT NULL
    );

	CREATE TABLE IF NOT EXISTS education (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
        description TEXT NOT NULL,
        issued_at DATETIME NOT NULL,
		profile_id INTEGER,
		FOREIGN KEY (profile_id) REFERENCES profile(id)
    );

	CREATE TABLE IF NOT EXISTS experience (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
        company TEXT NOT NULL,
		location TEXT NOT NULL,
		description TEXT NOT NULL,
        start_date DATETIME NOT NULL,
		end_date DATETIME,
		profile_id INTEGER,
		FOREIGN KEY (profile_id) REFERENCES profile(id)
    );

	CREATE TABLE IF NOT EXISTS skill (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
    );

	CREATE TABLE IF NOT EXISTS licence (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
        issuer TEXT NOT NULL,
		expires TEXT,
        issued_at DATETIME NOT NULL,
		profile_id INTEGER,
		FOREIGN KEY (profile_id) REFERENCES profile(id)
    );

	CREATE TABLE IF NOT EXISTS skill_experience (
        experience_id INTEGER,
		skill_id INTEGER,
        PRIMARY KEY (experience_id, skill_id),
		FOREIGN KEY (experience_id) REFERENCES experience(id),
    	FOREIGN KEY (skill_id) REFERENCES skill(id)
    );`

	_, err := DB.Exec(query)
	return err
}

// Populates the static DB
// Could be a JSON file
func fulfillTables() error {

	query := `INSERT OR IGNORE INTO profile (
	 	id,
        firstname,
        lastname,
		pronoun,
		email,
		location,
		postal_code,
        headline,
		about,
        birthdate
    ) VALUES (
		1,
		"Florent",
		"Maillard",
		"He/Him",
		"florent@maillard.icu",
		"Switzeralnd - Vaud",
		1867,
		"Integration Expert",
		"There is no subsitute for hard Work./n- Thomas Edison",
		"1990-08-21 00:00:00"
	);

	INSERT OR IGNORE INTO education (
        id,
		title,
        description,
        issued_at,
		profile_id
	) VALUES (
		1,
		"Université de Technologie de Compiègne (UTC)",
		"System and Network Engineer. Member of university sport team",
		"2014-09-01 00:00:00",
		1
	);

	INSERT OR IGNORE INTO experience (
        id,
		title ,
        company,
		location,
		description,
        start_date,
		end_date,
		profile_id
    ) VALUES (
		1,
		"Information Technology Engineer",
		"CoDEM Picardie",
		"Amiens",
		"Custom an ERP/CRM opensource software.\nR&D project management. Web app development and integration of 3 dimensions buildings data and point clouds",
		"2014-02-01 00:00:00",
		"2016-07-01 00:00:00",
		1
	), (
		2,
		"IT Project Manager / Lead developer",
		"French public services",
		"Beauvais",
		"Internal projects: Develop apps in a DevOps Team.\nExternal projects: Project specifications, budgets, plannings, leading IT service providers teams",
		"2016-08-01 00:00:00",
		"2021-09-01 00:00:00",
		1
	), (
		3,
		"Integration technical architect",
		"Pocalin Hydraulics",
		"Verberie",
		"Design and describe APIs using RAML, OpenAPI, AsyncAPI or GraphQL.\nDevelop Mule4 applications and achieve reliability implementing known design patterns and using platforms such as RabbitMQ / AnypointMQ or Apache Kafka.\nImplement CI/CD pipelines in Gitlab. Using Maven, mule cli and ansible.\nGovern APIs using Open Policy Agent or AMF. Enforce policies to ensure reliability, resilience\nand security (OAuth2, quotas, IP filtering, mTLS)\n\nOn-premise mule runtimes to Anypoint cloudhub 2.0 migration\n\nMulesoft connect 2022 speaker",
		"2021-09-01 00:00:00",
		"2023-02-01 00:00:00",
		1
	), (
		4,
		"Integration engineer",
		"Vaudoise Assurances",
		"Lausanne",
		"Build and run the existing wso2 clusters (API Manager and Identity Server)\n\nImprove or develop new custom wso2 components. From JWT token issuers and mediators to webapps\n\nFine tune the WSO2 engine\n\nStart the transition to the cloud for the API management. 6 environments created and designed to be fully managed using infrastructure as code",
		"2023-02-01 00:00:00",
		"2024-04-01 00:00:00",
		1
	);

	INSERT OR IGNORE INTO experience (
        id,
		title ,
        company,
		location,
		description,
        start_date,
		profile_id
    ) VALUES (
		5,
		"Interation expert",
		"Vaudoise Assurances",
		"Lausanne",
		"In addition to the preceding role\n\nLevel 3 support on all the integration platforms. Strong rise in skills on apache Kafka.\n\nManage the WSO2 platform migration project. Supporting teams during migration to the Cloud services.\n\nManage the lift and shift project to Confluent Cloud and Azure Kubernetes Service.\n\nImprove platforms logs and metrics in ELK. From custom Dashboards to watcher alerts.\n\nAutomate some admin common tasks with Azure DevOps pipelines.\n\nIntegrate all team's legacy projects into Jenkins and SonarQube\n\nPart of the the Vaudoise Azure community of practice. Defining standards and helping teams to achieve them.",
		"2024-04-01 00:00:00",
		1
	), (
		6,
		"Hhikig trail marker ",
		"Vaud Rando",
		"Lavey/Morcles",
		"Mark mountain hiking trails",
		"2025-03-01 00:00:00",
		1
	);

	INSERT OR IGNORE INTO skill (
        id,
		name
    ) VALUES (
	 	1,
		"PostgreSQL"
	),	(
		2,
		"Git"
	),	(
		3,
		"MySQL"
	),	(
		4,
		"Podman"
	),	(
		5,
		"Istio"
	),	(
		6,
		"MongoDB"
	),	(
		7,
		"Apache Kafka"
	),	(
		8,
		"Maven"
	),	(
		9,
		"OAuth2"
	),	(
		10,
		"SAML"
	),	(
		11,
		"Terraform"
	),	(
		12,
		"GraphQL"
	),	(
		13,
		"Ansible"
	),	(
		14,
		"RabbitMQ"
	),	(
		15,
		"OIDC"
	),	(
		16,
		"Mulesoft"
	),	(
		17,
		"API Management"
	),	(
		18,
		"Kerberos"
	),	(
		19,
		"Azure"
	),	(
		20,
		"AWS"
	),	(
		21,
		"GO"
	),	(
		22,
		"Jenkins"
	),	(
		23,
		"Azure DevOps"
	),	(
		24,
		"AKS"
	),	(
		25,
		"SonarQube"
	),	(
		26,
		"KSql"
	);

	INSERT OR IGNORE INTO skill_experience (
        experience_id,
		skill_id
	) VALUES (
	 	1,
		1
 	), (
		1,
		2
	), (
		1,
		3
	), (
		1,
		2
	), (
		2,
		1
	), (
		2,
		2
	), (
		2,
		3
	), (
		2,
		4
	), (
		2,
		5
	), (
		2,
		6
	), (
		3,
		7
	), (
		3,
		8
	), (
		3,
		1
	), (
		3,
		9
	), (
		3,
		10
	), (
		3,
		11
	), (
		3,
		12
	), (
		3,
		13
	), (
		3,
		14
	), (
		3,
		15
	), (
		3,
		16
	), (
		3,
		17
	), (
		4,
		2
	), (
		4,
		7
	), (
		4,
		8
	), (
		4,
		9
	), (
		4,
		10
	), (
		4,
		11
	), (
		4,
		13
	), (
		4,
		14
	), (
		4,
		15
	), (
		4,
		17
	), (
		4,
		18
	), (
		4,
		19
	), (
		5,
		26
	), (
		5,
		25
	), (
		5,
		24
	), (
		5,
		23
	), (
		5,
		22
	), (
		5,
		21
	), (
		5,
		20
	), (
		5,
		19
	);

	INSERT OR IGNORE INTO licence (
        id,
		title,
        issuer,
		expires,
        issued_at,
		profile_id
    ) VALUES (
		1,
		"Mulesoft Certified Platform Architect (MCPA)",
		"Mulesoft",
		"2022-01-01 00:00:00",
		"2024-04-01 00:00:00",
		1
	), (
		2,
		"Mulesoft Certified Integration Architect (MCIA)",
		"Mulesoft",
		"2022-01-01 00:00:00",
		"2024-04-01 00:00:00",
		1
	),(
		3,
		"Confluent Certified Developer for Apache Kafka",
		"Confluent",
		"2023-01-01 00:00:00",
		"2025-01-01 00:00:00",
		1
	),(
		4,
		"Microsoft Certified: Cybersecurity Architect Expert",
		"Mulesoft",
		"2024-02-01 00:00:00",
		"2025-02-01 00:00:00",
		1
	),(
		5,
		"Microsoft Certified : Azure Security Engineer Associate",
		"Mulesoft",
		"2024-02-01 00:00:00",
		"2025-02-01 00:00:00",
		1
	),(
		6,
		"CKAD",
		"The Linux Foundation",
		"2022-10-01 00:00:00",
		"2025-10-01 00:00:00",
		1
	);

	INSERT OR IGNORE INTO licence (
        id,
		title,
        issuer,
        issued_at,
		profile_id
    ) VALUES (
		8,
		"Scrum Basics",
		"Scrum INC",
		"2025-09-01 00:00:00",
		1
	),(
		7,
		"WSO2 Certified API Manager",
		"WSO2",
		"2022-12-01 00:00:00",
		1
	);`

	_, err := DB.Exec(query)
	return err
}
