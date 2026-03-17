package kumuh

// GET RESPONSE
type GetKumuhV1Response KawasanKumuh

// ADD RESPONSE
// Contains only fields returned by PostKumuhHandler.
type AddKumuhV1Response struct {
	ID string `json:"id"`
}

func ToGetKumuhV1Response(data KawasanKumuh) GetKumuhV1Response {
	return GetKumuhV1Response(data)
}

func ToGetKumuhV1Responses(data []KawasanKumuh) []GetKumuhV1Response {
	responses := make([]GetKumuhV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetKumuhV1Response(item)
	}

	return responses
}
