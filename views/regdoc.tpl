<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="../static/css/index.css" />
    <link rel="stylesheet" href="../static/bootstrap/css/bootstrap.css" />
    <link rel="stylesheet" href="../static/css/bootstrap-social.css" />
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.5.0/css/font-awesome.min.css" />
    <title>{{.FileName}}</title>

    <style type="text/css" media="screen">
        #editor {
            position: relative;
            width: 95%;
            height: 100%;
            left: 2.5%;
            border: 1px solid black;
        }
    </style>

    <script src="../static/js/logcontrol.js" type="text/javascript" ></script>
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
                <button onclick="changePrivacy()">Change Privacy Option...</button>
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

    <div id="dialogbg"></div>
    <div id="newdocdialog" style="display:none">
        <div> Please give the privacy option of the current file: </div>
        <div> Who can access this file? : <select id="privacy" onchange="privacyOptionChange()">
            <option value="E">Everyone can access</option>
            <option value="S">Some user can access</option>
            <option value="N">No other user can access</option>
        </select></div>
        <div> List for emails that can access this file:</div>
        <div> <textarea id="accessemails" placeholder="One line for every email."></textarea> </div>
        <div> <button onclick="handinPrivacy()">Submit</button> <button onclick="handinCancel()">Cancel</button>  </div>
    </div>

    <script src="../static/ace/ace.js" type="text/javascript" charset="utf-8"></script>
    <script type="text/javascript" src="../static/js/leaps.js"></script>
    <script type="text/javascript" src="../static/js/leapexample.js"></script>

</body>
</html>