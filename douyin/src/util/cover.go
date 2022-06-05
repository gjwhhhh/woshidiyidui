package util

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

// GetSnapshot 根据视频生成封面cover存储到本地
// frameNum 第几帧
func GetSnapshot(videoPath, coverPath string, frameNum int) (err error) {
	// 读取本地文件流
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter(
			"select",
			ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output(
			"pipe:",
			ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(
			buf,
			os.Stdout).
		Run()
	if err != nil {
		log.Fatal("read local file err ,", err)
		return err
	}

	// 将视频帧转换为cover
	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("decode buffer as cover err ,", err)
		return err
	}

	// 将cover存储到本地
	err = imaging.Save(img, coverPath)
	if err != nil {
		log.Fatal("cover store local err ,", err)
		return err
	}
	return nil
}
