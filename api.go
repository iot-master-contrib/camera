package camera

import (
	"github.com/gin-gonic/gin"
	"github.com/zgwit/iot-master/v4/api"
	"github.com/zgwit/iot-master/v4/pkg/log"
	"github.com/zgwit/iot-master/v4/web/curd"
	"io"
	"time"
)

func init() {
	api.Register("POST", "ipc/count", curd.ApiCount[Camera]())
	api.Register("POST", "ipc/search", curd.ApiSearch[Camera]())
	api.Register("POST", "ipc/create", curd.ApiCreate[Camera]())
	api.Register("POST", "ipc/:id", curd.ParseParamStringId, curd.ApiUpdate[Camera]())
	api.Register("GET", "ipc/:id", curd.ParseParamStringId, curd.ApiGet[Camera]())
	api.Register("GET", "ipc/:id/delete", curd.ParseParamStringId, curd.ApiDelete[Camera]())
	api.Register("POST", "ipc/:id/stream", func(ctx *gin.Context) {
		id := ctx.Param("id")

		cam, err := Ensure(id)
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		err = cam.Check()
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		c := NewClient()
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		answer, err := c.CreateAnswer(cam.rtsp.CodecData, string(buf))
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		curd.OK(ctx, answer)

		go func() {
			defer c.Close()
			finished := time.NewTimer(10 * time.Second)

			var videoStart bool
			for {
				select {
				case <-finished.C:
					return
				case pck := <-c.queue:
					if pck.IsKeyFrame {
						finished.Reset(10 * time.Second)
						videoStart = true
					}
					if videoStart {
						err = c.WritePacket(pck)
						if err != nil {
							log.Println("WritePacket", err)
							return
						}
					}
				}
			}

		}()
	})
}

func sendVideo(c *Client) {

}
