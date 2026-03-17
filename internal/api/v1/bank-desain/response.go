package bank_desain

// GET RESPONSE
type GetBankDesainV1Response BankDesain

// ADD RESPONSE
// Contains only fields returned by PostBankDesainHandler.
type AddBankDesainV1Response struct {
	ID        string   `json:"id"`
	ImageURLs []string `json:"image_urls"`
	FileURLs  []string `json:"file_urls"`
}

func ToGetBankDesainV1Response(data BankDesain) GetBankDesainV1Response {
	return GetBankDesainV1Response(data)
}

func ToGetBankDesainV1Responses(data []BankDesain) []GetBankDesainV1Response {
	responses := make([]GetBankDesainV1Response, len(data))

	for i, item := range data {
		responses[i] = ToGetBankDesainV1Response(item)
	}

	return responses
}
