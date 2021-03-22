package cmd

import (
	"fmt"
	"log"

	"github.com/linuxing3/gospider/prisma/db"
	"github.com/linuxing3/vpsman/util"
	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "管理数据库，调用prisma,...",
	Long: ` About usage of using database. For example: 
Prisma is aweseme next generation ORM tools
Denodb is tools to generate models.`,
	Run: func(cmd *cobra.Command, args []string) {
		databaseMenu()
	},
}

// TODO: InitPrismaClient 初始化PrismaClient, 不能成功调用,因为context不同
func InitPrismaClient() *db.PrismaClient{
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	return client
}

// TrojanMenu 控制TrojanMenu
func databaseMenu() {
exit:
	for {
		fmt.Println()
		fmt.Print(util.Cyan("请选择"))
		fmt.Println()
		loopMenu := []string{"创建表格", "Deno创建表格", "创建Prisma客户端", "自动生成Schema", "迁移数据库"}
		choice := util.LoopInput("回车退出:   ", loopMenu, false)
		switch choice {
		case 1:
			fmt.Println("create_table")
		case 2:
			fmt.Println("deno run -A --unstable ./create_table.ts")
		case 3:
			fmt.Println("go run github.com/prisma/prisma-client-go generate")
		case 4:
			fmt.Println("go run github.com/prisma/prisma-client-go introspect")
		case 5:
			fmt.Println("go run github.com/prisma/prisma-client-go db push --preview-feature")
		default:
			break exit
		}
	}
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
