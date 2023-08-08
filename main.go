package main

import (
	"DeliciousService/cmd"
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "deliciousService"}

func init() {
	// データ取得関数にアクセストークンを設定するフラグを追加
	getDataCmd.Flags().StringP("token", "t", "", "アクセストークン (必須)")
	getDataCmd.MarkFlagRequired("token") // アクセストークンを必須にする

	// データ表示関数にIDを設定するフラグを追加
	// getBreadsInfoCmd.Flags().StringP("id", "i", "", "取得したいID (任意)　未入力の場合、すべての情報を表示")

	// データ取得関数とデータ表示関数をルートコマンドに追加
	rootCmd.AddCommand(getDataCmd)
	rootCmd.AddCommand(getBreadsInfoCmd)
}

// contentful APIを利用し、FireBaseへ格納する
var getDataCmd = &cobra.Command{
	Use:   "getData",
	Short: "Contentful APIからデータを取得し、FireStoreへ格納します",
	Run: func(inputCmd *cobra.Command, args []string) {
		// 実行
		cmd.GetData(inputCmd)
	},
}

// graphQLを利用し、FireStore内の情報へアクセスする
var getBreadsInfoCmd = &cobra.Command{
	Use:   "getBreadsInfo",
	Short: "graphQLを利用し、FireStore内の情報へアクセスします",
	Run: func(inputCmd *cobra.Command, args []string) {
		// 実行
		cmd.ExecuteBreadsInfoCmd()
	},
}

func main() {

	// 処理を実行
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
