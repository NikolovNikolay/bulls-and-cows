import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { NameInitComponent } from './name-init.component';
import { PlayComponent } from './play.component';
import { RTPlayComponent } from './rt-play.component';

const routes: Routes = [
  { path: '', redirectTo: '/', pathMatch: 'full' },
  { path: 'init', component: NameInitComponent },
  { path: 'play', component: PlayComponent },
  { path: 'p2p', component: RTPlayComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }