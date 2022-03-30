package data

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"gitee.com/moyusir/compilation-center/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

const (
	// CLIENT_CODE_KEY 用户客户端代码hash的key
	CLIENT_CODE_KEY = "client_code"
)

// RedisRepo redis数据库操作对象，可以理解为dao
type RedisRepo struct {
	client *Data
}

// NewRedisRepo 实例化redis数据库操作对象
func NewRedisRepo(data *Data, logger log.Logger) biz.ClientCodeRepo {
	return &RedisRepo{
		client: data,
	}
}

// SaveClientCode 以zip文件的二进制数据的十六进制字符串形式保存客户端代码
func (r *RedisRepo) SaveClientCode(key string, files map[string]*bytes.Reader) error {
	// 将文件压缩为zip二进制数据形式
	result := bytes.NewBuffer(make([]byte, 0, 1024))
	zipWriter := zip.NewWriter(result)
	for k, v := range files {
		writer, err := zipWriter.Create(k)
		if err != nil {
			return err
		}
		_, err = v.WriteTo(writer)
		if err != nil {
			return err
		}
	}
	err := zipWriter.Close()
	if err != nil {
		return errors.Newf(
			500, "Repo_Error", "将客户端代码压缩成zip文件时发生了错误:%v", err)
	}

	// 以十六进制字符串的形式保存二进制数据
	value := fmt.Sprintf("%x", result.Bytes())
	err = r.client.HSetNX(context.Background(), CLIENT_CODE_KEY, key, value).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errors.Newf(
			500, "Repo_Error", "将客户端代码写入redis时发生了错误:%v", err)
	}

	return nil
}

// IsValid 通过判断数据库中是否存在客户端的代码信息，判断账号是否有效
func (r *RedisRepo) IsValid(username string) bool {
	exist, err := r.client.HExists(context.Background(), CLIENT_CODE_KEY, username).Result()
	if err != nil {
		return false
	}

	return exist
}
