import { useState } from 'react'

function TodoForm({ onAddTodo, disabled }) {
	const [inputValue, setInputValue] = useState('')

	const handleSubmit = (e) => {
		e.preventDefault()
		if (inputValue.trim()) {
			onAddTodo(inputValue.trim())
			setInputValue('')
		}
	}

	const handleInputChange = (e) => {
		setInputValue(e.target.value)
	}

	return (
		<section className="todo-input-section">
			<form onSubmit={handleSubmit} className="input-container">
				<input
					type="text"
					value={inputValue}
					onChange={handleInputChange}
					placeholder="Bugün ne yapacaksınız? (örn: süt al)"
					disabled={disabled}
					autoComplete="off"
					data-testid="todo-input"
					aria-label="Todo text input"
				/>
				<button
					type="submit"
					className="add-btn"
					disabled={disabled || !inputValue.trim()}
					data-testid="add-todo-button"
					aria-label="Add todo button"
				>
					<span>Ekle</span>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
						<line x1="12" y1="5" x2="12" y2="19"></line>
						<line x1="5" y1="12" x2="19" y2="12"></line>
					</svg>
				</button>
			</form>
		</section>
	)
}

export default TodoForm 