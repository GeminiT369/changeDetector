package middleware

import (
	"changeDetector/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"net/http"
	"strings"
	"time"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 24 * 7

var MySecret = []byte("Secret of change detector")

var AuthUser dto.User

// GenToken 生成JWT
func GenToken(username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "Change Detector",                          // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func AuthHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user dto.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Invalid Data",
		})
		return
	}
	// 校验用户名和密码是否正确
	if user.UserName == AuthUser.UserName && user.PassWord == AuthUser.PassWord {
		// 生成Token
		tokenString, _ := GenToken(user.UserName)
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
			"token":   tokenString,
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "Authorize Failed",
	})
	return
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware(authUser dto.User) func(c *gin.Context) {
	AuthUser = authUser
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("[WARN] AuthMiddleware failed, header is empty\n")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2001,
				"msg":  "Request header lacks token",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			log.Printf("[WARN] AuthMiddleware failed, format error\n")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2002,
				"msg":  "Request header format error",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])
		if err != nil {
			log.Printf("[WARN] AuthMiddleware failed, token parse error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 2003,
				"msg":  "Request header token parse error",
			})
			c.Abort()
			return
		}
		log.Printf("[INFO] AuthMiddleware success, request path: %s, user: %s\n", c.Request.URL.Path, mc.UserName)
		c.Set("username", mc.UserName)
		c.Next()
		// 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
