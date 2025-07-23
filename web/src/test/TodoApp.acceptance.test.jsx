import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { vi } from 'vitest'
import TodoApp from '../TodoApp'

// Mock axios for API calls
vi.mock('axios', () => ({
	default: {
		get: vi.fn(),
		post: vi.fn(),
		put: vi.fn(),
		delete: vi.fn(),
	}
}))

import axios from 'axios'

describe('TodoApp Acceptance Tests', () => {
	beforeEach(() => {
		vi.clearAllMocks()
	})

	test('User Story: Add todo item and see it in the list', async () => {
		// Given: Empty todo list (API returns empty array)
		axios.get.mockResolvedValue({ data: [] })
		axios.post.mockResolvedValue({
			data: { id: 1, text: 'süt al' }
		})

		const user = userEvent.setup()
		render(<TodoApp />)

		// Wait for initial load
		await waitFor(() => {
			expect(screen.getByText(/henüz görev yok/i)).toBeInTheDocument()
		})

		// When: User types "süt al" and clicks "Ekle" button
		const input = screen.getByPlaceholderText(/bugün ne yapacaksınız/i)
		const addButton = screen.getByRole('button', { name: /add todo button/i })

		await user.type(input, 'süt al')
		await user.click(addButton)

		// Then: "süt al" should appear in the todo list
		await waitFor(() => {
			expect(screen.getByText('süt al')).toBeInTheDocument()
		})

		// And: Input should be cleared
		expect(input).toHaveValue('')

		// And: API calls should be made correctly
		expect(axios.post).toHaveBeenCalledWith('/api/todos', {
			text: 'süt al'
		})
	})

	test('User Story: Edit existing todo item', async () => {
		// Given: Todo list with one item
		const mockTodos = [
			{ id: 1, text: 'süt al' }
		]
		axios.get.mockResolvedValue({ data: mockTodos })
		axios.put.mockResolvedValue({
			data: { id: 1, text: 'organik süt al' }
		})

		const user = userEvent.setup()
		render(<TodoApp />)

		// Wait for todos to load
		await waitFor(() => {
			expect(screen.getByText('süt al')).toBeInTheDocument()
		})

		// When: User clicks edit button
		const editButton = screen.getByRole('button', { name: /düzenle/i })
		await user.click(editButton)

		// Then: Todo should become editable (input field appears)
		const editInput = screen.getByDisplayValue('süt al')
		expect(editInput).toBeInTheDocument()

		// When: User changes text and saves
		await user.clear(editInput)
		await user.type(editInput, 'organik süt al')
		const saveButton = screen.getByRole('button', { name: /kaydet/i })
		await user.click(saveButton)

		// Then: Updated todo should appear in the list
		await waitFor(() => {
			expect(screen.getByText('organik süt al')).toBeInTheDocument()
		})

		// And: API call should be made correctly
		expect(axios.put).toHaveBeenCalledWith('/api/todos/1', {
			text: 'organik süt al'
		})
	})

	test('User Story: Delete existing todo item', async () => {
		// Given: Todo list with two items
		const mockTodos = [
			{ id: 1, text: 'süt al' },
			{ id: 2, text: 'alışveriş yap' }
		]
		axios.get.mockResolvedValue({ data: mockTodos })
		axios.delete.mockResolvedValue({})

		const user = userEvent.setup()
		render(<TodoApp />)

		// Wait for todos to load
		await waitFor(() => {
			expect(screen.getByText('süt al')).toBeInTheDocument()
			expect(screen.getByText('alışveriş yap')).toBeInTheDocument()
		})

		// When: User clicks delete button for first todo
		const deleteButtons = screen.getAllByRole('button', { name: /sil/i })
		await user.click(deleteButtons[0])

		// Then: Todo should be removed from the list
		await waitFor(() => {
			expect(screen.queryByText('süt al')).not.toBeInTheDocument()
		})

		// And: Other todo should still be visible
		expect(screen.getByText('alışveriş yap')).toBeInTheDocument()

		// And: API call should be made correctly
		expect(axios.delete).toHaveBeenCalledWith('/api/todos/1')
	})

	test('User Story: Cancel edit operation', async () => {
		// Given: Todo list with one item
		const mockTodos = [
			{ id: 1, text: 'süt al' }
		]
		axios.get.mockResolvedValue({ data: mockTodos })

		const user = userEvent.setup()
		render(<TodoApp />)

		// Wait for todos to load
		await waitFor(() => {
			expect(screen.getByText('süt al')).toBeInTheDocument()
		})

		// When: User clicks edit button
		const editButton = screen.getByRole('button', { name: /düzenle/i })
		await user.click(editButton)

		// And: User changes text but cancels
		const editInput = screen.getByDisplayValue('süt al')
		await user.clear(editInput)
		await user.type(editInput, 'organik süt al')

		const cancelButton = screen.getByRole('button', { name: /İptal/i })
		await user.click(cancelButton)

		// Then: Original todo should be restored
		await waitFor(() => {
			expect(screen.getByText('süt al')).toBeInTheDocument()
		})

		// And: No API call should be made
		expect(axios.put).not.toHaveBeenCalled()
	})

	test('User Story: Empty state is shown when no todos', async () => {
		// Given: API returns empty array
		axios.get.mockResolvedValue({ data: [] })

		render(<TodoApp />)

		// Then: Empty state should be visible
		await waitFor(() => {
			expect(screen.getByText(/henüz görev yok/i)).toBeInTheDocument()
			expect(screen.getByText(/ilk görevinizi ekleyin/i)).toBeInTheDocument()
		})
	})

	test('User Story: Multiple todos are displayed correctly', async () => {
		// Given: API returns multiple todos
		const mockTodos = [
			{ id: 1, text: 'süt al' },
			{ id: 2, text: 'alışveriş yap' }
		]
		axios.get.mockResolvedValue({ data: mockTodos })

		render(<TodoApp />)

		// Then: All todos should be visible
		await waitFor(() => {
			expect(screen.getByText('süt al')).toBeInTheDocument()
			expect(screen.getByText('alışveriş yap')).toBeInTheDocument()
		})

		// And: Todo count should be displayed
		expect(screen.getByText('2 görev')).toBeInTheDocument()
	})
}) 