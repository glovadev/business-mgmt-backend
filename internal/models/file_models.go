package models

type UploadUrlRequest struct {
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Purpose     string `json:"purpose"` // EXPENSE_RECEIPT, PURCHASE_INVOICE, etc.
}

type UploadUrlResponse struct {
	Data struct {
		UploadURL string `json:"upload_url"`
		ObjectKey string `json:"object_key"`
		ExpiresIn int    `json:"expires_in"`
	} `json:"data"`
}

type DownloadUrlResponse struct {
	Data struct {
		DownloadURL string `json:"download_url"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"data"`
}
