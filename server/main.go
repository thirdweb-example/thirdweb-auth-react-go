package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/spruceid/siwe-go"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Payload struct {
	Domain         string `json:"domain"`
	Address        string `json:"address"`
	Statement      string `json:"statement"`
	Version        string `json:"version"`
	Nonce          string `json:"nonce"`
	IssuedAt       string `json:"issued_at"`
	ExpirationTime string `json:"expiration_time"`
	InvalidBefore  string `json:"invalid_before"`
	ChainId        string `json:"chain_id"`
	URI            string `json:"uri"`
}

type LoginRequest struct {
	Signature string  `json:"signature"`
	Payload   Payload `json:"payload"`
}

func main() {
	godotenv.Load()

	router := gin.Default()

	// CORS Configuration
	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/login", func(c *gin.Context) {
		address := c.Query("address")
		chainId := c.Query("chainId")
		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "address is required"})
			return
		}

		chainIdInt, err := strconv.Atoi(chainId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chainId"})
			return
		}
		payload := generateLoginPayload(address, chainIdInt)

		c.JSON(http.StatusOK, gin.H{"payload": map[string]interface{}{
			"domain":          payload.GetDomain(),
			"address":         address,
			"statement":       payload.GetStatement(),
			"version":         payload.GetVersion(),
			"nonce":           payload.GetNonce(),
			"issued_at":       payload.GetIssuedAt(),
			"expiration_time": payload.GetExpirationTime(),
			"invalid_before":  payload.GetNotBefore(),
			"chain_id":        chainId,
			"uri":             payload.GetURI().Host,
		}})
	})

	router.POST("/login", func(c *gin.Context) {
		var loginReq LoginRequest

		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		signature := loginReq.Signature
		payload := loginReq.Payload

		message, err := siwe.InitMessage(payload.Domain, payload.Address, payload.URI, payload.Nonce, map[string]interface{}{
			"chainId":        payload.ChainId,
			"issuedAt":       payload.IssuedAt,
			"expirationTime": payload.ExpirationTime,
			"notBefore":      payload.InvalidBefore,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message"})
			return
		}

		result := verifySignature(message, signature, "")

		if result {
			token := generateJWT(payload.Address, payload.ChainId)
			c.SetCookie("jwt", token, 3600, "/", "", false, true)
		}

		c.JSON(http.StatusOK, gin.H{"result": result})
	})

	router.GET("/isLoggedIn", func(c *gin.Context) {
		token, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"result": false})
			return
		}

		result, err := verifyJWT(token)
		if err != nil || result == nil {
			c.JSON(http.StatusOK, gin.H{"result": false})
			return
		}

		c.JSON(http.StatusOK, gin.H{"result": true})
	})

	router.POST("/logout", func(c *gin.Context) {
		c.SetCookie("jwt", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"result": true})
	})

	router.Run("localhost:8080")
}

func generateJWT(address string, chainId string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": address,
		"chainId": chainId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func verifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			panic("Unexpected signing method")
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return token, nil
}

func generateLoginPayload(address string, chainId int) *siwe.Message {
	now := time.Now()
	expirationTime := 3600 * time.Second

	message, err := siwe.InitMessage(os.Getenv("DOMAIN"), address, os.Getenv("URI"), siwe.GenerateNonce(), map[string]interface{}{
		"chainId":        chainId,
		"issuedAt":       now,
		"expirationTime": now.Add(expirationTime),
		"notBefore":      now.Add(-expirationTime),
	})

	if err != nil {
		panic(err)
	}
	return message
}

func verifySignature(payload *siwe.Message, signature string, address string) bool {
	publicKey, err := payload.Verify(signature, nil, nil, nil)
	if err != nil {
		return false
	}

	if address != publicKey.Params().Name {
		return false
	}

	return true
}
