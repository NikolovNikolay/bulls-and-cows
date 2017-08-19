import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { HttpClientModule } from '@angular/common/http';
import { AppRoutingModule } from './app-routing.module';
import { CommonModule } from '@angular/common';

import { AppComponent } from './app.component';
import { NameInitComponent } from './name-init.component';
import { PlayComponent } from './play.component'
import { RTPlayComponent } from './rt-play.component'
import { GameDataService } from './game-data.service';

@NgModule({
  declarations: [
    AppComponent,
    NameInitComponent,
    PlayComponent,
    RTPlayComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    HttpModule,
    HttpClientModule,
    CommonModule
  ],
  providers: [GameDataService],
  bootstrap: [AppComponent]
})
export class AppModule { }
