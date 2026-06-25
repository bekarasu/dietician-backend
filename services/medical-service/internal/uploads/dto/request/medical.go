package request

type UploadType string

const (
	BloodTest UploadType = "BloodTest"
)

type FileData struct {
	FileName    string `validate:"required"`
	FileSize    int64  `validate:"required"`
	ContentType string `validate:"required,eq=application/pdf"`
	Data        []byte `validate:"required"`
}

type CreateUploadRequest struct {
	UploadType  UploadType `json:"uploadType" validate:"required,oneof=BloodTest"`
	File        *FileData  `json:"-" validate:"required"`
}
