package handler

import (
	"database/sql"
	"demoexam"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "github.com/mattn/go-sqlite3"
)

type Handler struct {
	db *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) LoadPartners() ([]demoexam.Partner, error) {
	var partners []demoexam.Partner
	rows, err := h.db.Query(`SELECT * FROM partners;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p demoexam.Partner
		err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Rating, &p.Address, &p.Director, &p.Phone, &p.Email)
		if err != nil {
			return nil, err
		}
		partners = append(partners, p)
	}

	return partners, err
}

func (h *Handler) ShowPartnerForm(w fyne.Window, p *demoexam.Partner) {
	nameEntry := widget.NewEntry()
	typeSelect := widget.NewSelect([]string{"Дистрибьютор", "Розничный", "Оптовый"}, nil)
	ratingEntry := widget.NewEntry()
	addressEntry := widget.NewEntry()
	directorEntry := widget.NewEntry()
	phoneEntry := widget.NewEntry()
	emailEntry := widget.NewEntry()

	if p != nil {
		nameEntry.SetText(p.Name)
		typeSelect.SetSelected(p.Type)
		ratingEntry.SetText(strconv.Itoa(p.Rating))
		addressEntry.SetText(p.Address)
		directorEntry.SetText(p.Director)
		phoneEntry.SetText(p.Phone)
		emailEntry.SetText(p.Email)
	}

	dialog.ShowForm("Данные партнера", "Сохранить", "Отмена", []*widget.FormItem{
		{Text: "Название:", Widget: nameEntry},
		{Text: "Тип:", Widget: typeSelect},
		{Text: "Рейтинг:", Widget: ratingEntry},
		{Text: "Адрес:", Widget: addressEntry},
		{Text: "Директор:", Widget: directorEntry},
		{Text: "Телефон:", Widget: phoneEntry},
		{Text: "Почта:", Widget: emailEntry},
	}, func(b bool) {
		if b {
			rating, err := strconv.Atoi(ratingEntry.Text)
			if err != nil || rating < 0 || rating >= 100 {
				dialog.ShowError(fmt.Errorf("поле 'Рейтинг' пустое или имеет неправильное значение\n Значение должно быть в диапазоне [1, 100]\n%w", err), w)
				return
			}
			if p == nil {
				_, err := h.db.Exec(
					`INSERT INTO partners
					(name, type, rating, address, director, phone, email)
					VALUES (?,?,?,?,?,?,?);`,
					nameEntry.Text, typeSelect.Selected, rating, addressEntry.Text, directorEntry.Text, phoneEntry.Text, emailEntry.Text,
				)

				if err != nil {
					dialog.ShowError(err, w)
				}
			} else {
				_, err = h.db.Exec(
					`UPDATE partners SET
					name = ?, type = ?, rating = ?, address = ?, director = ?, phone = ?, email = ? WHERE id = ?`,
					nameEntry.Text, typeSelect.Selected, rating, addressEntry.Text, directorEntry.Text, phoneEntry.Text, emailEntry.Text, p.ID,
				)
				if err != nil {
					dialog.ShowError(err, w)
				}
			}

			dialog.ShowInformation("Успех", "Данные сохранены", w)
			w.Content().Refresh()
		}
	}, w)
}
