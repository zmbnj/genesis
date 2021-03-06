package controllers

import (
	"encoding/json"
	"genesis/models"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// WeblogController oprations for Weblog
type WeblogController struct {
	beego.Controller
}

//URLMapping URLMapping
func (c *WeblogController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Prepare 拦截请求
func (c *WeblogController) Prepare() {
	token := c.Ctx.Request.Header.Get("Token")
	err := models.CheckSessionByToken(token)
	if err != nil {
		c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
		c.ServeJSON()
		c.StopRun()
	}
}

// Post Post
// @Title Post
// @Description create Weblog
// @Param	body		body 	models.Weblog	true		"body for Weblog content"
// @Success 201 {int} models.Weblog
// @Failure 403 body is empty
// @router / [post]
func (c *WeblogController) Post() {
	var v models.Weblog
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddWeblog(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = models.GetReturnData(0, "OK", v)
		} else {
			c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
		}
	} else {
		c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
	}
	c.ServeJSON()
}

// GetOne GetOne
// @Title Get
// @Description get Weblog by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Weblog
// @Failure 403 :id is empty
// @router /:id [get]
func (c *WeblogController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetWeblogByID(id)
	if err != nil {
		c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
	} else {
		c.Data["json"] = models.GetReturnData(0, "OK", v)
	}
	c.ServeJSON()
}

// GetAll GetAll
// @Title Get All
// @Description get Weblog
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Weblog
// @Failure 403
// @router / [get]
func (c *WeblogController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query map[string]string = make(map[string]string)
	var limit int64 = 10
	var offset int64 = 0

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.Split(cond, ":")
			if len(kv) != 2 {
				c.Data["json"] = models.GetReturnData(-1, "Error: invalid query key/value pair", nil)
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllWeblog(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
	} else {
		c.Data["json"] = models.GetReturnData(0, "OK", l)
	}
	c.ServeJSON()
}

// Put Put
// @Title Update
// @Description update the Weblog
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Weblog	true		"body for Weblog content"
// @Success 200 {object} models.Weblog
// @Failure 403 :id is not int
// @router /:id [put]
func (c *WeblogController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Weblog{ID: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateWeblogByID(&v); err == nil {
			c.Data["json"] = models.GetReturnData(0, "OK", nil)
		} else {
			c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
		}
	} else {
		c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
	}
	c.ServeJSON()
}

// PutReviewed PutReviewed
// @router /:id/reviewed [put]
func (c *WeblogController) PutReviewed() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Weblog{ID: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateWeblogReviewedByID(&v); err == nil {
			c.Data["json"] = models.GetReturnData(0, "OK", nil)
		} else {
			c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
		}
	} else {
		c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
	}
	c.ServeJSON()
}

// Delete Delete
// @Title Delete
// @Description delete the Weblog
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *WeblogController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteWeblog(id); err == nil {
		c.Data["json"] = models.GetReturnData(0, "OK", nil)
	} else {
		c.Data["json"] = models.GetReturnData(-1, err.Error(), nil)
	}
	c.ServeJSON()
}
