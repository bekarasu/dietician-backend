package response

import "dietician.local/services/medical-service/internal/medical/repository"

type UploadDetailResponse struct {
	Upload   repository.MedicalUpload         `json:"upload"`
	Metadata []repository.MedicalFileMetadata `json:"metadata"`
}
