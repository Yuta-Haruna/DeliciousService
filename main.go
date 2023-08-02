// main.go

package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "deliciousService"}

var accessToken string

func main() {
	rootCmd.AddCommand(getDataCmd)

	// アクセストークンをコマンドライン引数として追加
	rootCmd.PersistentFlags().StringVarP(&accessToken, "token", "t", "", "Contentful Access Token")
	// 必須入力化
	rootCmd.MarkPersistentFlagRequired("token")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

var spaceID = "2vskphwbz4oc"
var environmentIDs = []string{
	"6QRk7gQYmOyJ1eMG9H4jbB",
	"41RUO5w4oIpNuwaqHuSwEc",
	"4Li6w5uVbJNVXYVxWjWVoZ",
}

var getDataCmd = &cobra.Command{
	Use:   "getData",
	Short: "Contentful APIからデータを取得します",
	Run:   getData,
}

func getData(cmd *cobra.Command, args []string) {
	client := resty.New()

	for _, envID := range environmentIDs {
		url := fmt.Sprintf("https://cdn.contentful.com/spaces/%s/entries/%s", spaceID, envID)

		resp, err := client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
			SetHeader("Content-Type", "application/json").
			Get(url)

		if err != nil {
			fmt.Printf("Contentful APIからデータを取得中にエラーが発生しました（environmentID: %s）: %s\n", envID, err)
			continue
		}

		// レスポンスのボディをJSONデコード
		var data map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &data); err != nil {
			fmt.Printf("Contentful APIからのレスポンスの解析中にエラーが発生しました（environmentID: %s）: %s\n", envID, err)
			continue
		}

		// 必要な情報を取得して表示
		sys := data["sys"].(map[string]interface{})
		fields := data["fields"].(map[string]interface{})

		fmt.Printf("取得したデータ（environmentID: %s）:\n", envID)
		fmt.Println("ID:", sys["id"])
		fmt.Println("Name:", fields["name"])
		fmt.Println("CreatedAt:", sys["createdAt"])
	}
}
