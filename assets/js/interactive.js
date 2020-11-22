'use strict';

/*
 * Compatibility Checks
 */

if (!window.customElements) {
	console.warn('Custom Elements API is not available. Thus, interactivity is not available');
}

/*
 * Custom Elements
 */

const slotTemplate = document.createElement('template');
slotTemplate.innerHTML = `
	<slot></slot>
`;

// Icons

const iconSolidTrashTemplate = document.querySelector('#icon-solid-trash');
class IconSolidTrashElement extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(iconSolidTrashTemplate.content.cloneNode(true));
	}
}

const iconSolidPlusTemplate = document.querySelector('#icon-solid-plus');
class IconSolidPlusElement extends HTMLElement {
	constructor() {
		super();
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(iconSolidPlusTemplate.content.cloneNode(true));
	}
}

// Zone

const zoneTemplate = document.createElement('template');
zoneTemplate.innerHTML = `
	<style>
		#toolbar {
			/*
			position: absolute;
			top: 0;
			right: 0;
			*/
			padding: 0.5rem 0;
		}

		button {
			padding: 0.2em 0.5em;

			border: 0;
			border-radius: 4px;

			font-family: inherit;
			font-size: 1em;
		}
		button.delete {
			background: rgba(255, 0, 0, 0.125);
			color: rgba(255, 0, 0, 1);
		}
	</style>
	<slot></slot>
	<div id="toolbar">
		<button type="button" class="delete" title="Delete this zone"><icon-solid-trash /></button>
	</div>
`;

class ZoneElement extends HTMLElement {
	constructor() {
		super();
		
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(zoneTemplate.content.cloneNode(true));
	}
}

// ZoneName

const zoneNameTemplate = document.createElement('template');
zoneNameTemplate.innerHTML = `
	<!--
	<style>
		:host {
			display: flex;
			justify-content: flex-start;
		}
		.zonearea {
			flex: 0 1 auto;

			display: block;

			overflow: hidden;
			white-space: nowrap;
			text-overflow: ellipsis;
		}
		.zonecountry {
			flex: 0 0 auto;

			display: block;
		}
	</style>
	-->
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
	<style>
		:host {
			display: block;

			white-space: nowrap;
			opacity: 0.5;
		}
	</style>
	<slot></slot>
`;

class ZoneOffsetElement extends HTMLElement {
	constructor() {
		super();
		
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(zoneOffsetTemplate.content.cloneNode(true));
	}
}

// Date

class DateElement extends HTMLElement {
	constructor() {
		super();
		
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(slotTemplate.content.cloneNode(true));
	}
}

// Time

class TimeElement extends HTMLElement {
	constructor() {
		super();
		
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(slotTemplate.content.cloneNode(true));
	}
}

// SearchList

const zoneSearchTemplate = document.createElement('template');
zoneSearchTemplate.innerHTML = `
	<style>
		#layout {
			position: relative;
		}

		input[type=search] {
			display: block;
			width: 100%;
			padding: 0.75rem;
			border: 1px solid rgba(127, 127, 127, 1);
			background: none;
			color: inherit;

			font-family: inherit;
			font-size: inherit;
		}

		ul {
		/*
			position: absolute;
			left: 0;
			right: 0;
		*/
			padding: 0;
			margin: 0;

			list-style: none;
			box-shadow: 4px 12px 12px rgba(127, 127, 127, 0.5);
		}
		li {
			display: flex;
			align-items: center;

			border: 1px solid rgba(127, 127, 127, 1);
			border-top: none;
		}
		ul, li:last-of-type {
			border-bottom-right-radius: 6px;
			border-bottom-left-radius: 6px;
		}
		li .left {
			flex: 1 1 0;
			padding: 0.75rem;
		}
		li .right {
			padding: 0.75rem;
		}
		button {
			padding: 0.2em 0.5em;

			border: 0;
			border-radius: 4px;

			font-family: inherit;
			font-size: 1.2em;
		}
		button.add {
			background: rgba(0, 180, 0, 0.125);
			color: rgba(0, 180, 0, 1);
		}
	</style>

	<div id="layout">
		<input type="search" placeholder="Search for a timezone...">
		
		<template id="result-template">
			<li>
				<div class="left">
					<d-zonename>Zone Name</d-zonename><d-zoneoffset>+XX:XX</d-zoneoffset>
				</div>
				<div class="right">
					<button type="button" class="add" title="Add this zone"><icon-solid-plus /></button>
				</div>
			</li>
		</template>
		<ul id="results"></ul>
	</div>
`;

class ZoneSearchElement extends HTMLElement {
	constructor() {
		super();
		
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(zoneSearchTemplate.content.cloneNode(true));
		this.shadowRoot.querySelector('input[type=search]').addEventListener('input', this.input.bind(this));
	}

	input() {
		const text = this.shadowRoot.querySelector('input[type=search]').value;
		const e = new CustomEvent('searchinput', { detail: { text } });
		this.dispatchEvent(e);
	}

	show(results) {
		const resultsElement = this.shadowRoot.querySelector('#results');
		const template = this.shadowRoot.querySelector('#result-template').content;

		resultsElement.innerHTML = '';
		for (const result of results) {
			const resultElement = template.cloneNode(true);
			resultElement.querySelector('d-zonename').innerText = result.n;
			resultElement.querySelector('d-zoneoffset').innerText = '+to:do';
			resultElement.querySelector('.add').dataset.id = result.id;
			resultElement.querySelector('.add').addEventListener('click', (e) => {
				this.resultClick(result);
			});
			resultsElement.appendChild(resultElement);
		}
	}

	resultClick(zone) {
		console.debug(zone.id);
	}
}

// ZoneAdd

const zoneAddTemplate = document.createElement('template');
zoneAddTemplate.innerHTML = `
	<style>
		#layout {
			margin-left: 0.75rem;
			margin-right: 0.753rem;
			margin-top: 0.5rem;
			margin-bottom: 0.5rem;
		}
	</style>
	
	<div id="layout">
		<d-zoneadd-search></d-zoneadd-search>
	</div>
`;

class ZoneAddElement extends HTMLElement {
	constructor() {
		super();
		
		this.attachShadow({ mode: 'open' });
		this.shadowRoot.appendChild(zoneAddTemplate.content.cloneNode(true));

		this.searchElement = this.shadowRoot.querySelector('d-zoneadd-search');
		window.s = this.searchElement;
		this.searchElement.show([
			{ n: 'St. John\'s, Newfoundland and Labrador, CA', id: 'St_John\'s-Newfoundland_and_Labrador-CA' },
			{ n: 'Singapore, SG', id: 'Singapore-SG' },
		]);
		this.searchElement.addEventListener('searchinput', this.input.bind(this));
	}

	input({ detail: { text } }) {
		console.debug(text);
	}
}

// Definitions

window.customElements.define('icon-solid-trash', IconSolidTrashElement);
window.customElements.define('icon-solid-plus', IconSolidPlusElement);

window.customElements.define('d-zone', ZoneElement);
window.customElements.define('d-zonename', ZoneNameElement);
window.customElements.define('d-zoneoffset', ZoneOffsetElement);
window.customElements.define('d-date', DateElement);
window.customElements.define('d-time', TimeElement);
// TODO: Instead of relying on <slots>, render using attributes

window.customElements.define('d-zoneadd-search', ZoneSearchElement);
window.customElements.define('d-zoneadd', ZoneAddElement);
