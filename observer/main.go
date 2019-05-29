package main

import "fmt"

// update stores the updated measurements.
type update struct {
	temp, humid, pressure int
}

type observer interface {
	update(update)
}

type subject interface {
	registerObserver(o observer)
	removeObserver(o observer)
	notifyObservers(update)
}

type displayElement interface {
	display()
}

type currentConditionsDisplay struct {
	temp, humid, pressure int
	// implements observer & displayElement
}

type statDisplay struct {
	// implements observer & displayElement
	temp, humid, pressure             []int
	tempMean, humidMean, pressureMean float32
}

type forecastDisplay struct {
	// implements displayElement
	temp, humid, pressure int
}

type weatherData struct {
	// implements subject
	observers             []observer
	temp, humid, pressure int
}

func newcurrentConditionsDisplay() observer {
	return &currentConditionsDisplay{}
}

func (d currentConditionsDisplay) update(u update) {
	d.temp = u.temp
	d.humid = u.humid
	d.pressure = u.pressure
	d.display()
}

func (d currentConditionsDisplay) display() {
	fmt.Printf(
		"CURRENT conditions are\ntemperature %v\nhumidity %v\n pressure %v\n\n",
		d.temp, d.humid, d.pressure,
	)
}

func newstatDisplay() observer {
	return &statDisplay{
		temp: make([]int, 0),
	}
}

func (d *statDisplay) update(u update) {
	d.temp = append(d.temp, u.temp)
	d.humid = append(d.humid, u.humid)
	d.pressure = append(d.pressure, u.pressure)
	// calculate 'average' mean.
	d.tempMean = calculateMean(d.temp)
	d.humidMean = calculateMean(d.humid)
	d.pressureMean = calculateMean(d.pressure)

	d.display()
}

func calculateMean(data []int) float32 {
	var mean float32
	for _, v := range data {
		mean += (float32(v))
	}
	mean = float32(mean / float32(len(data)))

	return mean
}

func (d *statDisplay) display() {
	fmt.Printf(
		"MEAN conditions are\ntemperature %.2f\nhumidity %.2f\n pressure %.2f\n\n",
		d.tempMean, d.humidMean, d.pressureMean,
	)
}

func newForecastDisplay() observer {
	return &forecastDisplay{}
}
func (d *forecastDisplay) update(u update) {
	d.display()
}

func (d *forecastDisplay) display() {
	fmt.Printf("FORECAST is what you make of it\n\n")
}

func newWeatherData() *weatherData {
	return &weatherData{observers: make([]observer, 0)}
}

func (w *weatherData) registerObserver(o observer) {
	w.observers = append(w.observers, o)
}

func (w *weatherData) removeObserver(o observer) {
	for i, obs := range w.observers {
		if obs == o {
			w.observers = append(w.observers[:i], w.observers[i+1:]...)
		}
	}
}

func (w *weatherData) notifyObservers(u update) {
	for _, obs := range w.observers {
		obs.update(u)
	}
}

func (w *weatherData) setMeasurements(u update) {
	w.temp = u.temp
	w.humid = u.humid
	w.pressure = u.pressure
	w.notifyObservers(u)
}

func main() {
	// create the subject that has the changes the observers care about.
	w := newWeatherData()
	// create and register observers
	curr := newcurrentConditionsDisplay()
	w.registerObserver(curr)
	avg := newstatDisplay()
	w.registerObserver(avg)
	forecast := newForecastDisplay()
	w.registerObserver(forecast)

	w.setMeasurements(update{
		temp:     79,
		humid:    90,
		pressure: 30,
	})
	w.setMeasurements(update{
		temp:     80,
		humid:    90,
		pressure: 20,
	})
	w.setMeasurements(update{
		temp:     10,
		humid:    90,
		pressure: 20,
	})
}
