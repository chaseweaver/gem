document.addEventListener('astilectron-ready', function() {
	let helpText = 'sample';
	document.getElementById('console-box').innerHTML = helpText;

	astilectron.onMessage(function(message) {
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('file').addEventListener('click', function() {
		max = 0;
		astilectron.showOpenDialog({properties: ['openFile', 'singleSelection'],
			title: 'File to Encrypt/Decrypt'}, function(paths) {
				astilectron.sendMessage({name: 'open-file', payload: paths[0]});
		});
	});
});