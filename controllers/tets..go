// @APIVersion 1.0.0
// @Title API Test
// @Description API para responder con 'Pong' al hacer un GET a /ping
// @Contact email@example.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

package controllers

import (
	"github.com/astaxie/beego"
)

type TestController struct {
	beego.Controller
}

func (c *TestController) URLMapping() {
	c.Mapping("Ping", c.Ping)
}

// Ping responde con "Pong" a la solicitud GET en /ping
// @router /ping [get]
func (c *TestController) Ping() {
	c.Ctx.WriteString("Pong")
}
