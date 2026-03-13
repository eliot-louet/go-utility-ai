package curves

import "math"

type Linear struct {
	M float64
	B float64
}

func (l Linear) Apply(x float64) float64 {
	return l.M*x + l.B
}

type Logistic struct {
	M float64
	C float64
	K float64
	B float64
}

func (l Logistic) Apply(x float64) float64 {
	return l.K/(1+math.Exp(-l.M*(x-l.C))) + l.B
}

type Parabolic struct {
}

func (p Parabolic) Apply(x float64) float64 {
	return 4 * x * (1 - x)
}

type Polynomial struct {
	M float64
	C float64
	K float64
	B float64
}

func (p Polynomial) Apply(x float64) float64 {
	return p.M*math.Pow(x-p.C, p.K) + p.B
}

type Logit struct {
	M float64
	C float64
	K float64
	B float64
}

func (l Logit) Apply(x float64) float64 {
	epsilon := 1e-9
	return l.M*math.Log((x+epsilon)/(1-x+epsilon)) + l.B
}
