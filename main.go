// main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "deliciousService"}

var accessToken string
var BREADS_DOCUMENT_ID = "GyNVqdXn86W20lPCEJ0Y"
var YOUR_PROJECT_ID = "samplefirebaseproject-c2ebb"

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
	Short: "Contentful APIからデータを取得し、FireStoreへ格納します",
	Run:   fetchDataAndStore,
}

func fetchDataAndStore(cmd *cobra.Command, args []string) {
	var allData []map[string]interface{}

	for _, envID := range environmentIDs {
		// Contentful APIからデータを取得
		data, err := getData(envID)
		if err != nil {
			log.Fatalf("Contentfulからデータを取得中にエラーが発生しました（environmentID: %s）: %v", envID, err)
		}

		// 取得データの格納　後に一括DB保存をする
		allData = append(allData, data...)

		for _, item := range data {
			fmt.Printf("ID: %v, Name: %v, CreatedAt: %v\n",
				item["id"], item["name"], item["createdAt"])
		}

	}
	// Firestoreにデータを保存
	if err := storeDataInFirestore(allData); err != nil {
		log.Fatalf("Firestoreへのデータ保存中にエラーが発生しました: %v", err)
	} else {
		fmt.Printf("データをFirestoreに保存しました\n")
	}

	fmt.Printf("処理が完了しました\n")

}

func getData(environmentID string) ([]map[string]interface{}, error) {
	client := resty.New()

	url := fmt.Sprintf("https://cdn.contentful.com/spaces/%s/entries/%s", spaceID, environmentID)

	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetHeader("Content-Type", "application/json").
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("Contentful APIからデータを取得中にエラーが発生しました（environmentID: %s）: %s\n", environmentID, err)
	}

	// レスポンスのボディをJSONデコード
	var data map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &data); err != nil {
		return nil, fmt.Errorf("Contentful APIからのレスポンスの解析中にエラーが発生しました（environmentID: %s）: %s\n", environmentID, err)
	}

	// 必要な情報を取得して表示
	sys := data["sys"].(map[string]interface{})
	fields := data["fields"].(map[string]interface{})

	fmt.Printf("取得したデータ（environmentID: %s）:\n", environmentID)
	fmt.Println("ID:", sys["id"])
	fmt.Println("Name:", fields["name"])
	fmt.Println("CreatedAt:", sys["createdAt"])

	// 必要な情報を取得してリストに格納
	var result []map[string]interface{}
	docData := map[string]interface{}{
		"id":        sys["id"],
		"name":      fields["name"],
		"createdAt": sys["createdAt"],
	}
	result = append(result, docData)

	return result, nil
}

// storeDataInFirestore関数（引数にallDataを受け取る）
func storeDataInFirestore(data []map[string]interface{}) error {

	// Firestoreクライアントを初期化
	ctx := context.Background()
	client, err := initFirestoreClient(ctx, "credentials.json", YOUR_PROJECT_ID)
	if err != nil {
		return fmt.Errorf("Firestoreクライアントの初期化に失敗しました: %v", err)
	}
	defer client.Close()

	// 一括書き込みのためにBulkWriterを初期化
	bw := client.BulkWriter(ctx)

	// 一括で書き込むコレクションの参照を取得
	collectionRef := client.Collection("Breads")

	// ドキュメントを一括で書き込み
	for _, item := range data {
		docRef := collectionRef.Doc(BREADS_DOCUMENT_ID)
		bw.Set(docRef, item)
	}

	// 書き込みを実行
	bw.Flush()

	return nil
}

func initFirestoreClient(ctx context.Context, serviceAccountKey, projectID string) (*firestore.Client, error) {
	// サービスアカウントキーから認証情報を取得
	opt := option.WithCredentialsFile(serviceAccountKey)

	// Firestoreクライアントを初期化
	client, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		return nil, fmt.Errorf("Firestoreクライアントの初期化に失敗しました: %v", err)
	}

	return client, nil
}
