import face_recognition


def detect_face(image_path: str) -> dict:
    """
    얼굴 인식 함수

    Args:
        image_path: 로컬 이미지 경로
    Returns:
        dict: 얼굴 위치와 임베딩  example: [{"face_location": { "top": 100, "right": 200, "bottom": 300, "left": 400 }, "embedding": [0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0 ...]}]
    """

    image = face_recognition.load_image_file(image_path)
    # 얼굴 위치 tuple(top, right, bottom, left)
    face_locations = face_recognition.face_locations(image)
    # 얼굴 임베딩(128d)
    face_encodings = face_recognition.face_encodings(image, face_locations)

    result = []
    for face_location, face_encoding in zip(face_locations, face_encodings):
        result.append(
            {
                "face_location": {
                    "top": face_location[0],
                    "right": face_location[1],
                    "bottom": face_location[2],
                    "left": face_location[3],
                },
                "embedding": face_encoding.tolist(),
            }
        )
    return {"faces": result}
