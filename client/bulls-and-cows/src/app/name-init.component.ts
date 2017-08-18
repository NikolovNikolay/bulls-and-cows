import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { HttpHeaders } from '@angular/common/http';
import { GameTypes } from './game-types';

@Component({
    selector: 'name-init',
    templateUrl: './name-init.component.html',
    styleUrls: ['./name-init.component.css']
})
export class NameInitComponent implements OnInit {
    private static get initURL(): string { return 'http://localhost:8080/api/init'; }

    private http: HttpClient;
    private route: ActivatedRoute;
    private router: Router;
    name: string;
    gameType: number;

    constructor(
        http: HttpClient,
        route: ActivatedRoute,
        router: Router
    ) {
        this.http = http;
        this.route = route;
        this.router = router;
        this.name = "";
    }

    ngOnInit(): void {
        // Parsing the query game type parameter
        this.route.params.subscribe(
            (params: Params) => {
                this.gameType = params["id"];
                if (this.gameType == 2) {
                    this.initGame();
                }
            });
    }

    public initGame() {
        let gt = this.gameType + "";

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
                sessionStorage.setItem("gameID", data.p.gameID);
                this.proceedPostInit(data)
            },
            error => {
                console.log(error);
            });
    }

    private forInitBodyParams(): string {
        return `userName=${this.name}&gameType=${this.gameType}`;
    }

    private proceedPostInit(data) {
        sessionStorage.setItem("name", data != null ? data.p.name : this.name);
        sessionStorage.setItem("gameType", this.gameType + "");

        if (this.gameType == parseInt(GameTypes.CVC)) {
            sessionStorage.setItem("guess", data.p.guess);
        }

        if (this.gameType == parseInt(GameTypes.P2P)) {
            this.router.navigateByUrl('/p2p');
        } else {
            this.router.navigateByUrl('/play');
        }
    }
}
