'use strict';

var pkg = require('../package.json');
var spawn = require('child_process').spawn;
var os = require('os');
var path = require('path');
var ext = /^win/i.test(os.platform()) ? '.exe' : '';

// @scope/packagename => packagename
// { bin: { "packagename": "bin/runner" } } => "bin/runner"
var bin = path.resolve(__dirname, '..', pkg.bin[pkg.name.replace(/.*\//, '')]);

function spawner(args) {
	return new Promise(function(resolve, reject) {
		var runner = spawn(path.join(bin + ext), args, {
			windowsHide: true
		});
		runner.stdout.on('data', function(chunk) {
			console.info(chunk.toString('utf8'));
		});
		runner.stderr.on('data', function(chunk) {
			console.err(chunk.toString('utf8'));
		});
		runner.on('exit', function(code) {
			if (0 !== code) {
				reject(
					new Error("exited with non-zero status code '" + code + "'")
				);
				return;
			}
			resolve({ code: code });
		});
	});
}

module.exports = spawner;
