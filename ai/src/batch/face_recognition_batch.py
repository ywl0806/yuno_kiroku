"""
얼굴 인식 배치 공통 코어

S3 클라이언트, 이미지 다운로드, ML 추론 + resize-worker 결과 전달을 담당합니다.
job을 가져오는 방식(DB / SQS)은 각 진입점(local.py / sqs.py)에서 처리합니다.
"""

import logging
import os
import tempfile

import boto3
import requests
from botocore.config import Config

from src.service.face import detect_face

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger(__name__)

STORAGE_TYPE = os.environ.get("STORAGE_TYPE", "minio")
S3_ENDPOINT = os.environ.get("S3_ENDPOINT", "http://minio:9000")
MEDIA_BUCKET_NAME = os.environ.get("MEDIA_BUCKET_NAME", "my-bucket")
MINIO_ROOT_USER = os.environ.get("MINIO_ROOT_USER", "root")
MINIO_ROOT_PASSWORD = os.environ.get("MINIO_ROOT_PASSWORD", "password")
AWS_REGION = os.environ.get("AWS_REGION", "ap-northeast-1")
RESIZE_WORKER_URL = os.environ.get("RESIZE_WORKER_URL", "http://resize-worker:1325")


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
    ext = key.rsplit(".", 1)[-1] if "." in key else "jpg"
    tmp = tempfile.NamedTemporaryFile(suffix=f".{ext}", delete=False)
    s3_client.download_fileobj(MEDIA_BUCKET_NAME, key, tmp)
    tmp.close()
    return tmp.name


def process_job(s3_client, job: dict) -> bool:
    """ML 추론 후 resize-worker에 결과 전달. 성공 시 True, 실패 시 False 반환."""
    job_id = job["id"]
    view_key = job["view_storage_key"]

    logger.info(f"job={job_id} 처리 시작 (view_key={view_key})")

    tmp_path = None
    try:
        tmp_path = download_image(s3_client, view_key)

        result = detect_face(tmp_path)
        faces = result.get("faces", [])
        logger.info(f"job={job_id} 감지된 얼굴 수: {len(faces)}")

        resp = requests.post(
            f"{RESIZE_WORKER_URL}/face-recognition/complete",
            json={"job_id": job_id, "faces": faces},
            timeout=30,
        )
        resp.raise_for_status()
        logger.info(f"job={job_id} 처리 완료")
        return True

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
        return False
    finally:
        if tmp_path and os.path.exists(tmp_path):
            os.unlink(tmp_path)
