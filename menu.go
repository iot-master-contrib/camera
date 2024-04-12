package camera

import "github.com/zgwit/iot-master/v4/menu"

func init() {
	menu.Register("camera", &menu.Menu{
		Name:       "视频监控",
		Icon:       "video-camera",
		Domain:     []string{"admin"},
		Privileges: nil,
		Items: []*menu.Item{
			{Name: "摄像头", Url: "/$camera/camera", Type: "web"},
			{Name: "推流器", Url: "/$camera/streamer", Type: "web"},
		},
	})
}
