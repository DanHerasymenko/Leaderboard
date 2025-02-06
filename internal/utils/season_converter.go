package utils

import (
	"fmt"
	"time"
	_ "time"
)

func GetCurrentSeasonAndTime() (time.Time, string) {

	now := time.Now()
	year := now.Year()
	month := now.Month()

	var season string
	switch month {
	case 12, 1, 2:
		season = "Winter"
	case 3, 4, 5:
		season = "Spring"
	case 6, 7, 8:
		season = "Summer"
	case 9, 10, 11:
		season = "Autumn"
	}

	return now, fmt.Sprintf("%s%d", season, year)
}
