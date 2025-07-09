package util

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

// SendMail 发送验证码邮件
func SendMail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "greecenew@163.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "openlist-strm 验证码通知")
	body := fmt.Sprintf("【openlist-strm】您的验证码是：%s，有效期10分钟。如非本人操作请忽略。", code)
	m.SetBody("text/plain", body)

	log.Printf("[邮件发送] to: %s, code: %s\n", to, code)

	// ⚠️ 请将下面的密码替换为你的163邮箱授权码
	d := gomail.NewDialer("smtp.163.com", 465, "greecenew@163.com", "MEmGa38rUS6k8VEU")
	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("[邮件发送失败] to: %s, err: %v\n", to, err)
	} else {
		log.Printf("[邮件发送成功] to: %s\n", to)
	}
	return err
}
