import {Routes} from '@angular/router';
import {UnknownComponent} from "iot-master-smart";
import {IndexComponent} from "./index/index.component";
import {CameraComponent} from './camera/camera.component';
import {CameraEditComponent} from './camera-edit/camera-edit.component';
import {CameraDetailComponent} from './camera-detail/camera-detail.component';

export const routes: Routes = [
    {path: '', pathMatch: "full", component: IndexComponent},

    {path: 'camera', component: CameraComponent},
    {path: 'camera/create', component: CameraEditComponent},
    {path: 'camera/:id', component: CameraDetailComponent},
    {path: 'camera/:id/edit', component: CameraEditComponent},

    {path: '**', component: UnknownComponent},
];
