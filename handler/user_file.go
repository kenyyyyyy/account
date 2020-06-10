package handler

import (
	"context"
	"encoding/json"
	proto "filestore/service/account/proto"
	dbcli "filestore/service/dbproxy/client"
	"filestore/service/dbproxy/orm"
	"github.com/mitchellh/mapstructure"
)


// UserFiles : 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, res *proto.RespUserFile) error {
	dbResp, err := dbcli.QueryUserFileMetas(req.Username, int(req.Limit),int(req.Offset))
	if err != nil || !dbResp.Suc {
		res.Code = 0
		return err
	}
	files := []orm.UserFile{}
	mapstructure.Decode(dbResp.Data, &files)
	data, err := json.Marshal(files)
	if err != nil {
		res.Code = -1
		return nil
	}
	res.FileData = data
	return nil
}
// UserFiles : 用户文件重命名
func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.RespUserFileRename) error {
	dbResp, err := dbcli.FileRename(req.Username,req.Filehash,req.NewFileName)
	if err != nil || !dbResp.Suc {
		res.Code = 0
		return err
	}
	if err != nil {
		res.Code = -1
		return nil
	}
	return nil
}
// UserFileDelete : 删除用户文件
func (u *User) UserFileDelete(ctx context.Context, req *proto.ReqUserFileDelete, res *proto.RespUserFileDelete) error {
	dbResp, err := dbcli.FileDelete(req.Username,req.Filehash)
	if err != nil || !dbResp.Suc {
		res.Code = 0
		return err
	}
	if err != nil {
		res.Code = -1
		return nil
	}
	return nil
}