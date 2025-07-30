package jwt

import (
	"fmt"
	"go-42/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type JWT struct {
	PrivateKey []byte
	PublicKey  []byte
	Log        *zap.Logger
}

type customClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	IP    string `json:"ip"`
	jwt.StandardClaims
}

func NewJWT(privateKey []byte, publicKey []byte, log *zap.Logger) JWT {
	return JWT{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Log:        log,
	}
}

func (j *JWT) CreateToken(email, ip string, ID string) (string, error) {
	//prepare private key parsing
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.PrivateKey)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &customClaims{
		ID:             ID,
		Email:          email,
		IP:             ip,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

// JWT for API
func (j *JWT) AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware jwt")

		key, err := jwt.ParseRSAPublicKeyFromPEM(j.PublicKey)
		if err != nil {
			return
		}

		claims := &customClaims{}
		tokenValue := c.GetHeader("token")

		if len(tokenValue) == 0 {
			response.ResponseBadRequest(c, http.StatusUnauthorized, "token invalid")
			return
		}

		tkn, err := jwt.ParseWithClaims(string(tokenValue), claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				response.ResponseBadRequest(c, http.StatusUnauthorized, fmt.Sprintf("unexpected method: %s", token.Header["alg"]))
				return nil, err
			}

			return key, nil
		})

		if err != nil {
			response.ResponseBadRequest(c, http.StatusUnauthorized, "fail to validate signature or session expired")
			return
		}

		if !tkn.Valid {
			response.ResponseBadRequest(c, http.StatusUnauthorized, "invalid token")
			return
		}

		c.Next()
	}
}
