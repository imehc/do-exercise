package util

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

type EmailData struct {
	EmailTitle       string // 邮件标题
	VerificationType string // 验证类型
	GreetingText     string // 问候语
	VerificationCode string // 验证码
	Minute           int    // 有效分钟数
}

type Email struct {
	TemplatePath string
	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPass     string
}

// SendEmail 发送邮件
func (e *Email) SendEmail(from, to, subject string, data EmailData) error {
	tmpl, err := e.parseTemplate()
	if err != nil {
		return err
	}
	// 渲染模板到缓冲区
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	// 设置邮件服务器
	d := gomail.NewDialer(e.SmtpHost, e.SmtpPort, e.SmtpUser, e.SmtpPass)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// parseTemplate 解析模版文件
func (e *Email) parseTemplate() (*template.Template, error) {
	// 检查文件扩展名是否为.html
	if filepath.Ext(e.TemplatePath) != ".html" {
		return nil, fmt.Errorf("only .html files are supported")
	}
	// 检查文件是否存在
	if _, err := os.Stat(e.TemplatePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template file does not exist")
	}
	// 解析单个模板文件
	tmpl, err := template.ParseFiles(e.TemplatePath)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
