package sdr

import (
	"math"
	"math/cmplx"
)

type ComplexSource interface {
	Source() ([]complex64, error)
}

type RealSink interface {
	Sink([]float32) error
}

type ComplexFilter interface {
	Filter([]complex64) ([]complex64, error)
}

type Demodulator interface {
	Demodulate(input []complex64, output []float32) (int, error)
}

type Rotate90Filter struct {
	// currentAngle int
}

func (fi *Rotate90Filter) Filter(samples []complex64) ([]complex64, error) {
	// for i := 0; i < len(samples); i++ {
	// 	switch fi.currentAngle {
	// 	case 0:
	// 		// noop
	// 	case 1:
	// 		samples[i] = complex(-imag(samples[i]), real(samples[i]))
	// 	case 2:
	// 		samples[i] = -samples[i]
	// 	case 3:
	// 		samples[i] = complex(imag(samples[i]), -real(samples[i]))
	// 	}
	// 	fi.currentAngle = (fi.currentAngle + 1) & 3
	// }
	// return samples, nil
	for i := 0; i < len(samples); i += 4 {
		samples[i+1] = complex(-imag(samples[i+1]), real(samples[i+1]))
		samples[i+2] = -samples[i+2]
		samples[i+3] = complex(imag(samples[i+3]), -real(samples[i+3]))
	}
	return samples, nil
}

func PolarDiscriminator(a, b complex128) float64 {
	return cmplx.Phase(a*cmplx.Conj(b)) / math.Pi
}

func PolarDiscriminator32(a, b complex64) float32 {
	return FastPhase32(a * Conj32(b)) // / math.Pi
}

type FMDemodFilter struct {
	pre complex64
}

func (fi *FMDemodFilter) Demodulate(input []complex64, output []float32) (int, error) {
	for i := 0; i < len(input); i++ {
		pcm := PolarDiscriminator32(input[i], fi.pre)
		fi.pre = input[i]
		output[i] = pcm
	}
	return len(input), nil
}

type LowPassDownsampleRationalFilter struct {
	Fast, Slow int

	sum       float32
	prevIndex int
}

func (fi *LowPassDownsampleRationalFilter) Filter(samples []float32) ([]float32, error) {
	i2 := 0
	fastSlowRatio := float32(fi.Slow) / float32(fi.Fast)
	for i := 0; i < len(samples); i++ {
		fi.sum += samples[i]
		fi.prevIndex += fi.Slow
		if fi.prevIndex < fi.Fast {
			continue
		}
		i2++
		samples[i2] = fi.sum * fastSlowRatio
		fi.prevIndex -= fi.Fast
		fi.sum = 0.0
	}
	return samples[:i2], nil
}
