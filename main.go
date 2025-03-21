package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/heavi715/api-proxy/config"
	"github.com/heavi715/api-proxy/public/logger"
	"github.com/heavi715/api-proxy/public/response"
)

func init() {
	config.LoadConfig("config.json")

	logger.Info("Start Listen On ", config.Config.ServerAddr)

}

func main() {

	Start()
}

func Start() {

	router := gin.New()
	router.Use(logger.RequestLog(), logger.ResponseLog(), gin.Recovery(), func(c *gin.Context) {
		//获取请求ip
		ip := c.ClientIP()
		if !strings.HasPrefix(ip, "192.168.") &&
			!strings.HasPrefix(ip, "172.1") &&
			!strings.HasPrefix(ip, "127.0") &&
			!strings.HasPrefix(ip, "10.100") &&
			!strings.HasPrefix(ip, "223.71.41") {
			c.JSON(http.StatusForbidden, "非法请求")
			c.Abort()
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.Any("/health", func(c *gin.Context) {
		response.Success(
			c,
			"success",
		)
	})
	router.Any("/proxy/:platform/:source/*url", chatProxyHandler)
	router.Run(config.Config.ServerAddr)

	logger.Info("Start Listen On ", config.Config.ServerAddr)

}

func chatProxyHandler(c *gin.Context) {

	//todo
	//1.source+ key 授权v,whiteIp
	//2.token计数到source+限量
	//3.加密
	//source 获取get参数source
	source := c.Param("source")
	platform := c.Param("platform")
	fmt.Println("source:", source)
	fmt.Println("platform:", platform)

	//source 在不在config.Config.SourceList里
	if !config.Config.IsSource(source) {
		return
	}
	//targetURLMap 里有没有platform
	if _, ok := config.Config.PlatformList[platform]; !ok {
		return
	}

	u := c.Param("url")
	fmt.Println("url:", u)

	target := config.Config.PlatformList[platform].Url // 目标域名

	r := c.Request
	w := c.Writer
	// 过滤无效URL
	_, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println("Error parsing URL: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// 拼接目标URL
	newPath := strings.Replace(r.RequestURI, "/"+source, "", 1)
	newPath = strings.Replace(newPath, "/"+platform, "", 1)
	newPath = strings.Replace(newPath, "/proxy", "", 1)
	newPath = strings.Replace(newPath, "/v1/v1", "/v1", 1)

	targetURL := target + newPath

	fmt.Println("targetURL:", targetURL)
	// 创建代理HTTP请求
	proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		log.Println("Error creating proxy request: ", err.Error())
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	// 将原始请求头复制到新请求中
	for headerKey, headerValues := range r.Header {
		for _, headerValue := range headerValues {
			proxyReq.Header.Add(headerKey, headerValue)
			fmt.Println("headerKey:", headerKey, "headerValue:", headerValue)
		}
	}

	authorization := proxyReq.Header.Get("Authorization")
	fmt.Println("authorization:", authorization)
	isValidServerKey := false
	if config.Config.ServerKeyList == nil || len(config.Config.ServerKeyList) == 0 {
		isValidServerKey = true
	} else {
		for _, serverKey := range config.Config.ServerKeyList {
			if serverKey == authorization {
				isValidServerKey = true
				break
			}
		}
	}

	if isValidServerKey {
		targetKeyHeader := config.Config.PlatformList[platform].HeaderKey
		targetKeyList := config.Config.PlatformList[platform].HeaderValues
		if targetKeyHeader != "" && len(targetKeyList) > 0 {
			targetKey := targetKeyList[rand.Intn(len(targetKeyList))]
			proxyReq.Header.Set(targetKeyHeader, targetKey)
		}
	} else {
		// 如果serverKey不合法，则返回401
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 默认超时时间设置为300s（应对长上下文）
	client := &http.Client{
		Timeout: 300 * time.Second,
	}

	// 本地测试通过代理请求 OpenAI 接口
	if config.Config.ProxyURL != "" {
		proxyURL, _ := url.Parse(config.Config.ProxyURL) // 本地HTTP代理配置
		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// 向 OpenAI 发起代理请求
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Println("Error sending proxy request: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 将响应头复制到代理响应头中
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 将响应状态码设置为原始响应状态码
	w.WriteHeader(resp.StatusCode)

	// 将响应实体写入到响应流中（支持流式响应）
	buf := make([]byte, 1024)
	for {
		if n, err := resp.Body.Read(buf); err == io.EOF || n == 0 {
			return
		} else if err != nil {
			log.Println("error while reading respbody: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if _, err = w.Write(buf[:n]); err != nil {
				log.Println("error whilfe writing resp: ", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.(http.Flusher).Flush()
		}
	}

}
