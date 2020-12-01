学习笔记

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
根据老师建议的处理原则，需要在调用3方库的现场调用wrap，然后在api处理的最外层打印以及处理相关error。同时，需要在dao层提供判定error的相关方法。

代码如下
func GetUser() (*user, error) {
	return nil, errors.Wrap(sql.ErrNoRows, "get user failed")
}

func IsNotFound(err error) bool {
	if errors.Cause(err) == sql.ErrNoRows {
		return true
	}
	return false
}