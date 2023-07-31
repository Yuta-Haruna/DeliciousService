package cmd

import (
	"context"
	"fmt"

	"github.com/contentful-labs/contentful-go"
	"github.com/spf13/cobra"
	// "cloud.google.com/go/firestore"
	// "google.golang.org/api/option"
)

// ContentfulData 構造体はContentful APIから取得したパンの情報を格納するためのものです
type ContentfulData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A simple CLI application using Cobra",
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when the application is called without any subcommands or flags
		fmt.Println("Welcome to myapp!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	// 処理を記載する
	rootCmd.AddCommand(getCmd)
	// rootCmd.AddCommand(saveCmd)
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get data from Contentful API",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement the logic to get data from Contentful API and display it
		// For simplicity, we'll just print a message here
		fmt.Println("Getting data from Contentful API...")
		data, err := getDataFromContentfulAPI()
		if err != nil {
			fmt.Println("Error: Failed to get data from Contentful API:", err)
			return
		}
		fmt.Println("Data:", data)
	},
}

// var saveCmd = &cobra.Command{
// 	Use:   "save",
// 	Short: "Save data to the database",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		// Implement the logic to save data to the database
// 		// For simplicity, we'll just print a message here
// 		fmt.Println("Saving data to the database...")
// 		err := saveDataToDatabase()
// 		if err != nil {
// 			fmt.Println("Error: Failed to save data to the database:", err)
// 			return
// 		}
// 		fmt.Println("Data saved successfully!")
// 	},
// }

// これらの値は、実際のContentfulのスペースIDとアクセストークンに置き換えてください
const (
	spaceID     = "あなたのContentfulスペースID"
	accessToken = "あなたのContentfulアクセストークン"
)

func getDataFromContentfulAPI() ([]ContentfulData, error) {
	// Contentfulのクライアントを初期化します
	client := contentful.NewClient(
		contentful.SpaceID(spaceID),
		contentful.AccessToken(accessToken),
	)

	// Contentful APIのクエリパラメータを定義します
	query := &contentful.Query{
		ContentType: "コンテンツタイプのID",
	}

	// Contentful APIを呼び出してデータを取得します
	entries, err := client.Entries(context.Background(), query)
	if err != nil {
		return nil, err
	}

	// APIのレスポンスを処理し、必要なデータを抽出します
	var result []ContentfulData
	for _, entry := range entries.Items {
		// エントリーからフィールドを抽出します
		id := entry.Sys.ID
		name := entry.Fields["name"].(string)
		createdAt := entry.Sys.CreatedAt.Format("2006-01-02 15:04:05")

		// データを格納するための構造体を作成します
		data := ContentfulData{
			ID:        id,
			Name:      name,
			CreatedAt: createdAt,
		}

		result = append(result, data)
	}

	return result, nil
}

// // データベースのクライアント
// var dbClient *firestore.Client

// // FirebaseプロジェクトIDと秘密鍵ファイルのパス
// const (
// 	projectID     = "YOUR_FIREBASE_PROJECT_ID"
// 	credentialsPath = "path/to/your/credentials.json"
// )

// // 初期化関数でデータベースのクライアントを作成します
// func init() {
// 	ctx := context.Background()
// 	opt := option.WithCredentialsFile(credentialsPath)
// 	client, err := firestore.NewClient(ctx, projectID, opt)
// 	if err != nil {
// 		log.Fatalf("Failed to create Firestore client: %v", err)
// 	}
// 	dbClient = client
// }

// func saveDataToFirestore(data []ContentfulData) error {
// 	// ドキュメントのコレクションを指定します
// 	collection := dbClient.Collection("breads")

// 	// データをFireStoreに保存します
// 	for _, d := range data {
// 		docData := map[string]interface{}{
// 			"id":         d.ID,
// 			"name":       d.Name,
// 			"created_at": d.CreatedAt,
// 		}
// 		docRef, _, err := collection.Add(context.Background(), docData)
// 		if err != nil {
// 			log.Println("Failed to add document to Firestore:", err)
// 			continue
// 		}
// 		fmt.Printf("Document added with ID: %s\n", docRef.ID)
// 	}

// 	return nil
// }
