"""
얼굴 인식 배치 워커

face_recognition_jobs 테이블을 폴링하여 detect_face()를 실행하고
결과를 resize-worker API로 전달합니다.
비즈니스 로직(identity 매칭, 크롭, DB 저장)은 resize-worker가 담당합니다.
"""

import logging
import os
import signal
import tempfile
import time

import boto3
import psycopg2
import psycopg2.extras
import requests
from botocore.config import Config

from src.service.face import detect_face

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger(__name__)

# 환경변수
DATABASE_URL = os.environ.get("DATABASE_URL", "")
STORAGE_TYPE = os.environ.get("STORAGE_TYPE", "minio")
S3_ENDPOINT = os.environ.get("S3_ENDPOINT", "http://minio:9000")
STORAGE_BUCKET = os.environ.get("STORAGE_BUCKET", "yuno-media")
MINIO_ROOT_USER = os.environ.get("MINIO_ROOT_USER", "root")
MINIO_ROOT_PASSWORD = os.environ.get("MINIO_ROOT_PASSWORD", "password")
AWS_REGION = os.environ.get("AWS_REGION", "us-east-1")
RESIZE_WORKER_URL = os.environ.get("RESIZE_WORKER_URL", "http://resize-worker:1325")

POLL_INTERVAL = 5
BATCH_SIZE = 10

_running = True


def get_db_conn():
    return psycopg2.connect(DATABASE_URL)


def get_s3_client():
    if STORAGE_TYPE == "minio":
        return boto3.client(
            "s3",
            endpoint_url=S3_ENDPOINT,
            aws_access_key_id=MINIO_ROOT_USER,
            aws_secret_access_key=MINIO_ROOT_PASSWORD,
            region_name=AWS_REGION,
            config=Config(signature_version="s3v4"),
        )
    return boto3.client("s3", region_name=AWS_REGION)


def download_image(s3_client, key: str) -> str:
    """S3에서 이미지를 임시 파일로 다운로드하고 경로 반환"""
    ext = key.rsplit(".", 1)[-1] if "." in key else "jpg"
    tmp = tempfile.NamedTemporaryFile(suffix=f".{ext}", delete=False)
    s3_client.download_fileobj(STORAGE_BUCKET, key, tmp)
    tmp.close()
    return tmp.name


def fetch_pending_jobs(conn, limit: int) -> list:
    """pending job을 원자적으로 fetch (FOR UPDATE SKIP LOCKED)"""
    with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
        cur.execute(
            """
            UPDATE face_recognition_jobs
            SET status = '02',
                attempt_count = attempt_count + 1,
                updated_at = NOW()
            WHERE id IN (
                SELECT id FROM face_recognition_jobs
                WHERE status = '01'
                ORDER BY created_at ASC
                LIMIT %s
                FOR UPDATE SKIP LOCKED
            )
            RETURNING *
            """,
            (limit,),
        )
        jobs = cur.fetchall()
    conn.commit()
    return [dict(j) for j in jobs]


def process_job(s3_client, job: dict):
    """단일 face_recognition_job 처리: ML 추론 후 resize-worker에 결과 전달"""
    job_id = job["id"]
    view_key = job["view_storage_key"]

    logger.info(f"job={job_id} 처리 시작 (view_key={view_key})")

    tmp_path = None
    try:
        # 1. view 이미지 다운로드
        tmp_path = download_image(s3_client, view_key)

        # 2. 얼굴 감지 + 임베딩 계산 (ML 추론만)
        result = detect_face(tmp_path)
        faces = result.get("faces", [])
        logger.info(f"job={job_id} 감지된 얼굴 수: {len(faces)}")

        # 3. resize-worker에 결과 전달 (비즈니스 로직은 Go가 처리)
        resp = requests.post(
            f"{RESIZE_WORKER_URL}/face-recognition/complete",
            json={"job_id": job_id, "faces": faces},
            timeout=30,
        )
        resp.raise_for_status()
        logger.info(f"job={job_id} 처리 완료")

    except Exception as e:
        logger.error(f"job={job_id} 처리 실패: {e}")
        try:
            requests.post(
                f"{RESIZE_WORKER_URL}/face-recognition/fail",
                json={"job_id": job_id, "error": str(e)},
                timeout=10,
            )
        except Exception as fail_e:
            logger.error(f"job={job_id} fail 상태 업데이트 실패: {fail_e}")
    finally:
        if tmp_path and os.path.exists(tmp_path):
            os.unlink(tmp_path)


def run_batch_loop():
    logger.info("AI Batch Worker 시작")
    conn = get_db_conn()
    s3_client = get_s3_client()

    while _running:
        try:
            jobs = fetch_pending_jobs(conn, BATCH_SIZE)
            if jobs:
                logger.info(f"처리할 job: {len(jobs)}개")
                for job in jobs:
                    if not _running:
                        break
                    process_job(s3_client, job)
            else:
                time.sleep(POLL_INTERVAL)
        except psycopg2.OperationalError:
            logger.warning("DB 연결 재시도...")
            try:
                conn.close()
            except Exception:
                pass
            time.sleep(5)
            conn = get_db_conn()
        except Exception as e:
            logger.error(f"배치 루프 에러: {e}")
            time.sleep(POLL_INTERVAL)

    logger.info("AI Batch Worker 종료")
    conn.close()


def _handle_sigterm(sig, frame):
    global _running
    logger.info("SIGTERM 수신 - 현재 job 완료 후 종료")
    _running = False


if __name__ == "__main__":
    signal.signal(signal.SIGTERM, _handle_sigterm)
    signal.signal(signal.SIGINT, _handle_sigterm)
    run_batch_loop()
