package main

import (
	"github.com/codeflames/cizzors/models"
	"github.com/codeflames/cizzors/server"
)

// @title           Cizzors url shortener API
// @version         1.0
// @description     This is a simple url shortening api service.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Yereka
// @contact.url    https://github.com/codeflames/cizzors/issues
// @contact.email  yerekadonald@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3001

//@securityDefinitions.apikey Bearer
//@in header
//@name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	models.Setup()
	server.SetupAndListen()
}
