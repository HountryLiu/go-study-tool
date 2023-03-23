package ffmpeg

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/HountryLiu/go-study-tool/utils"
	"github.com/gin-gonic/gin"
)

// @Deprecated
// @Tags FFmpeg
// @Summary 测试swagger弃用接口语法
// @Description 测试swagger弃用接口语法
// @Accept multipart/form-data
// @Produce application/json
// @Success 200 {object} object{no=int,data=string,msg=string}
// @Router /api/ffmpeg [post]
func test(ctx *gin.Context) {
	utils.Success(ctx)
}

// @Tags FFmpeg
// @Summary FFmpeg api
// @Description FFmpeg api
// @Accept multipart/form-data
// @Produce application/json
// @Success 200 {object} object{no=int,data=string,msg=string}
// @Router /api/ffmpeg [get]
func Index(ctx *gin.Context) {
	url := "http://127.0.0.1:8099/api/ffmpeg/"
	method := "GET"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	//获取视频第一帧
	_ = writer.WriteField("arg", "-vframes 1 -f singlejpeg -")
	file, errFile2 := os.Open("data/video/test.mp4")
	defer file.Close()
	part2, errFile2 := writer.CreateFormFile("file", filepath.Base("data/video/test.mp4"))
	_, errFile2 = io.Copy(part2, file)
	if errFile2 != nil {
		utils.Error(ctx, utils.InternalServerError, errFile2)
	}
	err := writer.Close()
	if err != nil {
		utils.Error(ctx, utils.InternalServerError, err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		utils.Error(ctx, utils.InternalServerError, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		utils.Error(ctx, utils.InternalServerError, err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.Error(ctx, utils.InternalServerError, err)
	}
	err = utils.WriteFile("data/video/test.jpg", string(body))
	if err != nil {
		utils.Error(ctx, utils.InternalServerError, err)
	}
	utils.Success(ctx)
}
