"""
얼굴 인식 배치 공통 코어

S3 클라이언트, 이미지 다운로드, ML 추론 후 Go 바이너리 실행을 담당합니다.
job을 가져오는 방식(SQS)은 각 진입점(sqs.py)에서 처리합니다.
"""

import json
import logging
import os
import subprocess
import tempfile

import boto3
from botocore.config import Config
from botocore.exceptions import ClientError

from src.service.face import detect_face

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger(__name__)

STORAGE_TYPE = os.environ.get("STORAGE_TYPE", "s3")
S3_ENDPOINT = os.environ.get("S3_ENDPOINT", "http://minio:9000")
MEDIA_BUCKET_NAME = os.environ.get("MEDIA_BUCKET_NAME", "my-bucket")
MINIO_ROOT_USER = os.environ.get("MINIO_ROOT_USER", "root")
MINIO_ROOT_PASSWORD = os.environ.get("MINIO_ROOT_PASSWORD", "password")
AWS_REGION = os.environ.get("AWS_REGION", "ap-northeast-1")
FACE_RECOGNITION_WORKER = os.environ.get("FACE_RECOGNITION_WORKER", "/app/face-recognition-worker")


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
    """ML 추론 후 Go 바이너리로 identity 매칭·크롭·저장 처리. 성공 시 True 반환."""
    view_key = job["view_storage_key"]
    media_item_id = job["media_item_id"]
    family_id = job["family_id"]

    logger.info(f"media_item={media_item_id} 처리 시작 (view_key={view_key})")

    tmp_path = None
    try:
        try:
            tmp_path = download_image(s3_client, view_key)
        except ClientError as e:
            if e.response["Error"]["Code"] in ("404", "NoSuchKey"):
                logger.error(
                    f"media_item={media_item_id} S3 파일 없음 (key={view_key}), 재시도 불필요 - 메시지 삭제"
                )
                return True  # 재시도 불필요 → 즉시 메시지 삭제
            raise

        result = detect_face(tmp_path)
        faces = result.get("faces", [])
        logger.info(f"media_item={media_item_id} 감지된 얼굴 수: {len(faces)}")

        payload = json.dumps({
            "media_item_id": media_item_id,
            "family_id": family_id,
            "view_image_path": tmp_path,
            "faces": faces,
        }).encode()

        proc = subprocess.run(
            [FACE_RECOGNITION_WORKER],
            input=payload,
            capture_output=True,
        )
        if proc.returncode != 0:
            logger.error(f"media_item={media_item_id} face-recognition-worker 실패: {proc.stderr.decode()}")
            return False

        logger.info(f"media_item={media_item_id} 처리 완료")
        return True

    except Exception as e:
        logger.error(f"media_item={media_item_id} 처리 실패: {e}")
        return False
    finally:
        if tmp_path and os.path.exists(tmp_path):
            try:
                os.unlink(tmp_path)
            except OSError as e:
                logger.warning(f"임시 파일 삭제 실패 (path={tmp_path}): {e}")
