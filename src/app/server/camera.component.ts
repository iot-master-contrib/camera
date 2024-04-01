import {Component} from '@angular/core';
import {CommonModule} from '@angular/common';
import {
    ParamSearch, RequestService, SmartTableButton, SmartTableColumn,
    SmartTableComponent, SmartTableOperator
} from "iot-master-smart";

@Component({
    selector: 'app-camera',
    standalone: true,
    imports: [
        CommonModule,
        SmartTableComponent,
    ],
    templateUrl: './camera.component.html',
    styleUrls: ['./camera.component.scss'],
})
export class CameraComponent {
    datum: any[] = [];
    total = 0;
    loading = false;

    buttons: SmartTableButton[] = [
        {icon: "plus", label: "创建", link: () => `/camera/create`}
    ];

    columns: SmartTableColumn[] = [
        {key: "id", sortable: true, label: "ID", keyword: true, link: (data) => `/camera/${data.id}`},
        {key: "name", sortable: true, label: "名称", keyword: true},
        {key: "port", sortable: true, label: "端口", keyword: true},
        {key: "created", sortable: true, label: "创建时间", date: true},
    ];

    operators: SmartTableOperator[] = [
        {icon: 'edit', title: '编辑', link: data => `/camera/${data.id}/edit`},
        {
            icon: 'delete', title: '删除', confirm: "确认删除？", action: data => {
                this.rs.get(`camera/${data.id}/delete`).subscribe(res => this.refresh())
            }
        },
    ];

    constructor(private rs: RequestService) {
    }


    query!: ParamSearch

    refresh() {
        this.search(this.query)
    }

    search(query: ParamSearch) {
        //console.log('onQuery', query)
        this.query = query
        this.loading = true
        this.rs.post('camera/search', query).subscribe((res) => {
            this.datum = res.data;
            this.total = res.total;
        }).add(() => this.loading = false);
    }

}
