package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/polynetwork/poly-nft-bridge/models"
)

const (
	ErrCodeRequest int = 400
)

var errMap = map[int]string{
	ErrCodeRequest: "request parameter is invalid!",
}

func input(c *beego.Controller, req interface{}) error {
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		code := ErrCodeRequest
		c.Data["json"] = models.MakeErrorRsp(errMap[code])
		c.Ctx.ResponseWriter.WriteHeader(code)
		c.ServeJSON()
		return err
	}
	return nil
}

func output(c *beego.Controller, data interface{}) {
	c.Data["json"] = data
	c.ServeJSON()
}
