package mysql

import (
	"errors"
	"fmt"
	"web_app/models"

	"go.uber.org/zap"
)

func CheckUserExist(username string) (err error) {
	sqlStr := "select count(id) from user where name = ?"
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("user already exist")
	}
	return
}

func SaveUser(user *models.User) (err error) {
	sqlStr := "insert into user(id,name,password,status,mobile) values(?,?,?,?,?)"
	if _, err = db.Exec(sqlStr, user.Id, user.Name, user.Password, user.Status, user.Mobile); err != nil {
		zap.L().Error("save user failed", zap.Error(err))
		return err
	}
	return
}

func CheckUser(user *models.User) (error, *models.LoginReply) {
	r := new(models.LoginReply)
	sqlStr := "select id, role from user where mobile = ? and password = ? "
	if err := db.Get(r, sqlStr, user.Mobile, user.Password); err != nil {
		zap.L().Error("query mysql failed", zap.Error(err))
		return err, nil
	}
	if r == nil {
		return errors.New("mobile or pwd error"), nil
	}
	return nil, r
}

func GetUserList(pageNum, pageSize int64) (error, []*models.UserListReply) {
	res := []*models.UserListReply{}
	sqlStr := "select id, name, add_time, status, mobile from user limit ?, ? "
	rows, err := db.Query(sqlStr, (pageNum-1)*pageSize, pageSize)
	if err != nil {
		return errors.New("query db failed"), nil
	}
	for rows.Next() {
		var u models.UserListReply
		if e := rows.Scan(&u.Id, &u.Name, &u.AddTime, &u.Status, &u.Mobile); e != nil {
			fmt.Printf("scan failed, err:%v\n", e)
			continue
		}
		res = append(res, &u)
	}
	defer rows.Close()
	return nil, res
}
