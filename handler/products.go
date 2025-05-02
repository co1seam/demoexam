package handler

import (
	"demoexam"
	"strconv"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"
)

func (h *Handler) LoadProducts() ([]demoexam.Product, error) {
	var products []demoexam.Product
	rows, err := h.db.Query(`SELECT * FROM products;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product demoexam.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Article, &product.Description, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (h *Handler) ShowProductForm(w fyne.Window, product *demoexam.Product) {
	nameEntry := widget.NewEntry()
	articleEntry := widget.NewEntry()
	descEntry := widget.NewEntry()
	priceEntry := widget.NewEntry()

	if product != nil {
		nameEntry.SetText(product.Name)
		articleEntry.SetText(product.Article)
		descEntry.SetText(product.Description)
		priceEntry.SetText(strconv.FormatFloat(product.Price, 'g', -1, 64))
	}
	
	dialog.ShowForm("Данные продукта", "Сохранить", "Отмена", []*widget.FormItem{
		{Text: "Название", Widget:nameEntry},
		{Text: "Артикул", Widget:articleEntry},
		{Text: "Описание", Widget:descEntry},
		{Text: "Цена", Widget:priceEntry},
	}, func(b bool) {
		if b {

			price, err := strconv.ParseFloat(priceEntry.Text, 64)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			if product != nil {
				_, err := h.db.Exec(
					`UPDATE products SET
					product_name = ?, product_article = ?, description = ?, price = ? WHERE id = ?`,
					nameEntry.Text, articleEntry.Text, descEntry.Text, price, product.ID,
				)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
			} else {
				_, err = h.db.Exec(
					`INSERT INTO products
					(product_name, product_article, description, price)
					VALUES (?,?,?,?)`,
					nameEntry.Text, articleEntry.Text, descEntry.Text, price,
				)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
			}

			dialog.ShowInformation("Успех", "Данные сохранены", w)
			w.Content().Refresh()
		}
	}, w)
}