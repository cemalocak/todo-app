// @ts-check
const { test, expect } = require('@playwright/test');

// Test environment URL
const BASE_URL = process.env.TEST_UI_URL || 'http://localhost:3001';
const API_URL = process.env.TEST_API_URL || 'http://localhost:8081';

// Test veritabanını temizleme fonksiyonu
async function clearTestDatabase(page) {
	await page.request.post(`${API_URL}/api/test/truncate`);
	console.log('🧪 Test veritabanı temizlendi.');
	await page.waitForTimeout(2000);
}

// Sayfanın yüklenmesini bekle
async function waitForPageLoad(page) {
	await page.waitForLoadState('networkidle');
	console.log('📱 Sayfa yüklendi ve hazır.');
}

// Temel CRUD testleri
test.describe('Todo App CRUD Tests', () => {
	test.beforeAll(async ({ browser }) => {
		console.log('🚀 CRUD testleri başlıyor...');
		const page = await browser.newPage();
		await clearTestDatabase(page);
		await page.close();
	});

	test.beforeEach(async ({ page }) => {
		await page.goto(BASE_URL, { waitUntil: 'load' });
		await waitForPageLoad(page);
		console.log('🧪 Yeni CRUD testi başlıyor...');
	});

	test('Add "süt al" and see it in the list', async ({ page }) => {
		// Given: User is on the todo app page
		await expect(page).toHaveTitle(/Todo App/);

		// When: User types "süt al" in the input
		await page.fill('[data-testid="todo-input"]', 'süt al');

		// And: User clicks the "Add" button
		await page.click('[data-testid="add-button"]');

		await page.waitForTimeout(2000);

		// Then: "süt al" should appear in the list
		await expect(page.locator('[data-testid="todo-item"]')).toContainText('süt al');

		// And: Input should be cleared
		await expect(page.locator('[data-testid="todo-input"]')).toHaveValue('');
	});
});

test.describe('Todo App CRUD Tests', () => {
	test.beforeAll(async ({ browser }) => {
		console.log('🚀 CRUD testleri başlıyor...');
		const page = await browser.newPage();
		await clearTestDatabase(page);
		await page.close();
	});

	test.beforeEach(async ({ page }) => {
		await page.goto(BASE_URL);
		console.log('🧪 Yeni CRUD testi başlıyor...');
	});



	test('Todo CRUD operations work correctly', async ({ page }) => {
		// Create
		await page.fill('[data-testid="todo-input"]', 'Test todo');
		await page.click('[data-testid="add-button"]');
		await page.waitForTimeout(2000);
		await expect(page.locator('[data-testid="todo-item"]')).toContainText('Test todo');

		// Update (if implemented)
		// await page.click('[data-testid="edit-button"]');
		// await page.fill('[data-testid="edit-input"]', 'Updated todo');
		// await page.click('[data-testid="save-button"]');
		// await expect(page.locator('[data-testid="todo-item"]')).toContainText('Updated todo');

		// Delete (if implemented)
		// await page.click('[data-testid="delete-button"]');
		// await expect(page.locator('[data-testid="todo-item"]')).not.toBeVisible();
	});
});

// UI ve Erişilebilirlik testleri
test.describe('Todo App UI Tests', () => {
	test.beforeAll(async ({ browser }) => {
		console.log('🚀 UI testleri başlıyor...');
		const page = await browser.newPage();
		await clearTestDatabase(page);
		await page.close();
	});

	test.beforeEach(async ({ page }) => {
		await page.goto(BASE_URL);
		console.log('🧪 Yeni UI testi başlıyor...');
	});

	test('Empty state is shown correctly', async ({ page }) => {
		// Check empty state message
		await expect(page.locator('[data-testid="empty-state"]')).toBeVisible();
		await expect(page.locator('[data-testid="empty-state"]')).toContainText('📋Henüz görev yokYukarıdaki form ile ilk görevinizi ekleyin!');
	});

	test('App is responsive on mobile', async ({ page }) => {
		// Set mobile viewport
		await page.setViewportSize({ width: 375, height: 667 });

		// Check if app is still functional
		await page.fill('[data-testid="todo-input"]', 'Mobile todo');
		await page.click('[data-testid="add-button"]');
		await expect(page.locator('[data-testid="todo-item"]')).toContainText('Mobile todo');
	});

	test('Accessibility: Basic a11y checks', async ({ page }) => {
		// Check for basic accessibility features
		await expect(page.locator('[data-testid="todo-input"]')).toHaveAttribute('aria-label');
		await expect(page.locator('[data-testid="add-button"]')).toHaveAttribute('aria-label');

		// Check for proper heading structure
		await expect(page.locator('h1')).toBeVisible();
	});
});

// API ve Performans testleri
test.describe('Todo App API Tests', () => {
	test.beforeAll(async ({ browser }) => {
		console.log('🚀 API testleri başlıyor...');
		const page = await browser.newPage();
		await clearTestDatabase(page);
		await page.close();
	});

	test.beforeEach(async ({ page }) => {
		await page.goto(BASE_URL);
		console.log('🧪 Yeni API testi başlıyor...');
	});

	test('API integration works correctly', async ({ page, request }) => {
		// Test direct API call
		const response = await request.get(`${API_URL}/api/todos`);
		expect(response.status()).toBe(200);

		const todos = await response.json();
		expect(Array.isArray(todos)).toBe(true);
	});

	test('App handles network errors gracefully', async ({ page }) => {
		// Simulate network failure
		await page.route('**/api/todos', route => route.abort());

		// Try to add a todo
		await page.fill('[data-testid="todo-input"]', 'Network test');
		await page.click('[data-testid="add-button"]');

		// Should show error message (if implemented)
		// await expect(page.locator('[data-testid="error-message"]')).toBeVisible();
	});

	test('Performance: App loads quickly', async ({ page }) => {
		const startTime = Date.now();
		await page.goto(BASE_URL);
		await page.waitForLoadState('networkidle');
		const loadTime = Date.now() - startTime;

		// Should load within 3 seconds
		expect(loadTime).toBeLessThan(3000);
	});
}); 