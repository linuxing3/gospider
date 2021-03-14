echo "======================================================================"
echo "以下命令仅需要运行一次"
echo "======================================================================"
echo "数据库导入(optional)"
echo "go run github.com/prisma/prisma-client-go introspect"
echo "数据库迁移(optional)"
echo "go run github.com/prisma/prisma-client-go db push --preview-feature"

echo "2. 生成客户端"
echo "go run github.com/prisma/prisma-client-go generate"

echo "3. 运行"
echo "go run main.go"

echo "4. 安装gospider到$GOPATH/bin/"
echo "go install ."

echo "5. 初始化数据库, 链接如下："
echo "postgresql://spider:20090909@db:5432/spider?schema=public"
echo "deno run -A --unstable https://raw.githubusercontent.com/linuxing3/gospider/main/create_table.ts"

echo "======================================================================"
echo "6. 启动docker-compose"
echo "======================================================================"

docker-compose up -d --remove-orphans

echo "开始抓取"
echo "docker exec -it spider /root/go/bin/gospider spider"
