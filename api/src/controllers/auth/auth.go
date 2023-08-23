package auth

import (
	"binder_api/configuration"
	"binder_api/db"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthService struct {
	Configuration *configuration.AppConfiguration
	logger        *zap.Logger
}

func ProvideAuthService(config *configuration.AppConfiguration, logger *zap.Logger) *AuthService {
	return &AuthService{Configuration: config, logger: logger}
}

func (service AuthService) GenerateToken(user db.UserDTO) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(user.Id, 10),
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(3 * time.Hour)},
		Issuer:    "Binder-API",
	})

	signedToken, err := token.SignedString([]byte(service.Configuration.JwtSecretKey))
	if err != nil {
		service.logger.Error("GenerateToken() error", zap.Error(err))
	}

	return signedToken
}

func (service AuthService) validateToken(signedToken string) (int64, error) {
	parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(service.Configuration.JwtSecretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if parsedToken.Valid == false {
		return 0, fmt.Errorf("Token invalid")
	}

	subject, _ := parsedToken.Claims.GetSubject()
	id, _ := strconv.ParseInt(subject, 10, 64)

	return id, nil
}

func (service AuthService) AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	jwt := strings.Split(authHeader, " ")
	if len(jwt) != 2 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId, err := service.validateToken(jwt[1])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("userId", userId)
	c.Next()
}
