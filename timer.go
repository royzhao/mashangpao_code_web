package main

import (
//	"time"
)

func HotTimerList() {
	images := QueryImage()
	var hotList []int
	var hot int
	for i := 0; i < len(images); i++ {
		hot = images[i].Star + images[i].Fork
		hotList = append(hotList, hot)
	}
}

func Qsort(images []CRImage, hot []int, start int, end int) {
	if start < end {
		pivot := Partition(images, hot, start, end)
		Qsort(images, hot, start, pivot-1)
		Qsort(images, hot, pivot+1, end)
	}
}

func Partition(images []CRImage, hot []int, start int, end int) int {
	hotTmp := hot[start]
	imageTmp := images[start]
	i := start + 1
	j := end
	for i < j {
		for ; hot[j] < hot[i] && i < j; j-- {
			hot[i] = hot[j]
			images[i] = images[j]
		}
		for ; hot[j] < hot[i] && i < j; i++ {
			hot[j] = hot[i]
			images[j] = images[i]
		}
	}
	hot[i] = hotTmp
	images[i] = imageTmp
	return i
}
