package vestaboard

// 6×22 for Flagship, 3×15 for Note
type BoardLayout [][]int

type MessageResult struct {
	ID     string      `json:"id"`
	Layout BoardLayout `json:"layout"`
}

type SendResult struct {
	Status  string `json:"status"`
	ID      string `json:"id"`
	Created int64  `json:"created"`
}



//extracted from docs, unsused
type Transition string

const (
	TransitionClassic Transition = "classic"
	TransitionWave    Transition = "wave"
	TransitionDrift   Transition = "drift"
	TransitionCurtain Transition = "curtain"
)

type TransitionSpeed string

const (
	SpeedGentle TransitionSpeed = "gentle"
	SpeedFast   TransitionSpeed = "fast"
)

// TransitionInfo holds the current transition configuration.
type TransitionInfo struct {
	Transition      Transition      `json:"transition"`
	TransitionSpeed TransitionSpeed `json:"transitionSpeed"`
}
