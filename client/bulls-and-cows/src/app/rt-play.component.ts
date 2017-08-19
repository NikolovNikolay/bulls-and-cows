import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http'
import { HttpHeaders } from '@angular/common/http'
import { BCResolver } from "./bc-resolver";
import { GameTypes } from './game-types';
import { Socket } from './socket';
import * as io from 'socket.io-client';
import { GameDataService } from "./game-data.service";
import { PlayComponent } from "./play.component";

/**
 * Component initialized when real-time 
 * Peer 2 Peer game was selected
 */
@Component({
	selector: 'rtPlay',
	templateUrl: './rt-play.component.html',
	styleUrls: ['./rt-play.component.css']
})
export class RTPlayComponent extends PlayComponent {
	private socket: any;
	private ownNumberSet: boolean;

	greet: string;
	name: string
	gameType: number;
	join: boolean;
	ownNumber: string;
	roomName: string;
	playing: boolean;
	guess: string;
	turnMsg: string
	rivalName: string;
	onTurn: boolean;
	gameEnd: boolean;


	constructor(
		http: HttpClient,
		route: ActivatedRoute,
		router: Router,
		service: GameDataService
	) {
		super(http, route, router, service);
		this.name = this.gameDataService.userName;
		this.gameType = this.gameDataService.gameType;
		this.greet = `Welcome, ${this.name}`;
		this.configureSocketIO();
	}

	// Attaches the custom event listeners to
	// the socket instance
	private configureSocketIO() {
		let self = this;
		this.socket = io(Socket.defaultURL);
		this.socket.on(Socket.evtConnect, function () {
			self.socket = this;
			self.socket.emit(Socket.evtGetAvailableRooms, self.name, (rooms) => {
				self.appendRooms(rooms);
			});
			self.socket.on(Socket.evtJoinedMyGroup, function (rival) {
				self.rivalName = rival;
				self.greet = self.formActionMsg(self.name, self.rivalName);
				self.join = true;
			});
			self.socket.on(Socket.evtUpdateRooms, function (rooms) {
				self.appendRooms(rooms);
			});
			self.socket.on(Socket.evtConfirmJoin, function () {
				self.confirmJoin();
			});
			self.socket.on(Socket.evtStartP2P, function () {
				self.startP2P();
			});
			self.socket.on(Socket.evtGameEnd, function () {
				self.gameEnds();
			});
		});
	}

	// Called when the user hosts a game.
	// It receives the corresponding game ID
	// as an ack result
	public prepareHost() {
		this.roomName = this.name;
		this.socket.emit(Socket.evtCreateRoom, this.roomName, this.gameType, (data) => {
			console.log(data);
			if (data) {
				document.getElementById('p2p-btns').innerHTML = '';
				document.getElementById('active-games').innerHTML = '';
				this.greet = `${this.roomName}, you are waiting for someone to join`;
				this.gameDataService.gameId = data;
			} else {
				console.log('It seems that you have already hosted a game');
			}
		});
	}

	// Called when the user joins a game
	// It receives the corresponding game ID
	// that the user joined as an ack result
	public prepareJoin() {
		this.roomName = prompt("Please input the name of your rival");
		this.socket.emit(Socket.evtJoinRoom, this.roomName, this.name, (data) => {
			console.log(data);
			if (data) {
				document.getElementById('p2p-btns').innerHTML = '';
				document.getElementById('active-games').innerHTML = '';
				this.rivalName = this.roomName;
				this.greet = this.formActionMsg(this.name, this.rivalName);
				this.join = true;
				this.gameDataService.gameId = data;
			} else {
				console.log('Could not join the game. Sorry...');
			}
		});
	}

	// Called from server when a join was confirmed
	public confirmJoin() {
		this.join = true;
	}

	// Called when the user inputs his own guess number
	// and sends it to the server
	public inputGuess() {
		if (!this.ownNumberSet) {
			this.socket.emit(Socket.evtInputGuess, this.roomName, this.ownNumber, this.name, (success) => {
				this.greet += ` Your number is ${this.ownNumber}`;
				this.ownNumberSet = true;
				if (success) {
					this.startP2P()
				}
			});
		}
	}

	// Called from server when both players have set
	// their numbers to guess and the game can start
	public startP2P() {
		this.playing = true;
		this.join = false;
	}

	// Called when the player calls the server with a
	// number, trying to guess the other player's num
	public makeRTGuess() {
		if (!this.gameEnd) {
			this.socket.emit(
				Socket.evtMakeGuess,
				this.roomName,
				this.guess,
				this.name,
				this.gameDataService.gameId,
				(data) => {
					if (data != null) {
						this.handleGameResponse(data);
						if (data.win) {
							this.gameEnd = true
						}
					}
				});
		}
	}

	// Called from the server when someone guessed 
	// the other player's number
	public gameEnds() {
		document.getElementById("results").innerHTML = '';
		document.getElementById("play-guess-container").innerHTML = '';
		document.getElementById("play-guess-container").innerHTML = `Game ended! You have lost!`;
	}

	// Some helper methods to visualize data

	private appendRooms(rooms) {
		try {
			document.getElementById('active-games').innerHTML = '';
			rooms.forEach(r => {
				if (r === "") return;
				let newP = document.createElement('p')
				newP.innerHTML = r;
				document.getElementById('active-games').appendChild(newP);
			});
		} catch (e) { }
	}

	private formActionMsg(name: string, rival: string): string {
		return `Now we are playing, ${name}! Your rival is ${rival}.`;
	}
}