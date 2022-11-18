/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "展示目录结构",
	Long:  "以树形展示目录结构",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("请输入目录路径")
			return
		}
		l, err := cmd.Flags().GetInt("length")
		if err == nil {
			ShowPathTree(args[0], l)
		}

	},
}

var (
	levelFlag []bool // 路径级别标志
	fileCount,
	dirCount int
	level int
	len1  int
)

const (
	//space  = ""
	line   = "│  "
	last   = "└─ "
	middle = "├─ "
)

func ShowPathTree(path string, length int) {
	levelFlag = make([]bool, 0)
	len1 = length
	walk(path)
	fmt.Println(fmt.Sprintf("\n指定路径下有 %d 个目录，%d 个文件。", dirCount, fileCount))
}

// walk 递归遍历指定路径
func walk(path string) {
	levelFlag = append(levelFlag, true)
	if dir, err := os.ReadDir(path); err == nil {
		for k, file := range dir {
			absFile := filepath.Join(path, file.Name())
			// 判断是否当前级别下的最后一个节点
			isLast := false
			if k == len(dir)-1 {
				isLast = true
			}
			if level >= len1 {
				isLast = true
				levelFlag = append(levelFlag, !isLast)
				showLine(level, isLast, file)
				continue
			}
			// 不是当前级别的最后节点，则设置为上级节点未结束
			levelFlag = append(levelFlag, !isLast)
			showLine(level, isLast, file)
			if file.IsDir() {
				level++
				walk(absFile)
				level--
			}
		}
	} else {
		panic(err)
	}
}

// showLine 显示当前节点的输出行
func showLine(level int, isLast bool, file os.DirEntry) {
	preFix := buildPrefix(level)
	var out string
	fName := file.Name()
	//如果判断为文件夹则格式化名字
	if file.IsDir() {
		fName = fmt.Sprintf("<%s>", fName)
		dirCount++
	} else {
		fileCount++
	}
	if isLast {
		out = fmt.Sprintf("%s%s%s", preFix, last, fName)
	} else {
		out = fmt.Sprintf("%s%s%s", preFix, middle, fName)
	}
	fmt.Println(out)
}

// buildPrefix 根据levelFlag的标志，构造上级节点的关系线
func buildPrefix(level int) string {
	var result string
	for i := 0; i < level; i++ {
		if levelFlag[i] {
			result += line
		} else {
			//result += space
		}
	}
	return result
}

func init() {
	glsCmd.AddCommand(treeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// treeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// treeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	treeCmd.Flags().IntP("length", "L", 100, "目录树的深度")

}
