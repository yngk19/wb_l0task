package service

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	"log/slog"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) OrderIsValid(data []byte, log *slog.Logger) (bool, *models.Order) {
	var order models.Order
	err := json.Unmarshal(data, &order)
	if err != nil {
		log.Info("Validation: message decoding fail", err)
		return false, nil
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(order)
	if err != nil {
		log.Info("Validation: message is invalid: ", err)
		return false, nil
	}
	return true, &order
}
