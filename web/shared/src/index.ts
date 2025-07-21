import type { App } from 'vue';
import * as components from './components';

export default {
	install (app: App, _options = {}) {
		for (const key in components) {
			app.component(key, (components as any)[key]);
		}
	}
}
