package bll

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

const maxFileSize = 8 << 20 // 8MiB

var supportedFileTypeRegex = regexp.MustCompile(`.(jpg|jpeg|JPG|png|PNG|gif|GIF|ico|ICO)$`)

func UploadPublicFile(file *multipart.FileHeader) (entity.File, error) {
	// only support img now
	if !supportedFileTypeRegex.MatchString(file.Filename) {
		return entity.File{}, util.NewErrorWithForbidden("unsupported file type")
	}

	record := entity.File{
		FileName:   file.Filename,
		AccessType: int8(entity.FileTypePublic),
	}
	err := uploadFile(file, &record)
	return record, err
}

func UploadTodoFile(userID, todoID int64, file *multipart.FileHeader) (entity.File, error) {
	_, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.File{}, err
	}

	record := entity.File{
		FileName:   file.Filename,
		AccessType: int8(entity.FileTypeTodo),
		RelatedID:  todoID,
	}
	if err := uploadFile(file, &record); err != nil {
		return entity.File{}, err
	}

	if err := dal.InsertTodoFile(todoID, record.ID); err != nil {
		return entity.File{}, fmt.Errorf("fails to upload todo file: %w", err)
	}

	return record, nil
}

func uploadFile(file *multipart.FileHeader, record *entity.File) error {
	write := func(err error) error {
		return fmt.Errorf("fails to upload file: %w", err)
	}

	if file.Size > maxFileSize {
		return util.NewError(otodo.ErrRequestEntityTooLarge, "file too large")
	}

	if err := dal.InsertFile(record); err != nil {
		return write(err)
	}

	// TODO[pref]: avoid duplicate save, remove :id in template?
	record.FileServerID = otodo.Conf.Server.ID
	record.FilePath = applyFilePathTemplate(record)
	if err := util.SaveFile(file, record.FilePath); err != nil {
		return write(err)
	}

	if err := dal.SaveFile(record); err != nil {
		return write(err)
	}

	return nil
}

func GetFile(fileID int64) (*entity.File, error) {
	file, err := dal.SelectFile(fileID)
	return file, fmt.Errorf("fails to get file: %w", err)
}

// Get file path, auto
func GetFilePath(userID int64, fileID string) (string, error) {
	id, err := strconv.ParseInt(fileID, 10, 64)
	if err != nil {
		return GetPreSignFilePath(fileID)
	}

	file, err := OwnFile(userID, id)
	if err != nil {
		return "", err
	}

	return getFilePath(file), nil
}

func OwnFile(userID, fileID int64) (*entity.File, error) {
	file, err := GetFile(fileID)
	if err != nil {
		return nil, err
	}

	switch entity.FileAccessType(file.AccessType) {
	case entity.FileTypePublic:
		break

	case entity.FileTypeTodo:
		if _, err := OwnTodo(userID, file.RelatedID); err != nil {
			return nil, util.NewErrorWithForbidden("unable to get non-owned file: %w", err)
		}

	default:
		return nil, fmt.Errorf("invalid file access type: %v", file.AccessType)
	}

	return file, nil
}

func applyFilePathTemplate(file *entity.File) string {
	template := otodo.Conf.Server.FilePathTemplate
	template = strings.ReplaceAll(template, ":id", strconv.FormatInt(file.ID, 10))
	template = strings.ReplaceAll(template, ":ext", filepath.Ext(file.FileName))
	template = strings.ReplaceAll(template, ":name", file.FileName)
	template = strings.ReplaceAll(template, ":path", file.FilePath)
	template = strings.ReplaceAll(template, ":date", file.CreatedAt.Format("20060102"))
	return template
}

func getFilePath(file *entity.File) string {
	// TODO[feat]: If exist multi servers, how to get file? maybe we need redirect
	return file.FilePath
}

func ParseConfigFile(path string) (config otodo.ConfigSms, err error) {
	viper.SetConfigFile(path)
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("unmarshall conf failed:%s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被人修改了。。")
		if err := viper.Unmarshal(config); err != nil {
			panic(fmt.Errorf("Unmarshal config file failed: %s", err))
		}
	})
	return config, err
}
