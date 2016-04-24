/**
 * Created by jiyu on 07/04/16.
 */

"use strict";

function logout() {
    var xhr = new XMLHttpRequest()
    xhr.open('get', '/logout')
    xhr.onreadystatechange = function() {
        if ( xhr.readyState == 4 && xhr.status == 200) {
            window.location.pathname = "/" ;
        }
    }
    xhr.send()
}

function changePrivacy() {

    document.getElementById("dialogbg").style.display ="block";
    document.getElementById("newdocdialog").style.display ="block";

    var sWidth, sHeight;
    sWidth = screen.width;
    sWidth = document.body.offsetWidth;
    sHeight = document.body.offsetHeight;

    if( sHeight < screen.height) { sHeight = screen.height;}

    document.getElementById("dialogbg").style.width = sWidth + "px";
    document.getElementById("dialogbg").style.height = sHeight + "px";
    document.getElementById("dialogbg").style.display = "block";
    document.getElementById("dialogbg").style.right = document.getElementById("newdocdialog").offsetLeft + "px";

    loadPrivacyData();
}

function loadPrivacyData() {

    var filename = document.title;
    var json_upload = JSON.stringify( {documentname: filename})

    var xhr = new XMLHttpRequest()
    xhr.open('post', '/loadfileprivacy')
    xhr.setRequestHeader('Content-Type', 'application/json')
    xhr.onreadystatechange = function() {
        if( xhr.readyState == 4 && xhr.status == 200) {
            var jsonData = JSON.parse( xhr.responseText );

            var pri = jsonData["privacy"];
            var ae = jsonData["accessemails"];

            var priSelect = document.getElementById("privacy");
            priSelect.value = pri;

            var aeTextarea = document.getElementById("accessemails");
            aeTextarea.value = ae;

            privacyOptionChange();
        }
    }
    xhr.send(json_upload)

}

function privacyOptionChange() {
    var priSelect = document.getElementById("privacy");
    var aeTextarea = document.getElementById("accessemails");

    if( priSelect.value == "S") {
        aeTextarea.disabled = false;
    }
    else {
        aeTextarea.disabled = true;
    }
}

function handinPrivacy() {
    var filename = document.title;
    var pri = document.getElementById("privacy").value;
    var ae = document.getElementById("accessemails").value;

    var json_upload = JSON.stringify( {documentname: filename, privacy: pri, accessemails: ae})

    var xhr = new XMLHttpRequest()
    xhr.open('post', '/updatedocprivacy')
    xhr.setRequestHeader('Content-Type', 'application/json')
    xhr.onreadystatechange = function() {
        if( xhr.readyState == 4 && xhr.status == 200) {

        }
    }
    xhr.send(json_upload)


    document.getElementById("dialogbg").style.display ="none";
    document.getElementById("newdocdialog").style.display ="none";
}

function handinCancel() {
    document.getElementById("dialogbg").style.display ="none";
    document.getElementById("newdocdialog").style.display ="none";


}