package accessLevel

type AccessLevel int64

const (
	None AccessLevel = iota
	User
	Admin
)

func (al AccessLevel) String() string {
	switch al {
	case None:
		return "None"
	case User:
		return "User"
	case Admin:
		return "Admin"
	}
	return "unknown"
}
