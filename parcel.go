package main

import (
	"database/sql"
	"log"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	db, err := sql.Open("sqlite3", "./parcels.db")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO parcel (number, client, status, adress, created_at) VALUES (:number, :client, :status, :adress, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("adress", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		log.Println(err)
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	db, err := sql.Open("sqlite3", "./parcels.db")
	if err != nil {
		log.Println(err)
		return p, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	err = row.Scan(Parcel{})
	if err != nil {
		log.Println(err)
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	db, err := sql.Open("sqlite3", "./parcels.db")
	if err != nil {
		log.Println(err)
		return p, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(append(res, Parcel{}))
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	db, err := sql.Open("sqlite3", "./parcels.db")
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE status SET status = :status WHERE number = :id",
		sql.Named("status", status),
		sql.Named("id", number))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	db, err := sql.Open("sqlite3", "./parcels.db")
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	row := db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	var status string
	err = row.Scan(&status)
	if err != nil {
		log.Println(err)
		return err
	}

	if status == "registered" {
		_, err := db.Exec("UPDATE status SET address = :address WHERE number = :number",
			sql.Named("number", number),
			sql.Named("address", address))
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		log.Println("[SetAddress] parcel not registered:", status)
		return nil
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	db, err := sql.Open("sqlite3", "./parcels.db")
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	row := db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	var status string
	err = row.Scan(&status)
	if err != nil {
		log.Println(err)
		return err
	}

	if status == "registered" {
		_, err := db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
