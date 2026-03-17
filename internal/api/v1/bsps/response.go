package bsps

// GET RESPONSE
type GetBSPSV1Response BSPS

// ADD RESPONSE
// Contains only fields returned by PostBSPSHandler.
type AddBSPSV1Response struct {
	ID string `json:"id"`
}

func ToGetBSPSV1Response(data BSPS) GetBSPSV1Response {
	return GetBSPSV1Response(data)
}

func ToGetBSPSV1Responses(data []BSPS) []GetBSPSV1Response {
	responses := make([]GetBSPSV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetBSPSV1Response(item)
	}

	return responses
}
