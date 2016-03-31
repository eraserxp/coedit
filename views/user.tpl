<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>New Document</title>
    <link rel="stylesheet" href="../static/css/index.css" />
    <style type="text/css" media="screen">
        #editor {
            position: relative;
            width: 95%;
            height: 100%;
            left: 2.5%;
            border: 1px solid black;
        }
    </style>
    <script src="../static/js/main.js" type="text/javascript"></script>
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
        <div class="iconright"> <button>Logout</button></div>
    </div>
</div>
<div class="body">

    <div class="filelist">
        <Select id="filelist" class="files" size="2" name="files">
            {{ .Options }}
        </Select>

        <button onclick="createNewFile()">Create New</button>
    </div>

    <div class="reginfo">

        <div class="editor">
            <div id="editor">function foo(items) {
                var x = "All this is syntax highlighted";
                return x;}
            </div>
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
                <option value="twilight">Twilight</option>
                <option value="tomorrow">Tomorrow</option>
                <option value="chrome">Chrome</option>
            </select>
            </div>
            <div class="selectcolumn">

            </div>
        </div>
    </div>

    <div class="chatter">
        Current users:<br/>
        <textarea class="usertext" readonly>User1
User2</textarea>
        <br/><br/>
        Message log:<br/>
        <textarea class="logtext" readonly>Ace: sometext</textarea>
        <br/><br/>
        Send message: <br/>
        <textarea class="sendtext"></textarea>
    </div>


</div>

<script src="http://d1n0x3qji82z53.cloudfront.net/src-min-noconflict/ace.js" type="text/javascript" charset="utf-8"></script>
<script src="../static/js/aceedit.js" type="text/javascript">
</script>

</body>
</html>