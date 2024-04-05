package camera

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/zgwit/iot-master/v4/pkg/log"
	"github.com/zgwit/iot-master/v4/plugin"
	"github.com/zgwit/iot-master/v4/web"
	"github.com/zgwit/webrtc-streamer/signaling"
	"gopkg.in/yaml.v3"
	"net/http"
	"strings"
)

//go:embed all:www
var wwwFiles embed.FS

//go:embed manifest.yaml
var _manifest string
var manifest plugin.Manifest

func Manifest() *plugin.Manifest {
	return &manifest
}

var server signaling.Server

var upper = &websocket.Upgrader{
	//HandshakeTimeout: time.Second,
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	Subprotocols:    []string{"webrtc"},
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func init() {
	//前端静态文件
	web.Static.Put("/$camera", http.FS(wwwFiles), "www", "index.html")

	d := yaml.NewDecoder(strings.NewReader(_manifest))
	err := d.Decode(&manifest)
	if err != nil {
		log.Fatal(err)
	}

	manifest.Startup = Startup
	manifest.Shutdown = Shutdown

	plugin.Register(&manifest)
}

func Startup() error {
	//group := web.Engine.Group("/$gateway")
	//
	////注册前端接口
	//api.RegisterRoutes(group.Group("/api"))
	//
	////注册接口文档
	//web.RegisterSwaggerDocs(group, "gateway")

	web.Engine.GET("streamer/:id", func(ctx *gin.Context) {
		ws, err := upper.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Error(err)
			return
		}

		//注册
		server.ConnectStreamer(ctx.Param("id"), ws)
	})

	//这里没有鉴权了
	web.Engine.GET("streamer/:id/connect", func(ctx *gin.Context) {
		ws, err := upper.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Error(err)
			return
		}
		server.ConnectViewer(ctx.Param("id"), ws)
	})

	return nil
}

func Shutdown() error {
	//只关闭Web就行了，其他通过defer关闭
	return nil
}
