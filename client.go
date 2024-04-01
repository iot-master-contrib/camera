package camera

import (
	"bytes"
	"errors"
	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/codec/h265parser"
	"github.com/pion/interceptor"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"log"
	"time"
)

var (
	ErrorNotFound          = errors.New("WebRTC Stream Not Found")
	ErrorCodecNotSupported = errors.New("WebRTC Codec Not Supported")
	ErrorClientOffline     = errors.New("WebRTC Client Offline")
	ErrorNotTrackAvailable = errors.New("WebRTC Not Track Available")
	ErrorIgnoreAudioTrack  = errors.New("WebRTC Ignore Audio Track codec not supported WebRTC support only PCM_ALAW or PCM_MULAW")
)

type Client struct {
	streams   map[int8]*Stream
	status    webrtc.ICEConnectionState
	stop      bool
	pc        *webrtc.PeerConnection
	ClientACK *time.Timer
	StreamACK *time.Timer

	queue      chan *av.Packet
	candidates []*webrtc.ICECandidate
}

type Stream struct {
	codec av.CodecData
	track *webrtc.TrackLocalStaticSample
}

func consumeRtpSender(sender *webrtc.RTPSender) {
	rtcpBuf := make([]byte, 1500)
	for {
		if _, _, rtcpErr := sender.Read(rtcpBuf); rtcpErr != nil {
			return
		}
	}
}

func NewClient() *Client {
	m := Client{
		ClientACK: time.NewTimer(time.Second * 20),
		StreamACK: time.NewTimer(time.Second * 20),
		streams:   make(map[int8]*Stream),
		queue:     make(chan *av.Packet, 100),
	}
	//go m.WaitCloser()
	return &m
}

func (c *Client) NewPeerConnection(configuration webrtc.Configuration) (*webrtc.PeerConnection, error) {
	if len(configuration.ICEServers) == 0 {
		configuration.ICEServers = []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}}
	}

	me := &webrtc.MediaEngine{}
	if err := me.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}

	r := &interceptor.Registry{}
	if err := webrtc.RegisterDefaultInterceptors(me, r); err != nil {
		return nil, err
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(me), webrtc.WithInterceptorRegistry(r))
	return api.NewPeerConnection(configuration)
}

func (c *Client) CreateAnswer(streams []av.CodecData, sdp string) (string, error) {
	var success bool
	if len(streams) == 0 {
		return "", ErrorNotFound
	}

	offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp}

	pc, err := c.NewPeerConnection(webrtc.Configuration{SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback})
	if err != nil {
		return "", err
	}
	defer func() {
		if !success {
			_ = c.Close()
		}
	}()

	for i, stream := range streams {
		var track *webrtc.TrackLocalStaticSample
		if stream.Type().IsVideo() {
			if stream.Type() == av.H264 {
				track, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
					MimeType: webrtc.MimeTypeH264,
				}, "pion-rtsp-video", "pion-video")
				if err != nil {
					return "", err
				}
				if rtpSender, err := pc.AddTrack(track); err != nil {
					return "", err
				} else {
					go consumeRtpSender(rtpSender)
				}
			}
			if stream.Type() == av.H265 {
				track, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
					MimeType: webrtc.MimeTypeH265,
				}, "pion-rtsp-video", "pion-video")
				if err != nil {
					return "", err
				}
				if rtpSender, err := pc.AddTrack(track); err != nil {
					return "", err
				} else {
					go consumeRtpSender(rtpSender)
				}
			}
		} else if stream.Type().IsAudio() {
			AudioCodecString := webrtc.MimeTypePCMA
			switch stream.Type() {
			case av.PCM_ALAW:
				AudioCodecString = webrtc.MimeTypePCMA
			case av.PCM_MULAW:
				AudioCodecString = webrtc.MimeTypePCMU
			case av.OPUS:
				AudioCodecString = webrtc.MimeTypeOpus
			default:
				//log.Println(ErrorIgnoreAudioTrack)
				continue
			}
			track, err = webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{
				MimeType:  AudioCodecString,
				Channels:  uint16(stream.(av.AudioCodecData).ChannelLayout().Count()),
				ClockRate: uint32(stream.(av.AudioCodecData).SampleRate()),
			}, "pion-rtsp-audio", "pion-audio")
			if err != nil {
				return "", err
			}
			if rtpSender, err := pc.AddTrack(track); err != nil {
				return "", err
			} else {
				go consumeRtpSender(rtpSender)
			}
		}
		c.streams[int8(i)] = &Stream{track: track, codec: stream}
	}
	if len(c.streams) == 0 {
		return "", ErrorNotTrackAvailable
	}
	pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		c.status = connectionState
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			_ = c.Close()
		}
	})
	pc.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			c.ClientACK.Reset(5 * time.Second)
		})
	})
	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		c.candidates = append(c.candidates, candidate)
	})

	if err = pc.SetRemoteDescription(offer); err != nil {
		return "", err
	}

	gc := webrtc.GatheringCompletePromise(pc)

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		return "", err
	}

	if err = pc.SetLocalDescription(answer); err != nil {
		return "", err
	}
	c.pc = pc

	waitT := time.NewTimer(time.Second * 10)
	select {
	case <-waitT.C:
		return "", errors.New("GatheringCompletePromise wait")
	case <-gc:
		//Connected
	}

	resp := pc.LocalDescription()
	success = true

	return resp.SDP, nil
}

