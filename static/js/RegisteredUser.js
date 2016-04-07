/**
 * Created by jiyu on 06/04/16.
 */
"use strict";

document.onload = getFileList();


function createNewFile() {
    var pro = prompt("Please enter the name for the file", "example: code.js");
    //Todo: check the uniqueness of the filename
    if ( pro != null ) {

        var json_upload = JSON.stringify( {documentname: pro})

        var xhr = new XMLHttpRequest()
        xhr.open('post', '/addnewdoc')
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.onreadystatechange = function() {
            if( xhr.readyState == 4 && xhr.status == 200) {
                getFileList();
            }
        }
        xhr.send(json_upload)


    }
}

function getFileList() {
    var xhr = new XMLHttpRequest();

    xhr.open('get', '/requestuserlist')
    xhr.onreadystatechange = function() {
        if ( xhr.readyState == 4 && xhr.status == 200) {
            var jsonData = JSON.parse( xhr.responseText );

            var files = jsonData.toString().split(",");
            var filelist = document.getElementById('filelist');
            filelist.options.length = 0;

            for ( var index in files) {
                var opt = new Option(files[index], files[index]);
                filelist.options.add( opt );
            }
        }
    }
    xhr.send()
}

function openSelectedDoc() {
    var filelist = document.getElementById('filelist');

    var si = filelist.selectedIndex;
    var filename = filelist.options[si].value ;

    var json_upload = JSON.stringify( {documentname: filename})

    var xhr = new XMLHttpRequest()
    xhr.open('post', '/opendoc')
    xhr.setRequestHeader('Content-Type', 'application/json')
    xhr.onreadystatechange = function() {
        if( xhr.readyState == 4 && xhr.status == 200) {
            var re = xhr.responseText;

            window.open('/regdoc/' + re);
        }
    }
    xhr.send(json_upload)
}
