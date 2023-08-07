package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/graphql-go/graphql"
	"google.golang.org/api/option"
)

type BreadInfo struct {
	ID        string `firestore:"id"`
	Name      string `firestore:"name"`
	CreatedAt string `firestore:"createdAt"`
}

// graphQLを利用し、FireStore内の情報へアクセスする
func ExecuteBreadsInfoCmd() {
	// Firebaseの初期化
	opt := option.WithCredentialsFile("credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("Error initializing Firebase:", err)
		return
	}

	fmt.Println("FiraBaseの初期化成功")

	// Firestoreのクライアントを取得
	fireStoreClient, err := app.Firestore(context.Background())
	if err != nil {
		fmt.Println("Error initializing Firestore client:", err)
		return
	}
	defer fireStoreClient.Close()
	fmt.Println("GraphQL　★★")

	// スキーマの定義
	fields := graphql.Fields{
		"item": &graphql.Field{
			Type:        itemType,
			Description: "Get item by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: getItemByIDResolver(fireStoreClient), // Firestoreクライアントをリゾルバ関数に渡す
		},
	}
	fmt.Println("GraphQL　★★★★")

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Println("Error creating GraphQL schema:", err)
		return
	}

	fmt.Println("GraphQL　★★★★")

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		queryID := r.URL.Query().Get("id")
		if queryID != "" {
			fmt.Println("GraphQL　ID　Target")
			// IDが与えられた場合は該当する情報のみ取得する
			result := graphql.Do(graphql.Params{
				Context: r.Context(), // リクエストのコンテキストを渡す
				Schema:  schema,
				RequestString: fmt.Sprintf(`
					{
						item(id: "%s") {
							id,
							name,
							createdAt
						}
					}`, queryID),
			})
			json.NewEncoder(w).Encode(result)
		} else {
			fmt.Println("GraphQL　ALL")

			// IDが与えられない場合はすべての情報を取得する
			result := graphql.Do(graphql.Params{
				Context: r.Context(), // リクエストのコンテキストを渡す
				Schema:  schema,
				RequestString: `
					{
						item {
							id,
							name,
							createdAt
						}
					}`,
			})
			json.NewEncoder(w).Encode(result)
		}
	})

	fmt.Println("GraphQL server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

// GraphQLスキーマの定義
var itemType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Item",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// Firestoreからデータを取得するリゾルバ関数を返す関数
func getItemByIDResolver(firestoreClient *firestore.Client) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		// クエリからIDを取得
		// queryID, ok := p.Args["id"].(string)
		// if !ok {
		// 	// return nil, fmt.Errorf("★★★★★ID not provided")
		// }
		// Firestoreからデータを取得する処理を実装
		documentID := "GyNVqdXn86W20lPCEJ0Y" // スラッシュ（"/"）を含まない形式に修正
		doc, err := firestoreClient.Collection("Breads").Doc(documentID).Get(context.Background())
		// doc, err := firestoreClient.Collection("GyNVqdXn86W20lPCEJ0Y").Doc(queryID).Get(p.Context)
		if err != nil {
			// エラーが発生した場合はログに出力
			fmt.Println("Error fetching data from Firestore:", err)
			return nil, err
		}

		fmt.Println("aaaaaaaa")

		// データが見つからなかった場合はエラーを返す
		if !doc.Exists() {
			return nil, fmt.Errorf("item not found")
		}
		fmt.Println("bbbbbbbbbbbb")

		// 取得したデータを構造体にマッピングして返す
		var bread BreadInfo
		if err := doc.DataTo(&bread); err != nil {
			fmt.Println("★★Error converting Firestore data:", err)
			return nil, err
		}
		fmt.Println("cccccccccccc")

		return bread, nil
	}
}
