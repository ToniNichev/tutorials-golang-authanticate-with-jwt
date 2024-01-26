package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type requestBody struct {
	OperationName *string     `json:"operationName"`
	Query         *string     `json:"query"`
	Variables     interface{} `json:"variables"`
}

func validateSession(c *gin.Context) {
	if c.Request.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		body := requestBody{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return
		}

		// extract the token from the headers
		tokenStr := c.Request.Header.Get("X-Auth-Token")

		product := body.Variables.(map[string]interface{})["product"]

		var payload string
		var err error
		if product == "web" {
			payload, err = retreiveTokenWithSymmetrikKey(c, tokenStr)
		} else {
			payload, err = retreiveTokenWithAsymmetrikKey(c, tokenStr)
		}

		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Session token signature can't be confirmed!"})
		}

		if payload == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		c.Next()
	}

}

func retreiveTokenWithSymmetrikKey(c *gin.Context, tokenStr string) (string, error) {
	fmt.Println("retreive Token With Symmetric Key ...")

	tknStr := c.Request.Header.Get("X-Auth-Token")
	secretKey := "itsasecret123"

	token, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println("Error !")
		fmt.Println(err)
	} else {
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println("======================================")
		fmt.Println(claims)
		fmt.Println(claims["author"])
		fmt.Println(claims["data"])
		fmt.Println("======================================")
	}
	return "retreiveToken: OK 123", nil
}

func retreiveTokenWithAsymmetrikKey(c *gin.Context, tokenStr string) (string, error) {

	fmt.Println("retreive Token With Asymmetric Key ...")

	publicKeyPath := "key/public_key.pem"
	keyData, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		fmt.Println("Error reading the key")
		return "", errors.New("error reading public key")
	}

	var parsedToken *jwt.Token

	// parse token
	state, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		// ensure signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unknown signing method")
		}

		parsedToken = token

		// verify
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(keyData))
		if err != nil {
			return nil, errors.New("parsing key failed")
		}

		return key, nil
	})

	claims := state.Claims.(jwt.MapClaims)

	fmt.Println("======================================")
	fmt.Println("Header [alg]:", parsedToken.Header["alg"])
	fmt.Println("Header [expiresIn]:", parsedToken.Header["expiresIn"])
	fmt.Println("Claims [author]:", claims["author"])
	fmt.Println("Claims [data]:", claims["data"])
	fmt.Println("======================================")
	if !state.Valid {
		return "", errors.New("verification failed")
	}

	if err != nil {
		errors.New("unknown signing error")
	}

	return "retreiveToken: OK 123", nil
}

func returnData(c *gin.Context) {
	fmt.Println("Returning data ...")
	c.String(200, "TEST 123")
}
