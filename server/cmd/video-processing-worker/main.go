package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	workerVideo "github.com/ywl0806/yuno_kiroku/internal/worker/video"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

func main() {
	setting.SettingEnv()

	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	defer conn.Close()

	st := store.New(conn, db.New(conn))
	storageService := providers.NewStorageProvider().StorageService()

	svc := workerVideo.NewVideoProcessingService(st.MediaItem, storageService)

	// SQS 클라이언트 초기화
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("AWS 설정 로드 실패: %v", err)
	}
	if ep := viper.GetString("SQS_ENDPOINT_URL"); ep != "" {
		cfg.BaseEndpoint = aws.String(ep)
	}
	sqsClient := sqs.NewFromConfig(cfg)
	queueURL := viper.GetString("VIDEO_SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("VIDEO_SQS_QUEUE_URL 환경변수가 설정되지 않았습니다")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Video Processing Worker 시작")

	for {
		select {
		case <-ctx.Done():
			log.Println("Video Processing Worker 종료")
			return
		default:
		}

		out, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,  // Long Polling
			VisibilityTimeout:   900, // 15분
		})
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("SQS ReceiveMessage 실패: %v", err)
			continue
		}

		for _, msg := range out.Messages {
			var params services.VideoJobParams
			// S3-01: 파싱 실패 시 deleteMessage 호출 제거 → visibility timeout 후 재발행 → maxReceiveCount 소진 시 DLQ 이동
			if err := json.Unmarshal([]byte(aws.ToString(msg.Body)), &params); err != nil {
				log.Printf("[video-worker] [ERROR] 메시지 파싱 실패 (DLQ 경유 예정): %v", err)
				continue
			}

			// S3-05: 처리 시간이 15분을 초과할 수 있으므로 5분마다 visibility timeout을 10분으로 갱신
			extendCtx, extendCancel := context.WithCancel(ctx)
			go func(receiptHandle *string) {
				ticker := time.NewTicker(5 * time.Minute)
				defer ticker.Stop()
				for {
					select {
					case <-extendCtx.Done():
						return
					case <-ticker.C:
						if _, err := sqsClient.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
							QueueUrl:          aws.String(queueURL),
							ReceiptHandle:     receiptHandle,
							VisibilityTimeout: 600, // 10분으로 갱신
						}); err != nil {
							log.Printf("[video-worker] [WARN] visibility timeout 갱신 실패: %v", err)
						}
					}
				}
			}(msg.ReceiptHandle)

			log.Printf("비디오 처리 시작: media_item_id=%s", params.MediaItemID)
			if err := svc.ProcessVideo(ctx, params); err != nil {
				extendCancel()
				log.Printf("[video-worker] [WARN] 비디오 처리 실패 (재시도 예정): media_item_id=%s err=%v", params.MediaItemID, err)
				// DeleteMessage 하지 않으면 visibility timeout 후 재시도
				continue
			}

			extendCancel()
			deleteMessage(ctx, sqsClient, queueURL, msg.ReceiptHandle)
		}
	}
}

func deleteMessage(ctx context.Context, client *sqs.Client, queueURL string, receiptHandle *string) {
	if _, err := client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: receiptHandle,
	}); err != nil {
		log.Printf("SQS DeleteMessage 실패: %v", err)
	}
}
