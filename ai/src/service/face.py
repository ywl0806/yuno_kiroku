import logging

import cv2
from insightface.app import FaceAnalysis

logger = logging.getLogger(__name__)
_MAX_DIMENSION = 2048  # OOM 방지 다운샘플링 임계값 (px)

_face_analyzer = None


def get_face_analyzer():
    """ArcFace 모델 싱글톤 인스턴스 반환"""
    global _face_analyzer
    if _face_analyzer is None:
        _face_analyzer = FaceAnalysis(
            providers=["CPUExecutionProvider"],  # GPU 사용시: ['CUDAExecutionProvider']
            name="buffalo_l",
        )
        _face_analyzer.prepare(ctx_id=0, det_size=(640, 640))
    return _face_analyzer


def detect_face(image_path: str) -> dict:
    """
    ArcFace를 사용한 얼굴 인식 함수

    Args:
        image_path: 로컬 이미지 경로
    Returns:
        dict: 얼굴 위치와 임베딩 (512차원)
        example: [{"face_location": { "top": 100, "right": 200, "bottom": 300, "left": 400 }, "embedding": [0.1, 0.2, ...]}]
    """
    # 이미지 로드
    image = cv2.imread(image_path)
    if image is None:
        return {"faces": []}

    # RGB로 변환 (OpenCV는 BGR 사용)
    image_rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)

    # OOM 방지: 장축이 _MAX_DIMENSION 초과 시 다운샘플링
    h, w = image_rgb.shape[:2]
    if max(h, w) > _MAX_DIMENSION:
        scale = _MAX_DIMENSION / max(h, w)
        image_rgb = cv2.resize(image_rgb, (int(w * scale), int(h * scale)), interpolation=cv2.INTER_AREA)
        logger.warning(f"이미지 다운샘플링: {w}x{h} → {int(w*scale)}x{int(h*scale)} (path={image_path})")

    # ArcFace 모델로 얼굴 감지 및 임베딩 추출
    face_analyzer = get_face_analyzer()
    try:
        faces = face_analyzer.get(image_rgb)
    except Exception as e:
        logger.error(f"InsightFace 추론 실패 (path={image_path}): {e}", exc_info=True)
        return {"faces": []}

    result = []
    for face in faces:
        # 얼굴 위치 (bbox: [x1, y1, x2, y2])
        bbox = face.bbox.astype(int)

        # ArcFace 임베딩 (512차원, 정규화됨)
        embedding = face.normed_embedding.tolist()

        result.append(
            {
                "face_location": {
                    "top": int(bbox[1]),  # y1
                    "right": int(bbox[2]),  # x2
                    "bottom": int(bbox[3]),  # y2
                    "left": int(bbox[0]),  # x1
                },
                "embedding": embedding,
            }
        )

    return {"faces": result}
