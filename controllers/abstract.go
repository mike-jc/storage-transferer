package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"gitlab.com/24sessions/lib-go-logger/logger"
	"gitlab.com/24sessions/lib-go-logger/logger/services"
	"net/http"
	"strings"
	"sync"
)

type AbstractController struct {
	beego.Controller

	loggerLock sync.Mutex
	lg         *logger.Logger
}

func (c *AbstractController) getClientIp() string {
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	if len(ip) > 0 && ip[0] != "[" {
		return ip[0]
	}
	return "127.0.0.1"
}

func (c *AbstractController) GetLogger() *logger.Logger {
	c.loggerLock.Lock()
	defer c.loggerLock.Unlock()

	if c.lg == nil {
		c.lg = new(logger.Logger).
			SetSubject("anonymous", "").
			SetTraceId(uuid.NewV4String()).
			SetClientIp(c.getClientIp()).
			SetParent(LogMain)
	}
	return c.lg
}

func (c *AbstractController) addCors() {
	c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Output.Header("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE")
	c.Ctx.Output.Header("Access-Control-Allow-Headers", "Authorization,Content-type,Accept-Language")
}

func (c *AbstractController) SuccessResponse() {
	result := make(map[string]interface{})
	result["status"] = "success"

	c.Data["json"] = result
	c.addCors()
	c.ServeJSON()
}

func (c *AbstractController) ShowError(message string, status int, errorCode int, reason string, showReason bool) {
	if errorCode <= 0 {
		errorCode = status
	}

	// log the error
	logRow := logger.CreateError(fmt.Sprintf("%d: %s, %s", status, message, reason)).
		SetErrorCode(errorCode).
		AddData("request", c.MarshalRequest(c.Ctx.Request))

	if status < 500 {
		logRow.SetLevel(logger.LOG_LEVEL_NOTICE)
	}

	c.GetLogger().Log(logRow)

	// output the error
	result := make(map[string]interface{})
	result["status"] = "error"
	result["error"] = message
	if showReason {
		result["description"] = reason
	}
	if logRow.GetId() != "" {
		result["logId"] = logRow.GetId()
	}

	c.Data["json"] = result
	c.Ctx.Output.Status = status
	c.addCors()
	c.ServeJSON()
}

func (c *AbstractController) MarshalRequest(r *http.Request) string {
	res := make(map[string]interface{})
	res["method"] = r.Method
	res["url"] = r.URL

	bytes, _ := json.Marshal(res)
	return string(bytes)
}
