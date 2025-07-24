// @ts-check
const { defineConfig, devices } = require('@playwright/test');

/**
 * @see https://playwright.dev/docs/test-configuration
 */
module.exports = defineConfig({
	testDir: './',
	/* Run tests in files in parallel */
	fullyParallel: false,
	/* Fail the build on CI if you accidentally left test.only in the source code. */
	forbidOnly: !!process.env.CI,
	/* Retry on CI only */
	retries: process.env.CI ? 2 : 0,
	/* Opt out of parallel tests on CI. */
	workers: 1,
	/* Reporter to use. See https://playwright.dev/docs/test-reporters */
	reporter: [
		['html'],
		['json', { outputFile: 'test-results.json' }],
		['junit', { outputFile: 'test-results.xml' }]
	],
	/* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
	use: {
		/* Base URL to use in actions like `await page.goto('/')`. */
		baseURL: process.env.TEST_UI_URL || 'http://localhost:3000',

		/* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
		trace: 'on-first-retry',

		/* Take screenshot on failure */
		screenshot: 'only-on-failure',

		/* Record video on failure */
		video: 'retain-on-failure',

		/* Timeout for each action */
		actionTimeout: 30000,

		/* Navigation timeout */
		navigationTimeout: 30000,

		/* Wait until page load state */
		waitUntil: 'networkidle',
	},

	/* Configure projects for major browsers */
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] },
		},
	],

	/* Run your local dev server before starting the tests */
	webServer: process.env.CI ? undefined : {
		command: 'cd ../../ && make up',
		url: 'http://localhost:3000',
		timeout: 120 * 1000,
		reuseExistingServer: !process.env.CI,
	},
}); 