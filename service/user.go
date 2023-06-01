package service

import (
	"github.com/jinzhu/copier"
	"monitor/database"
	"monitor/model"
	"monitor/pkg/response"
	"monitor/pkg/utils"
)

type UserService struct {
	NickName string `json:"nick_name" form:"nick_name"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Avatar   string `json:"avatar" form:"avatar"`
	Key      string `json:"key" form:"key"` // 前端验证
}

func (s *UserService) Register() *response.Response {
	if s.Key == "" || len(s.Key) != 16 {
		return response.FailWithMessage("密钥错误")
	}

	// 对称加密
	utils.Encrypt.SetKey(s.Key)

	// 注册逻辑
	var user model.User
	// 1.判断该用户名是否已存在
	rowsAffected := database.Mysql.Table("user").Where("username = ?", s.Username).Find(&user).RowsAffected
	if rowsAffected > 0 {
		return response.FailWithMessage("用户名已存在")
	}

	user = model.User{
		Username: s.Username,
		NickName: s.NickName,
		Status:   model.Active,
		Avatar:   "imgs/avatar.jpg",
		Money:    utils.Encrypt.AesEncoding("10000"), // 对初始金额的加密
	}

	err := user.SetPassword(s.Password)
	if err != nil {
		return response.FailWithMessage("密码加密失败")
	}

	// 2.向数据库写入新用户
	err = database.Mysql.Table("user").Create(&user).Error
	if err != nil {
		return response.FailWithMessage(err.Error())
	}

	return response.Ok()
}

func (s *UserService) Login() (*response.Response, *model.User) {
	var user model.User

	// 验证账号密码
	affected := database.Mysql.Table("user").Where("username = ?", s.Username).
		Find(&user).RowsAffected
	if affected == 0 {
		return response.FailWithMessage("账号不存在"), nil
	}

	// 从数据库中获取hash后的密码
	hashedPasswordFromDB := user.PasswordDigest
	// 验证密码
	flag := user.CheckPassword(s.Password, hashedPasswordFromDB)
	if !flag {
		return response.FailWithMessage("密码错误"), nil
	}

	return response.OkWithData(user), &user
}

func (s *UserService) Update(userId int64) *response.Response {
	var user model.User
	_ = copier.Copy(&user, s)

	affected := database.Mysql.Table("user").Where("id = ?", userId).Updates(user).RowsAffected // gorm默认只会更新非空字段，0/false/""都将被忽略
	if affected == 0 {
		return response.FailWithMessage("用户信息更新失败")
	}

	// 将该用户查出用于返回
	database.Mysql.Table("user").Where("id = ?", userId).Find(&user)
	return response.OkWithData(user)
}
