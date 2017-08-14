import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http'
import { HttpHeaders } from '@angular/common/http'

@Component({
    selector: 'name-init',
    templateUrl: './name-init.component.html',
    styleUrls: ['./name-init.component.css']
})
export class NameInitComponent implements OnInit {
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
    }

    ngOnInit(): void {
        this.name = "";
        this.route.params.subscribe(
            (params: Params) => {
                this.gameType = params["id"];
                if (this.gameType == 2) {
                    this.initGame();
                }
            }
        );
    }

    initGame() {
        let gt = this.gameType + "";
        let s = new URLSearchParams();
        s.append("userName", this.name);
        s.append("gameType", gt);

        this.http.post(
            "http://localhost:8081/api/init",
            `userName=${this.name}&gameType=${this.gameType}`,
            { headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded') }
        )
            .subscribe(
            (data: any) => {
                console.log(data);
                sessionStorage.setItem("gameID", data.p.gameID);
                sessionStorage.setItem("name", data.p.name);
                sessionStorage.setItem("gameType", gt);
                if (this.gameType == 2) {
                    sessionStorage.setItem("guess", data.p.guess);
                }
                this.router.navigateByUrl('/play');
            },
            error => {
                console.log(error);
            });
    }
}
