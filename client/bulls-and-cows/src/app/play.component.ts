import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http'
import { HttpHeaders } from '@angular/common/http'

@Component({
    selector: 'play',
    templateUrl: './play.component.html',
    styleUrls: ['./play.component.css']
})
export class PlayComponent implements OnInit {
    private http: HttpClient;
    private route: ActivatedRoute;
    private router: Router;
    guess: string;
    name: string;

    constructor(
        http: HttpClient,
        route: ActivatedRoute,
        router: Router
    ) {
        this.name = sessionStorage.getItem("name");
        this.http = http;
        this.router = router;
        this.route = route;
    }

    ngOnInit(): void {
        this.http.get(
            `http://localhost:8081/api/game/${sessionStorage.getItem("gameID")}`
        )
            .subscribe(
            (data: any) => {
                this.handleGameData(data);
            },
            error => {
                this.handleError(error);
            });
    }

    makeGuess() {
        this.http.post(
            `http://localhost:8081/api/guess/${this.guess}`,
            `X-GameID=${sessionStorage.getItem('gameID')}`,
            { headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded') }
        )
            .subscribe(
            (data: any) => {
                this.handleGameData(data);
            },
            error => {
                this.handleError(error);
            });
    }

    private handleError(e) {
        if (e.error) {
            alert(e.error.e)
        }
    }

    private handleGameData(data) {
        document.getElementById("bandc").innerHTML = ''
        document.getElementById("history").innerHTML = '';

        let newP: HTMLParagraphElement = document.createElement("p");
        if (data.p.win === false && data.p.bc) {
            newP.innerHTML =
                `You got: ${data.p.bc.b} Bulls and ${data.p.bc.c} Cows`;
        } else if (data.p.win === true) {
            newP.innerHTML =
                `Congrats, you won! You have 4 Bulls with ${data.p.m[data.p.m.length - 1]}.
                        It took you ${data.p.t} seconds. You can try again if you like :)`;
        }
        document.getElementById("bandc").appendChild(newP);

        data.p.m.forEach(m => {
            let np = document.createElement("p");
            np.innerHTML = m;
            document.getElementById("history").appendChild(np);
        });
    }
}