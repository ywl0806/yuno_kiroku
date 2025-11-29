package models

type FaceLocation struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}

type FaceDetection struct {
	FaceLocation FaceLocation `json:"face_location"`
	Embedding    []float64    `json:"embedding"`
}
type FaceDetectionResponse struct {
	Faces []FaceDetection `json:"faces"`
}
