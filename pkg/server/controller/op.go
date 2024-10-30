package controller

import (
	"fmt"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
)

type CustomClaims struct {
	Sub   string   `json:"sub"`
	Ports []string `json:"ports"`
	jwt.RegisteredClaims
}

type OpController struct {
	secret []byte
}

func NewOpController(s []byte) *OpController {
	return &OpController{
		secret: s,
	}
}

func (c *OpController) Register(engine *gin.Engine) {
	engine.POST("/handler", MakeGinHandlerFunc(c.HandleLogin))
}

func verifyJWT(tokenString string, secretKey []byte) (*CustomClaims, error) {
	// Parse the token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Type-assert the token claims to our custom type
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (c *OpController) HandleLogin(ctx *gin.Context) (interface{}, error) {
	var r plugin.Request
	var content plugin.LoginContent
	r.Content = &content
	if err := ctx.BindJSON(&r); err != nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}

	var claims, err = verifyJWT(content.User, c.secret)

	var res plugin.Response
	if err != nil && claims != nil {
		// check port claim here
		var proxyRequest plugin.Request
		var proxyContent plugin.NewProxyContent

		proxyRequest.Content = &proxyContent

		fmt.Println("-------------Plugin: Allowed Ports--------------------")
		fmt.Printf("ProxyName: %s\tProxyType%s\t", proxyContent.ProxyName, proxyContent.ProxyType)
		if strings.ToLower(proxyContent.ProxyType) == "tcp" || strings.ToLower(proxyContent.ProxyType) == "udp" {
			fmt.Printf("RemotePort: %d\r\n", proxyContent.RemotePort)
		} else if strings.HasPrefix(proxyContent.ProxyType, "http") {
			fmt.Printf("CustomDomains%s\r\n", proxyContent.CustomDomains)
		} else {
			fmt.Println("Won't do validation for this type")
			res.Unchange = true
			return res, nil
		}

		subdomain := proxyContent.SubDomain
		remoteport := strconv.Itoa(proxyContent.RemotePort)

		if subdomain == "" && remoteport == "0" && len(proxyContent.CustomDomains) == 0 {
			res.Reject = true
			res.RejectReason = "Rejected due to misconfiguration of the client"
		}

		find := false
		for _, portAllowed := range claims.Ports {
			if portAllowed == remoteport || portAllowed == subdomain {
				find = true
			}

			if contains(proxyContent.CustomDomains, portAllowed) {
				find = true
			}
		}

		if !find {
			res.Reject = true
			res.RejectReason = "Client is not allowed => Port or subdomain false"
		}

		if !res.Reject {
			res.Unchange = true
		}

	} else {
		res.Reject = true
		res.RejectReason = "invalid token"
	}

	return res, nil
}
