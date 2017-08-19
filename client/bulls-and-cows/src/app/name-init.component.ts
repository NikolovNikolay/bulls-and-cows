import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { HttpHeaders } from '@angular/common/http';
import { GameTypes } from './game-types';

import { GameDataService } from './game-data.service';

@Component({
    selector: 'name-init',
    templateUrl: './name-init.component.html',
    styleUrls: ['./name-init.component.css'],
})
export class NameInitComponent {
    private static get initURL(): string { return 'http://localhost:8080/api/init'; }

    name: string;

    constructor(
        private gameDataService: GameDataService,
        private http: HttpClient,
        private route: ActivatedRoute,
        private router: Router
    ) {
        this.http = http;
        this.route = route;
        this.router = router;
        this.name = "";
    }

    public initGame() {
        this.gameDataService.userName = this.name;
        let gt = this.gameDataService.gameType;

        // If game types is peer 2 peer, then
        // navigate p2p root
        if (gt === GameTypes.P2P) {
            return this.proceedPostInit(null);
        }

        // Make an init request to server
        this.http.post(
            NameInitComponent.initURL,
            this.forInitBodyParams(),
            { headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded') }
        )
            .subscribe(
            (data: any) => {
                console.log(data);
                this.gameDataService.gameId = data.p.gameID;
                this.proceedPostInit(data)
            },
            error => {
                console.log(error);
            });
    }

    private forInitBodyParams(): string {
        return `userName=${this.gameDataService.userName}&gameType=${this.gameDataService.gameType}`;
    }

    private proceedPostInit(data) {

        if (this.gameDataService.gameType == GameTypes.CVC) {
            this.gameDataService.guess = data.p.guess;
        }

        if (this.gameDataService.gameType == GameTypes.P2P) {
            this.router.navigateByUrl('/p2p');
        } else {
            this.router.navigateByUrl('/play');
        }
    }
}
