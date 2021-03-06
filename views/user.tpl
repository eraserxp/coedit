<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Email}}'s main page</title>
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
    <script scr="https://code.jquery.com/jquery-2.2.3.min.js" type="text/javascript"></script>
    <script src="../static/js/RegisteredUser.js" type="text/javascript" ></script>
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

    <div class="filelist">
        <Select id="filelist" class="files" size="2" name="files" ondblclick="openSelectedDoc()">

        </Select>

        <div>
            <button onclick="showDocDialog()">Create New File</button>
            <button onclick="deleteFile()">Delete Selected File</button>
        </div>

    </div>

    <div id="dialogbg"></div>
    <div id="newdocdialog" style="display:none;width:400px">
        <div> Please give the name for the file: </div>
        <div> <input id="newdocname" placeholder="example: test.js" /> </div>
        <div> <button onclick="createNewFile()">Create</button>  <button onclick="CreateCancel()">Cancel</button></div>
    </div>


</div>

<script src="http://d1n0x3qji82z53.cloudfront.net/src-min-noconflict/ace.js" type="text/javascript" charset="utf-8"></script>
<script src="../static/js/aceedit.js" type="text/javascript"></script>

</body>
</html>