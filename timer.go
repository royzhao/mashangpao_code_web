package main

import (
	"encoding/json"
)

var (
	starWeight = 3
	forkWeight = 2
)

type HotImages struct {
	List []CRImage
}

func HotTimerList() {
	images := QueryImage()
	var hotList []int
	var hot int
	for i := 0; i < len(images); i++ {
		hot = images[i].Star*starWeight + images[i].Fork*forkWeight
		hotList = append(hotList, hot)
	}
	Qsort(images, hotList, 0, len(hotList)-1)
	//	fmt.Println(images)
	var key = "hotimage"
	//	buf, _ := json.Marshal(images)
	buf, _ := json.Marshal(HotImages{List: images})
	redis_client.Set(key, buf)
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
	i := start
	j := end + 1
	for i < j {
		j--
		for ; i < j && hot[j] <= hotTmp; j-- {
		}
		hot[i] = hot[j]
		images[i] = images[j]
		i++
		for ; i < j && hot[i] >= hotTmp; i++ {
		}
		hot[j] = hot[i]
		images[j] = images[i]
	}
	hot[j] = hotTmp
	images[j] = imageTmp
	return j
}
