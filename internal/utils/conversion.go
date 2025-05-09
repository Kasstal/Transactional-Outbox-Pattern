package utils

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"math/big"
)

// Преобразует строку в pgtype.Text
func ToText(val string) pgtype.Text {
	return pgtype.Text{String: val, Valid: val != ""}
}

// Преобразует float64 в pgtype.Numeric
func ToNumeric(val float64) pgtype.Numeric {
	b := new(big.Float).SetFloat64(val)
	i := new(big.Int)
	b.Mul(b, big.NewFloat(100)).Int(i) // Умножаем на 100 для хранения копеек
	return pgtype.Numeric{
		Int:   i,
		Exp:   -2, // Два знака после запятой (копейки)
		Valid: true,
	}
}

// Преобразует UUID в pgtype.UUID
func ToUUID(val uuid.UUID) pgtype.UUID {

	return pgtype.UUID{Bytes: val, Valid: true}
}

// Преобразует bool в pgtype.Bool
func ToBool(val bool) pgtype.Bool {
	return pgtype.Bool{Bool: val, Valid: true}
}
