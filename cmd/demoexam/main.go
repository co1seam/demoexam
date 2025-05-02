package main

import (
	"database/sql"
	"demoexam/handler"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered: ", r)
		}
	}()

	db, err := InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handlers := handler.New(db)

	a := app.New()
	w := a.NewWindow("Учет поставщиков")
	w.Resize(fyne.NewSize(800, 600))

	partners, err := handlers.LoadPartners()
	if err != nil {
		log.Fatal(err)
	}

	partnersList := widget.NewList(
		func() int { return len(partners) },
		func() fyne.CanvasObject { return widget.NewLabel("Шаблон") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(partners[id].Name)
		},
	)

	partnersList.OnSelected = func(id widget.ListItemID) {
		handlers.ShowPartnerForm(w, &partners[id])
	}

	addPartner := widget.NewButton("Добавить партнера", func() {
		handlers.ShowPartnerForm(w, nil)
	})

	products, err := handlers.LoadProducts()
	if err != nil {
		log.Fatal(err)
	}

	productsList := widget.NewList(
		func() int { return len(products) },
		func() fyne.CanvasObject { return widget.NewLabel("Продукты") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(products[id].Name)
		},
	)

	productsList.OnSelected = func(id widget.ListItemID) {
		handlers.ShowProductForm(w, &products[id])
	}

	addProduct := widget.NewButton("Добавить продукт", func() {
		handlers.ShowProductForm(w, nil)
	})

	sales, err := handlers.LoadSales()
	if err != nil {
		log.Fatal(err)
	}

	salesList := widget.NewList(
		func() int { return len(sales) },
		func() fyne.CanvasObject { return widget.NewLabel("Продажи") },
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(sales[id].SaleDate)
		},
	)

	salesList.OnSelected = func(id widget.ListItemID) {
		handlers.ShowSalesForm(w, &sales[id], partners, products)
	}

	addSales := widget.NewButton("Добавить продажу", func() {
		handlers.ShowSalesForm(w, nil, partners, products)
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("Поставщики", container.NewBorder(nil, addPartner, nil, nil, partnersList)),
		container.NewTabItem("Продажи", container.NewBorder(nil, addSales, nil, nil, salesList)),
		container.NewTabItem("Продукты", container.NewBorder(nil, addProduct, nil, nil, productsList)),
	)
	w.SetContent(tabs)
	go func() {
		for range time.Tick(time.Second) {
			fyne.Do(func() {
				partners, err = handlers.LoadPartners()
				if err != nil {
					log.Fatal(err)
				}
				partnersList.Refresh()

				products, err = handlers.LoadProducts()
				if err != nil {
					log.Fatal(err)
				}
				productsList.Refresh()

				sales, err = handlers.LoadSales()
				if err != nil {
					log.Fatal(err)
				}
				salesList.Refresh()

			})
		}
	}()
	w.Show()
	a.Run()
}

func InitDatabase() (*sql.DB, error) {
	var err error

	db, err := sql.Open("sqlite3", "./partners.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS partners 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			rating INTEGER NOT NULL,
			address TEXT NOT NULL,
			director TEXT NOT NULL,
			phone TEXT NOT NULL,
			email TEXT NOT NULL
		);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS products
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			product_name TEXT NOT NULL,
			product_article TEXT NOT NULL,
			description TEXT NOT NULL,
			price FLOAT NOT NULL
		);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS sales 
		(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			partner_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			sale_date TEXT NOT NULL,
			FOREIGN KEY (partner_id) REFERENCES partners(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		);`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
