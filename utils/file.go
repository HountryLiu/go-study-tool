package utils

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func WriteFile(path, content string) (err error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(content)
	return
}

func MkdirAll(path string) (err error) {
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return
	}
	return
}

func MvDir(source string, target string) (err error) {
	err = os.Rename(source, target)
	if err != nil {
		return
	}
	return
}

func RmDir(source string) (err error) {
	err = os.RemoveAll(source)
	if err != nil {
		return
	}
	return
}

func CopyDir(from, to string) error {
	var err error

	f, err := os.Stat(from)
	if err != nil {
		return err
	}

	fn := func(fromFile string) error {
		//复制文件的路径
		rel, err := filepath.Rel(from, fromFile)
		if err != nil {
			return err
		}
		toFile := filepath.Join(to, rel)

		//创建复制文件目录
		if err = os.MkdirAll(filepath.Dir(toFile), 0777); err != nil {
			return err
		}

		//读取源文件
		file, err := os.Open(fromFile)
		if err != nil {
			return err
		}

		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建复制文件用于保存
		out, err := os.Create(toFile)
		if err != nil {
			return err
		}

		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, err = io.Copy(out, bufReader)
		return err
	}

	//转绝对路径
	pwd, _ := os.Getwd()
	if !filepath.IsAbs(from) {
		from = filepath.Join(pwd, from)
	}
	if !filepath.IsAbs(to) {
		to = filepath.Join(pwd, to)
	}

	//复制
	if f.IsDir() {
		return filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				return fn(path)
			} else {
				if err = os.MkdirAll(path, 0777); err != nil {
					return err
				}
			}
			return err
		})
	} else {
		return fn(from)
	}
}

func ReadFile(path string) (body string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	line, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	body = string(line)
	return
}

func GetFileList(dirPath string) []string {
	var ls []string

	filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ls = append(ls, path)
		return nil
	})
	sort.Slice(ls, func(i, j int) bool {
		return ls[i] < ls[j]
	})

	return ls
}

func ZipFiles(filename string, files []string) error {
	//创建输出文件目录
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	//创建空的zip档案，可以理解为打开zip文件，准备写入
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	//打开要压缩的文件
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	//获取文件的描述
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	//FileInfoHeader返回一个根据fi填写了部分字段的Header，可以理解成是将fileinfo转换成zip格式的文件信息
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filename
	/*
	   预定义压缩算法。
	   archive/zip包中预定义的有两种压缩方式。一个是仅把文件写入到zip中。不做压缩。一种是压缩文件然后写入到zip中。默认的Store模式。就是只保存不压缩的模式。
	   Store   unit16 = 0  //仅存储文件
	   Deflate unit16 = 8  //压缩文件
	*/
	header.Method = zip.Deflate
	//创建压缩包头部信息
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	//将源复制到目标，将fileToZip 写入writer   是按默认的缓冲区32k循环操作的，不会将内容一次性全写入内存中,这样就能解决大文件的问题
	_, err = io.Copy(writer, fileToZip)
	return err
}

func Base64ToImg(file_path string, base64_code string) (err error) {
	i := strings.Index(base64_code, ",")
	dist, err := base64.StdEncoding.DecodeString(base64_code[i+1:])
	if err != nil {
		return
	}
	f, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(dist)
	if err != nil {
		return
	}
	return
}
