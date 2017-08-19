import { Injectable } from '@angular/core';

Injectable()
export class GameDataService {
    private _gameId: string;
    private _userName: string;
    private _gameType: number;
    private _guess: string;

    get gameId(): string {
        return this._gameId;
    }

    set gameId(gId: string) {
        this._gameId = gId;
    }

    get userName(): string {
        return this._userName;
    }

    set userName(un: string) {
        this._userName = un;
    }

    get gameType(): number {
        return this._gameType;
    }

    set gameType(gt: number) {
        this._gameType = gt;
    }

    get guess(): string {
        return this._guess;
    }

    set guess(g: string) {
        this._guess = g;
    }
}