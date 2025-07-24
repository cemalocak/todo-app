import { test, expect } from '@playwright/test';

test.describe('Todo App', () => {
	test.beforeEach(async ({ page }) => {
		await page.goto('/');
	});

	test('should add a new todo', async ({ page }) => {
		// Add a new todo
		await page.fill('[data-testid="todo-input"]', 'New Todo Item');
		await page.click('[data-testid="add-todo-button"]');

		// Verify todo is added
		const todoText = await page.textContent('[data-testid="todo-item"]:first-child');
		expect(todoText).toBe('New Todo Item');
	});

	test('should edit a todo', async ({ page }) => {
		// Add a todo first
		await page.fill('[data-testid="todo-input"]', 'Todo to Edit');
		await page.click('[data-testid="add-todo-button"]');

		// Edit the todo
		await page.click('[data-testid="edit-button"]:first-child');
		await page.fill('[data-testid="edit-input"]', 'Edited Todo');
		await page.click('[data-testid="save-button"]');

		// Verify todo is edited
		const todoText = await page.textContent('[data-testid="todo-item"]:first-child');
		expect(todoText).toBe('Edited Todo');
	});

	test('should delete a todo', async ({ page }) => {
		// Add a todo first
		await page.fill('[data-testid="todo-input"]', 'Todo to Delete');
		await page.click('[data-testid="add-todo-button"]');

		// Get initial todo count
		const initialCount = await page.locator('[data-testid="todo-item"]').count();

		// Delete the todo
		await page.click('[data-testid="delete-button"]:first-child');

		// Verify todo is deleted
		const finalCount = await page.locator('[data-testid="todo-item"]').count();
		expect(finalCount).toBe(initialCount - 1);
	});
}); 