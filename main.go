package main

import (
	"encoding/json"
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
`, os.Args[0])
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
	flag.StringVar(&files, "f", "", "邮件附件，多个文件使用英文逗号分割")
	flag.StringVar(&address, "a", "", "收件人邮箱，多个邮箱使用英文逗号分割")
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
	b, err := ioutil.ReadFile(".mailuser")
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
		err = ioutil.WriteFile(".mailuser", js, 0644)
	}
	return err
}

func main() {

	flag.Parse()
	flag.Usage()

	fmt.Println(subject)
	fmt.Println(text)
	fmt.Println(files)
	fmt.Println(address)
	user, err := getUser()
	if err != nil {
		fmt.Println("-----------------")
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

		user = Euser{pserver, pport, sserver, sport, uname, passwd}
		err = saveUser(user)
		if err != nil {
			return
		}
	}
	fmt.Println(user)

}
