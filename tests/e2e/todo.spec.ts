import { test, expect } from '@playwright/test';

test.describe('Todo App', () => {
	test.beforeEach(async ({ page, request }) => {
		// Clear test database before each test
		try {
			const response = await request.post('/api/test/truncate');
			if (!response.ok()) {
				console.warn('Failed to truncate test database:', await response.text());
			}
		} catch (error) {
			console.warn('Error truncating test database:', error.message);
		}

		await page.goto('/');
	});

	test('should add a new todo', async ({ page }) => {
		// Add a new todo
		await page.fill('[data-testid="todo-input"]', 'New Todo Item');
		await page.click('[data-testid="add-todo-button"]');

		// Verify todo is added
		const todoItem = page.locator('[data-testid="todo-item"]').first();
		await expect(todoItem.locator('.todo-text')).toContainText('New Todo Item');
	});

	test('should edit a todo', async ({ page }) => {
		// Add a todo to edit
		await page.fill('[data-testid="todo-input"]', 'Todo to Edit');
		await page.click('[data-testid="add-todo-button"]');

		// Edit the todo
		const todoItem = page.locator('[data-testid="todo-item"]').first();
		await todoItem.locator('[data-testid="edit-button"]').click();
		await page.fill('[data-testid="edit-input"]', 'Edited Todo');
		await page.click('[data-testid="save-button"]');

		// Verify todo is edited
		await expect(todoItem.locator('.todo-text')).toContainText('Edited Todo');
	});

	test('should delete a todo', async ({ page }) => {
		// Add a todo to delete
		await page.fill('[data-testid="todo-input"]', 'Todo to Delete');
		await page.click('[data-testid="add-todo-button"]');

		// Delete the todo
		const todoItem = page.locator('[data-testid="todo-item"]').first();
		await todoItem.locator('[data-testid="delete-button"]').click();

		// Verify todo is deleted
		const finalCount = await page.locator('[data-testid="todo-item"]').count();
		expect(finalCount).toBe(0);
	});
}); 