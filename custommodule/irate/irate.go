package irate

import (
	"math"
)

// PMT func
func PMT(rate float64, nper float64, pv float64, fv float64) float64 {
	if rate == 0 {
		return -(pv + fv) / nper
	}

	pvif := math.Pow(1+rate, nper)

	return rate / (pvif - 1) * -(pv*pvif + fv)
}

// IPMT func
func IPMT(pv float64, pmt float64, rate float64, per float64) float64 {
	tmp := math.Pow(1+rate, per)
	return 0 - (pv*tmp*rate + pmt*(tmp-1))
}

// PPMT func
func PIPMT(rate float64, per float64, nper float64, pv float64, fv float64) (float64, float64) {
	if per < 1 || per >= nper+1 {
		return 0, 0
	}

	pmt := PMT(rate, nper, pv, fv)
	ipmt := IPMT(pv, pmt, rate, per-1)

	return pmt - ipmt, ipmt
}

// FLATANNUAL func
func FLATANNUAL(rate float64, v float64, months float64) (monthlyloan float64, monthlyinterest float64, totalmonthlypay float64, accumulatedtotalpay float64) {
	monthlyloan = v / months
	monthlyinterest = (rate / 12) * v
	totalmonthlypay = monthlyloan + monthlyinterest
	accumulatedtotalpay = totalmonthlypay * months

	return monthlyloan, monthlyinterest, totalmonthlypay, accumulatedtotalpay
}

// ONETIMEPAYMENT func
func ONETIMEPAYMENT(rate float64, v float64, months float64) (monthlyloan float64, monthlyinterest float64, totalmonthlypay float64, accumulatedtotalpay float64) {
	monthlyloan = v / months
	monthlyinterest = (rate * v) / months
	totalmonthlypay = monthlyloan + monthlyinterest
	accumulatedtotalpay = totalmonthlypay * months

	return monthlyloan, monthlyinterest, totalmonthlypay, accumulatedtotalpay
}
