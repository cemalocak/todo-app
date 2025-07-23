import { useState, useEffect } from 'react'
import axios from 'axios'
import TodoForm from './components/TodoForm'
import TodoList from './components/TodoList'
import './TodoApp.css'

function TodoApp() {
	const [todos, setTodos] = useState([])
	const [loading, setLoading] = useState(false)

	// Fetch todos on component mount
	useEffect(() => {
		fetchTodos()
	}, [])

	const fetchTodos = async () => {
		try {
			setLoading(true)
			const response = await axios.get('/api/todos')
			setTodos(response.data)
		} catch (error) {
			console.error('Error fetching todos:', error)
		} finally {
			setLoading(false)
		}
	}

	const addTodo = async (text) => {
		if (!text.trim()) return

		try {
			setLoading(true)
			const response = await axios.post('/api/todos', { text })
			setTodos(prevTodos => [...prevTodos, response.data])
		} catch (error) {
			console.error('Error adding todo:', error)
		} finally {
			setLoading(false)
		}
	}

	return (
		<div className="todo-app">
			<header className="todo-header">
				<h1 data-testid="app-title">📝 ToDo Listesi</h1>
				<p>Günlük görevlerinizi organize edin</p>
			</header>

			<main className="todo-main">
				<TodoForm onAddTodo={addTodo} disabled={loading} />
				<TodoList todos={todos} loading={loading} />
			</main>

			<footer className="todo-footer">
				<p>A-TDD ile ❤️ kullanılarak geliştirildi</p>
			</footer>
		</div>
	)
}

export default TodoApp 