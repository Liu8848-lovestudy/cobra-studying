/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// glsCmd represents the gls command
var glsCmd = &cobra.Command{
	Use:   "gls",
	Short: "gls命令",
	Long:  "gls命令用于展示该目录下的所有文件",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := cmd.Flags().GetBool("ls")      //获取参数l的值默认为false
		a, err1 := cmd.Flags().GetBool("all")    //获取参数a的值默认为false
		s, err2 := cmd.Flags().GetString("sort") //获取排序方式
		x, err3 := cmd.Flags().GetString("xl")   //获取文件大小的显示方式
		if err == nil && err1 == nil && err2 == nil && err3 == nil {
			dir, err2 := os.ReadDir(args[0])

			if err2 != nil {
				fmt.Println("输入文件目录有误！")
			}
			var infos = make([]os.FileInfo, 0)
			for _, file := range dir {
				info, _ := file.Info()
				infos = append(infos, info)
			}
			infos = a_filter(a, infos) //获得是否过滤隐藏文件后的文件切片
			infos = s_filter(s, infos) //获得排序后的文件切片
			um, m := x_filter(x)       //获得文件大小的单位以及处理方式
			l_filter(l, infos, um, m)  //打印文件
		}
	},
}

func a_filter(a bool, infos []os.FileInfo) (infos1 []os.FileInfo) {
	if a {
		return infos
	} else {
		//var infos1 = make([]os.FileInfo, 5, 10)
		for _, info := range infos {
			if !CheckIsHidden(info) {
				infos1 = append(infos1, info)
			}
		}
		return infos1
	}
}

func s_filter(s string, infos []os.FileInfo) []os.FileInfo {
	switch s {
	case "":
		return infos
	case "name", "n", "Name", "N", "NAME":
		sort.SliceStable(infos, func(i, j int) bool {
			return infos[j].Name() < infos[i].Name() //按名称逆序排序
		})
		return infos
	case "size", "s", "Size", "SIZE":
		sort.SliceStable(infos, func(i, j int) bool {
			return infos[j].Size() > infos[i].Size() //按文件大小从小到大排序
		})
		return infos
	case "update", "Update", "u", "U":
		sort.SliceStable(infos, func(i, j int) bool {
			return infos[i].ModTime().After(infos[j].ModTime()) //按文件更新时间排序
		})
		return infos
	default:
		panic("输入的排序参数有误！")
	}
}

func x_filter(x string) (string, float64) {
	switch x {
	case "", "byte", "Byte", "BYTE", "by":
		return "byte", 1
	case "bit", "Bit", "BIT", "bi":
		return "bit", 1024
	case "kb", "KB", "Kb", "kB", "k", "K":
		return "kb", float64(1) / 1024
	case "MB", "m", "M", "Mb", "mb":
		return "m", float64(1) / (1024 * 1024)
	case "GB", "Gb", "gb", "g", "G":
		return "g", float64(1) / (1024 * 1024 * 1024)
	default:
		panic("输入的单位参数有误！")
	}
}

func l_filter(l bool, infos []os.FileInfo, um string, m float64) {
	var f_type string
	for _, info := range infos {
		if info.IsDir() {
			f_type = "文件夹"
		} else {
			arr := strings.Split(info.Name(), ".")
			f_type = arr[len(arr)-1]
		}
		if l {
			switch f_type {
			case "文件夹":
				fmt.Printf("%c[1;40;32m%v\t%v\t%.2f%v\t%v\t%c[0m\n", 0x1B, info.Name(), f_type, m*float64(info.Size()), um, info.ModTime().Format("2006-01-02 15:04:05"), 0x1B)

			case "txt", "Txt", "TXT":
				fmt.Printf("%c[1;40;34m%v\t%v\t%.2f%v\t%v\t%c[0m\n", 0x1B, info.Name(), f_type, m*float64(info.Size()), um, info.ModTime().Format("2006-01-02 15:04:05"), 0x1B)

			case "exe", "EXE":
				fmt.Printf("%c[1;40;35m%v\t%v\t%.2f%v\t%v\t%c[0m\n", 0x1B, info.Name(), f_type, m*float64(info.Size()), um, info.ModTime().Format("2006-01-02 15:04:05"), 0x1B)

			default:
				fmt.Printf("%c[1;40;33m%v\t%v\t%.2f%v\t%v\t%c[0m\n", 0x1B, info.Name(), f_type, m*float64(info.Size()), um, info.ModTime().Format("2006-01-02 15:04:05"), 0x1B)
			}
		} else {
			switch f_type {
			case "文件夹":
				fmt.Printf("%c[1;40;32m%v%c[0m\n", 0x1B, info.Name(), 0x1B)
			case "txt", "Txt", "TXT":
				fmt.Printf("%c[1;40;34m%v%c[0m\n", 0x1B, info.Name(), 0x1B)
			case "exe", "EXE":
				fmt.Printf("%c[1;40;35m%v%c[0m\n", 0x1B, info.Name(), 0x1B)
			default:
				fmt.Printf("%c[1;40;33m%v%c[0m\n", 0x1B, info.Name(), 0x1B)
			}
		}
	}
}

func CheckIsHidden(file os.FileInfo) bool {
	//"通过反射来获取Win32FileAttributeData的FileAttributes
	fa := reflect.ValueOf(file.Sys()).Elem().FieldByName("FileAttributes").Uint()
	bytefa := []byte(strconv.FormatUint(fa, 2))
	if bytefa[len(bytefa)-2] == '1' {
		//fmt.Println("隐藏目录:", file.Name())
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(glsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// glsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// glsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	glsCmd.Flags().BoolP("ls", "l", false, "显示文件详细信息")
	glsCmd.Flags().BoolP("all", "a", false, "显示全部文件")
	glsCmd.Flags().StringP("sort", "s", "", "文件夹排序")
	glsCmd.Flags().StringP("xl", "x", "", "文件大小单位")
}
