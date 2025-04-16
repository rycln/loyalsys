package auth

type PasswordAdapter struct {
	hashFunc    func(string) (string, error)
	compareFunc func(string, string) error
}

func NewPasswordAdapter(
	hashFunc func(string) (string, error),
	compareFunc func(string, string) error,
) *PasswordAdapter {
	return &PasswordAdapter{
		hashFunc:    hashFunc,
		compareFunc: compareFunc,
	}
}

func (a *PasswordAdapter) Hash(password string) (string, error) {
	return a.hashFunc(password)
}

func (a *PasswordAdapter) Compare(hashed, plain string) error {
	return a.compareFunc(hashed, plain)
}
