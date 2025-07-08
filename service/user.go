package service

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"CollabDoc-go/utils"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
}

func (userService *UserService) Logout(c *gin.Context) {
	uuid := utils.GetUUID(c)
	jwtStr := utils.GetRefreshToken(c)
	utils.ClearRefreshToken(c)
	global.Redis.Del(context.Background(), uuid.String())
	_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{Jwt: jwtStr})
}
func (userService *UserService) Register(user database.User) (database.User, error) {
	if !errors.Is(global.DB.Where("email = ?", user.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered, please check the information you filled in, or retrieve your password")
	}

	user.Password = utils.BcryptHash(user.Password)
	user.UUID = uuid.Must(uuid.NewV4())
	user.Avatar = "https://avatars.githubusercontent.com/u/132669442"
	user.Roles = database.JSONStringList{"common"}
	user.Permissions = database.JSONStringList{"*:*:*"}
	t := time.Date(2030, time.October, 30, 0, 0, 0, 0, time.Local)
	pt := &t
	user.Expires = pt
	if err := global.DB.Create(&user).Error; err != nil {
		return database.User{}, err
	}

	return user, nil
}
func (userService *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ?", u.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("incorrect email or password")
		}
		return user, nil
	}
	return database.User{}, err
}
func (userService *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}
func (userService *UserService) GetUser(user database.User) (database.User, error) {
	db := global.DB
	if err := db.Where("username=?", user.Username).First(&user).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}
