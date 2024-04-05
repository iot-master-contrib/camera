import {Component, OnInit, ViewChild} from '@angular/core';
import {NzButtonComponent} from 'ng-zorro-antd/button';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {NzMessageService} from 'ng-zorro-antd/message';
import {CommonModule} from '@angular/common';
import {Router} from '@angular/router';
import {NzCardComponent} from "ng-zorro-antd/card";
import {RequestService, SmartEditorComponent, SmartField} from "iot-master-smart";

@Component({
    selector: 'app-camera-edit',
    standalone: true,
    imports: [
        CommonModule,
        NzButtonComponent,
        RouterLink,
        NzCardComponent,
        SmartEditorComponent,
    ],
    templateUrl: './camera-edit.component.html',
    styleUrls: ['./camera-edit.component.scss'],
})
export class CameraEditComponent implements OnInit {
    id: any = '';

    @ViewChild('form') form!: SmartEditorComponent

    fields: SmartField[] = [
        {key: "id", label: "ID", type: "text", min: 2, max: 30, placeholder: "选填"},
        {key: "name", label: "名称", type: "text", required: true, default: '新摄像头'},
        {
            key: "url", label: "链接", type: "text", required: true,
            default: 'rtsp://admin:admin@192.168.1.20:554/cam/realmonitor?channel=1&subtype=1'
        },
        {key: "project_id", label: "项目ID", type: "text"},
        {key: "streamer_id", label: "推流器ID", type: "text"},
        {key: "audio", label: "音频", type: "switch"},
        {key: "disabled", label: "禁用", type: "switch"},
        {key: "description", label: "说明", type: "textarea"},
    ]

    values: any = {}


    constructor(private router: Router,
                private msg: NzMessageService,
                private rs: RequestService,
                private route: ActivatedRoute
    ) {
    }

    ngOnInit(): void {
        if (this.route.snapshot.paramMap.has('id')) {
            this.id = this.route.snapshot.paramMap.get('id');
            this.load()
        }
    }

    load() {
        this.rs.get(`camera/` + this.id).subscribe((res) => {
            this.values = res.data
        });
    }

    onSubmit() {
        if (!this.form.valid) {
            this.msg.error('请检查数据')
            return
        }

        let url = `camera/${this.id || 'create'}`
        this.rs.post(url, this.form.value).subscribe((res) => {
            this.router.navigateByUrl('/camera/' + res.data.id);
            this.msg.success('保存成功');
        });
    }
}
