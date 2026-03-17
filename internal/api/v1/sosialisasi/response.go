package sosialisasi

// GET RESPONSE
type GetSosialisasiV1Response Sosialisasi

// ADD RESPONSE
// Contains only fields returned by PostSosialisasiHandler.
type AddSosialisasiV1Response struct {
	ID        string   `json:"id"`
	ImageURLs []string `json:"image_urls"`
}

func ToGetSosialisasiV1Response(data Sosialisasi) GetSosialisasiV1Response {
	return GetSosialisasiV1Response(data)
}

func ToGetSosialisasiV1Responses(data []Sosialisasi) []GetSosialisasiV1Response {
	responses := make([]GetSosialisasiV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetSosialisasiV1Response(item)
	}

	return responses
}
