(function (){
    var ws = new WebSocket("ws://0.0.0.0:9000/ws");
    var stream = document.getElementById("messages")

    ws.onopen = function(e) {
        console.log("connection opened")
    }
    ws.onmessage = function (e) {
        console.log(e.data)
        stream.value = stream.value + '<div class="message">' + e.data + '</div>';
    };
    ws.onerror = function() {
        console.log("failed to connect over websocket")
    }
    ws.onclose = function(e){
        if (e.wasClean) {
            console.log("closed correctly")
            } else {
            console.log("broken connection")
            }
            console.log("code: "+ e.code + " reason: " + e.reason)
    };
})()