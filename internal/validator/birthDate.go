package validator

import (
	"auth/internal/constants"
	"time"
)

type BirthDate struct{}

func newBirthDate() *BirthDate {
	return &BirthDate{}
}

func (b *BirthDate) IsValid(birthDate string) bool {
	return notMinor(birthDate)
}

func notMinor(birthDate string) bool {
	birthdate, _ := time.Parse("02/01/2006", birthDate)
	currentAge := calculateAge(birthdate, time.Now())
	return currentAge >= constants.CUSTOMER_MIN_AGE_TO_DEFINE_MINOR
}

func calculateAge(birthdate, today time.Time) int {
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}
