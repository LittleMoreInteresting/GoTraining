## socket 粘包的解包方式

|解包方式|说明|优缺点|应用|
|:---- |:---- |:---- |:----|
| fix length | 固定缓冲区大小：控制服务器端和客户端发送和接收字节的长度固定 | 逻辑简单，但可能会增加不必要的数据传输 |服务器心跳检测|
| delimiter based | 以特殊的分隔符结尾（如：\n） | 简单易用  | 按行来读、写和区分消息 |
| length field based | 数据头+数据正文 | 数据包结构更加灵活，方便扩充，实现逻辑相对复杂  | IM通讯 |