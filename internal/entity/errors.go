package entity

import "errors"

var (
	// ErrNoID -.
	ErrNoID = errors.New("no such id")

	// ErrOrderExists -.
	ErrOrderExists = errors.New("order already exists")

	// ErrOrderNoExists -.
	ErrOrderNoExists = errors.New("no such order id")

	// ErrNotEnoughMoney -.
	ErrNotEnoughMoney = errors.New("not enough money")

	// ErrOrderMismatch -.
	ErrOrderMismatch = errors.New("order data doesnt match order data from db")

	// ErrCantChangeStatus -.
	ErrCantChangeStatus = errors.New("cant update status of committed/canceled order")

	// ErrEmptyReport -.
	ErrEmptyReport = errors.New("no any operations in this month")

	// ErrEmptyPage -.
	ErrEmptyPage = errors.New("this page is empty")

	// ErrNoService -.
	ErrNoService = errors.New("no such service")
)
