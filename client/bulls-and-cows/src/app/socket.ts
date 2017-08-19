
// Socket is a wrapper over the socket.io lib
export class Socket {

	public static get defaultURL(): string { return 'http://localhost:8080/socket.io'; }

	public static get evtConnect(): string { return 'connect'; }
	public static get evtDisconnect(): string { return 'disconnect'; }
	public static get evtError(): string { return 'error'; }


	public static get evtGetAvailableRooms(): string { return 'getavr'; }
	public static get evtJoinedMyGroup(): string { return 'joinmy'; }
	public static get evtConfirmJoin(): string { return 'confjoin'; }
	public static get evtUpdateRooms(): string { return 'updater'; }
	public static get evtCreateRoom(): string { return 'creater'; }
	public static get evtJoinRoom(): string { return 'joinr'; }
	public static get evtInputGuess(): string { return 'inputguess'; }
	public static get evtStartP2P(): string { return 'startp2p'; }
	public static get evtMakeGuess(): string { return 'makeguess'; }
	public static get evtGameEnd(): string { return 'gameend'; }
}