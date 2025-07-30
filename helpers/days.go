package helpers

import "fieldreserve/constants"

func DayIntToName(day int) string {
	switch day {
	case constants.Sunday:
		return "Sunday"
	case constants.Monday:
		return "Monday"
	case constants.Tuesday:
		return "Tuesday"
	case constants.Wednesday:
		return "Wednesday"
	case constants.Thursday:
		return "Thursday"
	case constants.Friday:
		return "Friday"
	case constants.Saturday:
		return "Saturday"
	default:
		return "Day Invalid"
	}
}
