package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 用法提示
func usage() {
	fmt.Fprintf(os.Stderr, `用法: gomail [-s 主题] [-t 正文] [-f 文件名] [-u 邮箱地址]

选项:
`)
	flag.PrintDefaults()
}

// 全局参数
var (
	subject string
	text    string
	files   string
	address string
	recv    bool
	help    bool
)

func init() {
	// 命令行参数定义
	flag.StringVar(&subject, "s", "", "邮件主题，默认为空")
	flag.StringVar(&text, "t", "", "邮件文本内容，默认为空")
	flag.StringVar(&files, "f", "", "邮件附件，多个文件使用英文逗号分割")
	flag.StringVar(&address, "u", "", "收件人邮箱，多个邮箱使用英文逗号分割")
	flag.BoolVar(&recv, "r", false, "收取未读邮件，自动下载附件")
	flag.BoolVar(&help, "h", false, "使用帮助")
	flag.Usage = usage
}

type Euser struct {
	PopServer  string
	PopPort    int
	SmtpServer string
	SmtpPort   int
	UserName   string
	Passwd     string
}

// 获取账户
func getUser() (u Euser, err error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	f := filepath.Join(dir, ".mailuser")
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(b), &u)
	return u, err
}

// 保存账户
func saveUser(euser Euser) (err error) {
	js, err := json.Marshal(euser)
	if err == nil {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println(err)
		}
		f := filepath.Join(dir, ".mailuser")
		err = ioutil.WriteFile(f, js, 0644)
	}
	return err
}

// 发送邮件
func sendmail(euser Euser) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", euser.UserName)
	u := strings.Split(address, ",")
	m.SetHeader("To", u...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", text)
	f := strings.Split(files, ",")
	for _, i := range f {
		m.Attach(i)
	}
	d := gomail.NewDialer(euser.SmtpServer, euser.SmtpPort, euser.UserName, euser.Passwd)

	// Send the email to Bob, Cora and Dan.
	err = d.DialAndSend(m)

	return err
}

// 添加用户
func addUser() Euser {
	var (
		pserver string
		pport   int
		sserver string
		sport   int
		uname   string
		passwd  string
	)

	fmt.Println("Please enter the email pop3 server: ")
	_, _ = fmt.Scan(&pserver)
	fmt.Println("Please enter the email pop3 port: ")
	_, _ = fmt.Scan(&pport)
	fmt.Println("Please enter the email smtp server: ")
	_, _ = fmt.Scan(&sserver)
	fmt.Println("Please enter the email smtp port: ")
	_, _ = fmt.Scan(&sport)
	fmt.Println("Please enter the email username: ")
	_, _ = fmt.Scan(&uname)
	fmt.Println("Please enter the email password: ")
	_, _ = fmt.Scan(&passwd)
	return Euser{pserver, pport, sserver, sport, uname, passwd}

}

func main() {
	user, err := getUser()
	if err != nil {
		user = addUser()
		err = saveUser(user)
		if err != nil {
			return
		}
	}
	fmt.Println(user)

	flag.Parse()
	if len(os.Args) == 1 {
		flag.Usage()
	} else if help {
		flag.Usage()
	} else if recv {
		fmt.Println("recv...")
	} else if address == "" {
		flag.Usage()
	} else {
		fmt.Println("send...")
		err = sendmail(user)
		if err == nil {
			fmt.Println("success!")
		} else {
			fmt.Println(err)
		}
	}
}
