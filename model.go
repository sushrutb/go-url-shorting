package main

import (
	"database/sql"
	"log"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type short_url struct {
	ID          int    `json:"id"`
	Destination string `json:"destination"`
	Shortcode   string `json:"shortcode"`
}

func (s *short_url) getShortUrl(db *sql.DB) error {
	return db.QueryRow("SELECT id, destination FROM short_urls WHERE shortcode=?", s.Shortcode).Scan(&s.ID, &s.Destination)
}

func (s *short_url) createShortUrl(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO short_urls(destination, shortcode) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	res, err := stmt.Exec(s.Destination, s.Shortcode)
	if err != nil {
		log.Fatal(err.Error())
		return err
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		} else {
			s.ID = int(id)
		}
	}
	return nil

}
func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=?", p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) deleteProduct(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM products WHERE id=?")
	if err == nil {
		_, err := stmt.Exec(p.ID)
		if err != nil {
			return err
		}
	}
	return err
}

func (p *product) updateProduct(db *sql.DB) error {
	stmt, err := db.Prepare("UPDATE products SET name=?, price=? WHERE id=?")
	if err == nil {
		_, err := stmt.Exec(p.Name, p.Price, p.ID)
		if err != nil {
			return err
		}
	}
	return err
}

func (p *product) createProduct(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO products(name, price) VALUES(?, ?)")

	if err != nil {
		return err
	}

	res, err := stmt.Exec(p.Name, p.Price)
	if err != nil {
		return err
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			return err
		} else {
			p.ID = int(id)
		}
	}

	return nil
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products LIMIT 10 OFFSET 0")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
