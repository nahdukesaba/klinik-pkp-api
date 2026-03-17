package rusun

// GET RESPONSE
type GetRusunV1Response Rusun

// ADD RESPONSE
// Contains only fields returned by PostRusunHandler.
type AddRusunV1Response struct {
	ID        string   `json:"id"`
	ImageURLs []string `json:"image_urls"`
}

func ToGetRusunV1Response(data Rusun) GetRusunV1Response {
	return GetRusunV1Response(data)
}

func ToGetRusunV1Responses(data []Rusun) []GetRusunV1Response {
	responses := make([]GetRusunV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetRusunV1Response(item)
	}

	return responses
}
