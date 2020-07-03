package typefile

//structを複数ファイルに跨って使用したいので、専用のファイルを作成

//Token そのトークンのワードと次へのポインタ、そのトークンの種類を所持する
type Token struct {
	Word  string
	Next  *Token
	Tkind string
}
