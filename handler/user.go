package handler

import (
	"context"
	proto "filestore/service/account/proto"
	dbcli "filestore/service/dbproxy/client"
	"filestore/service/dbproxy/orm"
	"filestore/util"
	"filestore/util/jwt"
	"github.com/mitchellh/mapstructure"
)
// User : 用于实现UserServiceHandler接口的对象

type User struct{}

// Signup : 处理用户注册请求
func (u *User) Signup(ctx context.Context, req *proto.ReqSignup, res *proto.RespSignup) error {
	username := req.Username
	passwd := req.Password
	// 从文件表查找记录
	// 参数简单校验
	if len(username) < 3 || len(passwd) < 5 {
		res.Code = -1
		res.Message = "注册参数无效"
		return nil
	}
	// 对密码进行加盐及取Sha1值加密
	salt:=util.GetSalt()
	encPasswd :=util.MD5(passwd +salt)
	// 将用户信息注册到用户表中
	dbResp,err := dbcli.UserSignup(username, encPasswd,salt)
	//filehash:=
	//dbResp, err := dbcli.GetFileMeta(filehash)
	if err == nil && dbResp.Suc {
		res.Code = 0
		res.Message = "注册成功"
	} else {
		res.Code = -1
		res.Message = "注册失败"
	}
	return nil
}
// Signin : 处理用户登录请求
func (u *User) Login(ctx context.Context, req *proto.ReqLogin, res *proto.RespLogin) error {
	username := req.Username
	password := req.Password
	// 1. 校验用户名及密码
	dbResp, err := dbcli.UserLogin(username, password)
	if err != nil || dbResp.Code!=0 {
		res.Code = -1
		return nil
	}
	user:=orm.User{}
	mapstructure.Decode(dbResp.Data,&user)
	// 生成访问凭证(token)
	token,err:= jwt.GetToken(user.UserName,user.Id)
	if err != nil {
		res.Code = -1
		return nil
	}
	// 3. 登录成功, 返回token
	res.Code = 0
	res.Token = token
	return nil
}
// UserInfo : 获取用户信息
func (u *User) UserInfo(ctx context.Context, req *proto.ReqUserInfo, res *proto.RespUserInfo) error {
	username := req.Username
	dbResp, err := dbcli.GetUserInfo(username)
	if err != nil  {
		res.Code = -1
		res.Message = "服务错误"
		return nil
	}
	//SQL执行错误
	if !dbResp.Suc {
		res.Code = -2
		res.Message = "用户不存在"
		return nil
	}
	user:=orm.User{}
	mapstructure.Decode(dbResp.Data,&user)
	// 3. 组装并且响应用户数据
	res.Code = 0
	res.Id=int32(user.Id)
	res.Username = user.UserName
	return nil
}