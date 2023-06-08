package main


import (
	"crypto/sha1"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)
func wxCallbackHandler(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echoStr := c.Query("echostr")
	fmt.Println("收到消息",echoStr)
	if signature == "" || timestamp == "" || nonce == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 实际应用中可添加额外的逻辑进行安全校验，例如判断请求的来源IP是否合法等

	token := "your_token"
	if !checkSignature(token, timestamp, nonce, signature) {
		c.AbortWithStatusJSON(403, gin.H{"error": "Invalid signature"})
		return
	}

	c.String(200, echoStr)
}

func checkSignature(token, timestamp, nonce, signature string) bool {
	// 将 token、timestamp、nonce 进行字典序排序
	params := []string{token, timestamp, nonce}
	sort.Strings(params)

	// 将排序后的三个参数拼接成一个字符串
	str := strings.Join(params, "")

	// 对拼接后的字符串进行 sha1 加密
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(str))
	hashed := sha1Hash.Sum(nil)

	// 将 sha1 加密后的结果转换为字符串形式
	signatureCalculated := fmt.Sprintf("%x", hashed)

	// 将加密后的结果与微信传递的签名进行比对，如果一致，则返回 true，否则返回 false
	return signature == signatureCalculated
}

var Router * gin.Engine
func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello world!",
		})
	})
	r.GET("/wx/callback", wxCallbackHandler)
	r.Run()
}
