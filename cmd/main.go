package main

import (
	"github.com/iot-master-contrib/camera"
	_ "github.com/iot-master-contrib/webui"
	master "github.com/zgwit/iot-master/v4"
	"github.com/zgwit/iot-master/v4/log"
	"github.com/zgwit/iot-master/v4/web"
)

func main() {

	err := master.Startup()
	if err != nil {
		log.Fatal(err)
	}

	_ = camera.Startup()

	err = web.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
