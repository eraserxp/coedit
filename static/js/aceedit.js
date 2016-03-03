/**
 * Created by jiyu on 11/02/16.
 */
"use strict";


var editor = ace.edit("editor");
editor.setTheme("ace/theme/twilight");
editor.getSession().setMode("ace/mode/javascript");

editor.resize();

function changeLan() {
    var e = document.getElementById('mode');
    var lang = e.options[e.selectedIndex].value;
    editor.getSession().setMode( "ace/mode/" + lang);
}

function changeTheme() {
    var e = document.getElementById('theme');
    var theme = e.options[e.selectedIndex].value;
    editor.setTheme( "ace/theme/" + theme);
}

/*document.getElementById('mode').addEventListener('change',
    function() {
        var e = document.getElementById('mode');
        var lang = e.options[e.selectedIndex].text;
        editor.getSession().setMode( "ace/mode/" + lang);
    }
    , false);*/
