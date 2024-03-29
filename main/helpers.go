package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "gioui.org/layout"
    "gioui.org/text"
    "gioui.org/unit"
    "gioui.org/widget"
    "gioui.org/widget/material"
    "golang.org/x/exp/shiny/materialdesign/icons"
    "image/color"
    "io"
    "log"
    "net/http"
    "strconv"
)

// ----- Window Stuff -----
func outerSongListWrapper(gtx layout.Context, f func(songList []*Song) []layout.FlexChild, songList []*Song) layout.Dimensions {
    l := layout.Flex{
        // Vertical alignment, from top to bottom
        Axis: layout.Vertical,
        // Empty space is left at the start, i.e. at the top
        Spacing: layout.SpaceStart,
    }.Layout(gtx,
        f(songList)...,
    )
    return l
}

// ----- Theme Stuff ------
// Copied from widget.material.theme.NewTheme
func rgb(c uint32) color.NRGBA {
    return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
    return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
func mustIcon(ic *widget.Icon, err error) *widget.Icon {
    if err != nil {
        panic(err)
    }
    return ic
}

func CreateTheme(fontCollection []text.FontFace) *material.Theme {
    t := &material.Theme{
        Shaper: text.NewCache(fontCollection),
    }
    t.Palette = material.Palette{
        Fg:         rgb(0x000000),
        Bg:         rgb(0x000000),
        ContrastBg: rgb(0x717EDD),
        ContrastFg: rgb(0xffffff),
    }
    t.TextSize = unit.Sp(16)

    t.Icon.CheckBoxChecked = mustIcon(widget.NewIcon(icons.ToggleCheckBox))
    t.Icon.CheckBoxUnchecked = mustIcon(widget.NewIcon(icons.ToggleCheckBoxOutlineBlank))
    t.Icon.RadioChecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonChecked))
    t.Icon.RadioUnchecked = mustIcon(widget.NewIcon(icons.ToggleRadioButtonUnchecked))

    // 38dp is on the lower end of possible finger size.
    t.FingerSize = unit.Dp(38)

    return t
}

// -------- End Theme -----------

// --------- Styling --------------
func songLineMargins(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
    margins := layout.Inset{
        Top:    unit.Dp(2),
        Right:  unit.Dp(1),
        Bottom: unit.Dp(3),
        Left:   unit.Dp(1),
    }
    return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return d
    })
}

func songFieldsMargins(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
    margins := layout.Inset{
        Top:    unit.Dp(0.5),
        Right:  unit.Dp(3),
        Bottom: unit.Dp(0.5),
        Left:   unit.Dp(3),
    }
    return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return d
    })
}

func headerFieldsMargins(gtx layout.Context, d layout.Dimensions) layout.Dimensions {
    margins := layout.Inset{
        Top:    unit.Dp(0),
        Right:  unit.Dp(0),
        Bottom: unit.Dp(5),
        Left:   unit.Dp(0),
    }
    return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
        return d
    })
}

func generateDurationString(d int) string {
    // Create zero padded strings for the duration of songs up to the hours level.
    durationSeconds := d % 60
    durationMinutes := d / 60
    durationHours := 0

    if durationMinutes >= 60 {
        durationHours = durationMinutes / 60
        durationMinutes = durationMinutes % 60
        return fmt.Sprintf(
            "%s:%s:%s",
            generateTimePaddedStrings(durationHours),
            generateTimePaddedStrings(durationMinutes),
            generateTimePaddedStrings(durationSeconds),
        )
    }
    return fmt.Sprintf(
        "%s:%s",
        generateTimePaddedStrings(durationMinutes),
        generateTimePaddedStrings(durationSeconds),
    )
}

func generateTimePaddedStrings(i int) string {
    if i < 10 {
        // Pad with an extra zero if needed
        return fmt.Sprintf("0%s", strconv.Itoa(i))
    }
    return strconv.Itoa(i)
}

// --------- HTTP --------------
func Get(url string, responseValue interface{}, auth bool) error {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }
    req.Header.Add("Content-Type", "application/json;")
    if auth {
        req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    err = json.Unmarshal(body, &responseValue)
    if err != nil {
        return err
    }
    return nil
}

func Post(url string, requestBody, responseValue interface{}, auth bool) error {
    jsonBody, err := json.Marshal(requestBody)
    if err != nil {
        // TODO - figure out how to display a banner or something
        log.Printf("CredentialMarshallingFailure: %s\n", err)
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
    if err != nil {
        return err
    }
    if auth {
        req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    // TODO - see what happens with an empty response? - import
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    err = json.Unmarshal(body, responseValue)
    if err != nil {
        return err
    }
    return nil
}
