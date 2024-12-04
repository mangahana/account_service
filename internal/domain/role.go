package domain

type RoleID int

const (
	UserRole RoleID = iota + 1
	AuthorRole
	ModeratorRole
	AdminRole
	OwnerRole
)

type Role struct {
	ID          RoleID
	Name        string
	Permissions []string
}
