package model

type GoogleDriveGetFilesResponse struct {
	Items []GoogleDriveGetFilesItemsResponse
}

type GoogleDriveGetFilesItemsResponse struct {
	Id       string `json:"id"`
	Filename string `json:"title"`
}
