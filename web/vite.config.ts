import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import fs from 'fs';

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		https: {
			key: fs.readFileSync('./certs/localhost+2-key.pem'),
			cert: fs.readFileSync('./certs/localhost+2.pem')
		},
		host: 'localhost',
		port: 5173
	}
});
