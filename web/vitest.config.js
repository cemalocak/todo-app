import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
	plugins: [react()],
	test: {
		environment: 'jsdom',
		globals: true,
		setupFiles: ['./src/test/setup.js'],
		coverage: {
			provider: 'v8',
			reporter: ['text', 'json', 'html'],
			reportsDirectory: './coverage',
			exclude: [
				'node_modules/',
				'src/test/',
				'**/*.{test,spec}.{js,jsx}',
				'vite.config.js',
				'vitest.config.js',
			],
			branches: 80,
			functions: 80,
			lines: 80,
			statements: 80,
		},
		reporters: ['default', 'junit'],
		outputFile: {
			junit: './test-results/junit.xml'
		}
	},
}); 