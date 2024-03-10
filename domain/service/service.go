package service

import (
	"context"
)

type GenerateIDService interface {
	GenerateID(ctx context.Context) string
}

type EncryptPasswordService interface {
	EncryptPassword(ctx context.Context, text string) string
}
