package service

type Request struct {
	Query     string
	Shardid   uint64
	Sharduid  string
	UserToken string
}
