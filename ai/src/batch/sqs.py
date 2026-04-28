"""
프로덕션 진입점 (SQS)

SQS 큐에서 job을 수신하여 처리합니다.
메시지 본문은 {"media_item_id": int, "family_id": int, "view_storage_key": string} 형식의 JSON이어야 합니다.
처리 성공/실패 여부와 무관하게 메시지를 삭제합니다(재시도는 resize-worker/DB가 관리).
"""

import json
import logging
import os
import signal

import boto3

from src.batch.face_recognition_batch import get_s3_client, logger, process_job

SQS_QUEUE_URL = os.environ.get("SQS_QUEUE_URL", "")
AWS_REGION = os.environ.get("AWS_REGION", "ap-northeast-1")

# SQS long-polling 대기 시간 (초, 최대 20)
WAIT_TIME_SECONDS = 20
# 한 번에 수신할 최대 메시지 수 (SQS 최대 10)
MAX_MESSAGES = 10

_running = True


def get_sqs_client():
    return boto3.client("sqs", region_name=AWS_REGION)


def run_sqs_batch_loop():
    logger.info("AI Batch Worker (production/SQS) 시작")
    sqs = get_sqs_client()
    s3_client = get_s3_client()

    while _running:
        try:
            response = sqs.receive_message(
                QueueUrl=SQS_QUEUE_URL,
                MaxNumberOfMessages=MAX_MESSAGES,
                WaitTimeSeconds=WAIT_TIME_SECONDS,
                AttributeNames=["All"],
            )
            messages = response.get("Messages", [])
            if not messages:
                continue

            logger.info(f"SQS 메시지 수신: {len(messages)}개")
            for msg in messages:
                if not _running:
                    break

                receipt_handle = msg["ReceiptHandle"]
                try:
                    job = json.loads(msg["Body"])
                except json.JSONDecodeError as e:
                    logger.error(f"메시지 파싱 실패: {e} body={msg['Body']}")
                    return False
                    continue

                success = process_job(s3_client, job)
                if not success:
                    return False

                # 성공 시 메시지 삭제
                _delete_message(sqs, receipt_handle)
                return True

        except Exception as e:
            logger.error(f"SQS 루프 에러: {e}")
            return False

    logger.info("AI Batch Worker (production/SQS) 종료")


def _delete_message(sqs, receipt_handle: str):
    try:
        sqs.delete_message(QueueUrl=SQS_QUEUE_URL, ReceiptHandle=receipt_handle)
    except Exception as e:
        logger.error(f"SQS 메시지 삭제 실패: {e}")


def _handle_signal(sig, frame):
    global _running
    logger.info("종료 신호 수신 - 현재 job 완료 후 종료")
    _running = False


if __name__ == "__main__":
    signal.signal(signal.SIGTERM, _handle_signal)
    signal.signal(signal.SIGINT, _handle_signal)
    run_sqs_batch_loop()
