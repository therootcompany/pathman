#!/usr/bin/env node

'use strict';

var pkg = require('../package.json');
var os = require('os');
var path = require('path');
var fs = require('fs');
var manager = require('../index.js');

if (require.main === module) {
	run();
}

function run() {
	var ext = /^win/i.test(os.platform()) ? '.exe' : '';
	//var homedir = require('os').homedir();
	//var bindir = path.join(homedir, '.local', 'bin');
	var bindir = path.resolve(__dirname, '..', 'bin');
	var name = pkg.name.replace(/.*\//, '');
	if ('.exe' === ext) {
		winpmstall(pkg.name, name, bindir);
	}

	return manager.install(
		name,
		[bindir],
		'pathman version',
		function parseVersion(stdout) {
			return (stdout || '').split(' ')[0];
		},
		'https://rootprojects.org/pathman/dist/{{ .Platform }}/{{ .Arch }}/pathman{{ .Ext }}'
	);
}

function winpmstall(pkgname, name, bindir) {
	var dd = /\//.test(pkgname) ? '../' : '';
	var pkgpath = pkgname.replace(/@/g, '\\@');

	try {
		fs.writeFileSync(
			path.join(bindir, name),
			[
				'#!/usr/bin/env bash',
				'"$(dirname "$0")/' + name + '.exe" "$@"',
				'exit $?'
			].join('\n')
		);
	} catch (e) {
		// ignore
	}

	// because bugs in npm + git bash oddities, of course
	// https://npm.community/t/globally-installed-package-does-not-execute-in-git-bash-on-windows/9394
	try {
		fs.writeFileSync(
			path.join(__dirname, dd + '../../.bin', name),
			[
				'#!/bin/sh',
				'# manual bugfix patch for npm on windows',
				'basedir=$(dirname "$(echo "$0" | sed -e \'s,\\\\,/,g\')")',
				'"$basedir/../' + pkgpath + '/bin/' + name + '"   "$@"',
				'exit $?'
			].join('\n')
		);
	} catch (e) {
		// ignore
	}
	try {
		fs.writeFileSync(
			path.join(__dirname, dd + '../../..', name),
			[
				'#!/bin/sh',
				'# manual bugfix patch for npm on windows',
				'basedir=$(dirname "$(echo "$0" | sed -e \'s,\\\\,/,g\')")',
				'"$basedir/node_modules/' +
					pkgname +
					'/bin/' +
					name +
					'"   "$@"',
				'exit $?'
			].join('\n')
		);
	} catch (e) {
		// ignore
	}
	// end bugfix
}
