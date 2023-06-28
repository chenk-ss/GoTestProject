package dao

import (
	"time"
)

type File struct {
	Id         int64 `gorm:"primaryKey"`
	Name       string
	Type       string
	BasePath   string
	Path       string
	FullPath   string
	Size       int64
	CreateTime time.Time
}

func (f *File) getTableName() string {
	return "file"
}

func (f *File) AddBatch(files []File) error {
	if len(files) == 0 {
		return nil
	}
	mySqlDB.Table(f.getTableName()).CreateInBatches(&files, 100)
	return nil
}

func (f *File) Random() (file File) {
	mySqlDB.Table(f.getTableName()).Select("*").Order("RAND()").Take(&file)
	return file
}
