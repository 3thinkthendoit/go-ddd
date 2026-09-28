package valueobject

import "fmt"

// UserType 用户类型
type UserType int

const (
	UserTypeToB UserType = 1 // TO_B
	UserTypeToC UserType = 2 // TO_C
)

var userTypeDesc = map[UserType]string{
	UserTypeToB: "TO_B",
	UserTypeToC: "TO_C",
}

// Code 返回用户类型编码
func (t UserType) Code() int { return int(t) }

// Desc 返回用户类型描述
func (t UserType) Desc() string { return userTypeDesc[t] }

func (t UserType) String() string {
	if desc, ok := userTypeDesc[t]; ok {
		return desc
	}
	return fmt.Sprintf("UserType(%d)", int(t))
}

// UserTypeOf 按编码查找用户类型
func UserTypeOf(code int) (UserType, bool) {
	t := UserType(code)
	_, ok := userTypeDesc[t]
	return t, ok
}

// UserInfo 下单用户信息（值对象，不可变）
type UserInfo struct {
	userId   int64
	username string
	userType UserType
}

// NewUserInfo 创建用户信息值对象
func NewUserInfo(userId int64, username string, userType UserType) UserInfo {
	return UserInfo{userId: userId, username: username, userType: userType}
}

// UserId 用户 ID
func (u UserInfo) UserId() int64 { return u.userId }

// Username 用户名称
func (u UserInfo) Username() string { return u.username }

// UserType 用户类型
func (u UserInfo) UserType() UserType { return u.userType }

func (u UserInfo) String() string {
	return fmt.Sprintf("UserInfo{userId=%d, username=%s, userType=%s}", u.userId, u.username, u.userType)
}
