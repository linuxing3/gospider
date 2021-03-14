package cmd

import (
	"fmt"
	"os/exec"

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

// TrojanMenu 控制TrojanMenu
func databaseMenu() {
exit:
	for {
		fmt.Println()
		fmt.Print(util.Cyan("请选择"))
		fmt.Println()
		loopMenu := []string{"创建表格", "创建Prisma客户端", "自动生成Schema", "迁移数据库"}
		choice := util.LoopInput("回车退出:   ", loopMenu, false)
		switch choice {
		case 1:
			remoteScript := "https://raw.githubusercontent.com/linuxing3/gospider/main/create_table.ts"
			fmt.Println("deno run -A --unstable https://raw.githubusercontent.com/linuxing3/gospider/main/create_table.ts")
			if c, err := exec.Command("cmd", "/C", "deno", "run", "-A", "--unstable", remoteScript).CombinedOutput(); err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Printf("%s ", c)
			}
		case 2:
			fmt.Println("go run github.com/prisma/prisma-client-go generate")
			if c, err := exec.Command("cmd", "/C", "go", "run", "github.com/prisma/prisma-client-go", "generate").CombinedOutput(); err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Printf("%s\n", c)
			}
		case 3:
			// util.ExecCommand("go run github.com/prisma/prisma-client-go introspect")
			if c, err := exec.Command("cmd", "/C", "go", "run", "github.com/prisma/prisma-client-go", "introspect").CombinedOutput(); err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Printf("%s ", c)
			}
		case 4:
			// util.ExecCommand("go run github.com/prisma/prisma-client-go db push --preview-feature")
			if c, err := exec.Command("cmd", "/C", "go", "run", "github.com/prisma/prisma-client-go", "db", "push", "--preview-feature").CombinedOutput(); err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Printf("%s ", c)

			}
		default:
			break exit
		}
	}
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
