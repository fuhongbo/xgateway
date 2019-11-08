# xGateway 

这是一个api网关的项目，可以做负载均衡等。

当前支持如下功能：

1. 负载均衡支持

    -> RoundRobin

    -> Random

    -> Consistent Hashing
    
    -> Bounded Consistent Hashing

2. Basic 认证设置

3. Rate Limit

4. IP过滤

5. JWT 认证

6. 发送中间请求

7. 断言

8. 变量计算

9. 健康检查

通过上面的一些功能，可以在请求送达后端服务器之前做一些预判断及处理，做认证等操作。详细可参考 configbak.yaml

