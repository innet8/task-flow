package util

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/tealeg/xlsx"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GetDatasFromXlsx 从xlsx文件获取数据
func GetDatasFromXlsx(file *os.File) ([][]string, error) {
	if file == nil {
		return nil, fmt.Errorf("file值为nil")
	}
	xlFile, err := xlsx.OpenFile(file.Name())
	if err != nil {
		return nil, err
	}
	datas, err := xlFile.ToSlice()
	if err != nil {
		return nil, fmt.Errorf("文件数据转换成数组失败:%s", err.Error())
	}
	return datas[0], nil
}

// GetDatasFromCSV 从csv文件获取数据
func GetDatasFromCSV(file *os.File) ([][]string, error) {
	if file == nil {
		return nil, fmt.Errorf("file值为nil")
	}
	csvfile, err := os.Open(file.Name())
	if err != nil {
		return nil, fmt.Errorf("打开csv文件:%s", err.Error())
	}
	defer csvfile.Close()
	r := transform.NewReader(bufio.NewReader(csvfile), simplifiedchinese.GBK.NewDecoder())
	reader := csv.NewReader(r)
	if reader == nil {
		return nil, fmt.Errorf("解析csv文件:reader为空")
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("从csv读取文件:%s", err.Error())
	}
	return records, nil
}

// UploadFile 上传文件
func UploadFile(file multipart.File, header *multipart.FileHeader, uploadDir string) (string, error) {
	// 检查文件大小
	if header.Size > 5*1024*1024 {
		return "", errors.New("file size exceeds the limit")
	}
	// 创建上传目录
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}
	// 生成文件名
	ext := filepath.Ext(header.Filename)
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	filename := filepath.Join(uploadDir, hex.EncodeToString(hash.Sum(nil))+ext)
	// 创建文件
	out, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer out.Close()
	// 复制文件内容
	if _, err = file.Seek(0, 0); err != nil {
		return "", err
	}
	if _, err = io.Copy(out, file); err != nil {
		return "", err
	}
	return filename, nil
}
