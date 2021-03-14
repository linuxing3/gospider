rem 1.1. 数据库导入(optional)
rem go run github.com/prisma/prisma-client-go introspect
rem 1.2. 数据库迁移(optional)
rem go run github.com/prisma/prisma-client-go db push --preview-feature
rem 2. 生成客户端
go run github.com/prisma/prisma-client-go generate
rem 3. 运行
rem go run main.go
rem 4. 安装
go build .
go install .

rem 5. 设置数据库
deno run -A --unstable https://raw.githubusercontent.com/linuxing3/gospider/main/create_table.ts 

rem headless-shell onwardlinux
docker run -d -p 9222:9222 --rm --name headless-shell chromedp/headless-shell