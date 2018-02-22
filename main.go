package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `
DROP TABLE IF EXISTS cities;
CREATE TABLE cities (id SERIAL PRIMARY KEY NOT NULL, name varchar(80), lat real);
INSERT INTO cities(name, lat) VALUES ('San Francisco', '194.4');
INSERT INTO cities(name, lat) VALUES ('Los Angeles', '220.4');
INSERT INTO cities(name, lat) VALUES ('Paris', '80.4');`

type City struct {
	Id   int64
	Name string
	Lat  float64
}

type Mapper func(map[string]interface{}) interface{}

func CityRowMapper(row map[string]interface{}) interface{} {
	return &City{
		Id:   row["id"].(int64),
		Name: row["name"].(string),
		Lat:  row["lat"].(float64),
	}
}

func GetDbSlice(db *sqlx.DB, fn Mapper, query string, args ...interface{}) []interface{} {

	var result []interface{}
	rows, err := db.Queryx(query, args...)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		row := make(map[string]interface{})
		err = rows.MapScan(row)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, fn(row))
	}
	rows.Close()

	return result
}

func main() {

	db, err := sqlx.Connect("postgres", "user=docker dbname=docker password=docker sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(schema)

	list := GetDbSlice(db, CityRowMapper, "SELECT * FROM cities WHERE id=$1", 1)
	fmt.Println(list[0])

	list = GetDbSlice(db, CityRowMapper, "SELECT * FROM cities LIMIT 100")
	fmt.Println(list[0], list[1], list[2])

}
