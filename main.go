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
	Id   uint32
	Name string
	Lat  float32
}

func main() {

	db, err := sqlx.Connect("postgres", "user=docker dbname=docker password=docker sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(schema)

	cities := []City{}
	db.Select(&cities, "SELECT * FROM cities LIMIT 100")
	//san, los := cities[0], cities[1]
	//fmt.Println(san, los)

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

	fmt.Println(list[0]["id"], list[0]["name"], list[0]["lat"])
	fmt.Println(list[1]["id"], list[1]["name"], list[1]["lat"])
	fmt.Println(list[2]["id"], list[2]["name"], list[2]["lat"])

	/*
		db := pg.Connect(&pg.Options{
			Addr:     "localhost:5432",
			Database: "docker",
			User:     "docker",
			Password: "docker",
		})
		defer db.Close()

		cities, err := GetCitiesSQL(db)
		if err != nil {
			panic(err)
		}

		fmt.Println(cities)

		city, err := GetCityByObject(db)
		if err != nil {
			panic(err)
		}

		fmt.Println(city)
		GetCitiesClassic()
	*/

}
