package request

type CreateUploadRequest struct {
	UploadType  string `json:"uploadType"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
