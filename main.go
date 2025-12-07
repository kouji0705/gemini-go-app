package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("🚀 Dockerコンテナ内で処理を開始します...")
	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ エラー: GEMINI_API_KEY 環境変数が設定されていません")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// PDFファイル名
	fileName := "sample.pdf"

	// ファイルの存在確認
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("❌ エラー: %s が見つかりません。同じフォルダに置いてください。", fileName)
	}
	defer f.Close()

	fmt.Printf("1. ファイルをアップロード中: %s\n", fileName)
	uploadFile, err := client.UploadFile(ctx, "", f, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   完了 -> URI: %s\n", uploadFile.URI)

	fmt.Println("2. Google側での処理待ち...")
	for {
		fileInfo, err := client.GetFile(ctx, uploadFile.Name)
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.State == genai.FileStateActive {
			fmt.Println("   準備OK！")
			break
		} else if fileInfo.State == genai.FileStateFailed {
			log.Fatal("❌ 処理に失敗しました")
		}
		fmt.Print(".")
		time.Sleep(2 * time.Second)
	}

	// モデル設定
	model := client.GenerativeModel("gemini-1.5-flash")

	fmt.Println("\n3. AIに質問中...")
	resp, err := model.GenerateContent(ctx,
		genai.FileData{URI: uploadFile.URI},
		genai.Text("この資料の要点を、技術者向けに箇条書きで3点にまとめてください。"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				fmt.Println("\n--- 回答 ---")
				fmt.Println(txt)
				fmt.Println("------------")
			}
		}
	} else {
		fmt.Println("回答が得られませんでした。")
	}
}
