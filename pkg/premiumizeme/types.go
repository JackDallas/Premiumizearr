package premiumizeme

type SimpleResponse struct {
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

type CreateFolderResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

type ListFoldersResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Content  []Item `json:"content"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
	FolderID string `json:"folder_id"`
}

type Item struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	CreatedAt  int    `json:"created_at"`
	MimeType   string `json:"mime_type"`
	Link       string `json:"link"`
	StreamLink string `json:"stream_link"`
}
type FolderItems struct {
	Status   string `json:"status"`
	Contant  []Item `json:"content"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
	FolderID string `json:"folder_id"`
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

const (
	ERROR_FOLDER_ALREADY_EXISTS = "This folder already exists."
)
