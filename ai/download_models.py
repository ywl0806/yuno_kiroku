#!/usr/bin/env python3
"""
insightface 모델을 로컬에 다운로드하는 스크립트
"""

import os
from insightface.app import FaceAnalysis

# 모델 저장 디렉토리 설정
MODELS_DIR = os.path.join(os.path.dirname(__file__), "models")
os.makedirs(MODELS_DIR, exist_ok=True)

# 모델 다운로드
print("Downloading buffalo_s model...")
app = FaceAnalysis(
    providers=["CPUExecutionProvider"],
    # name="buffalo_s",
    name="buffalo_l",
    root=MODELS_DIR,  # 모델 저장 경로 지정
)
app.prepare(ctx_id=0, det_size=(640, 640))
print(f"Model downloaded to: {MODELS_DIR}")
