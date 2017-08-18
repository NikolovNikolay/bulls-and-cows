# bulls-and-cows

Server code is in **'server'** folder

Client code is in **'client/bulls-and-cows'** folder

Server is configured to run on localhost, port **8080**

Client is configured to run on localhost, port **4200**

If some port configuration is required, please refer to:
- change server port: `server/cmd/srv/main.go` and update `servePort` constant
- change client port: `client/bulls-and-cows/protractor.conf.js` and update `baseUrl` property

The server app depends on MongoDB with default configuration, listening on port `27017`

Know issues:
- The auto reolve algorithm in **Browser vs. Computer** mode crashes sometimes for still some unclear circumstances
