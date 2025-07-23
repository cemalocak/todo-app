function TodoList({ todos, loading }) {
	const todoCount = todos.length

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
							<div className="todo-content">
								<span className="todo-text">{todo.text}</span>
								<span className="todo-id">#{todo.id}</span>
							</div>
						</li>
					))}
				</ul>
			)}
		</section>
	)
}

export default TodoList 