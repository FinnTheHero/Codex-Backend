package worker

import "mime/multipart"

type ProcessEPUBArgs struct {
	File *multipart.FileHeader `json:"file"`
}

func (ProcessEPUBArgs) Kind() string { return "process_epub" }
