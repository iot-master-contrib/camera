package camera

import (
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/zgwit/iot-master/v4/lib"
	"github.com/zgwit/iot-master/v4/pkg/db"
	"time"
)

func init() {
	db.Register(new(Camera))
}

type Camera struct {
	Id string `json:"id" xorm:"pk"`

	ProjectId string `json:"project_id,omitempty" xorm:"index"`
	Project   string `json:"project,omitempty" xorm:"<-"`

	Name        string    `json:"name"`
	Url         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	Audio       bool      `json:"audio,omitempty"`
	Disabled    bool      `json:"disabled,omitempty"`
	Created     time.Time `json:"created,omitempty" xorm:"created"`

	//RTSP连接
	rtsp    *rtspv2.RTSPClient
	clients lib.Map[Client]
}

func (c *Camera) Check() error {
	if c.rtsp != nil {
		return nil
	}

	var err error
	c.rtsp, err = rtspv2.Dial(rtspv2.RTSPClientOptions{
		URL:              c.Url,
		DialTimeout:      3 * time.Second,
		ReadWriteTimeout: 3 * time.Second,
		DisableAudio:     !c.Audio,
	})
	if err != nil {
		return err
	}

	go c.receive()

	return nil
}

func (c *Camera) receive() {
	defer func() {
		c.rtsp.Close()
		c.rtsp = nil
	}()

	autoClose := time.NewTimer(20 * time.Second)
	for {
		select {
		case <-autoClose.C:
			return
		case signals := <-c.rtsp.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				//Config.coAd(name, c.rtsp.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return
			}
		case pkt := <-c.rtsp.OutgoingPacketQueue:
			if pkt.IsKeyFrame {
				autoClose.Reset(20 * time.Second)
			}
			c.clients.Range(func(_ string, client *Client) bool {
				if len(client.queue) < cap(client.queue) {
					client.queue <- pkt
				}
				return true
			})
		}
	}
}
