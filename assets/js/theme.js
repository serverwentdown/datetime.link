'use strict';

const OPTIONS = {
	'dark': {
		name: 'theme-dark',
	},
	'light': {
		name: 'theme-light',
	},
	'system': {
		name: '',
	},
};

function setTheme(theme) {
	for (let key in OPTIONS) {
		const cls = OPTIONS[key];
		if (cls.name) {
			document.body.classList.remove(cls.name);
		}
	}
	const cls = OPTIONS[theme];
	if (cls === undefined) {
		return;
	}
	if (cls.name) {
		document.body.classList.add(cls.name);
	}
}

function toggleTheme() {
	let theme = localStorage.getItem('theme');
	if (theme === 'dark') {
		theme = 'light';
	} else if (theme === 'light') {
		theme = 'system';
	} else if (theme === 'system') {
		theme = 'dark';
	} else {
		// Default is system -> dark
		theme = 'dark';
	}
	localStorage.setItem('theme', theme);
	setTheme(theme);
}

const theme = localStorage.getItem('theme');
if (theme) {
	setTheme(theme);
}

