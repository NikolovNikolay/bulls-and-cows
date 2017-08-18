import { Component, OnInit, Input } from '@angular/core';
import { ActivatedRoute, ParamMap, Params, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http'
import { HttpHeaders } from '@angular/common/http'
import { BCResolver } from "./bc-resolver";
import { GameTypes } from './game-types';
import { Socket } from './socket';
import * as io from 'socket.io-client';

@Component({
	selector: 'rtPlay',
	templateUrl: './rt-play.component.html',
	styleUrls: ['./rt-play.component.css']
})
export class RTPlayComponent {
	socket: any;
	greet: string;
	name: string
	gameType: string;
	join: boolean;
	guess: string;

	constructor() {
		this.name = sessionStorage.getItem("name");
		this.gameType = sessionStorage.getItem("gameType");
		this.greet = `Welcome, ${this.name}`;
		this.configureSocketIO();
	}

	public prepareHost() {
		this.socket.emit(Socket.evtCreateRoom, this.name, (data) => {
			console.log(data);
			if (data) {
				document.getElementById('p2p-btns').innerHTML = '';
				document.getElementById('active-games').innerHTML = '';
				this.greet = `${this.name}, you are waiting for someone to join`;
			} else {
				console.log('It seems that you have already hosted a game');
			}
		});
	}

	public prepareJoin() {
		let roomName = prompt("Please input the name of your rival");
		this.socket.emit(Socket.evtJoinRoom, roomName, this.name, (data) => {
			console.log(data);
			if (data) {
				document.getElementById('p2p-btns').innerHTML = '';
				document.getElementById('active-games').innerHTML = '';
				this.greet = this.formActionMsg(this.name, roomName);
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
		this.socket.emit(Socket.evtInputGuess, this.guess, this.name, (data) => {
			// console.log(data);
			// if (data) {
			//     document.getElementById('play-guess-container').innerHTML = '';
			//     this.greet = `${this.name}, you are waiting for someone to join`;
			// } else {
			//     console.log('It seems that you have already hosted a game');
			// }
		});
	}

	private configureSocketIO() {
		let self = this;
		this.socket = io(Socket.defaultURL);
		this.socket.on(Socket.evtConnect, function () {
			self.socket = this;
			self.socket.emit(Socket.evtGetAvailableRooms, self.name, (rooms) => {
				self.appendRooms(rooms);
			});
			self.socket.on(Socket.evtJoinedMyGroup, function (rival) {
				self.greet = self.formActionMsg(self.name, rival);
				self.join = true;
			});
			self.socket.on(Socket.evtUpdateRooms, function (rooms) {
				self.appendRooms(rooms);
			});
			self.socket.on(Socket.evtConfirmJoin, function () {
				self.confirmJoin();
			});
		});
	}

	private appendRooms(rooms) {
		rooms.forEach(r => {
			if (r === "") return;
			let newP = document.createElement('p')
			newP.innerHTML = r;
			document.getElementById('active-games').appendChild(newP)
		});
	}

	private formActionMsg(name: string, rival: string): string {
		return `Now we are playing, ${name}! Your rival is ${rival}.`;
	}
}