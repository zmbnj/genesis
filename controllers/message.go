package controllers

import (
	_ "errors"
	"fmt"
	"genesis/models"
	_ "strconv"
	_ "strings"

	"github.com/astaxie/beego"
)

// oprations for Message
type MessageController struct {
	beego.Controller
}

func (c *MessageController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
}

// @Title Post
// @Description create Message
// @Param	body		body 	models.Message	true		"body for Message content"
// @Success 201 {int} models.Message
// @Failure 403 body is empty
// @router / [post]
func (c *MessageController) Post() {
	var v models.Message
	fmt.Printf(string(c.Ctx.Input.RequestBody))
	if _, err := models.HandleMessage(c.Ctx.Input.RequestBody); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["xml"] = v
	} else {
		c.Data["xml"] = err.Error()
	}
	c.ServeXML()
}

//GET /?signature=d01007dcff994c555bc51d22e154956ccdc61ec5×tamp=1418970951&nonce=484765335&echostr=qwe1235

// @Title Get
// @Description get Message by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Message
// @Failure 403 :id is empty
// @router / [get]
func (c *MessageController) Get() {

	signature := c.GetString("signature")
	timestamp := c.GetString("timestamp")
	nonce := c.GetString("nonce")
	echostr := c.GetString("echostr")

	//v, err1 := models.GetUserByOpenid("oh6Xgw0qeS4Uk-D-CeDPpgtZ2Ss0")
	//fmt.Println("------------------------------", err1)
	_, err2 := models.GetReplayMessage()

	fmt.Println(err2)
	//v.subscribe
	err := models.CheckMessageInterface(signature, timestamp, nonce, echostr)
	if err != nil {
		c.Ctx.Output.Body([]byte(err.Error()))
	} else {
		c.Ctx.Output.Body([]byte(echostr))
	}
}
