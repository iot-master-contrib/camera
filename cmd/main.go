package main

import (
	"github.com/iot-master-contrib/camera"
	master "github.com/zgwit/iot-master/v4"
	"github.com/zgwit/iot-master/v4/pkg/log"
	"github.com/zgwit/iot-master/v4/plugin"
	"github.com/zgwit/iot-master/v4/web"
)

func main() {

	_ = camera.Startup()

	err := master.Startup()
	if err != nil {
		log.Fatal(err)
	}

	plg := camera.Manifest()
	err = plg.Startup()
	if err != nil {
		log.Fatal(err)
	}

	//注册插件
	plugin.Register(plg)

	err = web.Serve()
	if err != nil {
		log.Fatal(err)
	}

	err = plg.Shutdown()
	if err != nil {
		log.Fatal(err)
	}

}
