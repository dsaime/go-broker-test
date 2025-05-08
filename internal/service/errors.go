package service

import "errors"

var (
	ErrRequiredWorkerID  = errors.New("workerID обязателен")
	ErrRequiredAccountID = errors.New("accountID обязателен")
)
