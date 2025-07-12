package dari

type Key interface {
	Key() (string, string)
}

type KeyName string
type KeyValues [2]string
type KeySet map[KeyName]KeyValues

type Keys interface {
	Keys() KeySet
}
