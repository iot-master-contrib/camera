package camera

import (
	"github.com/zgwit/iot-master/v4/api"
	"github.com/zgwit/iot-master/v4/pkg/db"
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
}
