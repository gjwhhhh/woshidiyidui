package oss

import (
	"encoding/base64"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
	"strings"
	"time"
)

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
	return fmt.Sprintf("%s%s", UrlPathPrefix, urlPath), nil
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
	// 设置为时间文件夹
	builder.WriteString(fmt.Sprintf("%d/%d/%d/", time.Now().Year(), time.Now().Month(), time.Now().Day()))
	// 设置文件名 BASE
	split := strings.Split(fileName, ".")
	// 文件url进行base64编码
	builder.WriteString(base64.URLEncoding.EncodeToString([]byte(split[0])))
	builder.WriteByte('.')
	// 文件后缀
	builder.WriteString(split[1])
	return builder.String()
}
