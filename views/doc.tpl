<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Collaborative editing platform</title>
	<link href="static/css/main.css" type="text/css" rel="stylesheet"></link>
    <style type="text/css" media="screen">
        #editor {
            position: relative;
            width: 800px;
            height: 400px;
            margin-top: 100px;
            border: 1px solid black;
        }
    </style>
</head>
<body>

<div class="area">
    <div id="editor"></div>

    <div id="misc" class="func">
        <select id="mode" onchange="changeLan()">
            <option>javascript</option>
            <option>xml</option>
            <option>python</option>
            <option>java</option>
            <option>c</option>
        </select>
    </div>
</div>


<script src="http://cdn.jsdelivr.net/ace/1.1.8/min/ace.js" type="text/javascript"
        charset="utf-8"></script>
<script type="text/javascript" src="static/js/leaps.js"></script>
<script type="text/javascript" src="static/js/leapexample.js"></script>


</body>
</html>