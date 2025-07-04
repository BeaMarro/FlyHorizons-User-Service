package enums

type AccountType int

const (
	Admin AccountType = 0
	User  AccountType = 1
)

func AccountTypeFromInt(value int) AccountType {
	switch value {
	case 0:
		return Admin
	case 1:
		return User
	default:
		return User
	}
}
