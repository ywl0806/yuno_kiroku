from http import HTTPStatus
import os
import uuid
import requests
import logging
from fastapi import FastAPI, File, Form, Response, UploadFile
from fastapi.exceptions import HTTPException
from fastapi.requests import Request
from src.service.face import detect_face

app = FastAPI()

TMP_DIR = "tmp/files"
logger = logging.getLogger(__name__)


@app.get("/")
def read_root():
    return {"message": "Hello, World!"}


@app.post("/face-detection/file")
def detect_face_endpoint(request: Request, file: UploadFile = File(...)):
    """
    파일로 부터 얼굴 인식

    Args:
        file: 이미지 파일
    Returns:
        dict: 얼굴 정보
    """
    if not file.filename:
        raise HTTPException(status_code=400, detail="File name is required")
    if not file.filename.split(".")[-1]:
        raise HTTPException(status_code=400, detail="File extension is required")
    logger.info(f"request: {request.headers}")
    file_path = f"tmp/files/{uuid.uuid4()}.{file.filename.split('.')[-1]}"
    _check_tmp_dir()
    with open(file_path, "wb") as f:
        f.write(file.file.read())
    face_detection = detect_face(file_path)
    _clean_tmp_files()
    if len(face_detection["faces"]) > 0:
        return face_detection
    else:
        return Response(status_code=HTTPStatus.NO_CONTENT)


@app.post("/face-detection/url")
def detect_face_url_endpoint(url: str = Form(...)):
    """
    URL로 부터 얼굴 인식

    Args:
        url: 이미지 URL
    Returns:
        dict: 얼굴 정보
    """
    if not url:
        raise HTTPException(status_code=400, detail="URL is required")
    if not url.startswith("http"):
        raise HTTPException(status_code=400, detail="URL must start with http")

    file_path = f"tmp/files/{uuid.uuid4()}.{url.split('.')[-1]}"
    _check_tmp_dir()
    with open(file_path, "wb") as f:
        f.write(requests.get(url).content)

    face_detection = detect_face(file_path)
    _clean_tmp_files()
    if len(face_detection["faces"]) > 0:
        return face_detection
    else:
        return Response(status_code=HTTPStatus.NO_CONTENT)


def _check_tmp_dir():
    if not os.path.exists(TMP_DIR):
        os.makedirs(TMP_DIR, exist_ok=True)


def _clean_tmp_files():
    for file in os.listdir(TMP_DIR):
        os.remove(os.path.join(TMP_DIR, file))
