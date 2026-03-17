package village

// GET RESPONSE
type GetVillageV1Response Village

// ADD RESPONSE
// Contains only fields returned by PostVillageHandler.
type AddVillageV1Response struct {
	ID string `json:"id"`
}

func ToGetVillageV1Response(data Village) GetVillageV1Response {
	return GetVillageV1Response(data)
}

func ToGetVillageV1Responses(data []Village) []GetVillageV1Response {
	responses := make([]GetVillageV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetVillageV1Response(item)
	}

	return responses
}
