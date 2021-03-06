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
        <div class="profileinfo">
            Welcome, {{.Email}}
        </div>
        <div class="iconright"> <button onclick="logout()">Logout</button></div>
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
        <textarea class="nametext">Ace</textarea>
        <br/><br/>
        Message log:<br/>
        <textarea class="logtext" readonly>Ace: sometext</textarea>
        <br/><br/>
        Send message: <br/>
        <textarea class="sendtext"></textarea>
    </div>

    <script src="../static/ace/ace.js" type="text/javascript" charset="utf-8"></script>
    <script type="text/javascript" src="../static/js/leaps.js"></script>
    <script type="text/javascript" src="../static/js/leapexample.js"></script>

</body>
</html>