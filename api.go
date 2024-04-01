package camera

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pion/webrtc/v3"
	"github.com/zgwit/iot-master/v4/api"
	"github.com/zgwit/iot-master/v4/pkg/log"
	"github.com/zgwit/iot-master/v4/web/curd"
	"io"
	"net/http"
	"time"
)

func init() {
	api.Register("POST", "camera/count", curd.ApiCount[Camera]())
	api.Register("POST", "camera/search", curd.ApiSearch[Camera]())
	api.Register("POST", "camera/create", curd.ApiCreateHook[Camera](curd.GenerateID[Camera](), nil))
	api.Register("POST", "camera/:id", curd.ParseParamStringId, curd.ApiUpdate[Camera]())
	api.Register("GET", "camera/:id", curd.ParseParamStringId, curd.ApiGet[Camera]())
	api.Register("GET", "camera/:id/delete", curd.ParseParamStringId, curd.ApiDelete[Camera]())

	api.Register("GET", "camera/:id/getice/:uid", func(ctx *gin.Context) {
		id := ctx.Param("id")

		cam := Get(id)
		if cam == nil {
			curd.Fail(ctx, "找不到设备")
			return
		}
		client := cam.clients.Load(ctx.Param("uid"))
		if client == nil {
			curd.Fail(ctx, "找不到设备")
			return
		}
		curd.OK(ctx, client.candidates)
	})

	api.Register("POST", "camera/:id/addice/:uid", func(ctx *gin.Context) {
		id := ctx.Param("id")

		cam := Get(id)
		if cam == nil {
			curd.Fail(ctx, "找不到设备")
			return
		}
		client := cam.clients.Load(ctx.Param("uid"))
		if client == nil {
			curd.Fail(ctx, "找不到设备")
			return
		}
		var candidate webrtc.ICECandidateInit
		err := ctx.BindJSON(&candidate)
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		err = client.pc.AddICECandidate(candidate)
		if err != nil {
			curd.Error(ctx, err)
			return
		}

		curd.OK(ctx, nil)
	})

	api.Register("POST", "camera/:id/stream", func(ctx *gin.Context) {
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

		//生成ID
		uid := uuid.NewString()
		cam.clients.Store(uid, c)

		//curd.OK(ctx, answer)

		ctx.JSON(http.StatusOK, gin.H{
			"data": answer,
			"uuid": uid,
		})

		go sendVideo(c)
	})
}

func sendVideo(c *Client) {
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
				err := c.WritePacket(pck)
				if err != nil {
					log.Println("WritePacket", err)
					return
				}
			}
		}
	}
}
