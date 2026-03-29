package constants

// UserRole define os papéis disponíveis no sistema
type UserRole string

const (
	AdminRole UserRole = "admin"
	RoleUser  UserRole = "user"
)

// Query maps roles to string values for database storage
func (r UserRole) String() string {
	return string(r)
}

// IsAdmin verifica se a role é admin
func (r UserRole) IsAdmin() bool {
	return r == AdminRole
}

// IsValid valida se a role é válida
func (r UserRole) IsValid() bool {
	return r == AdminRole || r == RoleUser
}

// ValidRoles retorna a lista de roles válidas
func ValidRoles() []UserRole {
	return []UserRole{AdminRole, RoleUser}
}
