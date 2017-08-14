import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http'
import { HttpHeaders } from '@angular/common/http'
import { BCResolver } from "./bc-resolver";

@Component({
    selector: 'play',
    templateUrl: './play.component.html',
    styleUrls: ['./play.component.css']
})
export class PlayComponent implements OnInit {
    private http: HttpClient;
    private route: ActivatedRoute;
    private router: Router;
    private engine: any;
    guess: string;
    name: string;
    greet: string;
    gameType: string;

    constructor(
        http: HttpClient,
        route: ActivatedRoute,
        router: Router
    ) {
        this.name = sessionStorage.getItem("name");
        this.gameType = sessionStorage.getItem("gameType");
        this.http = http;
        this.router = router;
        this.route = route;
        this.greet = `Now we are playing, ${this.name}!`;

        if (this.gameType === "2") {
            this.greet += ` Your browser is trying to guess ${sessionStorage.getItem("guess")}.`;
            this.startAutoPlay();
        }
    }

    ngOnInit(): void {
        this.http.get(
            `http://localhost:8081/api/game/${sessionStorage.getItem("gameID")}`
        )
            .subscribe(
            (data: any) => {
                this.handleGameResponse(data);
            },
            error => {
                this.handleError(error);
            });
    }

    makeGuess(): Promise<{ bulls, cows, win }> {
        let __this = this;
        if (this.guess.length < 4) {
            this.guess = this.genPrepZeroes(this.guess);
        }
        return new Promise(
            (resolve: (res: { bulls, cows, win }) => void, reject: (res: boolean) => void) => {
                __this.http.post(
                    `http://localhost:8081/api/guess/${this.guess}`,
                    `X-GameID=${sessionStorage.getItem('gameID')}`,
                    { headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded') }
                )
                    .subscribe(
                    (data: any) => {
                        this.handleGameResponse(data);
                        resolve({ bulls: data.p.bc.b, cows: data.p.bc.c, win: data.p.win });
                    },
                    error => {
                        reject(error);
                        this.handleError(error);
                    });
            });

    }

    private startAutoPlay() {
        let _this = this;
        let bcResolver = new BCResolver();
        const turnInterval = 750;
        _this.engine = setInterval(playTurn, turnInterval);
        function playTurn() {
            _this.guess = bcResolver.makeGuess();
            _this.makeGuess()
                .then((res: { bulls, cows, win }) => {
                    if (res.win === true) {
                        clearInterval(_this.engine);
                        bcResolver = null;
                        return;
                    }
                    clearInterval(_this.engine);
                    bcResolver.prune(_this.guess, res.bulls, res.cows);
                    _this.engine = setInterval(playTurn, turnInterval);
                })
                .catch((err: Error) => {
                    console.error(err);
                    alert(`${err.message}`);
                    clearInterval(_this.engine);
                });
        };
    }

    private handleError(e) {
        if (e.error) {
            alert(e.error.e)
        }
    }

    private handleGameResponse(data) {
        try {
            document.getElementById("win").innerHTML = ''

            if (data.p.win === true) {
                let newP: HTMLParagraphElement = document.createElement("p");
                newP.innerHTML =
                    `Game won! It took ${data.p.t} seconds and ${data.p.m.length} tries.`;
                document.getElementById("win").appendChild(newP);
            }

            let np = document.createElement("p");
            var style = document.createElement('style');
            style.type = 'text/css';
            style.innerHTML = '.strong { color: red; }';
            np.innerHTML = `<strong>${this.guess}</strong> got you <strong>${data.p.bc.b}</strong> bulls and <strong>${data.p.bc.c}</strong> cows`;
            np.className = "play-guess-res"
            document.getElementById("history").appendChild(np);
        } catch (e) {

        }
    }

    private genPrepZeroes(guess) {
        let pz = "";
        for (let i = 0; i < 4 - guess.length; i++) {
            pz += "0";
        }

        return pz + guess;
    }
}