package handler

import (
	"demoexam"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (h *Handler) LoadSales() ([]demoexam.Sales, error) {
	var sales []demoexam.Sales
	rows, err := h.db.Query(`SELECT * FROM sales;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s demoexam.Sales
		err := rows.Scan(&s.ID, &s.PartnerID, &s.ProductID, &s.Quantity, &s.SaleDate)
		if err != nil {
			return nil, err
		}

		sales = append(sales, s)
	}


	return sales, nil
}

func (h *Handler) ShowSalesForm(w fyne.Window, sales *demoexam.Sales, partners []demoexam.Partner, products []demoexam.Product) error {
	partnersMap := make(map[string]int, len(partners))
	partnerNames := make([]string, 0, len(partners))
	productsMap := make(map[string]int, len(products))
	productNames := make([]string, 0, len(products))
	var productName string
	var partnerName string

	for _, partner := range partners {
		partnersMap[partner.Name] = partner.ID
		partnerNames = append(partnerNames, partner.Name)
		if sales != nil && partnersMap[partner.Name] == sales.PartnerID {
			partnerName = partner.Name
		}
	}

	for _, product := range products {
		productsMap[product.Name] = product.ID
		productNames = append(productNames, product.Name)
		if sales != nil && productsMap[product.Name] == sales.ProductID {
			partnerName = product.Name
		}
	}

	partnerNameSelect := widget.NewSelect(partnerNames, nil)  
	productNameSelect := widget.NewSelect(productNames, nil)
	quantityEntry := widget.NewEntry()
	saleDateEntry := widget.NewEntry()

	if sales != nil {
		partnerNameSelect.SetSelected(partnerName)
		productNameSelect.SetSelected(productName)
		quantityEntry.SetText(strconv.Itoa(sales.Quantity))
		saleDateEntry.SetText(sales.SaleDate)
	}

	dialog.ShowForm("Данные продаж", "Сохранить", "Отмена", []*widget.FormItem{
		{Text:"Партнер", Widget: partnerNameSelect},
		{Text:"Название продукта", Widget: productNameSelect},
		{Text:"Количество", Widget: quantityEntry},
		{Text:"Дата", Widget: saleDateEntry},
	}, func(b bool) {
		if b {
			quantity, err := strconv.Atoi(quantityEntry.Text)
			if err != nil || quantity <= 0{
				dialog.ShowError(fmt.Errorf("поле 'Цена' пустое или имеет неправильное значение\n Значение должно быть больше нуля\n%w", err), w)
				return
			}
			partnerID := partnersMap[partnerNameSelect.Selected]
			productID := productsMap[productName]

			if _, err := time.Parse("01.02.2006", saleDateEntry.Text); err != nil {
				dialog.ShowError(fmt.Errorf("некорректный формат даты. Используйте ДД.ММ.ГГГГ"), w)
				return
			}

			if sales != nil {
				if _, err = h.db.Exec(
					`UPDATE sales SET
					partner_id = ?, product_id = ?, quantity = ?, sale_date = ? WHERE id = ?;`,
					partnerID, productID, quantity, saleDateEntry.Text, sales.ID,
				); err != nil {
					dialog.ShowError(err, w)
					return
				}
			} else {
				_, err = h.db.Exec(
					`INSERT INTO sales
					(partner_id, product_id, quantity, sale_date)
					VALUES (?,?,?,?)`,
					partnerID, productID, quantityEntry.Text, saleDateEntry.Text,
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


	return nil
}