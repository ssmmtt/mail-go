package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

// 用法提示
func usage() {
	fmt.Fprintf(os.Stderr, `用法: %s [-s 主题] [-t 正文] [-f 文件名] [-a 邮箱地址]

选项:
  无
	不加任何参数收取未读邮件
` , os.Args[0])
	flag.PrintDefaults()
}

// 全局参数
var (
	subject string
	text    string
	files   string
	address string
)

func init() {
	// 命令行参数定义
	flag.StringVar(&subject, "s", "", "邮件主题，默认为空")
	flag.StringVar(&text, "t", "", "邮件文本内容，默认为空")
	flag.StringVar(&files, "f", "", "邮件附件，多个文件使用英文逗号分割，不要添加空格")
	flag.StringVar(&address, "a", "", "收件人邮箱，多个邮箱使用英文逗号分割，不要加空格")
	flag.Usage = usage
}

// 获取账户
func getUser() (u string, err error) {
	b, err := ioutil.ReadFile(".mailuser")
	if err != nil {
		fmt.Print(err)
		return u, err
	}
	fmt.Println("==============")
	fmt.Println(b)
	u = string(b)
	fmt.Println(u)
	return u, err
}

// 保存账户
func saveUser(euser Euser) (err error) {
	fmt.Println(euser)
	d1, err := Encode(euser)
	if err != nil{
		fmt.Println(err)
	}else {
		err = ioutil.WriteFile(".mailuser", d1, 0644)
		if err != nil {
			panic(err)
		}
	}

	return err
}



func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}


type Euser struct {
	popServer string
	popPort int
	smtpServer string
	smtpPort int
	userName string
	passwd string
}


func main() {

	//flag.Parse()
	//
	//fmt.Println(subject)
	//fmt.Println(text)
	//fmt.Println(files)
	//fmt.Println(address)

	u := Euser{"smtp.qq.com", 993, "imap.qq.com", 587, "xxxxx@qq.com", "xxxxxx"}
	err := saveUser(u)
	fmt.Println(err)
	c, err := getUser()
	fmt.Println(c, err)



}
