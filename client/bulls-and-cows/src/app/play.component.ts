import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http'
import { HttpHeaders } from '@angular/common/http'
import { BCResolver } from "./bc-resolver";
import { GameTypes } from './game-types'
import * as io from 'socket.io-client';
import { GameDataService } from "./game-data.service";

/**
 * Component initialized when Player vs. Computer
 * and Computer vs. Computer was selected 
 */
@Component({
    selector: 'play',
    templateUrl: './play.component.html',
    styleUrls: ['./play.component.css']
})
export class PlayComponent implements OnInit {
    private static get guessURL(): string { return `http://localhost:8080/api/guess/`; }
    private static get gameDataURL(): string { return `http://localhost:8080/api/game/`; }
    private static get autoPlayTurnInterval(): number { return 1500; }

    private engine: any;
    private lastGuess: string;

    guess: string;
    name: string;
    greet: string;
    gameType: number;

    constructor(
        private http: HttpClient,
        private route: ActivatedRoute,
        private router: Router,
        protected gameDataService: GameDataService
    ) {
        this.name = this.gameDataService.userName;
        this.gameType = this.gameDataService.gameType;
        this.guess = this.gameDataService.guess;
        this.http = http;
        this.router = router;
        this.route = route;

        if (this.gameType === GameTypes.PVC) {
            this.greet = `Now we are playing, ${this.name}!`;
        } else if (this.gameType === GameTypes.CVC) {
            this.greet = `${this.name}, your browser is trying to guess ${this.guess}.`;
            this.startAutoPlay();
        }
    }

    ngOnInit(): void {
        this.http.get(this.formGameDataURL())
            .subscribe(
            (data: any) => {
                this.handleGameResponse(data.p);
            },
            error => {
                this.handleError(error);
            });
    }

    // Called when a player makes a guess (player or browser bot)
    public makeGuess(guess: string): Promise<{ bulls, cows, win }> {
        let self = this;

        // If the guess is empty - don't do anything
        if (guess == null || guess === "") {
            return;
        }

        // make a valid guess if less chars are detected
        if (guess.length < 4) {
            guess = this.guess = this.genPrepZeroes(guess);
        }

        return new Promise(
            (resolve: (res: { bulls, cows, win }) => void, reject: (res: boolean) => void) => {
                self.http.put(
                    this.formMakeGuessURL(guess),
                    this.formMakeGuessParams(),
                    { headers: new HttpHeaders().set('Content-Type', 'application/x-www-form-urlencoded') })
                    .subscribe(
                    (data: any) => {
                        this.handleGameResponse(data.p);
                        resolve({ bulls: data.p.bc.b, cows: data.p.bc.c, win: data.p.win });
                    },
                    error => {
                        reject(error);
                        this.handleError(error);
                    });
            });
    }

    // Auto called if Computer vs. Computer game
    // was selected
    private startAutoPlay() {
        let self = this;
        let bcResolver = new BCResolver();

        self.engine =
            setInterval(playTurn, PlayComponent.autoPlayTurnInterval, bcResolver);

        function playTurn(resolver: BCResolver) {
            self.guess = resolver.makeGuess();
            if (self.guess == self.lastGuess) {
                clearInterval(self.engine);
                alert("Something messed up. Please try again!");
                return;
            }

            self.makeGuess(self.guess)
                .then((res: { bulls, cows, win }) => {
                    if (res.win === true) {
                        clearInterval(self.engine);
                        resolver = null;
                        return;
                    }

                    clearInterval(self.engine);
                    resolver.prune(self.guess, res.bulls, res.cows);
                    self.engine = setInterval(playTurn, PlayComponent.autoPlayTurnInterval, bcResolver);
                    self.lastGuess = self.guess;
                })
                .catch((err: Error) => {
                    console.error(err);
                    alert(`${err.message}`);
                    clearInterval(self.engine);
                });
        };
    }

    // Some helper methods for visualizing data
    
    protected handleGameResponse(data: any) {
        try {
            document.getElementById("win").innerHTML = '';

            if (data.win === true) {
                let newP: HTMLParagraphElement = document.createElement('p');
                newP.innerHTML =
                    `Game won! It took ${data.t} seconds and ${data.m.length} tries.`;
                document.getElementById('win').appendChild(newP);
            }

            if (data.bc == null) {
                data.m.forEach(g => {
                    this.formNewHistoryDoc(g.g, g);
                });
            } else {
                this.formNewHistoryDoc(this.guess, data);
            }
        } catch (e) {
            this.handleError(e)
        }
    }

    private formGameDataURL() {
        return `${PlayComponent.gameDataURL}${this.getGameId()}`;
    }

    private formMakeGuessURL(guess: string): string {
        return `${PlayComponent.guessURL}${this.guess}`;
    }

    private formMakeGuessParams(): string {
        return `gameID=${this.getGameId()}`;
    }

    private getGameId(): string {
        return this.gameDataService.gameId;
    }

    protected formBCString(guess: string, bulls: number, cows: number): string {
        return `<strong>${guess}</strong> got you <strong>${bulls}</strong> bulls and <strong>${cows}</strong> cows`;
    }

    private formNewHistoryDoc(guess: string, payload: any) {
        let np = this.genHistoryElement();
        np.innerHTML = this.formBCString(guess, payload.bc.b, payload.bc.c);
        document.getElementById("history").appendChild(np);
    }

    protected handleError(e) {
        if (e.error) {
            alert(e.error.e)
        }
    }

    protected genHistoryElement(): HTMLParagraphElement {
        let np = document.createElement("p");
        np.className = "play-guess-res";

        return np;
    }

    private genPrepZeroes(guess: string): string {
        let pz = "";
        for (let i = 0; i < 4 - guess.length; i++) {
            pz += "0";
        }

        return pz + guess;
    }
}