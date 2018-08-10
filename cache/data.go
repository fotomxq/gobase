package cache

//数据结构
type CacheGoData struct {
	//标识
	Mark string
	//过期时间戳
	ExpireTime int64
	//数据内容
	Content []byte
}
