package utils

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func ValidateCoordinate(coordinate Coordinate) bool {
	if coordinate.Latitude < -90 || coordinate.Latitude > 90 {
		return false
	}

	if coordinate.Longitude < -180 || coordinate.Longitude > 180 {
		return false
	}

	return true
}
