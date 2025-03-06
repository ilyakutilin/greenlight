package data

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Define a Permissions slice, which will be used to hold the permission codes (like
// "movies:read" and "movies:write") for a single user.
type Permissions []string

// Add a helper method to check whether the Permissions slice contains a specific
// permission code.
func (p Permissions) Include(code string) bool {
	return slices.Contains(p, code)
}

// Define the PermissionModel type.
type PermissionModel struct {
	DB *sql.DB
}

// Returns all permission codes for a specific user in a Permissions slice.
func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
        SELECT permissions.code
        FROM permissions
        INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
        INNER JOIN users ON users_permissions.user_id = users.id
        WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// Add the provided permission codes for a specific user.
func (m PermissionModel) AddForUser(userID int64, codes ...string) error {
	// In this query the $1 parameter is the user’s ID, and the $2 parameter is a
	// PostgreSQL array of the permission codes that we want to add for the user,
	// like {'movies:read', 'movies:write'}. The SELECT ... statement on the second line
	// creates an ‘interim’ table with rows made up of the user ID and the corresponding
	// IDs for the permission codes in the array. Then we insert the contents of this
	// interim table into our user_permissions table.
	query := `
        INSERT INTO users_permissions
        SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(codes))
	return err
}
