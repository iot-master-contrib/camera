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
    uid: any = ''

    @ViewChild("video") video!: ElementRef<any>

    config: RTCConfiguration = {
        iceServers: [{
            urls: ["stun:stun.l.google.com:19302"]
        }]
    };

    pc!: RTCPeerConnection

    stream = new MediaStream();



    constructor(private rs: RequestService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.id = this.route.snapshot.paramMap.get("id") || this.route.snapshot.queryParamMap.get("id")
        this.play().then()
    }


    async play() {
        console.log('play', this.id)

        this.pc = new RTCPeerConnection(this.config)

        this.pc.onnegotiationneeded = async (event) => {
            console.log('onnegotiationneeded', event)
            let offer = await this.pc.createOffer();
            await this.pc.setLocalDescription(offer);
            this.connect(offer);
        }

        this.pc.onicecandidate = (event) => {
            console.log('onicecandidate', event)
            //this.onIceCandidate(evt)
            this.rs.post(`camera/${this.id}/addice/${this.uid}`, event.candidate).subscribe((res: any) => {
            })
        }
        this.pc.onicegatheringstatechange = (event) => {
            console.log('onicegatheringstatechange', this.pc.iceConnectionState)
        };

        this.pc.ontrack = (event) => {
            console.log('ontrack', event)
            this.stream.addTrack(event.track);
            this.video.nativeElement.srcObject = this.stream;
        }

        this.pc.oniceconnectionstatechange = (eevent) => {
            console.log('oniceconnectionstatechange', this.pc.iceConnectionState)
        }

        let offer = await this.pc.createOffer({offerToReceiveAudio: true, offerToReceiveVideo: true});
        console.log('createOffer', offer)

        await this.pc.setLocalDescription(offer);
        this.connect(offer);
    }

    connect(offer: any) {
        this.rs.post(`camera/${this.id}/stream`, offer).subscribe((res: any) => {
            console.log('get stream', res)
            this.pc.setRemoteDescription(new RTCSessionDescription({type: "answer", sdp: res.data}))
            this.uid = res.uuid
        })
    }

}
