/***
    PING PONG GAME:
**/

const canvas = document.getElementById("pingpongWindow");
const ctx = canvas.getContext('2d');
drawGame();

function main(){
    /*
        Get player id and gameid after submit
        hide the inputs
    */
    var playerid = document.getElementById("playerid").value;
    var gameid = document.getElementById("gameid").value;
    //connectToServer(gameid, playerid);
}

function connectToServer(gameid, playerid) {
    
    var ws;
    var loc = window.location, new_uri;
    if (loc.protocol === "https:") {
        new_uri = "wss:";
    }
    else {
        new_uri = "ws:";
    }
    new_uri += "//" + loc.host;
    new_uri += loc.pathname + "ws?gameid=" + gameid + "&playerid=" + playerid;
    ws = new WebSocket(new_uri);
    ws.onopen = function (evt) {
        console.log("OPEN");
    };
    ws.onclose = function (evt) {
        console.log("CLOSE");
        ws = null;
    };
    ws.onmessage = function (evt) {
        console.log("RESPONSE: " + evt.data);
    };
    ws.onerror = function (evt) {
        console.log("ERROR: " + evt.data);
    };
}



// draw circle, will be used to draw the ball
function drawArc(x, y, r, color){
    ctx.fillStyle = color;
    ctx.beginPath();
    ctx.arc(x,y,r,0,Math.PI*2,true);
    ctx.closePath();
    ctx.fill();
}


function drawPlayer(color, x,y,w,h){
    ctx.fillStyle = color;
    ctx.fillRect(x, y, w, h);
}

function drawField(color, w,h){

}

function drawGame(){
    /****
    * Draw players and the game board
    *  
    * */ 
   var playerHeight = 10;
   var playerWidth = 100;

   drawPlayer("blue", 0,0,playerHeight,playerWidth);
   drawPlayer("red", canvas.width-playerHeight, 0,playerHeight,playerWidth);
    
}
