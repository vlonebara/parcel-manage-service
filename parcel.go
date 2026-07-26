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
	res, err := s.db.Exec("INSERT INTO parcel (number, client, status, address, created_at) VALUES (:number, :client, :status, :address, :created_at)",
		sql.Named("number", p.Number),
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
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

	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		log.Println(err)
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Parcel
		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))

	var status string
	err := row.Scan(&status)
	if err != nil {
		log.Println(err)
		return err
	}

	if status == "registered" {
		_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
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
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	var status string
	err := row.Scan(&status)
	if err != nil {
		log.Println(err)
		return err
	}

	if status == "registered" {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}
