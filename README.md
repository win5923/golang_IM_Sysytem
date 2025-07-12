# V1: Build basic server

```
go build -o server main.go server.go
./server
```

Create another terminal and run the following command:
```
nc 127.0.0.1 8888
```

# V2: User online and broadcast

```
go build -o server main.go server.go user.go
./server
```


# V3: User's messages broadcast
```
go build -o server main.go server.go user.go
./server
```

Create multiple terminals and run the following command:
```
nc 127.0.0.1 8888
test
```

# V4: Refactor the user's logic

Ref: 【8小时转职Golang工程师(如果你想低成本学习Go语言)】 https://www.bilibili.com/video/BV1gf4y1r79E/?share_source=copy_web&vd_source=cbd019538643d10ff09a29aebb9b1099