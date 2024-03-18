package handler

import (
	"github.com/gleb-korostelev/short-url.git/internal/db"
	"github.com/gleb-korostelev/short-url.git/internal/service"
)

type APIService struct {
	data db.DatabaseI
}

func NewAPIService(data db.DatabaseI) service.APIServiceI {
	return &APIService{
		data: data,
	}
}
