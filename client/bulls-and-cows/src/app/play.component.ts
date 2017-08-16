import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http'
import { HttpHeaders } from '@angular/common/http'
import { BCResolver } from "./bc-resolver";
import * as io from 'socket.io-client';

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
    private socket: any;
    private clientSock: any;

    pattern = /^\d+$/;
    guess: string;
    name: string;
    greet: string;
    gameType: string;
    doubleJoin: boolean;
    join: boolean;

    constructor(
        http: HttpClient,
        route: ActivatedRoute,
        router: Router
    ) {
        let __this = this;
        this.name = sessionStorage.getItem("name");
        this.gameType = sessionStorage.getItem("gameType");
        this.http = http;
        this.router = router;
        this.route = route;
        this.doubleJoin = false;
        this.join = false;

        if (this.gameType === "1") {
            this.greet = `Now we are playing, ${this.name}!`;
        } else if (this.gameType === "2") {
            this.greet += ` Your browser is trying to guess ${sessionStorage.getItem("guess")}.`;
            this.startAutoPlay();
        } else if (this.gameType === "3") {
            this.greet = `Welcome, ${this.name}`;
            this.socket = io('http://localhost:8080/socket.io');
            this.socket.on('connect', function () {
                __this.socket = this;
                __this.socket.emit("getavr", __this.name, (rooms) => {
                    __this.appendRooms(rooms);
                });
                __this.socket.on("joinmy", function (rival) {
                    __this.greet = `Now we are playing, ${__this.name}! Your rival is ${rival}.`;
                    __this.join = true;
                });
                __this.socket.on("updater", function (rooms) {
                    __this.appendRooms(rooms);
                });
                __this.socket.on("confjoin", function () {
                    __this.confirmJoin();
                });
            });
        }
    }

    ngOnInit(): void {
        if (this.gameType !== "3") {
            this.http.get(
                `http://localhost:8080/api/game/${sessionStorage.getItem("gameID")}`
            )
                .subscribe(
                (data: any) => {
                    this.handleGameResponse(data);
                },
                error => {
                    this.handleError(error);
                });
        }
    }

    private appendRooms(rooms) {
        rooms.forEach(r => {
            if (r === "") return;
            let newP = document.createElement('p')
            newP.innerHTML = r;
            document.getElementById('active-games').appendChild(newP)
        });
    }

    makeGuess(): Promise<{ bulls, cows, win }> {
        let __this = this;
        if (this.guess.length < 4) {
            this.guess = this.genPrepZeroes(this.guess);
        }

        return new Promise(
            (resolve: (res: { bulls, cows, win }) => void, reject: (res: boolean) => void) => {
                __this.http.put(
                    `http://localhost:8080/api/guess/${this.guess}`,
                    `gameID=${sessionStorage.getItem('gameID')}`,
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

    public prepareHost() {
        this.socket.emit("creater", this.name, (data) => {
            console.log(data);
            if (data) {
                document.getElementById('play-guess-container').innerHTML = '';
                this.greet = `${this.name}, you are waiting for someone to join`;
            } else {
                console.log('It seems that you have already hosted a game');
            }
        });
    }

    public prepareJoin() {
        let roomName = prompt("Please input the name of your rival");
        this.socket.emit("joinr", roomName, this.name, (data) => {
            console.log(data);
            if (data) {

                document.getElementById('p2p-btns').innerHTML = '';
                document.getElementById('active-games').innerHTML = '';
                this.greet = `Now we are playing, ${this.name}! Your rival is ${roomName}`;
                this.join = true;
            } else {
                console.log('Could not join the game. Sorry...');
            }
        });
    }

    public confirmJoin() {
        this.join = true;
    }

    public inputGuess() {
        this.socket.emit("inputguess", this.guess, this.name, (data) => {
            // console.log(data);
            // if (data) {
            //     document.getElementById('play-guess-container').innerHTML = '';
            //     this.greet = `${this.name}, you are waiting for someone to join`;
            // } else {
            //     console.log('It seems that you have already hosted a game');
            // }
        });
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

            if (data.p.bc == null) {
                data.p.m.forEach(g => {
                    let np = this.genHistoryElement();
                    np.innerHTML = `<strong>${g.g}</strong> got you <strong>${g.bc.b}</strong> bulls and <strong>${g.bc.c}</strong> cows`;
                    document.getElementById("history").appendChild(np);
                });
            } else {
                let np = this.genHistoryElement();
                np.innerHTML = `<strong>${this.guess}</strong> got you <strong>${data.p.bc.b}</strong> bulls and <strong>${data.p.bc.c}</strong> cows`;
                document.getElementById("history").appendChild(np);
            }
        } catch (e) {

        }
    }

    private genHistoryElement(): HTMLParagraphElement {
        let np = document.createElement("p");
        np.className = "play-guess-res";

        return np;
    }

    private genPrepZeroes(guess) {
        let pz = "";
        for (let i = 0; i < 4 - guess.length; i++) {
            pz += "0";
        }

        return pz + guess;
    }
}