func (c *Client) WritePacket(pkt *av.Packet) (err error) {
	var WritePacketSuccess bool
	defer func() {
		if !WritePacketSuccess {
			_ = c.Close()
		}
	}()

	if c.stop {
		return ErrorClientOffline
	}
	if c.status == webrtc.ICEConnectionStateChecking {
		WritePacketSuccess = true
		return nil
	}
	if c.status != webrtc.ICEConnectionStateConnected {
		return nil
	}
	if stream, ok := c.streams[pkt.Idx]; ok {
		c.StreamACK.Reset(10 * time.Second)
		if len(pkt.Data) < 5 {
			return nil
		}
		switch stream.codec.Type() {
		case av.H264:
			nalus, _ := h264parser.SplitNALUs(pkt.Data)
			for _, nalu := range nalus {
				naltype := nalu[0] & 0x1f
				if naltype == 5 {
					codec := stream.codec.(h264parser.CodecData)
					err = stream.track.WriteSample(media.Sample{
						Data: append([]byte{0, 0, 0, 1},
							bytes.Join([][]byte{codec.SPS(), codec.PPS(), nalu}, []byte{0, 0, 0, 1})...),
						Duration: pkt.Duration})
				} else {
					err = stream.track.WriteSample(media.Sample{
						Data:     append([]byte{0, 0, 0, 1}, nalu...),
						Duration: pkt.Duration})
				}
				if err != nil {
					return err
				}
			}
			WritePacketSuccess = true
			return
		case av.H265:
			nalus, _ := h265parser.SplitNALUs(pkt.Data)
			for _, nalu := range nalus {
				naltype := (nalu[0] & 0x7e) >> 1
				if naltype == 5 {
					codec := stream.codec.(h265parser.CodecData)
					err = stream.track.WriteSample(media.Sample{
						Data: append([]byte{0, 0, 0, 1},
							bytes.Join([][]byte{codec.VPS(), codec.SPS(), codec.PPS(), nalu}, []byte{0, 0, 0, 1})...),
						Duration: pkt.Duration})
				} else {
					err = stream.track.WriteSample(media.Sample{
						Data:     append([]byte{0, 0, 0, 1}, nalu...),
						Duration: pkt.Duration})
				}
				if err != nil {
					return err
				}
			}
			WritePacketSuccess = true
			return
		case av.PCM_ALAW:
		case av.OPUS:
		case av.PCM_MULAW:
		case av.AAC:
			//TODO: NEED ADD DECODER AND ENCODER
			return ErrorCodecNotSupported
		case av.PCM:
			//TODO: NEED ADD ENCODER
			return ErrorCodecNotSupported
		default:
			return ErrorCodecNotSupported
		}
		err = stream.track.WriteSample(media.Sample{Data: pkt.Data, Duration: pkt.Duration})
		if err == nil {
			WritePacketSuccess = true
		}
		return err
	} else {
		WritePacketSuccess = true
		return nil
	}
}

func (c *Client) Transport() {
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

func (c *Client) Close() error {
	c.stop = true
	if c.pc != nil {
		err := c.pc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
