/*
math-engine

--------

数学表达式解析计算引擎库

使用 Go 实现的数学表达式解析计算引擎库，它小巧，无任何依赖，具有扩展性(比如可以注册自己的函数到引擎中)，比较完整的完成了数学表达式解析执行，包括词法分析、语法分析、构建AST、运行。

`go get -u github.com/dengsgo/math-engine`

能够处理的表达式样例：
- `1+127-21+(3-4)*6/2.5`
- `(88+(1+8)*6)/2+99`
- `123_345_456 * 1.5 - 2 ^ 4`
- `-4 * 6 + 2e2 - 1.6e-3`
- `sin(pi/2)+cos(45-45*1)+tan(pi/4)`
- `99+abs(-1)-ceil(88.8)+floor(88.8)`
- `max(min(2^3, 3^2), 10*1.5-7)`
- `double(6) + 3` , `double`是一个自定义的函数

*/
package libMath
