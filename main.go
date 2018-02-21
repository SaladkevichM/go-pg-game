package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const schema = `
DROP TABLE cities;
CREATE TABLE cities (id SERIAL PRIMARY KEY NOT NULL, name varchar(80), lat real, number bigint);
INSERT INTO cities(name, lat, number) VALUES ('San Francisco', '194.4', '0');
INSERT INTO cities(name, lat, number) VALUES ('Los Angeles', '220.4', '0');
INSERT INTO cities(name, lat, number) VALUES ('Paris', '80.4', '0');`

type City struct {
	Id     int64
	Name   string
	Lat    float64
	Number int64
}

type Mapper func(map[string]interface{}) interface{}

func CityMapper(row map[string]interface{}) interface{} {
	return &City{Id: row["id"].(int64), Name: row["name"].(string), Lat: row["lat"].(float64), Number: row["number"].(int64)}
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

	list := GetDbSlice(db, CityMapper, "SELECT * FROM cities WHERE id=$1", 1)

	fmt.Println(list[0])

	/*
		rows, err := db.Queryx("SELECT * FROM cities LIMIT 100")
		var list []map[string]interface{}

		for rows.Next() {
			row := make(map[string]interface{})
			err = rows.MapScan(row)
			if err != nil {
				log.Fatalln(err)
			}
			list = append(list, row)
		}

		rows.Close()
	*/

}
