import { useState } from 'react'

function TodoList({ todos, loading, onUpdateTodo, onDeleteTodo }) {
	const [editingId, setEditingId] = useState(null)
	const [editingText, setEditingText] = useState('')
	const todoCount = todos.length

	const handleEditStart = (todo) => {
		setEditingId(todo.id)
		setEditingText(todo.text)
	}

	const handleEditSave = async (id) => {
		if (editingText.trim()) {
			await onUpdateTodo(id, editingText.trim())
		}
		setEditingId(null)
		setEditingText('')
	}

	const handleEditCancel = () => {
		setEditingId(null)
		setEditingText('')
	}

	const handleDelete = async (id) => {
		await onDeleteTodo(id)
	}

	if (loading) {
		return (
			<section className="todo-list-section">
				<div className="loading-state">
					<div className="spinner"></div>
					<p>Yükleniyor...</p>
				</div>
			</section>
		)
	}

	return (
		<section className="todo-list-section">
			<div className="list-header">
				<h2>Görevler</h2>
				<span className="todo-count">
					{todoCount} görev
				</span>
			</div>

			{todoCount === 0 ? (
				<div className="empty-state" data-testid="empty-state">
					<div className="empty-icon">📋</div>
					<h3>Henüz görev yok</h3>
					<p>Yukarıdaki form ile ilk görevinizi ekleyin!</p>
				</div>
			) : (
				<ul className="todo-list">
					{todos.map(todo => (
						<li key={todo.id} className="todo-item" data-testid="todo-item">
							{editingId === todo.id ? (
								// Edit Mode
								<div className="todo-edit-mode">
									<input
										type="text"
										value={editingText}
										onChange={(e) => setEditingText(e.target.value)}
										className="edit-input"
										autoFocus
									/>
									<div className="edit-actions">
										<button
											onClick={() => handleEditSave(todo.id)}
											className="save-btn"
											aria-label="Kaydet"
											disabled={!editingText.trim()}
										>
											✓ Kaydet
										</button>
										<button
											onClick={handleEditCancel}
											className="cancel-btn"
											aria-label="İptal"
										>
											✕ İptal
										</button>
									</div>
								</div>
							) : (
								// View Mode
								<div className="todo-content">
									<div className="todo-text-section">
										<span className="todo-text">{todo.text}</span>
										<span className="todo-id">#{todo.id}</span>
									</div>
									<div className="todo-actions">
										<button
											onClick={() => handleEditStart(todo)}
											className="edit-btn"
											aria-label="Düzenle"
										>
											✏️ Düzenle
										</button>
										<button
											onClick={() => handleDelete(todo.id)}
											className="delete-btn"
											aria-label="Sil"
										>
											🗑️ Sil
										</button>
									</div>
								</div>
							)}
						</li>
					))}
				</ul>
			)}
		</section>
	)
}

export default TodoList 