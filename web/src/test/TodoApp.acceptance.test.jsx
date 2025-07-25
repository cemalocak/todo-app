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
			expect(screen.getByText('alışveriş yap')).toBeInTheDocument()
			expect(screen.getByText('süt al')).toBeInTheDocument()
		})

		// And: Todo count should be displayed
		expect(screen.getByText('2 görev')).toBeInTheDocument()
	})
}) 