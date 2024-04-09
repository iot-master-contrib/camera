import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {RequestService} from "iot-master-smart";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-player',
    standalone: true,
    imports: [],
    templateUrl: './player.component.html',
    styleUrl: './player.component.scss'
})
export class PlayerComponent implements OnInit {

    id: any = ''
    camera: any;

    @ViewChild("video") video!: ElementRef<any>

    config: RTCConfiguration = {
        iceServers: [{
            urls: ["stun:stun.l.google.com:19302"]
        }]
    };

    //webrtc-streamer连接的参数
    cid: any = ''
    ws!: WebSocket
    pc: RTCPeerConnection = new RTCPeerConnection()
    stream = new MediaStream();

    constructor(private rs: RequestService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.id = this.route.snapshot.paramMap.get("id") || this.route.snapshot.queryParamMap.get("id")
        this.load()
    }

    load() {
        this.rs.get("camera/" + this.id).subscribe(res => {
            this.camera = res.data
            this.connect()
        })
    }

    send(type: string, data: any) {
        console.log('[SEND] ===>', type, data)

        if (typeof data === "object") {
            data = JSON.stringify(data)
        }
        let text = JSON.stringify({id: this.cid, type, data})
        this.ws.send(text)
    }

    connect() {
        //ws://localhost:8080/streamer/test/connect
        let url = location.protocol.replace("http", "ws")
            + "//" + location.host + "/streamer/" + this.camera.streamer_id + "/connect"
        this.ws = new WebSocket(url)
        this.ws.onerror = console.error

        this.ws.onopen = (event) => {
            console.log("websocket onopen")
            this.send("connect", {url: this.camera.url})
        }

        this.ws.onmessage = async (event) => {
            //console.log("<---", event.data)
            let msg = JSON.parse(event.data)
            console.log('[RECV] <===', msg.type, msg.data)

            this.cid = msg.id
            switch (msg.type) {
                case "offer":
                    await this.pc.setRemoteDescription(new RTCSessionDescription({type: 'offer', sdp: msg.data}))
                    let answer = await this.pc.createAnswer()
                    this.send("answer", answer.sdp)
                    await this.pc.setLocalDescription(answer)
                    break
                case "answer":
                    await this.pc.setRemoteDescription(new RTCSessionDescription({type: 'answer', sdp: msg.data}))
                    break
                case "candidate":
                    if (msg.data)
                        await this.pc.addIceCandidate(new RTCIceCandidate(JSON.parse(msg.data)))
                    break
                case "error":
                    console.error("streamer error", msg.data)
                    break
            }
        }

        this.pc.onnegotiationneeded = async function () {
            console.log("onnegotiationneeded")

            // let offer = await pc.createOffer()
            // await pc.setLocalDescription(offer)
            // send("offer", offer.sdp)
        };

        this.pc.ontrack = (event) => {
            console.log("ontrack", event.streams)

            this.stream.addTrack(event.track);
            // let videoElem = document.getElementById("video")
            // videoElem.srcObject = stream;
            this.video.nativeElement.srcObject = this.stream
        }

        this.pc.onicecandidate = (event) => {
            console.log("candidate", event.candidate)
            if (event.candidate)
                this.send("candidate", event.candidate.toJSON())
        }

        this.pc.oniceconnectionstatechange = (event) => {
            console.log("oniceconnectionstatechange", this.pc.iceConnectionState)
        }

        this.pc.addTransceiver("video", {'direction': 'sendrecv'})
    }

}
