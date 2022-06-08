package oss

import (
	"bufio"
	"douyin/src/global"
	"encoding/base64"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// UploadVideo 上传视频，返回视频url和封面url
func UploadVideo(fileName string) (string, string, error) {
	// TODO 使用协程优化
	// 上传视频到云端
	videoUrl, err := FileUpLoad(fileName)
	if err != nil {
		return "", "", err
	}
	// 上传封面到云端
	coverUrl := getUrlPath(fmt.Sprintf("%s.jpg", strings.Split(fileName, ".")[0]))
	err = NetFileDump(fmt.Sprintf("%s?x-oss-process=video/snapshot,t_1,f_jpg,w_0,h_0,m_fast", videoUrl), coverUrl)
	if err != nil {
		return "", "", err
	}
	return videoUrl, fmt.Sprintf("%s%s", global.OssSetting.UrlPathPrefix, coverUrl), nil
}

// NetFileDump 网络文件转存
func NetFileDump(srcUrl, dstUrl string) error {
	defaultBucketInstance, err := GetDefaultBucketInstance()
	if err != nil {
		return err
	}
	res, err := http.Get(srcUrl)
	if err != nil {
		return err
	}
	// defer后的为延时操作，通常用来释放相关变量
	defer res.Body.Close()
	reader := bufio.NewReader(res.Body)
	err = defaultBucketInstance.PutObject(dstUrl, reader)
	if err != nil {
		return err
	}
	return nil
}

// FileUpLoad 上传文件到默认bucket
func FileUpLoad(fileName string) (string, error) {
	defaultBucketInstance, err := GetDefaultBucketInstance()
	if err != nil {
		return "", err
	}
	return fileUpLoad(fileName, defaultBucketInstance)
}

// FileUpLoadByBucket 上传文件到指定BucketName
func FileUpLoadByBucket(fileName, bucketName string) (string, error) {
	// 加载指定Bucket
	bucketInstance, err := GetBucketInstance(bucketName)
	if err != nil {
		return "", err
	}
	return fileUpLoad(fileName, bucketInstance)
}

// fileUpLoad 上传文件到指定Bucket
func fileUpLoad(fileName string, bucket *oss.Bucket) (string, error) {
	// 生成url路径
	urlPath := getUrlPath(fileName)
	// 检查本地文件
	localPath := LocalFilePathPrefix + fileName
	err := checkLocalFile(localPath)
	if err != nil {
		log.Println("LocalFile err:", localPath)
		return "", err
	}
	// 上传文件
	err = bucket.PutObjectFromFile(urlPath, localPath)
	if err != nil {
		log.Println("Upload err:", err)
		// TODO 文件上传云端失败处理
		return "", err
	}
	// 删除本地文件
	err = os.Remove(localPath)
	if err != nil {
		// TODO 如果引入多线程，删除文件可能出现问题
		return "", err
	}
	// 返回url
	return fmt.Sprintf("%s%s", global.OssSetting.UrlPathPrefix, urlPath), nil
}

// checkLocalFile 检查本地文件是否正常
func checkLocalFile(localPath string) error {
	if _, err := os.Stat(localPath); err != nil {
		return err
	}
	return nil
}

// getUrlPath 根据文件名设置urlPath
func getUrlPath(fileName string) string {
	builder := &strings.Builder{}
	// 设置文件名 BASE
	split := strings.Split(fileName, ".")
	// 设置为时间文件夹
	builder.WriteString(fmt.Sprintf("%d/%d/%d/%s/", time.Now().Year(), time.Now().Month(), time.Now().Day(), split[1]))
	// 文件url进行base64编码
	builder.WriteString(base64.URLEncoding.EncodeToString([]byte(split[0])))
	builder.WriteByte('.')
	// 文件后缀
	builder.WriteString(split[1])
	return builder.String()
}
