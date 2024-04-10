package camera

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/iot-master/v4/api"
	"github.com/zgwit/iot-master/v4/db"
	"github.com/zgwit/iot-master/v4/log"
	"github.com/zgwit/iot-master/v4/web/curd"
	"time"
)

type Streamer struct {
	Id          string    `json:"id" xorm:"pk"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Disabled    bool      `json:"disabled,omitempty"`
	Created     time.Time `json:"created,omitempty" xorm:"created"`
}

func init() {
	db.Register(new(Streamer))

	api.Register("POST", "streamer/count", curd.ApiCount[Streamer]())
	api.Register("POST", "streamer/search", curd.ApiSearch[Streamer]())
	api.Register("POST", "streamer/create", curd.ApiCreateHook[Streamer](curd.GenerateID[Streamer](), nil))
	api.Register("POST", "streamer/:id", curd.ParseParamStringId, curd.ApiUpdate[Streamer]())
	api.Register("GET", "streamer/:id", curd.ParseParamStringId, curd.ApiGet[Streamer]())
	api.Register("GET", "streamer/:id/delete", curd.ParseParamStringId, curd.ApiDelete[Streamer]())

	api.Register("GET", "streamer/:id/connect", func(ctx *gin.Context) {
		ws, err := upper.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Error(err)
			return
		}
		server.ConnectViewer(ctx.Param("id"), ws)
	})
}
