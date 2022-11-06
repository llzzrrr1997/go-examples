package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"time"
)

func main() {
	defaultPath := "/Library/Containers/com.tencent.xinWeChat/Data/Library/Application Support/com.tencent.xinWeChat/2.0b4.0.9/26a529dbdfa6e0fbc07029799ed92861/Message/test"
	//获取删除的文件目录
	u, _ := user.Current()
	if u != nil {
		defaultPath = u.HomeDir + defaultPath
	} else {
		defaultPath = "/Users/ws" + defaultPath
	}
	deletePath := ""
	flag.StringVar(&deletePath, "p", "", "删除的目录路径")
	flag.Parse()
	if deletePath == "" {
		deletePath = defaultPath
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	//执行删除操作
	if IsDir(deletePath) {
		err := os.RemoveAll(deletePath)
		//输出删除结果
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println("delete dir success!")
	} else if IsFile(deletePath) {
		err := os.RemoveAll(deletePath)
		//输出删除结果
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println("delete file success!")
	}
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()

}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
