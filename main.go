package main

/*

db = docker
user = docker@docker

CREATE TABLE cities (id SERIAL PRIMARY KEY NOT NULL, name varchar(80), lat real);

INSERT INTO cities(name, lat) VALUES ('San Francisco', '194.4');
INSERT INTO cities(name, lat) VALUES ('Los Angeles', '220.4');
INSERT INTO cities(name, lat) VALUES ('Paris', '80.4');

*/

import (
	"fmt"

	"database/sql"

	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
)

type City struct {
	Id   uint32
	Name string
	Lat  float32
}

func GetCitiesSQL(db *pg.DB) ([]City, error) {
	var cities []City
	_, err := db.Query(&cities, `SELECT * FROM cities`)
	return cities, err
}

func GetCityByObject(db *pg.DB) (City, error) {
	с := City{
		Id: 2,
	}
	err := db.Select(&с)
	return с, err
}

func GetCitiesClassic() {

	connStr := "host=127.0.0.1 port=32772 user=docker password=docker dbname=docker sslmode=disable"
	db, err := sql.Open("postgres", connStr) // driver name
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("select * from cities")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	cities := []City{}

	for rows.Next() {
		c := City{}
		err := rows.Scan(&c.Id, &c.Name, &c.Lat)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cities = append(cities, c)
	}
	for _, c := range cities {
		fmt.Println(c.Id, c.Name, c.Lat)
	}

}

func main() {

	db := pg.Connect(&pg.Options{
		Addr:     "localhost:32772",
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

}
