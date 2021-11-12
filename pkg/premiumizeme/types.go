package premiumizeme

type DeleteTransferResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GenerateZipResponse struct {
	Status   string `json:"status"`
	Location string `json:"location"`
}

type ListTransfersResponse struct {
	Status    string     `json:"status"`
	Transfers []Transfer `json:"transfers"`
}

type CreateTransferResponse struct {
	Status  string `json:"status"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Transfer struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Message  string  `json:"message"`
	Status   string  `json:"status"`
	Progress float64 `json:"progress"`
	Src      string  `json:"src"`
	FolderID string  `json:"folder_id"`
	FileID   string  `json:"file_id"`
}
