package utils

import (
	"math"
)

func ComputeMean(values []float64) float64 {
	var sum float64
 	for _, value := range values {
		sum = sum + value
	}
 	return (sum / float64(len(values)))
}

func ComputeStdDeviation(values []float64) float64 {
	var (
		mean         float64
		vp           []float64
		stdDiv, temp float64
	)
 	vp = make([]float64, len(values))
 	/*Calculate mean of the data points*/
	mean = ComputeMean(values)
	/*Calculate standard deviation of individual data points*/
	for i, v := range values {
		temp = v - mean
		vp[i] = (temp * temp)
	}
 	/* Finally, Compute standard Deviation of data points
	 * by taking mean of individual std. Deviation.
	 */
	stdDiv = ComputeMean(vp)
	return float64(math.Sqrt(float64(stdDiv)))
}

