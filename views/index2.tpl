<html>
	<head>
		<script type="text/javascript" src="ace/ace.js"></script>
		<script>
			if (!window.ace) {
				console.log("injecting ace from cdn");
				document.write('<script src="http://cdn.jsdelivr.net/ace/1.1.8/min/ace.js"><\/script>');
			}
		</script>
		<script type="text/javascript" src="/leaps.js"></script>
		<script type="text/javascript" src="/leapexample.js"></script>
		<link href="/style.css" type="text/css" rel="stylesheet"></link>
	</head>
	<body>
		<h1>Leaps Ace Editor Example</h1>
		<div class="container">
			<div id="editor">console.log("This shouldn't be here");</div>
		</div>
	</body>
</html>
