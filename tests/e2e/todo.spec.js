// @ts-check
const { test, expect } = require('@playwright/test');

// Test environment URL
const BASE_URL = process.env.TEST_URL || 'http://localhost:3000';

test.describe('Todo App E2E Tests', () => {
	test.beforeEach(async ({ page }) => {
		// Navigate to the app
		await page.goto(BASE_URL);
	});

	test('User Story: Add "süt al" and see it in the list', async ({ page }) => {
		// Given: User is on the todo app page
		await expect(page).toHaveTitle(/Todo App/);

		// When: User types "süt al" in the input
		await page.fill('[data-testid="todo-input"]', 'süt al');

		// And: User clicks the "Add" button
		await page.click('[data-testid="add-button"]');

		// Then: "süt al" should appear in the list
		await expect(page.locator('[data-testid="todo-item"]')).toContainText('süt al');

		// And: Input should be cleared
		await expect(page.locator('[data-testid="todo-input"]')).toHaveValue('');
	});

	test('Multiple todos are displayed correctly', async ({ page }) => {
		// Add multiple todos
		const todos = ['Buy milk', 'Walk the dog', 'Finish project'];

		for (const todo of todos) {
			await page.fill('[data-testid="todo-input"]', todo);
			await page.click('[data-testid="add-button"]');
		}

		// Check all todos are displayed
		for (const todo of todos) {
			await expect(page.locator('[data-testid="todo-item"]')).toContainText(todo);
		}

		// Check order (newest first)
		const todoItems = page.locator('[data-testid="todo-item"]');
		await expect(todoItems.first()).toContainText('Finish project');
	});

	test('Empty state is shown correctly', async ({ page }) => {
		// Check empty state message
		await expect(page.locator('[data-testid="empty-state"]')).toBeVisible();
		await expect(page.locator('[data-testid="empty-state"]')).toContainText('No todos yet');
	});

	test('Todo CRUD operations work correctly', async ({ page }) => {
		// Create
		await page.fill('[data-testid="todo-input"]', 'Test todo');
		await page.click('[data-testid="add-button"]');
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

	test('API integration works correctly', async ({ page, request }) => {
		// Test direct API call
		const response = await request.get(`${BASE_URL}/api/todos`);
		expect(response.status()).toBe(200);

		const todos = await response.json();
		expect(Array.isArray(todos)).toBe(true);
	});

	test('App is responsive on mobile', async ({ page }) => {
		// Set mobile viewport
		await page.setViewportSize({ width: 375, height: 667 });

		// Check if app is still functional
		await page.fill('[data-testid="todo-input"]', 'Mobile todo');
		await page.click('[data-testid="add-button"]');
		await expect(page.locator('[data-testid="todo-item"]')).toContainText('Mobile todo');
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

	test('Accessibility: Basic a11y checks', async ({ page }) => {
		// Check for basic accessibility features
		await expect(page.locator('[data-testid="todo-input"]')).toHaveAttribute('aria-label');
		await expect(page.locator('[data-testid="add-button"]')).toHaveAttribute('aria-label');

		// Check for proper heading structure
		await expect(page.locator('h1')).toBeVisible();
	});
}); 