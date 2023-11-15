# README 

## 【 LOG 】

    来源于 B 站的视频： https://www.bilibili.com/video/BV1nZ4y1S7LZ?p=50&spm_id_from=pageDriver&vd_source=fd3a715c532ca990e4a8a6902ac1478c

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
    6. DB lab 相关的内容
        1. https://github.com/danvergara/dblab 
        2. 成功运行本地测试数据库📊 
            github.com/basic/business/data/store/user (master*) » go test -v

    7. upgrade dependency & software version 
    8. finish the project
    9. hey , update some trafic load, 模拟网络流量
    10. ultimate-go: https://www.bilibili.com/video/BV12341137Fo/?is_story_h5=false&p=1&share_from=ugc&share_medium=android&share_plat=android&share_session_id=72080e30-da82-45a3-926a-d251d9433f50&share_source=COPY&share_tag=s_i&timestamp=1664991916&unique_k=aBYklAq&vd_source=fd3a715c532ca990e4a8a6902ac1478c 
    11. 13 hours 的视频 & 基本✅
