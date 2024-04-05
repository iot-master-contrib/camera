import {Routes} from '@angular/router';
import {UnknownComponent} from "iot-master-smart";

import {IndexComponent} from "./index/index.component";
import {CamerasComponent} from "./camera/cameras/cameras.component";
import {CameraEditComponent} from './camera/camera-edit/camera-edit.component';
import {CameraDetailComponent} from './camera/camera-detail/camera-detail.component';
import {PlayerComponent} from "./player/player.component";
import {StreamersComponent} from "./streamer/streamers/streamers.component";
import {StreamerEditComponent} from "./streamer/streamer-edit/streamer-edit.component";
import {StreamerDetailComponent} from "./streamer/streamer-detail/streamer-detail.component";

export const routes: Routes = [
    {path: '', pathMatch: "full", component: IndexComponent},

    {path: 'camera', component: CamerasComponent},
    {path: 'camera/create', component: CameraEditComponent},
    {path: 'camera/:id', component: CameraDetailComponent},
    {path: 'camera/:id/edit', component: CameraEditComponent},

    {path: 'streamer', component: StreamersComponent},
    {path: 'streamer/create', component: StreamerEditComponent},
    {path: 'streamer/:id', component: StreamerDetailComponent},
    {path: 'streamer/:id/edit', component: StreamerEditComponent},

    {path: 'play', component: PlayerComponent},
    {path: 'play/:id', component: PlayerComponent},

    {path: '**', component: UnknownComponent},
];
