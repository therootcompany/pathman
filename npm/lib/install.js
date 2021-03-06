'use strict';

var pkg = require('../package.json');
var path = require('path');
var os = require('os');

// https://nodejs.org/api/os.html#os_os_arch
// 'arm', 'arm64', 'ia32', 'mips', 'mipsel', 'ppc', 'ppc64', 's390', 's390x', 'x32', and 'x64'
var arch = os.arch(); // process.arch

// https://nodejs.org/api/os.html#os_os_platform
// 'aix', 'darwin', 'freebsd', 'linux', 'openbsd', 'sunos', 'win32'
var platform = os.platform(); // process.platform
var ext = /^win/i.test(platform) ? '.exe' : '';

// This is _probably_ right. It's good enough for us
// https://github.com/nodejs/node/issues/13629
if ('arm' === arch) {
	arch += 'v' + process.config.variables.arm_version;
}

var map = {
	// arches
	armv6: 'armv6',
	armv7: 'armv7',
	arm64: 'armv8',
	ia32: '386',
	x32: '386',
	x64: 'amd64',
	// platforms
	darwin: 'darwin',
	linux: 'linux',
	win32: 'windows'
};

arch = map[arch];
platform = map[platform];

var pkg = require('../package.json');
var newVer = pkg.version;
var fs = require('fs');
var request = require('@root/request');
var mkdirp = require('@root/mkdirp');

function install(name, bindirs, getVersion, parseVersion, urlTpl) {
	if (!arch || !platform) {
		console.error(
			"'" +
				os.platform() +
				"' on '" +
				os.arch() +
				"' isn't supported yet."
		);
		console.error(
			'Please open an issue at https://git.rootprojects.org/root/pathman/issues'
		);
		process.exit(1);
	}

	var url = urlTpl
		.replace(/{{ .Version }}/g, newVer)
		.replace(/{{ .Platform }}/g, platform)
		.replace(/{{ .Arch }}/g, arch)
		.replace(/{{ .Ext }}/g, ext);

	console.info('Installing from', url);
	return request({ uri: url, encoding: null }, function(err, resp) {
		if (err) {
			console.error(err);
			return;
		}

		//console.log(resp.body.byteLength);
		//console.log(typeof resp.body);
		var bin = name + ext;
		function next() {
			if (!bindirs.length) {
				return;
			}
			var bindir = bindirs.pop();
			return mkdirp(bindir, function(err) {
				if (err) {
					console.error(err);
					return;
				}

				var localbin = path.join(bindir, bin);
				return fs.writeFile(localbin, resp.body, function(err) {
					next();
					if (err) {
						console.error(err);
						return;
					}
					fs.chmodSync(localbin, parseInt('0755', 8));
					console.info('Wrote', bin, 'to', bindir);
				});
			});
		}
		next();
	});
}

function shouldUpdate(oldVer, newVer) {
	// "v1.0.0-pre" is BEHIND "v1.0.0"
	newVer = newVer
		.replace(/^v/, '')
		.split(/[\.\-\+]/)
		.filter(Boolean);
	oldVer = oldVer
		.replace(/^v/, '')
		.split(/[\.\-\+]/)
		.filter(Boolean);

	if (!oldVer.length) {
		return true;
	}

	// ex: v1.0.0-pre vs v1.0.0
	if (newVer[3] && !oldVer[3]) {
		// don't install beta over stable
		return false;
	}

	// ex: old is v1.0.0-pre
	if (oldVer[3]) {
		if (oldVer[2] > 0) {
			oldVer[2] -= 1;
		} else if (oldVer[1] > 0) {
			oldVer[2] = 999;
			oldVer[1] -= 1;
		} else if (oldVer[0] > 0) {
			oldVer[2] = 999;
			oldVer[1] = 999;
			oldVer[0] -= 1;
		} else {
			// v0.0.0
			return true;
		}
	}

	// ex: v1.0.1 vs v1.0.0-pre
	if (newVer[3]) {
		if (newVer[2] > 0) {
			newVer[2] -= 1;
		} else if (newVer[1] > 0) {
			newVer[2] = 999;
			newVer[1] -= 1;
		} else if (newVer[0] > 0) {
			newVer[2] = 999;
			newVer[1] = 999;
			newVer[0] -= 1;
		} else {
			// v0.0.0
			return false;
		}
	}

	// ex: v1.0.1 vs v1.0.0
	if (oldVer[0] > newVer[0]) {
		return false;
	} else if (oldVer[0] < newVer[0]) {
		return true;
	} else if (oldVer[1] > newVer[1]) {
		return false;
	} else if (oldVer[1] < newVer[1]) {
		return true;
	} else if (oldVer[2] > newVer[2]) {
		return false;
	} else if (oldVer[2] < newVer[2]) {
		return true;
	} else if (!oldVer[3] && newVer[3]) {
		return false;
	} else if (oldVer[3] && !newVer[3]) {
		return true;
	} else {
		return false;
	}
}

/*
// Same version
console.log(false === shouldUpdate('0.5.0', '0.5.0'));
// No previous version
console.log(true === shouldUpdate('', '0.5.1'));
// The new version is slightly newer
console.log(true === shouldUpdate('0.5.0', '0.5.1'));
console.log(true === shouldUpdate('0.4.999-pre1', '0.5.0-pre1'));
// The new version is slightly older
console.log(false === shouldUpdate('0.5.0', '0.5.0-pre1'));
console.log(false === shouldUpdate('0.5.1', '0.5.0'));
*/

function checkVersion(getVersion, parseVersion) {
	var exec = require('child_process').exec;

	return new Promise(function(resolve) {
		exec(getVersion, { windowsHide: true }, function(err, stdout) {
			var oldVer = parseVersion(stdout);
			resolve(oldVer);
			/*
			//console.log('old:', oldVer, 'new:', newVer);
			if (!shouldUpdate(oldVer, newVer)) {
				console.info(
					'Current ' + name + ' version is new enough:',
					oldVer,
					newVer
				);
				return;
				//} else {
				//	console.info('Current version is older:', oldVer, newVer);
			}
      */
		});
	});
}

module.exports = install;
module.exports._shouldUpdate = shouldUpdate;
module.exports._checkVersion = checkVersion;
