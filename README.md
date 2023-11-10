# README 

## 【 LOG 】

    1. 4000/debug/pprof/ up and runing on my local machine
    2. gracefully shutdown , chatper4 
        
        {"started","host":"0.0.0.0:4000"}
        {"level":"info","ts":"2023-11-01T09:14:42.796Z","caller":"sales-api/main.go:151","msg":"startup","service":"SALES-API","status":"initializing V1 API support"}
        {"level":"info","ts":"2023-11-01T09:14:42.796Z","caller":"sales-api/main.go:174","msg":"startup","service":"SALES-API","status":"api router started","host":"0.0.0.0:3000"}
        {"level":"info","ts":"2023-11-01T09:15:18.921Z","caller":"sales-api/main.go:187","msg":"shutdown","service":"SALES-API","status":"shutdown started","signal":"terminated"}
        {"level":"info","ts":"2023-11-01T09:15:18.923Z","caller":"sales-api/main.go:201","msg":"shutdown","service":"SALES-API","status":"shutdown complete","signal":"terminated"}
        {"level":"info","ts":"2023-11-01T09:15:18.923Z","caller":"sales-api/main.go:201","msg":"shutdown complete","service":"SALES-API"}
        rpc error: code = NotFound desc = an error occurred when try to find container "7f8e2cb38aaa7be9777d02af5874c3025fa9e66e24a90bf575788f1961ddc3e4": not found%

    3. 增加 https://github.com/dimfeld/httptreemux 相关依赖
    4. 重构 APIMux， 变成洋葱式和其他方式的切片架构。
    5. 增加 JWT https://jwt.io/ 
    6. 

