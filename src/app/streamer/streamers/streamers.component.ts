import {Component} from '@angular/core';
import {CommonModule} from '@angular/common';
import {
    ParamSearch, RequestService, SmartTableButton, SmartTableColumn,
    SmartTableComponent, SmartTableOperator
} from "iot-master-smart";

@Component({
    selector: 'app-streamers',
    standalone: true,
    imports: [
        CommonModule,
        SmartTableComponent,
    ],
    templateUrl: './streamers.component.html',
    styleUrls: ['./streamers.component.scss'],
})
export class StreamersComponent {
    datum: any[] = [];
    total = 0;
    loading = false;

    buttons: SmartTableButton[] = [
        {icon: "plus", label: "创建", link: () => `/streamer/create`}
    ];

    columns: SmartTableColumn[] = [
        {key: "id", sortable: true, label: "ID", keyword: true, link: (data) => `/streamer/${data.id}`},
        {key: "name", sortable: true, label: "名称", keyword: true},
        {key: "created", sortable: true, label: "创建时间", date: true},
    ];

    operators: SmartTableOperator[] = [
        {icon: 'play-square', title: '播放', link: data => `/play/${data.id}`},
        {icon: 'edit', title: '编辑', link: data => `/streamer/${data.id}/edit`},
        {
            icon: 'delete', title: '删除', confirm: "确认删除？", action: data => {
                this.rs.get(`streamer/${data.id}/delete`).subscribe(res => this.refresh())
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
        this.rs.post('streamer/search', query).subscribe((res) => {
            this.datum = res.data;
            this.total = res.total;
        }).add(() => this.loading = false);
    }

}
