package main

import "database/sql"

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1", p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)

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
