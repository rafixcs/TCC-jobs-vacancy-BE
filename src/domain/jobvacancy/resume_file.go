package jobvacancy

import "mime/multipart"

type JobVacancyResumeFile struct {
	File   multipart.File
	Header *multipart.FileHeader
}
