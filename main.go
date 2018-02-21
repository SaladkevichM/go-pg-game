package main

import (
	"fmt"
	"log"
	"math"

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
	Id     uint32
	Name   string
	Lat    float32
	Number uint64
}

func main() {

	db, err := sqlx.Connect("postgres", "user=docker dbname=docker password=docker sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(schema)

	/*
		cities := []City{}
		db.Select(&cities, "SELECT * FROM cities LIMIT 100")
		san, los := cities[0], cities[1]
		fmt.Println(san, los)
	*/

	// http://jmoiron.github.io/sqlx/#altScanning

	r, err := db.Exec("INSERT INTO cities (name, lat, number) VALUES ($1, $2, $3)", "Dublin", 332.5, math.MaxInt64)
	if err != nil {
		log.Fatalln(err)
	}
	r.LastInsertId()

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

	fmt.Println(list[0]["id"], list[0]["name"], list[0]["lat"])
	fmt.Println(list[1]["id"], list[1]["name"], list[1]["lat"])
	fmt.Println(list[2]["id"], list[2]["name"], list[2]["lat"])

}
