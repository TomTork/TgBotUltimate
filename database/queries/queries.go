package queries

import (
	"TgBotUltimate/database/queries/helper"
	"fmt"
	"strings"
)

const CreateUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	tg_id BIGINT PRIMARY KEY,
	username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
	phone_number VARCHAR(12),
	email VARCHAR(255),
    ex_project_name VARCHAR(255) DEFAULT '',
    ex_building_liter VARCHAR(7) DEFAULT '',
    ex_floor_min VARCHAR(7) DEFAULT '',
    ex_floor_max VARCHAR(7) DEFAULT '',
    ex_rooms_amount_min VARCHAR(7) DEFAULT '',
    ex_rooms_amount_max VARCHAR(7) DEFAULT '',
    ex_square_min VARCHAR(7) DEFAULT '',
    ex_square_max VARCHAR(7) DEFAULT '',
    ex_cost_min VARCHAR(7) DEFAULT '',
    ex_cost_max VARCHAR(7) DEFAULT '',
    uoffset INTEGER DEFAULT 0
);
`

const CreateMessagesTable = `
CREATE TABLE IF NOT EXISTS messages (
	id SERIAL PRIMARY KEY,
	tg_id BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT now(),
	message TEXT NOT NULL,
    project_name VARCHAR(255) DEFAULT '',
    building_liter VARCHAR(7) DEFAULT '',
    floor_min VARCHAR(7) DEFAULT '',
    floor_max VARCHAR(7) DEFAULT '',
    rooms_amount_min VARCHAR(7) DEFAULT '',
    rooms_amount_max VARCHAR(7) DEFAULT '',
    square_min VARCHAR(7) DEFAULT '',
    square_max VARCHAR(7) DEFAULT '',
    cost_min VARCHAR(7) DEFAULT '',
    cost_max VARCHAR(7) DEFAULT '',
	CONSTRAINT fk_user
		FOREIGN KEY (tg_id)
		REFERENCES users(tg_id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
);
`

const CreateProjectsTable = `
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    code VARCHAR(63) UNIQUE NOT NULL,
    name VARCHAR(255),
    city VARCHAR(63),
    district VARCHAR(63),
    address_office VARCHAR(255)
);
`

const CreateProjectsInfoTable = `
CREATE TABLE IF NOT EXISTS info(
    id SERIAL PRIMARY KEY,
    code VARCHAR(63) UNIQUE NOT NULL,
    project_code VARCHAR(63),
    type VARCHAR(63),
    name VARCHAR(255),
    CONSTRAINT fk_project
    	FOREIGN KEY (project_code)
    	REFERENCES projects(code)
    	ON DELETE CASCADE
    	ON UPDATE CASCADE
);
`

const CreateBuildingsTable = `
CREATE TABLE IF NOT EXISTS buildings (
    id SERIAL PRIMARY KEY,
    project_code VARCHAR(63),
    code VARCHAR(63) UNIQUE NOT NULL,
    name VARCHAR(255),
    liter VARCHAR(3),
    delivery_date DATE,
    building_address VARCHAR(255),
    CONSTRAINT fk_project
    	FOREIGN KEY (project_code)
    	REFERENCES projects(code)
    	ON DELETE CASCADE
    	ON UPDATE CASCADE
);
`

const CreateSectionsTable = `
CREATE TABLE IF NOT EXISTS sections (
    id SERIAL PRIMARY KEY,
    code VARCHAR(63) UNIQUE NOT NULL,
    building_code VARCHAR(63),
    section_num VARCHAR(3),
    section_liter VARCHAR(3),
    CONSTRAINT fk_building
    	FOREIGN KEY (building_code)
    	REFERENCES buildings(code)
    	ON DELETE CASCADE
    	ON UPDATE CASCADE
);
`

const CreateFlatsTable = `
CREATE TABLE IF NOT EXISTS flats (
    id SERIAL PRIMARY KEY,
    code VARCHAR(63) UNIQUE NOT NULL,
    building_code VARCHAR(63),
    section_code VARCHAR(63),
    flat_number INTEGER,
    rooms_amount INTEGER,
    floor INTEGER,
    total_square NUMERIC(4, 2),
    living_square NUMERIC(4, 2),
    cost NUMERIC(10, 2),
    flat_img VARCHAR(255),
    floor_img VARCHAR(255),
    path VARCHAR(4095),
    status SMALLINT,
    place_type VARCHAR(63),
    CONSTRAINT fk_building
    	FOREIGN KEY (building_code)
	    REFERENCES buildings(code)
    	ON DELETE CASCADE
    	ON UPDATE CASCADE,
    CONSTRAINT fk_section
    	FOREIGN KEY (section_code)
    	REFERENCES sections(code)
    	ON DELETE CASCADE
    	ON UPDATE CASCADE
);
`

const CreateTagsTable = `
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    code VARCHAR(63) UNIQUE NOT NULL,
    flat_code VARCHAR(63),
    name VARCHAR(255),
    CONSTRAINT fk_tag
    	FOREIGN KEY (flat_code)
	    REFERENCES flats(code)
    	ON DELETE CASCADE
    	ON UPDATE CASCADE
);
`

const CreateMessagesTgIdCreatedAtIndex = `CREATE INDEX IF NOT EXISTS idx_messages_tg_id_created_at ON messages(tg_id, created_at ASC);`
const CreateBuildingsIndex = `CREATE INDEX IF NOT EXISTS idx_buildings_project_code ON buildings(project_code);`
const CreateFlatsIndex = `CREATE INDEX IF NOT EXISTS idx_flats_building_code ON flats(building_code);`
const CreateTagsIndex = `CREATE INDEX IF NOT EXISTS idx_tags_flat_code ON tags(flat_code);`
const CreateInfoIndex = `CREATE INDEX IF NOT EXISTS idx_info_project_code ON info(project_code);`

const FlatsQuery = `
SELECT
p.name as project_name,
p.city,
p.district,
p.address_office,
b.building_address,
b.name as building_name,
f.flat_number,
f.rooms_amount,
f.floor,
f.total_square,
f.living_square,
f.cost,
f.flat_img,
f.floor_img,
f.path,
f.place_type
FROM flats f
LEFT JOIN buildings b ON b.code = f.building_code
LEFT JOIN projects p ON p.code = b.project_code
LEFT JOIN tags t ON t.flat_code = f.code
LEFT JOIN info i ON i.project_code = p.code
`

func GetAll(tableName string) string {
	return fmt.Sprintf(`SELECT * FROM %s`, tableName)
}

func Get(tableName string, idName string, id uint64) string {
	return fmt.Sprintf(`
		SELECT * FROM %s WHERE %s = %d;
	`, tableName, idName, id)
}

func GetSort(tableName string, idName string, id uint64, sortName string, direction string) string {
	return fmt.Sprintf(`
		SELECT * FROM %s WHERE %s = %d ORDER BY %s %s;
	`, tableName, idName, id, sortName, direction)
}

func GetS(tableName string, idName string, id string) string {
	return fmt.Sprintf(`
		SELECT * FROM %s WHERE %s = '%s';
	`, tableName, idName, id)
}

func GetOneByMinValue(tableName string, idName string, minValue string) string {
	return fmt.Sprintf(`
		SELECT DISTINCT ON (%s)
    		*
		FROM %s
		ORDER BY %s, %s ASC LIMIT 1;
	`, idName, tableName, idName, minValue)
}

func Create(tableName string, fields []string, values []interface{}) string {
	return fmt.Sprintf(`
		INSERT INTO %s (%s) VALUES (%s);
	`, tableName, strings.Join(fields, ","), helper.ConvertValuesToSQLCreate(values))
}

func Update(tableName string, idName string, id uint64, fields []string, values []interface{}) string {
	return fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE %s = %d;
	`, tableName, helper.ConvertValuesToSQLUpdate(fields, values), idName, id)
}

func UpdateS(tableName string, idName string, id string, fields []string, values []interface{}) string {
	return fmt.Sprintf(`
		UPDATE %s
		SET %s
		WHERE %s = '%s';
	`, tableName, helper.ConvertValuesToSQLUpdate(fields, values), idName, id)
}

func Delete(tableName string, idName string, id uint64) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE %s = %d;`, tableName, idName, id)
}

func Count(tableName string, idName string, id uint64) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s = %d;`, tableName, idName, id)
}
