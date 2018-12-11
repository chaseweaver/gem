document.addEventListener('astilectron-ready', function() {
	let peers = [""];

	astilectron.onMessage(function(message) {
		if (message.name === 'ip') {
			document.getElementById('user-ip').innerHTML = 'IP: ' + message.payload;
		} else if (message.name === 'port') {
			document.getElementById('user-port').innerHTML = 'Port: ' + message.payload;
		} else if (message.name === 'receive') {
			let arr = JSON.parse(message.payload);
			let them = document.createElement('li');
			let container = document.createElement('div');
			let p = document.createElement('p');
			let time = new Date();
			let details = arr['name'] + ' at ' + time.toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });

			container.className = 'content-container';
			them.className = 'them';
			
			them.innerHTML = arr['msg'];
			p.innerHTML = details;

			container.appendChild(them);
			container.appendChild(p);
			document.getElementById('message-container').appendChild(container);
		}
	});

	document.getElementById('content-box').addEventListener('keyup', function(event) {
		event.preventDefault();
		if (event.keyCode === 13) {

			if (this.value.length === 0) {
				return;
			} 

			let pwd = document.getElementById('pwd').value;
			let nickname = document.getElementById('nickname').value;

			if (pwd.length === 0 ) {
				astilectron.showErrorBox('Error!', 'Please enter a password!');
			} else if (nickname.length === 0) {
				astilectron.showErrorBox('Error!', 'Please enter a nickname!');
			} else {
				let me = document.createElement('li');
				let container = document.createElement('div');

				container.className = 'content-container';
				me.className = 'me';
				me.innerHTML = this.value;
				container.appendChild(me);
				document.getElementById('message-container').appendChild(container);
				astilectron.sendMessage({name: 'send', payload: [this.value, pwd, nickname]});
				this.value = '';
			}
		}
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('connect').addEventListener('click', function() {
		let ip = document.getElementById('peer-ip').value;
		let port = document.getElementById('peer-port').value;
		let pwd = document.getElementById('pwd').value;

		if (peers.contains('tcp://'+ip+':'+port)) {
			astilectron.showErrorBox('Error!', 'You are already connected to that peer!');
			return
		} else {
			peers.push('tcp://'+ip+':'+port);
		}

		if (ip.length === 0) {
			astilectron.showErrorBox('Error!', 'Please enter a peer IP!');
		} else if (port.length === 0) {
			astilectron.showErrorBox('Error!', 'Please enter a peer Port! (i.e. 3000)');
		} else if (pwd.length === 0 ) {
			astilectron.showErrorBox('Error!', 'Please enter a password!');
		} else {
			astilectron.sendMessage({name: 'connect', payload: [ip, port, pwd]});
			astilectron.showMessageBox({message: 'Connecting to: tcp://' + ip + ':' + port, title: 'GEM : Go Encryption Messenger'});
		} 
	});
});