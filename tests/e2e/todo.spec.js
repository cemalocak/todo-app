// @ts-check
const { test, expect } = require('@playwright/test');

// Test environment URL
const BASE_URL = process.env.TEST_URL || 'http://localhost:3001';

test.describe('Todo App E2E Tests', () => {
	test.beforeEach(async ({ page, request }) => {
		// Clear test database before each test
		try {
			const response = await request.post(`${BASE_URL}/api/test/truncate`);
			if (!response.ok()) {
				console.warn('Failed to truncate test database:', await response.text());
			}
		} catch (error) {
			console.warn('Error truncating test database:', error.message);
		}

		// Navigate to the app
		await page.goto(BASE_URL);
	});

	test('User Story: Add "süt al" and see it in the list', async ({ page }) => {
		// Given: User is on the todo app page
		await expect(page).toHaveTitle(/Todo App/);

		// When: User types "süt al" in the input
		await page.fill('[data-testid="todo-input"]', 'süt al');

		// And: User clicks the "Add" button
		await page.click('[data-testid="add-todo-button"]');

		// Then: "süt al" should appear in the list
		const todoItem = page.locator('[data-testid="todo-item"]').first();
		await expect(todoItem.locator('.todo-text')).toContainText('süt al');

		// And: Input should be cleared
		await expect(page.locator('[data-testid="todo-input"]')).toHaveValue('');
	});

	test('Empty state is shown correctly', async ({ page }) => {
		// Check empty state message
		await expect(page.locator('[data-testid="empty-state"]')).toBeVisible();
		await expect(page.locator('[data-testid="empty-state"] h3')).toContainText('Henüz görev yok');
	});

	test('Todo CRUD operations work correctly', async ({ page }) => {
		// Create
		await page.fill('[data-testid="todo-input"]', 'Test todo');
		await page.click('[data-testid="add-todo-button"]');

		// Get the first todo item since we know it's the only one
		const todoItem = page.locator('[data-testid="todo-item"]').first();
		await expect(todoItem.locator('.todo-text')).toContainText('Test todo');

		// Update 
		await todoItem.locator('[data-testid="edit-button"]').click();
		await page.fill('[data-testid="edit-input"]', 'Updated todo');
		await page.click('[data-testid="save-button"]');
		await expect(todoItem.locator('.todo-text')).toContainText('Updated todo');

		// Delete 
		await todoItem.locator('[data-testid="delete-button"]').click();
		await expect(todoItem).not.toBeVisible();
	});

	test('API integration works correctly', async ({ page, request }) => {
		// Test direct API call
		const response = await request.get(`${BASE_URL}/api/todos`);
		expect(response.status()).toBe(200);

		const todos = await response.json();
		expect(Array.isArray(todos)).toBe(true);
		expect(todos.length).toBe(0); // Should be empty after beforeEach cleanup
	});

	test('App is responsive on mobile', async ({ page }) => {
		// Set mobile viewport
		await page.setViewportSize({ width: 375, height: 667 });

		// Check if app is still functional
		await page.fill('[data-testid="todo-input"]', 'Mobile todo');
		await page.click('[data-testid="add-todo-button"]');

		// Check the first (and only) todo item
		const todoItem = page.locator('[data-testid="todo-item"]').first();
		await expect(todoItem.locator('.todo-text')).toContainText('Mobile todo');
	});

	test('Accessibility: Basic a11y checks', async ({ page }) => {
		// Check for basic accessibility features
		await expect(page.locator('[data-testid="todo-input"]')).toHaveAttribute('aria-label', 'Todo text input');
		await expect(page.locator('[data-testid="add-todo-button"]')).toHaveAttribute('aria-label', 'Add todo button');

		// Check for proper heading structure
		await expect(page.locator('h2')).toContainText('Görevler');
	});
}); 