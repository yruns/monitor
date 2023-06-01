package service

import (
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"monitor/conf"
	"monitor/database"
	"monitor/model"
	"monitor/pkg/response"
	"monitor/pkg/utils"
	"strconv"
	"time"
)

type WarningEmailService struct {
	Detail string `json:"detail" form:"detail"`
}

type EmailService struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	// 1:绑定邮箱  2:解绑邮箱	 3:更改密码
	OperationType int64 `json:"operation_type" form:"operation_type"`
}

var EMAILNOTICE = [3]string{
	"您正在绑定邮箱",
	"您正在解绑邮箱",
	"您正在更改密码",
}

func (s *EmailService) SendEmailCode(userId int64) *response.Response {

	verifyCode := utils.GenerateEmailVerifyCode()

	// 发送邮件
	notice := fmt.Sprintf("%s，验证码为：<h2>%s</h2>若非本人操作请忽略该信息。", EMAILNOTICE[s.OperationType], verifyCode)
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(conf.SmtpUser, "Monitor"))
	m.SetHeader("To", s.Email)
	m.SetHeader("Subject", "Monitor邮箱验证")
	m.SetBody("text/html", notice)
	d := gomail.NewDialer(conf.SmtpHost, 465, conf.SmtpUser, conf.SmtpPass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err.Error())
		return response.FailWithMessage("邮件发送失败，请重试")
	}

	// 写入redis
	database.Redis.SetEx(context.Background(), "user:verify:"+strconv.Itoa(int(userId)), verifyCode, time.Minute*3)

	return response.Ok()
}

func (s *WarningEmailService) SendWarningEmail(userId int64) *response.Response {

	// 从用户表中查询用户的Email
	var user model.User
	err := database.Mysql.Table("user").Where("id = ?", userId).Find(&user).Error
	if err != nil || user.Email == "" {
		return response.FailWithMessage("请先绑定邮箱")
	}

	toEmail := user.Email

	// 发送邮件
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(conf.SmtpUser, "Monitor"))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Monitor警告邮件")
	m.SetBody("text/html", s.Detail)
	d := gomail.NewDialer(conf.SmtpHost, 465, conf.SmtpUser, conf.SmtpPass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err.Error())
		return response.FailWithMessage("邮件发送失败，请重试")
	}

	return response.Ok()
}
