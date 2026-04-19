# vestaboard-go

A Go wrapper for the [Vestaboard Read/Write API](https://docs.vestaboard.com/docs/read-write-api/introduction) and [VBML](https://docs.vestaboard.com/docs/vbml/) composition service.

## Installation

```bash
go get github.com/jam941/vestaboard-go
```

## Authentication

Generate an API token from the Vestaboard mobile or web app under **Settings → Developer**. Store it in an environment variable — never hardcode it.

## Quick Start

```go
import "github.com/jam941/vestaboard-go/vestaboard"

client := vestaboard.New(os.Getenv("VBOARD_TOKEN"))       // Flagship (6×22)
// client := vestaboard.NewNote(os.Getenv("VBOARD_TOKEN")) // Vestaboard Note (3×15)
```

## API Reference

### Reading the board

```go
msg, err := client.GetMessage()
// msg.ID     — UUID of the current message
// msg.Layout — [][]int grid of character codes
```

### Sending messages

**Plain text** — the simplest option, no layout control:
```go
result, err := client.SendText("HELLO WORLD", false)
```

**Raw character codes** — full manual control over every cell:
```go
layout := vestaboard.BoardLayout{
    {vestaboard.CharH, vestaboard.CharI, 0, 0, ...},
    // one []int per row
}
result, err := client.SendCharacters(layout, false)
```

Pass `forced: true` to either method to bypass configured quiet hours.

### VBML composition

**Format** — convert a string to a board layout:
```go
layout, err := client.FormatMessage("HELLO WORLD")
```

**Compose** — full VBML layout with alignment, positioning, props, and multi-component support:
```go
layout, err := client.Compose(vestaboard.ComposeRequest{
    Components: []vestaboard.Component{
        {
            Template: "VESTABOARD",
            Style: &vestaboard.ComponentStyle{
                Justify: "center",
                Align:   "center",
            },
        },
    },
})
```

Board dimensions (`style.height` / `style.width`) are filled automatically from the client's board type — no manual configuration needed.

**ComposeAndSend** — compose and send in one call:
```go
result, err := client.ComposeAndSend(vestaboard.ComposeRequest{
    Components: []vestaboard.Component{
        {Template: "{{title}}", Style: &vestaboard.ComponentStyle{Justify: "center"}},
        {Template: "{{subtitle}}"},
    },
    Props: map[string]string{
        "title":    "SCORE",
        "subtitle": "42 - 17",
    },
}, false)
```

#### Component options

| Field | Type | Description |
|-------|------|-------------|
| `Template` | `string` | Text with optional `{{prop}}` variables and `{N}` character codes |
| `RawCharacters` | `[]int` | Direct character code array, overrides Template |
| `Style.Height` | `int` | Component height in rows |
| `Style.Width` | `int` | Component width in columns |
| `Style.Justify` | `string` | `left`, `right`, `center`, `justified` |
| `Style.Align` | `string` | `top`, `bottom`, `center`, `justified` |
| `Style.AbsolutePosition` | `*Position` | Exact `{X, Y}` placement |

### Transitions

```go
ti, err := client.GetTransition()
// ti.Transition      — current transition type
// ti.TransitionSpeed — current speed

ti, err = client.SetTransition(vestaboard.TransitionWave, vestaboard.SpeedGentle)
```

Available transitions: `TransitionClassic`, `TransitionWave`, `TransitionDrift`, `TransitionCurtain`  
Available speeds: `SpeedGentle`, `SpeedFast`

> Transitions are supported on Flagship and Vestaboard Note only, not Note Arrays.

### Character codes

Named constants are provided for all supported characters:

| Constant | Code | Character |
|----------|------|-----------|
| `CharSpace` | 0 | _(blank)_ |
| `CharA` | 1 | A |
| `CharB` | 2 | B |
| `CharC` | 3 | C |
| `CharD` | 4 | D |
| `CharE` | 5 | E |
| `CharF` | 6 | F |
| `CharG` | 7 | G |
| `CharH` | 8 | H |
| `CharI` | 9 | I |
| `CharJ` | 10 | J |
| `CharK` | 11 | K |
| `CharL` | 12 | L |
| `CharM` | 13 | M |
| `CharN` | 14 | N |
| `CharO` | 15 | O |
| `CharP` | 16 | P |
| `CharQ` | 17 | Q |
| `CharR` | 18 | R |
| `CharS` | 19 | S |
| `CharT` | 20 | T |
| `CharU` | 21 | U |
| `CharV` | 22 | V |
| `CharW` | 23 | W |
| `CharX` | 24 | X |
| `CharY` | 25 | Y |
| `CharZ` | 26 | Z |
| `Char1` | 27 | 1 |
| `Char2` | 28 | 2 |
| `Char3` | 29 | 3 |
| `Char4` | 30 | 4 |
| `Char5` | 31 | 5 |
| `Char6` | 32 | 6 |
| `Char7` | 33 | 7 |
| `Char8` | 34 | 8 |
| `Char9` | 35 | 9 |
| `Char0` | 36 | 0 |
| `CharExclamation` | 37 | ! |
| `CharAt` | 38 | @ |
| `CharHash` | 39 | # |
| `CharDollar` | 40 | $ |
| `CharLeftParen` | 41 | ( |
| `CharRightParen` | 42 | ) |
| `CharHyphen` | 44 | - |
| `CharPlus` | 46 | + |
| `CharAnd` | 47 | & |
| `CharEquals` | 48 | = |
| `CharSemicolon` | 49 | ; |
| `CharColon` | 50 | : |
| `CharApostrophe` | 52 | ' |
| `CharQuote` | 53 | " |
| `CharPercent` | 54 | % |
| `CharComma` | 55 | , |
| `CharPeriod` | 56 | . |
| `CharSlash` | 59 | / |
| `CharQuestion` | 60 | ? |
| `CharHeart` | 62 | ♥ |
| `CharRed` | 63 | _(solid red)_ |
| `CharOrange` | 64 | _(solid orange)_ |
| `CharYellow` | 65 | _(solid yellow)_ |
| `CharGreen` | 66 | _(solid green)_ |
| `CharBlue` | 67 | _(solid blue)_ |
| `CharPurple` | 68 | _(solid purple)_ |
| `CharWhite` | 69 | _(solid white)_ there may be problems with this if used with local api|
| `CharBlack` | 70 | _(solid black)_ there may be problems with this if used with local api|
| `CharFilled` | 71 | _(filled)_ there may be problems with this if used with local api|

## Board types

| Constant | Dimensions | Constructor |
|----------|-----------|-------------|
| `BoardFlagship` | 6 rows × 22 columns | `vestaboard.New(token)` |
| `BoardNote` | 3 rows × 15 columns | `vestaboard.NewNote(token)` |

## Rate limiting

Vestaboard recommends no more than **one message per 15 seconds** to avoid dropped messages.

## Running the example

```bash
VBOARD_TOKEN=your_token_here go run ./cmd/example/
```

## Running tests

No API token required — tests use an in-process HTTP server.

```bash
go test ./vestaboard/...
```
