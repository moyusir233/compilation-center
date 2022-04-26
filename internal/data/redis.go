package data

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"gitee.com/moyusir/compilation-center/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"io"
	"time"
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

func (r *RedisRepo) SaveExe(key string, reader io.ReadCloser, expire time.Duration) error {
	defer reader.Close()
	// 将二进制文件经过gzip压缩后再保存到redis中
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	if err != nil {
		return errors.Newf(
			500, "Save_Exe_Error", "将可执行文件进行gzip压缩时发生了错误:%s", err)
	}

	_, err = io.Copy(gzipWriter, reader)
	if err != nil {
		return errors.Newf(
			500, "Save_Exe_Error", "将可执行文件进行gzip压缩时发生了错误:%s", err)
	}

	err = gzipWriter.Close()
	if err != nil {
		return errors.Newf(
			500, "Save_Exe_Error", "将可执行文件进行gzip压缩时发生了错误:%s", err)
	}

	err = r.client.SetEX(context.Background(), key, string(buffer.Bytes()), expire).Err()
	if err != nil {
		return errors.Newf(
			500, "Save_Exe_Error", "将可执行文件进行gzip压缩时发生了错误:%s", err)
	}

	return nil
}

func (r *RedisRepo) GetExe(key string) ([]byte, error) {
	result, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, errors.Newf(
			500, "Save_Exe_Error", "查询可执行文件的缓存时发生了错误:%s", err)
	}

	// 解压
	reader, err := gzip.NewReader(bytes.NewReader([]byte(result)))
	if err != nil {
		return nil, errors.Newf(
			500, "Save_Exe_Error", "将可执行文件进行gzip解压时发生了错误:%s", err)
	}

	exe, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Newf(
			500, "Save_Exe_Error", "将可执行文件进行gzip解压时发生了错误:%s", err)
	}

	return exe, nil
}
