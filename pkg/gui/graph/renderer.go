package graph

import (
	"fmt"
	"image/color"
	"math"
	"physicsGUI/pkg/function"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type GraphRenderer struct {
	graph   *GraphCanvas
	objects []fyne.CanvasObject

	// size of canvas
	size *fyne.Size

	// margin for labels etc.
	margin float32
}

type GraphScope struct {
	Min function.Coordinate
	Max function.Coordinate
}
type GraphRange struct {
	Min float64
	Max float64
}

// returns the minimum size needed for the graph
func (r *GraphRenderer) MinSize() fyne.Size {
	return fyne.NewSize(500, 200)
}

// initializes the base strcuture for every graph
func (r *GraphRenderer) base() {
	// background
	r.graph.background.Resize(*r.size)
	r.graph.background.Move(fyne.NewPos(0, 0))
	r.AddObject(r.graph.background)

	// title
	title := &canvas.Text{
		Text:      r.graph.config.Title,
		Color:     titleColor,
		TextSize:  16,
		TextStyle: fyne.TextStyle{Bold: true},
	}
	title.Move(fyne.NewPos(r.size.Width/2-float32(len(title.Text)*4), 0))
	r.AddObject(title)

	// x-axis
	r.AddObject(&canvas.Line{
		StrokeColor: axesColor,
		StrokeWidth: 1,
		Position1:   fyne.NewPos(r.margin-2, r.size.Height-r.margin+2),
		Position2:   fyne.NewPos(r.size.Width-r.margin/2, r.size.Height-r.margin+2),
	})

	// y-axis
	r.AddObject(&canvas.Line{
		StrokeColor: axesColor,
		StrokeWidth: 1,
		Position1:   fyne.NewPos(r.margin-2, r.size.Height-r.margin+2),
		Position2:   fyne.NewPos(r.margin-2, 0.5*r.margin),
	})
}

// draws the whole graph
func (r *GraphRenderer) Layout(size fyne.Size) {
	// clear objects
	r.objects = make([]fyne.CanvasObject, 0)

	// size of the graph
	r.size = &size

	// About layout when size is zero and therefor component is not visible
	if r.size.Width < r.MinSize().Width || r.size.Height < r.MinSize().Height {
		return
	}

	// set the base for the canvas
	r.base()

	// calculate the maximum scope
	scope := function.GetMaximumScope(append(r.graph.functions, r.graph.loadedData...)...)
	if scope == nil {
		r.DrawErrorMessage("Scope error")
		return
	}

	if len(r.graph.functions) == 0 || r.graph.functions[0].GetDataCount() < 1 && len(r.graph.loadedData) == 0 {
		r.DrawErrorMessage("No data available")
		return
	}

	// Add Remove Buttons
	r.DrawRemoveButtons()

	// draw model lines
	if r.graph.config.IsLog {
		for _, f := range r.graph.functions {
			f.Range(r.graph.config.DisplayRange.Min, r.graph.config.DisplayRange.Max)
			points, iPoints := f.Model(r.graph.config.Resolution, true)
			r.DrawGraphLog(scope, points, iPoints, pointColor)
		}
		for i, d := range r.graph.loadedData {
			d.Range(r.graph.config.DisplayRange.Min, r.graph.config.DisplayRange.Max)
			points, iPoints := d.Model(r.graph.config.Resolution, true)
			dataColor := DataTrackColors[i%len(DataTrackColors)]
			r.DrawGraphLog(scope, points, iPoints, dataColor)
		}
		r.DrawGridLog(scope)
		return
	}

	for _, f := range r.graph.functions {
		points, iPoints := f.Model(r.graph.config.Resolution, false)
		r.DrawGraphLinear(scope,
			points.Filter(r.graph.config.DisplayRange.Min, r.graph.config.DisplayRange.Max),
			iPoints.Filter(r.graph.config.DisplayRange.Min, r.graph.config.DisplayRange.Max),
			pointColor)
	}
	for i, d := range r.graph.loadedData {
		points, iPoints := d.Model(r.graph.config.Resolution, false)
		dataColor := DataTrackColors[i%len(DataTrackColors)]
		r.DrawGraphLinear(scope,
			points.Filter(r.graph.config.DisplayRange.Min, r.graph.config.DisplayRange.Max),
			iPoints.Filter(r.graph.config.DisplayRange.Min, r.graph.config.DisplayRange.Max),
			dataColor)
	}
	r.DrawGridLinear(scope)
}

// display remove buttons at the right border
func (r *GraphRenderer) DrawRemoveButtons() {
	offsetY := float32(0)
	startY := float32(0)
	startX := r.size.Width - r.margin
	for _, d := range r.graph.dataRemoveButtons {
		offsetY += d.Size().Height + RemoveButtonTopPadding
		for _, o := range d.Objects {
			o.Move(fyne.NewPos(startX, startY+offsetY))
			r.AddObject(o)
		}
	}
}

// draw an error message onto the graph
func (r *GraphRenderer) DrawErrorMessage(message string) {
	errorMsg := &canvas.Text{
		Text:     message,
		Color:    titleColor,
		TextSize: 16,
	}
	errorMsg.Move(fyne.NewPos(
		r.size.Width/2-float32(len(errorMsg.Text)*4),
		r.size.Height/2-errorMsg.Size().Height/2-r.margin/2))

	r.AddObject(errorMsg)
}

// draw a linear graph
func (r *GraphRenderer) DrawGraphLinear(scope *function.Scope, points, iPoints function.Points, pointColor color.Color) {
	// calc available space
	availableWidth := r.size.Width - (1.5 * r.margin)
	availableHeight := r.size.Height - (1.5 * r.margin)

	// complete range
	xRange := math.Abs(scope.MaxX - scope.MinX)
	yRange := math.Abs(scope.MaxY - scope.MinY)

	oX, oY := float32(0), float32(0)

	// draw line based on interpolated (resolution) points
	for i, point := range iPoints {
		// scale x value to available width
		x := float32((point.X-scope.MinX)/xRange) * availableWidth
		y := float32((point.Y-scope.MinY)/yRange) * availableHeight

		// first point is the origin
		if i == 0 {
			oX, oY = r.normalize(x, y)
			continue
		}

		xt, yt := r.normalize(x, y)

		// draw line
		r.AddObject(&canvas.Line{
			StrokeColor: pointColor,
			StrokeWidth: 1,
			Position1:   fyne.NewPos(oX, oY),
			Position2:   fyne.NewPos(xt, yt),
		})

		oX, oY = xt, yt
	}

	// draw data points
	for _, point := range points {
		// scale x value to available width
		x := float32((point.X-scope.MinX)/xRange) * availableWidth
		y := float32((point.Y-scope.MinY)/yRange) * availableHeight

		xt, yt := r.normalize(x, y)

		// error correction
		yE1 := float32((point.Y+point.Error-scope.MinY)/yRange) * availableHeight
		yE2 := float32((point.Y-point.Error-scope.MinY)/yRange) * availableHeight

		_, e1 := r.normalize(x, yE1)
		_, e2 := r.normalize(x, yE2)

		r.DrawError(xt, e1, e2, errorColor)
		r.DrawPoint(xt, yt, pointColor)
	}
}

// draw the graph in logarithmic scale
func (r *GraphRenderer) DrawGraphLog(scope *function.Scope, points, iPoints function.Points, pointColor color.Color) {
	// calc available space
	availableWidth := r.size.Width - (1.5 * r.margin)
	availableHeight := r.size.Height - (1.5 * r.margin)

	// Calculate shifts if needed for negative values
	xShift := 0.0
	if scope.MinX <= 0 {
		xShift = math.Abs(scope.MinX) + 1
	}
	yShift := 0.0
	if scope.MinY <= 0 {
		yShift = math.Abs(scope.MinY) + 2
	}

	// Calculate log ranges
	logMinX := math.Log10(scope.MinX + xShift)
	logMinX = float64(int(logMinX) - 1)
	logMaxX := math.Log10(scope.MaxX + xShift)
	logMaxX = float64(int(logMaxX))
	logMinY := math.Log10(scope.MinY + yShift)
	logMinY = float64(int(logMinY) - 1)
	logMaxY := math.Log10(scope.MaxY + yShift)
	logMaxY = float64(int(logMaxY))
	xRange := math.Abs(logMaxX - logMinX)
	yRange := math.Abs(logMaxY - logMinY)

	oX, oY := float32(0), float32(0)

	// draw line based on interpolated (resolution) points
	for i, point := range iPoints {
		// scale x and y values logarithmically
		logX := math.Log10(point.X + xShift)
		logY := math.Log10(point.Y + yShift)

		x := float32((logX-logMinX)/xRange) * availableWidth
		y := float32((logY-logMinY)/yRange) * availableHeight

		if i == 0 {
			oX, oY = r.normalize(x, y)
			continue
		}

		xt, yt := r.normalize(x, y)
		r.AddObject(&canvas.Line{
			StrokeColor: pointColor,
			StrokeWidth: 1,
			Position1:   fyne.NewPos(oX, oY),
			Position2:   fyne.NewPos(xt, yt),
		})
		oX, oY = xt, yt
	}

	// draw data points
	for _, point := range points {
		// scale x and y values logarithmically
		logX := math.Log10(point.X + xShift)
		logY := math.Log10(point.Y + yShift)

		x := float32((logX-logMinX)/xRange) * availableWidth
		y := float32((logY-logMinY)/yRange) * availableHeight

		xt, yt := r.normalize(x, y)

		// error correction (also logarithmic)
		yE1 := float32((math.Log10(point.Y+point.Error+yShift)-logMinY)/yRange) * availableHeight
		yE2 := float32((math.Log10(point.Y-point.Error+yShift)-logMinY)/yRange) * availableHeight
		_, e1 := r.normalize(x, yE1)
		_, e2 := r.normalize(x, yE2)

		r.DrawError(xt, e1, e2, pointColor)
		r.DrawPoint(xt, yt, errorColor)
	}
}

// draw grid lines and labels for linear scale
func (r *GraphRenderer) DrawGridLinear(scope *function.Scope) {
	// horizontal grid-lines + y-labels
	yGridCount := int(r.size.Height / 25)
	yStep := (scope.MaxY - scope.MinY) / float64(yGridCount)

	for i := 0; i <= yGridCount; i++ {
		value := scope.MinY + yStep*float64(i)
		yPos := r.size.Height - r.margin - float32(i)*float32(r.size.Height-1.5*r.margin)/float32(yGridCount)

		if i > 0 {
			r.DrawGridLine(fyne.NewPos(r.margin, yPos), false, false)
		}

		text := fmt.Sprintf("%.3f", value)
		if value < 0.01 {
			text = fmt.Sprintf("%.1e", value)
		}

		// label
		label := &canvas.Text{
			Text:     text,
			Color:    legendColor,
			TextSize: 12,
		}
		label.Move(fyne.NewPos(r.margin-45, yPos-10))
		r.AddObject(label)
	}

	// vertical grid-lines + x-labels
	xGridCount := int(r.size.Width / 25)
	xStep := math.Abs(scope.MaxX-scope.MinX) / float64(xGridCount)

	for i := 0; i <= xGridCount; i++ {
		xPos := r.margin + float32(i)*float32(r.size.Width-1.5*r.margin)/float32(xGridCount)

		if i > 0 {
			r.DrawGridLine(fyne.NewPos(xPos, r.margin/2), true, false)
		}

		// only draw every second label to prevent overlapping
		if i%4 == 0 {
			v := scope.MinX + xStep*float64(i)
			text := fmt.Sprintf("%.3f", v)
			if v < 0.01 {
				text = fmt.Sprintf("%.1e", v)
			}

			// label
			label := &canvas.Text{
				Text:     text,
				Color:    legendColor,
				TextSize: 12,
			}

			label.Move(fyne.NewPos(xPos-15, r.size.Height-r.margin+10))
			r.AddObject(label)
		}
	}
}

// TODO: draw grid lines and labels for logarithmic scale
func (r *GraphRenderer) DrawGridLog(scope *function.Scope) {
	// Horizontal grid-lines + y-labels (logarithmic)
	minLogY := math.Log10(math.Max(scope.MinY, 1e-10))
	minLogY = float64(int(minLogY) - 1)
	maxLogY := math.Log10(scope.MaxY)
	maxLogY = float64(int(maxLogY))
	yGridCount := int(maxLogY - minLogY)

	for i := 0; i <= yGridCount; i++ {
		// Calculate logarithmic value
		value := math.Pow(10, minLogY+float64(i))

		// Convert log space to screen space
		logVal := math.Log10(value)
		yPos := r.size.Height - r.margin -
			float32((logVal-minLogY)/(maxLogY-minLogY))*
				float32(r.size.Height-1.5*r.margin)

		if i > 0 {
			r.DrawGridLine(fyne.NewPos(r.margin, yPos), false, false)

			// Add minor grid lines between major decades
			for j := 2; j < 10; j++ {
				minorValue := value * float64(j)
				if minorValue < math.Pow(10, maxLogY) {
					minorLogVal := math.Log10(minorValue)
					minorYPos := r.size.Height - r.margin -
						float32((minorLogVal-minLogY)/(maxLogY-minLogY))*
							float32(r.size.Height-1.5*r.margin)

					r.DrawGridLine(fyne.NewPos(r.margin, minorYPos), false, true)
				}
			}
		}

		text := fmt.Sprintf("%.3f", value)
		if value < 0.01 {
			text = fmt.Sprintf("%.0e", value)
		}

		// Label for major grid lines
		label := &canvas.Text{
			Text:     text,
			Color:    legendColor,
			TextSize: 12,
		}
		label.Move(fyne.NewPos(r.margin-45, yPos-10))
		r.AddObject(label)
	}

	// Vertical grid-lines + x-labels (logarithmic)
	minLogX := math.Log10(math.Max(scope.MinX, 1e-10))
	minLogX = float64(int(minLogX) - 1)
	maxLogX := math.Log10(scope.MaxX)
	maxLogX = float64(int(maxLogX))
	xGridCount := int(maxLogX - minLogX)

	// TODO: make this more dynamic
	labelSkip := (xGridCount / (int(r.size.Width) / 50)) + 1

	for i := 0; i <= xGridCount; i++ {
		// Calculate logarithmic value
		value := math.Pow(10, minLogX+float64(i))

		// Convert log space to screen space
		logVal := math.Log10(value)
		xPos := r.margin +
			float32((logVal-minLogX)/(maxLogX-minLogX))*
				float32(r.size.Width-1.5*r.margin)

		if i > 0 {
			r.DrawGridLine(fyne.NewPos(xPos, r.margin/2), true, false)

			// Add minor grid lines between major decades
			for j := 2; j < 10; j++ {
				minorValue := value * float64(j)
				if minorValue < math.Pow(10, maxLogX) {
					minorLogVal := math.Log10(minorValue)
					minorXPos := r.margin +
						float32((minorLogVal-minLogX)/(maxLogX-minLogX))*
							float32(r.size.Width-1.5*r.margin)
					r.DrawGridLine(fyne.NewPos(minorXPos, r.margin/2), true, true)
				}
			}
		}

		// Label for major grid lines
		if i%labelSkip == 0 {
			text := fmt.Sprintf("%.3f", value)
			if value < 0.01 {
				text = fmt.Sprintf("%.0e", value)
			}
			label := &canvas.Text{
				Text:     text,
				Color:    legendColor,
				TextSize: 12,
			}
			label.Move(fyne.NewPos(xPos-25, r.size.Height-r.margin+10))
			r.AddObject(label)
		}
	}
}

// normalizes the coodinates from the bottom left of the canvas
func (r *GraphRenderer) normalize(x float32, y float32) (float32, float32) {
	return x + r.margin, r.size.Height - r.margin - y
}

// TODO: fix the small gap between points and lines
// TODO: if points are the same size as lines
// TODO: -> points are on the bottom of the lines

// draw a grid point
func (r *GraphRenderer) DrawPoint(x float32, y float32, pointColor color.Color) {
	r.AddObject(&canvas.Circle{
		FillColor: pointColor,
		Position1: fyne.NewPos(x-pointRadius, y-pointRadius),
		Position2: fyne.NewPos(x+pointRadius, y+pointRadius),
	})
}

// draw error correction lines
func (r *GraphRenderer) DrawError(x, y1, y2 float32, errorColor color.Color) {
	r.AddObject(&canvas.Line{
		StrokeColor: errorColor,
		StrokeWidth: 1,
		Position1:   fyne.NewPos(x, y1),
		Position2:   fyne.NewPos(x, y2),
	})
}

// helper function for the grid lines
func (r *GraphRenderer) DrawGridLine(pos fyne.Position, isVertical, isMinor bool) {
	line := &canvas.Line{
		StrokeColor: gridColor,
		StrokeWidth: 1,
		Position1:   pos,
	}

	if isVertical {
		line.Position2 = fyne.NewPos(pos.X, pos.Y+r.size.Height-1.5*r.margin)
	} else {
		line.Position2 = fyne.NewPos(pos.X+r.size.Width-1.5*r.margin, pos.Y)
	}

	if isMinor {
		line.StrokeColor = gridMinorColor
		line.StrokeWidth = 0.5
	}

	r.AddObject(line)
}

// returns the objects of the graph
func (r *GraphRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// destroy function (needs to be here to satisfy the interface)
func (r *GraphRenderer) Destroy() {}

// refresh function (needs to be here to satisfy the interface)
func (r *GraphRenderer) Refresh() {}

// add an object to the graph renderer
func (r *GraphRenderer) AddObject(object fyne.CanvasObject) {
	r.objects = append(r.objects, object)
}
