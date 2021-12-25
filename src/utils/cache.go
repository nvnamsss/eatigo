package utils

import "fmt"

func GetRestaurantKey(place, cursor string) string {
	return fmt.Sprintf("restaurant_%v_%v", place, cursor)
}
