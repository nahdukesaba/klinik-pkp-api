package region

// GET RESPONSE
type GetRegionV1Response Region

// ADD RESPONSE
// Contains only fields returned by PostRegionHandler.
type AddRegionV1Response struct {
	ID string `json:"id"`
}

func ToGetRegionV1Response(data Region) GetRegionV1Response {
	return GetRegionV1Response(data)
}

func ToGetRegionV1Responses(data []Region) []GetRegionV1Response {
	responses := make([]GetRegionV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetRegionV1Response(item)
	}

	return responses
}
