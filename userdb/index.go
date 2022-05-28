package userdb

type Index string

const (
	NicknameIndex Index = "nickname_idx"
)

func (index Index) String() string {
	return string(index)
}
