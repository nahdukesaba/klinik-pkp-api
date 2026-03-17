package district

// GET RESPONSE
type GetDistrictV1Response District

// ADD RESPONSE
// Contains only fields returned by PostDistrictHandler.
type AddDistrictV1Response struct {
	ID string `json:"id"`
}

func ToGetDistrictV1Response(data District) GetDistrictV1Response {
	return GetDistrictV1Response(data)
}

func ToGetDistrictV1Responses(data []District) []GetDistrictV1Response {
	responses := make([]GetDistrictV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetDistrictV1Response(item)
	}

	return responses
}
