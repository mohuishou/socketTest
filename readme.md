
golang socket 转发测试
====

>  用于接收client端的数据转发到app端

## 数据发送约定

1. 所有的数据采用json字符串发送
2. 建立连接使用TCP协议
3. 连接建立成功之后，需要先发送一次客户确认数据
4. 客户数据发送之后，client才将数据发送给server端

## 数据格式

1. 客户确认json

```json
{
  "id":1
  "type":"client" //或者是 user
}
```

2. 数据json

```json
{
  "lat":"11.3"
  "lon":"111.1"
  "fall":"1" //1：摔倒，0:正常
}
```

