package request

type FileData struct {
	FileName    string
	FileSize    int64
	ContentType string
	Data        []byte
}

type CreateUploadRequest struct {
	UploadType  string    `json:"uploadType"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	File        *FileData `json:"-"`
}
