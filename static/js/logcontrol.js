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