package postgres

const (
	queryInsertTodo = `
		INSERT INTO todos (id, title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	queryGetTodoByID = `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1`

	queryListTodos = `
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC`

	queryUpdateTodo = `
		UPDATE todos
		SET title = $2, description = $3, completed = $4, updated_at = $5
		WHERE id = $1`

	queryDeleteTodo = `
		DELETE FROM todos WHERE id = $1`

	queryCompleteAll = `
		UPDATE todos
		SET completed = TRUE, updated_at = NOW()
		WHERE completed = FALSE`
)
