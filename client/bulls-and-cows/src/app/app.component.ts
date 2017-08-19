import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { GameDataService } from './game-data.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: []
})
export class AppComponent {
  title = 'Bulls and cows game';
  service: GameDataService;

  constructor(
    gameDataService: GameDataService,
    private router: Router
  ) {
    this.service = gameDataService;
  }

  public selectGame(gameType: number) {
    this.service.gameId = "";
    this.service.guess = "";
    this.service.userName = "";
    this.service.gameType = gameType;
    this.router.navigateByUrl(`/init`);
  }

  public reset() {
    this.service.gameId = "";
    this.service.guess = "";
    this.service.userName = "";
    this.router.navigateByUrl(`/`);
  }
}
