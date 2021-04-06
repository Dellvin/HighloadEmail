DataBase Models

type Letter struct {
	Id            uint64
	Sender        string
	Receiver      string
	Theme         string
	Text          string
	DateTime      int64
	IsWatched     bool
	DirectoryRecv uint64
	DirectorySend uint64
	Spam          bool
	Box           bool
}

type Folder struct {
	Id   uint64
	Uid  uint64
	Type string
	Name string
}

type User struct {
	Id       uint64
	Name     string
	Surname  string
	Email    string
	Password string
	Img      string
}

type Session struct {
	Id     string
	UserId int64 `pg:"on_delete:RESTRICT,on_update: CASCADE"`
	User   *User `pg:"rel:has-one"`
}
