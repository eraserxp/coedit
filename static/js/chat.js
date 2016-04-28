/**
 * Created by yaoyaolin on 4/27/16.
 */
var app = require('express')();
var server = require('http').Server(app);
var io = require('socket.io')(server);

//app.get('/', function(req, res){
//    res.sendfile('index.html');
//});

var connectionList ={};
var chatRooms = {};


io.sockets.on('connection',function(socket){
    var socketId = socket.id;

    connectionList[socketId]= {
        socket: socket
    }

    socket.on('join',function(data){
        console.log("one user join " + data.username);
        socket.username = data.username;
        socket.roomUrl = data.roomUrl;

        if (!chatRooms[data.roomUrl]){
            chatRooms[data.roomUrl] = {};
        }
        chatRooms[data.roomUrl][socketId] = socket;
        console.log( chatRooms[data.roomUrl]);
        for ( var sId in chatRooms[data.roomUrl]){
            chatRooms[data.roomUrl][sId].emit('broadcast_join', {
               username : chatRooms[data.roomUrl][sId].username
            });
        }
    });

    socket.on('disconnect',function() {
        console.log("one user quit");
        for (var sId in chatRooms[socket.roomUrl]) {
            chatRooms[socket.roomUrl][sId].emit('broadcast_quit', {
                username: chatRooms[socket.roomUrl][sId].username
            });
        }

        if (chatRooms[socket.roomUrl]) {
            delete chatRooms[socket.roomUrl][socket.id];
        }
    });

    socket.on('say',function(data){
        console.log("one user say " + data.roomUrl );
        for ( var sId in chatRooms[data.roomUrl]){
            console.log(sId + " is saying " + data.text);
            chatRooms[data.roomUrl][sId].emit('broadcast_say', {
                username : data.username,
                text: data.text
            });
        }
    });
});


server.listen(8110, function(){
    console.log('linstening on 8110');
});




