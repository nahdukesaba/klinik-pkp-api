package province

// GET RESPONSE
type GetProvinceV1Response Province

// ADD RESPONSE
// Contains only fields returned by PostProvinceHandler.
type AddProvinceV1Response struct {
	ID string `json:"id"`
}

func ToGetProvinceV1Response(data Province) GetProvinceV1Response {
	return GetProvinceV1Response(data)
}

func ToGetProvinceV1Responses(data []Province) []GetProvinceV1Response {
	responses := make([]GetProvinceV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetProvinceV1Response(item)
	}

	return responses
}
