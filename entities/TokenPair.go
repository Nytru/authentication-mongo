package entities

type TokenPair struct {
	AccessToken  string
	RefreshToken *RefreshToken
}

func (tp TokenPair) IsValid() bool {
	return tp.AccessToken != "" && tp.RefreshToken != nil
}
