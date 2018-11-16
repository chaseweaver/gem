document.addEventListener('astilectron-ready', function() {
	astilectron.onMessage(function(message) {
		if (message.name === 'ip') {
			document.getElementById('user-ip').innerHTML = 'IP: ' + message.payload;
		} else if (message.name === 'port') {
			document.getElementById('user-port').innerHTML = 'Port: ' + message.payload;
		} else if (message.name === 'receive') {
			let them = document.createElement('li');
			them.className = 'them';
			them.innerHTML = message.payload;
			document.getElementById('message-container').appendChild(them);
		}
	});

	document.getElementById('content-box').addEventListener('keyup', function(event) {
		event.preventDefault();
		if (event.keyCode === 13) {
			let me = document.createElement('li');
			me.className = 'me';
			me.innerHTML = this.value;
			document.getElementById('message-container').appendChild(me);
			astilectron.sendMessage({name: 'send', payload: this.value});
			this.value = '';
		}
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('connect').addEventListener('click', function() {
		let ip = document.getElementById('peer-ip').value;
		let port = document.getElementById('peer-port').value;

		if (ip.length == 0) {
			astilectron.showErrorBox('Error!', 'Please enter a peer IP!');
		} else if (port.length == 0) {
			astilectron.showErrorBox('Error!', 'Please enter a peer Port! (i.e. 3000)')
		} else {
			astilectron.sendMessage({name: 'connect', payload: [ip, port]});
			astilectron.showMessageBox({message: 'Connecting to: tcp://' + ip + ':' + port, title: 'GEM : Go Encryption Messenger'});
		} 
	});
});