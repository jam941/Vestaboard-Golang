package vestaboard

type BoardType int

const (
	BoardFlagship BoardType = iota // 6 rows × 22 columns
	BoardNote                      // 3 rows × 15 columns
)

func (b BoardType) Rows() int {
	if b == BoardNote {
		return 3
	}
	return 6
}

func (b BoardType) Cols() int {
	if b == BoardNote {
		return 15
	}
	return 22
}

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


type ComposeRequest struct {
	Props      map[string]string `json:"props,omitempty"`
	Style      *BoardStyle       `json:"style,omitempty"`
	Components []Component       `json:"components"`
}

type BoardStyle struct {
	Height int `json:"height,omitempty"`
	Width  int `json:"width,omitempty"`
}

type Component struct {
	Template      string          `json:"template,omitempty"`
	RawCharacters []int           `json:"rawCharacters,omitempty"`
	Style         *ComponentStyle `json:"style,omitempty"`
}

type ComponentStyle struct {
	Height           int       `json:"height,omitempty"`
	Width            int       `json:"width,omitempty"`
	Justify          string    `json:"justify,omitempty"`  // left, right, center, justified
	Align            string    `json:"align,omitempty"`    // top, bottom, center, justified
	AbsolutePosition *Position `json:"absolutePosition,omitempty"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
