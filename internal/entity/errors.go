package entity

import "errors"

var (
	ErrNoID             = errors.New("no such id")
	ErrOrderExists      = errors.New("order already exists")
	ErrOrderNoExists    = errors.New("no such order id")
	ErrNotEnoughMoney   = errors.New("not enough money")
	ErrOrderMismatch    = errors.New("order data doesnt match order data from db")
	ErrCantChangeStatus = errors.New("cant update status of committed/canceled order")
	ErrEmptyReport      = errors.New("no any operations in this month")
	ErrEmptyPage        = errors.New("this page is empty")
)
