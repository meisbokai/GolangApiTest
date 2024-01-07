package jwt

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userId string, username string, isAdmin bool, email string, password string) (t string, err error)
	ParseToken(tokenString string) (claims JwtCustomClaim, err error)
}

type JwtCustomClaim struct {
	UserID   string
	Username string
	IsAdmin  bool
	Email    string
	Password string
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
	expired   int
}

func NewJWTService(secretKey, issuer string, expired int) JWTService {
	return &jwtService{
		issuer:    issuer,
		secretKey: secretKey,
		expired:   expired,
	}
}

func (j *jwtService) GenerateToken(userID string, username string, isAdmin bool, email string, password string) (t string, err error) {
	claims := &JwtCustomClaim{
		userID,
		username,
		isAdmin,
		email,
		password,
		jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * time.Duration(j.expired))},
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString([]byte(j.secretKey))
	return
}

func (j *jwtService) ParseToken(tokenString string) (claims JwtCustomClaim, err error) {
	if token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	}); err != nil || !token.Valid {
		// return JwtCustomClaim{}, errors.New("token cannot be parsed")
		if err.Error() == "token has invalid claims: token is expired" {
			parts := strings.Split(tokenString, ".")
			// header, _ := base64.StdEncoding.DecodeString(parts[0])
			payload, _ := base64.StdEncoding.DecodeString(parts[1])
			payload = append(payload, '}')

			var payloadJson JwtCustomClaim

			json.Unmarshal(payload, &payloadJson)

			// logger.ErrorF(fmt.Sprintf("Header: %s [end]", header), logrus.Fields{constants.LoggerCategory: constants.LoggerCategory})
			// logger.ErrorF(fmt.Sprintf("Payload: %s [end]", payload), logrus.Fields{constants.LoggerCategory: constants.LoggerCategory})
			// logger.ErrorF(fmt.Sprintf("Payload username: %s [end]", payloadJson.Username), logrus.Fields{constants.LoggerCategory: constants.LoggerCategory})
			// logger.ErrorF(fmt.Sprintf("Payload isAdmin: %t [end]", payloadJson.IsAdmin), logrus.Fields{constants.LoggerCategory: constants.LoggerCategory})

			return payloadJson, nil
		}

		return JwtCustomClaim{}, err

	}

	return
}
