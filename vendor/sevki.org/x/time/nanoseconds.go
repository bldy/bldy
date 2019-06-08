// Copyright 2015 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package time // import "sevki.org/x/time"
import (
	"fmt"
	"time"
)

// NsReadable takes in a nanosecond duration and converts it to it's least decimalled unit
// there is also a lot of rounding going on here
func NsReadable(ns int64) string {

	t := float64(ns)
	unit := "t"
	switch {
	case t > float64(time.Minute.Nanoseconds())/100:
		unit = "m"
		return fmt.Sprintf("%.2f%s", t/float64(time.Minute.Nanoseconds()), unit)
	case t > float64(time.Second.Nanoseconds())/100:
		unit = "s"
		return fmt.Sprintf("%.2f%s", t/float64(time.Second.Nanoseconds()), unit)
	case t > 10000:
		unit = "ms"
		return fmt.Sprintf("%.2f%s", t/10000, unit)
	case t > 10:
		unit = "μs"
		return fmt.Sprintf("%.2f%s", t/10, unit)
	default:
		return fmt.Sprintf("%.2f%s", t, unit)
	}
}

const lssthnd = "less than %d %s"
const lssthns = "less than a %s"
const aboutnd = "about %d %s"
const day time.Duration = 86400000000000
const month time.Duration = 2628000000001209
const year time.Duration = 31535999999964780

// InWords returns time in words after subtracting it from now.
func InWords(t time.Time) string {

	now := time.Now()
	d := now.Sub(t)
	if d >= time.Second && d <= (time.Second*4) {
		return fmt.Sprintf(lssthnd, 5, "seconds")
	} else if d >= (time.Second*5) && d < (time.Second*10) {
		return fmt.Sprintf(lssthnd, 10, "seconds")
	} else if d >= (time.Second*10) && d < (time.Second*20) {
		return fmt.Sprintf(lssthnd, 20, "seconds")
	} else if d >= (time.Second*20) && d < (time.Second*40) {
		return "half a minute"
	} else if d >= (time.Second*40) && d < (time.Second*60) {
		return fmt.Sprintf(lssthns, "minute")
	} else if d >= (time.Second*60) && d < time.Minute+(time.Second*30) {
		return "1 minute"
	} else if d >= time.Minute+(time.Second*30) && d < (time.Minute*44)+(time.Second*30) {
		return fmt.Sprintf("%d minutes", (d / time.Minute))
	} else if d >= (time.Minute*44)+(time.Second*30) && d < (time.Minute*89)+(time.Second*30) {
		return fmt.Sprintf(aboutnd, d/time.Hour, "hour")
	} else if d >= (time.Minute*89)+(time.Second*30) && d < (time.Hour*29)+(time.Minute*59)+(time.Second*30) {
		return fmt.Sprintf(aboutnd, (d / time.Hour), "hours")
	} else if d >= (time.Hour*23)+(time.Minute*59)+(time.Second*30) && d < (time.Hour*41)+(time.Minute*59)+(time.Second*30) {
		return "1 day"
	} else if d >= (time.Hour*41)+(time.Minute*59)+(time.Second*30) && d < (day*29)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) {
		return fmt.Sprintf("%d days", d/(time.Hour*24))
	} else if d >= (day*29)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) && d < (day*59)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) {
		return fmt.Sprintf(aboutnd, 1, "month")
	} else if d >= (day*59)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) && d < (year) {
		return fmt.Sprintf(aboutnd, d/month+1, "months")
	} else if d >= year && d < year+(3*month) {
		return fmt.Sprintf(aboutnd, 1, "year")
	} else if d >= year+(3*month) && d < year+(9*month) {
		return "over 1 year"
	} else if d >= year+(9*month) && d < (year*2) {
		return "almost 2 years"
	} else {
		return fmt.Sprintf(aboutnd, d/year, "years")
	}

}

// DurationInWords returns duration in words.
func DurationInWords(d time.Duration) string {
	if d >= time.Second && d <= (time.Second*4) {
		return fmt.Sprintf(lssthnd, 5, "seconds")
	} else if d >= (time.Second*5) && d < (time.Second*10) {
		return fmt.Sprintf(lssthnd, 10, "seconds")
	} else if d >= (time.Second*10) && d < (time.Second*20) {
		return fmt.Sprintf(lssthnd, 20, "seconds")
	} else if d >= (time.Second*20) && d < (time.Second*40) {
		return "half a minute"
	} else if d >= (time.Second*40) && d < (time.Second*60) {
		return fmt.Sprintf(lssthns, "minute")
	} else if d >= (time.Second*60) && d < time.Minute+(time.Second*30) {
		return "1 minute"
	} else if d >= time.Minute+(time.Second*30) && d < (time.Minute*44)+(time.Second*30) {
		return fmt.Sprintf("%d minutes", (d / time.Minute))
	} else if d >= (time.Minute*44)+(time.Second*30) && d < (time.Minute*89)+(time.Second*30) {
		return fmt.Sprintf(aboutnd, d/time.Hour, "hour")
	} else if d >= (time.Minute*89)+(time.Second*30) && d < (time.Hour*29)+(time.Minute*59)+(time.Second*30) {
		return fmt.Sprintf(aboutnd, (d / time.Hour), "hours")
	} else if d >= (time.Hour*23)+(time.Minute*59)+(time.Second*30) && d < (time.Hour*41)+(time.Minute*59)+(time.Second*30) {
		return "1 day"
	} else if d >= (time.Hour*41)+(time.Minute*59)+(time.Second*30) && d < (day*29)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) {
		return fmt.Sprintf("%d days", d/(time.Hour*24))
	} else if d >= (day*29)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) && d < (day*59)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) {
		return fmt.Sprintf(aboutnd, 1, "month")
	} else if d >= (day*59)+(time.Hour*23)+(time.Minute*59)+(time.Second*30) && d < (year) {
		return fmt.Sprintf(aboutnd, d/month+1, "months")
	} else if d >= year && d < year+(3*month) {
		return fmt.Sprintf(aboutnd, 1, "year")
	} else if d >= year+(3*month) && d < year+(9*month) {
		return "over 1 year"
	} else if d >= year+(9*month) && d < (year*2) {
		return "almost 2 years"
	} else {
		return fmt.Sprintf(aboutnd, d/year, "years")
	}

}
