// @Author : detaohe
// @File   : rewrite.go
// @Description:
// @Date   : 2022/9/8 21:07

package middleware

/*
rewrite和redirect的区别
redirect是在客户端的角度，客户端发送到服务器，服务器返回301和重定向的地址，客户端自动请求新地址
rewrite 是在服务器的角度 /resource 重写到 /different-resource 时，客户端会请求 /resource ，
并且服务器会在内部提取 /different-resource 处的资源。尽管客户端可能能够检索已重写URL处的资源
*/
