package handler

import (
	"fmt"
	"goTestProject/logic/dao"
	"goTestProject/tools"
	"log"
	"os"
	"strings"
	"sync"

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

var wg sync.WaitGroup

func FileCollect(c *gin.Context) {
	FilsInPath(baseFilePath)
	wg.Wait()
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
	wg.Add(1)
	go addFile(fileList)

}

func addFile(files []dao.File) {
	addErr := fileDao.AddBatch(files)
	if addErr != nil {
		log.Fatal("err02...")
	}
	defer wg.Done()
}

// Random file
// @Summary      Random file
// @Description  Random file
// @Tags         file
// @Accept       json
// @Produce      json
// @Success      200 string json "{"code":200,"data":{},"message":"success"}"
// @Router       /file/random [get]
func Random(c *gin.Context) {
	file := fileDao.Random()
	tools.SuccessWithMsg(c, "Query random file ok!", file)
}

var ch = make(chan int)

type ChannelTestParam struct {
	Id int `form:"id" json:"id" binding:"required"`
}

func ChannelTest(c *gin.Context) {
	// ch := make(chan int)
	var param ChannelTestParam
	if err := c.ShouldBind(&param); err != nil {
		tools.FailWithMsg(c, err.Error())
		return
	}
	go func(num int) {
		ch <- num
	}(param.Id)

}

func ChannelConsume() {
	for {
		fmt.Println("Channel consume start.........")
		fmt.Println(<-ch)
		fmt.Println("Channel consume end...........")
	}
}

var dst = "./upload/"

func UploadFile(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, dst+file.Filename)
	tools.SuccessWithMsg(c, fmt.Sprintf("'%s' uploaded!", file.Filename), nil)
}
