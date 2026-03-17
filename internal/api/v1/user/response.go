package user

// GET RESPONSE
type GetUserV1Response User

func ToGetUserV1Response(data User) GetUserV1Response {
	return GetUserV1Response(data)
}

func ToGetUserV1Responses(data []User) []GetUserV1Response {
	responses := make([]GetUserV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetUserV1Response(item)
	}

	return responses
}
