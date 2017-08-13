import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { NameInitComponent } from './name-init.component';
import { PlayComponent } from './play.component'

const routes: Routes = [
  { path: '', redirectTo: '/', pathMatch: 'full' },
  { path: 'init/:id', component: NameInitComponent },
  { path: 'play', component: PlayComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }