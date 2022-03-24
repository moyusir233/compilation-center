package data

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"gitee.com/moyusir/compilation-center/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
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
		return err
	}

	// 以十六进制字符串的形式保存二进制数据
	value := fmt.Sprintf("%x", result.Bytes())
	return r.client.HSetNX(context.Background(), CLIENT_CODE_KEY, key, value).Err()
}

// IsValid 通过判断数据库中是否存在客户端的代码信息，判断账号是否有效
func (r *RedisRepo) IsValid(username string) bool {
	exist, err := r.client.HExists(context.Background(), CLIENT_CODE_KEY, username).Result()
	if err != nil {
		return false
	}

	return exist
}