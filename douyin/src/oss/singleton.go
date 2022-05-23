package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"sync"
)

var clientSingleton *oss.Client
var defaultBucketSingleton *oss.Bucket
var lock sync.Mutex

//GetClientInstance 用于获取OssClient单例对象
func GetClientInstance() (*oss.Client, error) {
	var err error
	if clientSingleton == nil {
		lock.Lock()
		if clientSingleton == nil {
			clientSingleton, err = oss.New(
				Endpoint,
				AccessKeyId,
				AccessKeySecret)
		}
		lock.Unlock()
	}
	if err != nil {
		log.Println("Oss client crete error:", err)
		return nil, err
	}
	return clientSingleton, nil
}

//GetDefaultBucketInstance 用于获取默认Bucket单例
func GetDefaultBucketInstance() (*oss.Bucket, error) {
	var err error
	client, _ := GetClientInstance()
	if defaultBucketSingleton == nil {
		lock.Lock()
		if defaultBucketSingleton == nil {
			defaultBucketSingleton, err = client.Bucket(BucketName)
		}
		lock.Unlock()
	}
	if err != nil {
		log.Println("Default bucket init error:", err)
		return nil, err
	}
	return defaultBucketSingleton, nil
}

//GetBucketInstance 获取指定的Bucket，非单例
func GetBucketInstance(bucketName string) (*oss.Bucket, error) {
	client, _ := GetClientInstance()
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Println("Bucket init error:", err)
		return nil, err
	}
	return bucket, nil
}
