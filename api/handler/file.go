package handler

import (
	"fmt"
	"goTestProject/logic/dao"
	"goTestProject/tools"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var baseFilePath = "/Volumes/main"

var fileDao = dao.File{}

type QueryFileInfoParam struct {
	Id string `form:"id" json:"id" binding:"required"`
}

func FileInfo(c *gin.Context) {
	var param QueryFileInfoParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	tools.SuccessWithMsg(c, "Query file info ok!", param.Id)
}

func FileCollect(c *gin.Context) {
	FilsInPath(baseFilePath)
}

func FilsInPath(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(path + "err...")
		return
	}
	fileList := []dao.File{}
	for _, file := range files {
		if file.IsDir() {
			FilsInPath(path + "/" + file.Name())
		} else if file.Name() != ".DS_Store" && file.Name() != "Thumbs.db" {
			// fmt.Println(path + file.Name())
			fileInfo, e := file.Info()
			if e != nil {
				log.Fatal("err01...")
			}
			split := strings.Split(fileInfo.Name(), ".")
			fileList = append(fileList, dao.File{
				Name:       fileInfo.Name(),
				Type:       split[len(split)-1],
				BasePath:   baseFilePath,
				Path:       path[len(baseFilePath):],
				FullPath:   path,
				Size:       fileInfo.Size(),
				CreateTime: fileInfo.ModTime(),
			})
		}
	}
	addErr := fileDao.AddBatch(fileList)
	fmt.Println(len(files))
	if addErr != nil {
		log.Fatal("err02...")
	}
}

func Random(c *gin.Context) {
	file := fileDao.Random()
	tools.SuccessWithMsg(c, "Query random file ok!", file)
}
