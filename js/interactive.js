'use strict';

// Date formatting
//
// Import Luxon DateTime formatting wrapper.

async function importDateTime() {
	return luxon.DateTime;
}

// Zone data
//
// This maps an IANA zone into readable strings with bcp47. Sadly, Intl doesn't provide such strings.

async function importZoneData() {
	const zoneData = {};

	const res = await fetch("/js/bcp47timezone.json");
	const data = await res.json();
	for (const timezone of data) {
		for (const alias of timezone.aliases) {
			zoneData[alias] = timezone;
		}
	}

	return zoneData;
}

// Start

Promise.all([
	importDateTime(),
	importZoneData(),
]).then(dependencies => {
	start(...dependencies);
});

function start(DateTime, zoneData) {

	// Datetime translation
	//
	// This maps a datetime and zone into a zone name, date and time in the current locale.

	function translateDatetime(datetime, zone) {
		const dt = DateTime.fromISO(datetime).setZone(zone);
		const dateString = dt.toLocaleString(DateTime.DATE_MED_WITH_WEEKDAY);
		const timeString = dt.toLocaleString(DateTime.TIME_SIMPLE);
		const zoneObject = zoneData[zone];

		return {
			zone: zoneObject,
			date: dateString,
			time: timeString,
		};
	}


	// URL parser
	//
	// This maps the URL path into various parts. The server-side parser should match this code exactly.

	function parsePath(path) {
		let cleanup = path;
		cleanup = cleanup.replace(/^\/+/, ""); // Remove start slashes
		cleanup = cleanup.replace(/\/+$/, ""); // Remove end slashes
		let parts = cleanup.split("/");

		// Simple format: iso_time
		if (parts.length == 1) {
			return {
				datetime: parts[0],
				zones: ['local'],
			};
		}

		// Simple format: iso_time/csv_zones
		if (parts.length == 2) {
			return {
				datetime: parts[0],
				zones: parsePathZones(parts[1]),
			};
		}

		return null;
	}

	function parsePathZones(zones) {
		let parts = zones.split(",");

		let zs = parts.map(zone => {
			return zone.replace("+", "/");
		});
		return zs;
	}

	// Zone

	const zoneTemplate = document.createElement('template');
	zoneTemplate.innerHTML = `
		<slot></slot>
	`;

	class ZoneElement extends HTMLElement {
		constructor() {
			super();

			this.attachShadow({ mode: 'open' });
			this.shadowRoot.appendChild(zoneTemplate.content.cloneNode(true));
		}
	}
	
	// ZoneInfo

	const zoneInfoTemplate = document.createElement('template');
	zoneInfoTemplate.innerHTML = `
		<slot></slot>
	`;

	class ZoneInfoElement extends HTMLElement {
		constructor() {
			super();

			this.attachShadow({ mode: 'open' });
			this.shadowRoot.appendChild(zoneInfoTemplate.content.cloneNode(true));
		}
	}

	// ZoneName

	const zoneNameTemplate = document.createElement('template');
	zoneNameTemplate.innerHTML = `
		<slot></slot>
	`;

	class ZoneNameElement extends HTMLElement {
		constructor() {
			super();

			this.attachShadow({ mode: 'open' });
			this.shadowRoot.appendChild(zoneNameTemplate.content.cloneNode(true));
		}
	}

	// ZoneOffset

	const zoneOffsetTemplate = document.createElement('template');
	zoneOffsetTemplate.innerHTML = `
		<slot></slot>
	`;

	class ZoneOffsetElement extends HTMLElement {
		constructor() {
			super();

			this.attachShadow({ mode: 'open' });
			this.shadowRoot.appendChild(zoneOffsetTemplate.content.cloneNode(true));
		}
	}

	// Datetime

	const datetimeTemplate = document.createElement('template');
	datetimeTemplate.innerHTML = `
		<slot></slot>
	`;

	class DatetimeElement extends HTMLElement {
		constructor() {
			super();

			this.attachShadow({ mode: 'open' });
			this.shadowRoot.appendChild(datetimeTemplate.content.cloneNode(true));
		}
	}

	// Date

	const dateTemplate = document.createElement('template');
	dateTemplate.innerHTML = `
		<style>
			#editor-input {
				padding: 0;
				border: 0;

				width: 100%;

				background: none;
				color: inherit;
				font-size: inherit;
				font-weight: inherit;
				font-family: inherit;
				text-align: right;
			}
			@media (prefers-color-scheme: dark) {
				::-webkit-calendar-picker-indicator {
					filter: invert(1);
				}
			}
		</style>
		<input id="editor-input" type="date" required>
	`;

	class DateElement extends HTMLElement {
		constructor() {
			super();

			this.attachShadow({ mode: 'open' });
			this.shadowRoot.appendChild(dateTemplate.content.cloneNode(true));
		}
		connectedCallback() {
			const editorInput = this.shadowRoot.querySelector('#editor-input');

			const value = this.getAttribute('date');
			editorInput.value = value;

			editorInput.addEventListener('blur', () => {
				// TODO
			});
		}
	}

	// Time

	const timeTemplate = document.createElement('template');
	timeTemplate.innerHTML = `
		<style>
			#editor-input {
				padding: 0;
				border: 0;

				width: 100%;

				background: none;
				color: inherit;
				font-size: inherit;
				font-weight: inherit;
				font-family: inherit;
				text-align: right;
			}
			@media (prefers-color-scheme: dark) {
				::-webkit-calendar-picker-indicator {
					filter: invert(1);
				}
			}
		</style>
		<input id="editor-input" type="time" required pattern="[0-9]{2}:[0-9]{2}">
	`;

	class TimeElement extends HTMLElement {
		constructor() {
			super();

			this.attachShadow({ mode: 'open' });
			this.shadowRoot.appendChild(timeTemplate.content.cloneNode(true));
		}
		connectedCallback() {
			const editorInput = this.shadowRoot.querySelector('#editor-input');

			const value = this.getAttribute('time');
			editorInput.value = value;

			editorInput.addEventListener('blur', () => {
				// TODO
			});
		}
	}

	customElements.define('datetime-zone', ZoneElement);

	customElements.define('datetime-zoneinfo', ZoneInfoElement);
	customElements.define('datetime-zonename', ZoneNameElement);
	customElements.define('datetime-zoneoffset', ZoneOffsetElement);

	customElements.define('datetime-datetime', DatetimeElement);
	customElements.define('datetime-date', DateElement);
	customElements.define('datetime-time', TimeElement);

	// Page events

	const path = window.location.pathname;

}
