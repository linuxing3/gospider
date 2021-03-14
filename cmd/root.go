package cmd

import (
	"fmt"
	"os"

	"github.com/linuxing3/vpsman/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var dbPath string

// Any 别名
type Any map[string]interface{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gospider",
	Short: "Simple vps manager",
	Long: `A simple vps manager that can control
system services backend by sqlite.`,
	Run: func(cmd *cobra.Command, args []string) {
		mainMenu()
	},
}

func mainMenu() {
exit:
	for {
		fmt.Println()
		fmt.Println(util.Cyan("欢迎使用管理程序"))
		fmt.Println()
		menuList := []string{"爬虫", "数据"}
		switch util.LoopInput("请选择: ", menuList, false) {
		case 1:
			spiderMenu()
		case 2:
			spiderMenu()
		default:
			break exit
		}
	}
}

// Execute 执行rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// 加载配置文件
	cobra.OnInitialize(initConfig)
	// 通用选项
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gospider.yaml)")
	// 本地选项
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.SetConfigName(".gospider") 
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// CreateDefaultInitConfig 创建默认的配置文件
func CreateDefaultInitConfig() {
	path, _ := homedir.Expand(".gospider.yaml")
	util.EnsureFileExists(path)

	initConfig()

	defaultConf := Any{
		"db": Any {
			"sqlite": Any {
				"path": "./gospider.sqlite",
			},
			"leveldb": Any {
				"path": "./gospider.leveldb",
			},
			"jsondb": Any {
				"path": "./gospider.json",
			},
		},
	}
	viper.SetDefault("main", defaultConf)
}
