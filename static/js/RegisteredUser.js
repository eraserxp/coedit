/**
 * Created by jiyu on 06/04/16.
 */
"use strict";

document.onload = getFileList();


function showDocDialog() {
    document.getElementById("dialogbg").style.display ="block";
    document.getElementById("newdocdialog").style.display ="block";
    document.getElementById("newdocname").value = "";

    var sWidth, sHeight;
    sWidth = screen.width;
    sWidth = document.body.offsetWidth;
    sHeight = document.body.offsetHeight;

    if( sHeight < screen.height) { sHeight = screen.height;}

    document.getElementById("dialogbg").style.width = sWidth + "px";
    document.getElementById("dialogbg").style.height = sHeight + "px";
    document.getElementById("dialogbg").style.display = "block";
    document.getElementById("dialogbg").style.right = document.getElementById("newdocdialog").offsetLeft + "px";

}

function CreateCancel() {
    document.getElementById("dialogbg").style.display ="none";
    document.getElementById("newdocdialog").style.display ="none";
}

function createNewFile() {
    var text = document.getElementById("newdocname").value;

    //Todo: check the uniqueness of the filename
    if ( text != null && text != "") {

        var json_upload = JSON.stringify( {documentname: text})

        var xhr = new XMLHttpRequest()
        xhr.open('post', '/addnewdoc')
        xhr.setRequestHeader('Content-Type', 'application/json')
        xhr.onreadystatechange = function() {
            if( xhr.readyState == 4 && xhr.status == 200) {

                if( xhr.responseText == "OK")
                {
                    getFileList();
                }
                else
                {
                    alert("A file has the same name already exists or the name if an invalid name !");
                }

            }
        }
        xhr.send(json_upload)

        document.getElementById("dialogbg").style.display ="none";
        document.getElementById("newdocdialog").style.display ="none";
    }
}

function deleteFile() {



    var filelist = document.getElementById('filelist');

    var sel = filelist.options[filelist.selectedIndex].value;

    var json_upload = JSON.stringify( {documentname: sel} )
    var xhr = new XMLHttpRequest();

    xhr.open('post', '/deletedoc');

    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = function() {
        if ( xhr.readyState == 4 && xhr.status == 200) {
            getFileList();
        }
    }
    xhr.send(json_upload);


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
