'use strict';

var spawner = require('./lib/exec.js');
module.exports = spawner;
module.exports.install = require('./lib/install.js');

if (require.main === module) {
	// node index.js --foo bar => [ '--foo', 'bar' ]
	spawner(process.argv.slice(2)).catch(function(err) {
		console.error(err.message);
	});
}
