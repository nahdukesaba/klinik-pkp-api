package utils

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func ValidateCoordinate(coord Coordinate) bool {
	if coord.Latitude < -90 || coord.Latitude > 90 {
		return false
	}

	if coord.Longitude < -180 || coord.Longitude > 180 {
		return false
	}

	return true
}
