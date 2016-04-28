<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Collaborative editing platform</title>
    <link rel="stylesheet" href="../static/css/index.css" />
    <link rel="stylesheet" href="../static/bootstrap/css/bootstrap.css" />
    <link rel="stylesheet" href="../static/css/bootstrap-social.css" />
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.5.0/css/font-awesome.min.css" />

    <style type="text/css" media="screen">
        #editor {
            position: relative;
            width: 95%;
            height: 100%;
            left: 2.5%;
            border: 1px solid black;
        }
    </style>
</head>
<body>
<div class="header">
    <div class="icon iconleft">
        COEDITOR
    </div>
    <div class="iconright">
        <a href="/auth/gplus" class="btn btn-lg btn-social-icon btn-google">
            <span class="fa fa-google"></span>
        </a>

        <a class="btn btn-lg btn-social-icon btn-facebook">
            <span class="fa fa-facebook"></span>
        </a>

        <a class="btn btn-lg btn-social-icon btn-github">
            <span class="fa fa-github"></span>
        </a>

        <a class="btn btn-lg btn-social-icon btn-dropbox">
            <span class="fa fa-dropbox"></span>
        </a>
    </div>
</div>
<div class="body">
    <div class="info">
        <div class="editor">
            <div id="editor"></div>
        </div>


        <div class="funccolumn">
            <div class="selectcolumn">
                Language: <select id="mode" onchange="changeLan()">
                <option value="javascript">JavaScript</option>
                <option value="xml">XML</option>
                <option value="python">Python</option>
                <option value="csharp">C#</option>
            </select>
            </div>
            <div class="selectcolumn">
                Theme: <select id="theme" onchange="changeTheme()">
                <option value="tomorrow">Tomorrow</option>
                <option value="twilight">Twilight</option>
                <option value="chrome">Chrome</option>
            </select>
            </div>
            <div class="selectcolumn">

            </div>
        </div>
    </div>
    <div class="chatter">
        Your name:<br/>
        <textarea id="nametext">Anon</textarea>
        <br/><br/>
        Message log:<br/>
        <textarea id="messages" rows="15" cols="15"></textarea>
        <br/><br/>
        Send message: <br/>
        <form action="">
            <input id="sendtext" autocomplete="off" /><button>Send</button>
        </form>
    </div>
</div>
<script src="../static/ace/ace.js" type="text/javascript" charset="utf-8"></script>
<script type="text/javascript" src="../static/js/leaps.js"></script>
<script type="text/javascript" src="../static/js/leapexample.js"></script>
<script src="https://cdn.socket.io/socket.io-1.4.5.js"></script>
<script src="http://code.jquery.com/jquery-1.11.1.js"></script>

<script>

    var socket = io('192.168.0.11:8110').connect();
    var roomUrl = window.location.href;
    socket.emit('join',{
        username: $('#nametext').val(),
        roomUrl: roomUrl
    });

    socket.on('broadcast_join', function(data){
        console.log(data.username + " join");
        $('#messages').append("One user join \n");
    });

    socket.on('broadcast_quit', function(data){
        console.log(data.username + " leave");
        $('#messages').append(data.username + " left \n");
    });

    $('form').submit(function(){
        var username = $('#nametext').val();
        var message = $('#sendtext').val();

        socket.emit('say', {
            username: username,
            text : message,
            roomUrl : roomUrl
        });

        $('#sendtext').val('');
        return false;
    });

    socket.on('broadcast_say', function(msg){
        console.log("someone says");
        var dt = new Date();
        var time = dt.getHours() + ":" + dt.getMinutes() + ":" + dt.getSeconds();
        $('#messages').append(msg.username + ": " + msg.text + "\n");
    });



    $(document).ready(function(){
        $('#messages').scrollBottom($('#messages')[0].scrollHeight);
    });
</script>
</body>
</html>