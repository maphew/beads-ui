// htmx handles all filtering dynamically - no custom JS needed!

// Live reload WebSocket connection
(function() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = protocol + '//' + window.location.host + '/ws';
    const ws = new WebSocket(wsUrl);

    ws.onmessage = function(event) {
        if (event.data === 'reload') {
            window.location.reload();
        }
    };

    ws.onclose = function() {
        // Optionally reconnect after a delay
        setTimeout(function() {
            window.location.reload();
        }, 1000);
    };
})();
