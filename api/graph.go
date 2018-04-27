package api

import (
	"bytes"
	"encoding/base64"
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/frozzare/go-httpapi"
	"github.com/frozzare/statscoll/stat"
	"github.com/wcharczuk/go-chart"
)

const graphHTMLPage = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Statscoll</title>
		<style>
			* { margin: 0; padding: 0; }
			html, body {
				background: url(data:image/png;base64,{{ .img }});
				background-repeat:no-repeat;
				background-position: center center;
				min-height:100%;
			}
		</style>
	</head>
	<body>
	</body>
</html>
`

func generateGraphHTMLPage(c string) ([]byte, error) {
	tmpl, err := template.New("main").Parse(graphHTMLPage)
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"img": base64.StdEncoding.EncodeToString([]byte(c)),
	}); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func (h *Handler) handleGraph(w http.ResponseWriter, r *http.Request, ps httpapi.Params) {
	var stats []*stat.Stat

	query, err := h.statsQuery(r, ps)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Execute query and find any errors.
	if err := query.Find(&stats).Error; err != nil || len(stats) == 0 {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Sort stats so the last one is listed first.
	sort.Slice(stats, func(i, j int) bool {
		return time.Unix(stats[i].Timestamp, 0).After(time.Unix(stats[j].Timestamp, 0))
	})

	// Genrate x and y values.
	var xvalues []time.Time
	var yvalues []float64
	for _, s := range stats {
		yvalues = append(yvalues, float64(s.Value))
		xvalues = append(xvalues, time.Unix(s.Timestamp, 0))
	}

	// Get graph title.
	title := r.URL.Query().Get("title")
	if len(title) == 0 {
		title = ps.ByName("metric")
	}

	graph := chart.Chart{
		Title:      title,
		TitleStyle: chart.StyleShow(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 50,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 1.0,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorAlternateBlue,
					FillColor:   chart.ColorAlternateBlue.WithAlpha(64),
				},
				XValues: xvalues,
				YValues: yvalues,
			},
		},
	}

	if len(r.URL.Query().Get("image")) > 0 {
		w.Header().Set("Content-Type", "image/png")
		if err := graph.Render(chart.PNG, w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/html")
	buf := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buf); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	b, err := generateGraphHTMLPage(buf.String())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
