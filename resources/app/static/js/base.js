document.addEventListener('astilectron-ready', function() {
	let peers = [''];

	astilectron.onMessage(function(message) {
		if (message.name === 'ip') {
			document.getElementById('user-ip').value = message.payload;
		} else if (message.name === 'port') {
			document.getElementById('user-port').value = message.payload;
		} else if (message.name === 'receive') {
			let arr = JSON.parse(message.payload);
			let them = document.createElement('li');
			let container = document.createElement('div');
			let p = document.createElement('p');
			let time = new Date();
			let details = arr['name'] + ' at ' + time.toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });

			container.className = 'content-container left';
			container.id = 'content-container';
			them.className = 'them';
			
			them.innerHTML = arr['msg'];
			p.innerHTML = details;

			container.appendChild(them);
			container.appendChild(p);
			document.getElementById('message-container').appendChild(container);
		} else {
			sendResponse(message.payload, message.name);
		}
	});

	document.getElementById('content-box').addEventListener('keyup', function(event) {
		event.preventDefault();
		if (event.keyCode === 13) {

			if (this.value.length === 0) {
				return;
			} 

			if (peers.length === 0) {
				sendResponse('Please connect to a peer first!', 'error');
				return;
			}

			let pwd = document.getElementById('pwd').value;
			let nickname = document.getElementById('nickname').value;

			if (nickname.length === 0) {
				sendResponse('Please enter a nickname!', 'error');
			} else if (pwd.length === 0 ) {
				sendResponse('Please enter a password!', 'error');
			} else {
				let me = document.createElement('li');
				let container = document.createElement('div');
				let p = document.createElement('p');
				let time = new Date();
				let details = nickname + ' at ' + time.toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });
				
				container.className = 'content-container right';
				container.id = 'content-container';
				me.className = 'me';
				
				me.innerHTML = this.value;
				p.innerHTML = details;

				container.appendChild(me);
				container.appendChild(p);
				document.getElementById('message-container').appendChild(container);
				
				astilectron.sendMessage({name: 'send', payload: [this.value, pwd, nickname]});
				this.value = '';
			}
		}
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('log').addEventListener('click', function() {
		let li = document.getElementById('message-container').getElementsByTagName('li');
		let p = document.getElementById('message-container').getElementsByTagName('p');
		let content = '';

		for (let i = 0; i < li.length; i++) {
			content += '[' + p[i].innerHTML + ']' + ' ' + li[i].innerHTML + '\n';
		}

		if (content.length === 0) {
			sendResponse('There are no messages to save!', 'error');
			return
		}

		astilectron.showSaveDialog({title: 'Save Chat Logs'}, function(filename) {
			astilectron.sendMessage({name: 'save', payload: [filename, content]});
		})
	});

	document.getElementById('check').addEventListener('click', function() {
		let port = document.getElementById('user-port').value;

		if (port.length === 0) {
			sendResponse('Please enter a valid port for yourself!', 'error');
			return
		} else if (port < 1 || port > 65535) {
			sendResponse('Input port must be between 1-65535!', 'error');
			return
		} else {
			peers = [''];
			astilectron.sendMessage({name: 'change-port', payload: port});
		}
	});

	document.getElementById('connect').addEventListener('click', function() {
		let ip = document.getElementById('peer-ip').value;
		let port = document.getElementById('peer-port').value;
		let pwd = document.getElementById('pwd').value;

		if (ip.length === 0) {
			sendResponse('Please enter a Peer IP!', 'error');
		} else if (port.length === 0) {
			sendResponse('Please enter a peer Port! (i.e. 3000)!', 'error');
		} else if (pwd.length === 0 ) {
			sendResponse('Please enter a password!', 'error');
		} else if (peers.includes('tcp://'+ip+':'+port)) {
			sendResponse('You are already connected to that peer!', 'error');
		} else {
			peers.push('tcp://'+ip+':'+port);
			astilectron.sendMessage({name: 'connect', payload: [ip, port, pwd]});

			sendResponse('Adding peer: tcp://' + ip + ':' + port, 'success');

			ip.value = '';
			port.value = '';
		} 
	});
});

function sendResponse(msg, type) {
	let container = document.createElement('div');
	let li = document.createElement('li');
	let p = document.createElement('p');

	let time = new Date();
	let details = type.toUpperCase() + ' at ' + time.toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true });

	container.className = 'content-container';
	li.className = type;
	li.innerHTML = msg;
	p.innerHTML = details;

	container.appendChild(li);
	container.appendChild(p);
	
	document.getElementById('message-container').appendChild(container);
}