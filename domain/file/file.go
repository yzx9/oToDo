package file

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

const maxFileSize = 8 << 20 // 8MiB

type File struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	FileName     string
	FileServerID string
	FilePath     string
	AccessType   FileAccessType // FileAccessType
	RelatedID    int64          // Depend on access type
}

var Notfound = fmt.Errorf("file not found")

func GetFile(fileID int64) (*File, error) {
	f, err := FileRepository.Find(fileID)

	if err != nil {
		return nil, Notfound
	}

	return f, nil
}

func UploadFile(
	accessType FileAccessType,
	relatedID int64,
	file *multipart.FileHeader,
) (File, error) {
	record := File{
		FileName:   file.Filename,
		AccessType: accessType,
		RelatedID:  relatedID,
	}

	if err := record.upload(file); err != nil {
		return File{}, err
	}

	return record, nil
}

func (record *File) upload(file *multipart.FileHeader) error {
	write := func(err error) error {
		return fmt.Errorf("fails to upload file: %w", err)
	}

	if file.Size > maxFileSize {
		return util.NewError(errors.ErrRequestEntityTooLarge, "file too large")
	}

	if err := FileRepository.Save(record); err != nil {
		return write(err)
	}

	// TODO[pref]: avoid duplicate save, remove :id in template?
	record.FileServerID = config.Server.ID
	record.FilePath = record.newFilePath()
	if err := util.SaveFile(file, record.FilePath); err != nil {
		return write(err)
	}

	if err := FileRepository.Save(record); err != nil {
		return write(err)
	}

	return nil
}

func (file *File) GetFilePath() string {
	// TODO[feat]: If exist multi servers, how to get file? maybe we need redirect
	return file.FilePath
}

func (file *File) newFilePath() string {
	template := config.Server.FilePathTemplate
	template = strings.ReplaceAll(template, ":id", strconv.FormatInt(file.ID, 10))
	template = strings.ReplaceAll(template, ":ext", filepath.Ext(file.FileName))
	template = strings.ReplaceAll(template, ":name", file.FileName)
	template = strings.ReplaceAll(template, ":path", file.FilePath)
	template = strings.ReplaceAll(template, ":date", file.CreatedAt.Format("20060102"))
	return template
}
