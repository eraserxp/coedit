/**
 * Created by jiyu on 11/02/16.
 */
"use strict";

var editor = ace.edit("editor");
editor.setTheme("ace/theme/tomorrow");
editor.getSession().setMode("ace/mode/javascript");

editor.resize();

function changeLan() {
    var e = document.getElementById('mode');
    var lang = e.options[e.selectedIndex].text;
    editor.getSession().setMode( "ace/mode/" + lang);
}

/*document.getElementById('mode').addEventListener('change',
    function() {
        var e = document.getElementById('mode');
        var lang = e.options[e.selectedIndex].text;
        editor.getSession().setMode( "ace/mode/" + lang);
    }
    , false);*/
