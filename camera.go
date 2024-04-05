package camera

import (
	"github.com/zgwit/iot-master/v4/api"
	"github.com/zgwit/iot-master/v4/pkg/db"
	"github.com/zgwit/iot-master/v4/web/curd"
	"time"
)

type Camera struct {
	Id string `json:"id" xorm:"pk"`

	ProjectId string `json:"project_id,omitempty" xorm:"index"`
	Project   string `json:"project,omitempty" xorm:"<-"`

	StreamerId string `json:"streamer_id,omitempty" xorm:"index"`
	Streamer   string `json:"streamer,omitempty" xorm:"<-"`

	Name        string    `json:"name"`
	Url         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	Audio       bool      `json:"audio,omitempty"`
	Disabled    bool      `json:"disabled,omitempty"`
	Created     time.Time `json:"created,omitempty" xorm:"created"`
}

func init() {
	db.Register(new(Camera))

	api.Register("POST", "camera/count", curd.ApiCount[Camera]())
	api.Register("POST", "camera/search", curd.ApiSearch[Camera]())
	api.Register("POST", "camera/create", curd.ApiCreateHook[Camera](curd.GenerateID[Camera](), nil))
	api.Register("POST", "camera/:id", curd.ParseParamStringId, curd.ApiUpdate[Camera]())
	api.Register("GET", "camera/:id", curd.ParseParamStringId, curd.ApiGet[Camera]())
	api.Register("GET", "camera/:id/delete", curd.ParseParamStringId, curd.ApiDelete[Camera]())
}
