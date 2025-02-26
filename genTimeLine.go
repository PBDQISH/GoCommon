package test0226

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"time"
)

func minutesFromNow(givenTime time.Time) int {
	now := time.Now()
	timeDiff := now.Sub(givenTime).Minutes()
	return int(timeDiff / 10)
}

func formatDateTime(date time.Time) map[string]string {
	year, month, day := date.Date()
	hours, minutes, seconds := date.Clock()
	return map[string]string{
		"date":    fmt.Sprintf("%d-%02d-%02d", year, month, day),
		"year":    fmt.Sprintf("%d", year),
		"month":   fmt.Sprintf("%02d", month),
		"day":     fmt.Sprintf("%02d", day),
		"hours":   fmt.Sprintf("%02d", hours),
		"minutes": fmt.Sprintf("%02d", minutes),
		"seconds": fmt.Sprintf("%02d", seconds),
	}
}

func getViewList(list []map[string]interface{}, typeStr string) []map[string]interface{} {
	var arrZero []int
	for i, item := range list {
		rowSpanNum := 1
		for j := i + 1; j < len(list); j++ {
			if typeStr == "pre" {
				if fmt.Sprintf("%s%s", item["date"], item["hours"]) == fmt.Sprintf("%s%s", list[j]["date"], list[j]["hours"]) {
					arrZero = append(arrZero, j)
					rowSpanNum++
				} else {
					break
				}
			} else {
				if item["date"] == list[j]["date"] {
					arrZero = append(arrZero, j)
					rowSpanNum++
				} else {
					break
				}
			}
		}
		if contains(arrZero, i) {
			rowSpanNum = 0
		}
		list[i]["rowSpan"] = rowSpanNum
		list[i]["largeW"] = rowSpanNum == 1
	}
	return list
}

func contains(slice []int, element int) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

func CreateTimeArray(replayTime *time.Time) map[string]interface{} {
	var inputDateTime time.Time
	if replayTime != nil {
		inputDateTime = *replayTime
	} else {
		inputDateTime = time.Now()
	}

	var past48Hours []map[string]interface{}
	var future48Hours []map[string]interface{}

	// 生成过去48小时的时间数组
	for i := 0; i < 24*60/5; i++ {
		dateTime := inputDateTime.Add(-time.Duration(i*5) * time.Minute)
		minutes := (dateTime.Minute() / 5) * 5
		dateTime = dateTime.Truncate(time.Hour).Add(time.Duration(minutes) * time.Minute)
		past48Hours = append([]map[string]interface{}{
			{
				"date":       formatDateTime(dateTime)["date"],
				"year":       formatDateTime(dateTime)["year"],
				"month":      formatDateTime(dateTime)["month"],
				"day":        formatDateTime(dateTime)["day"],
				"hours":      formatDateTime(dateTime)["hours"],
				"minutes":    formatDateTime(dateTime)["minutes"],
				"seconds":    formatDateTime(dateTime)["seconds"],
				"type":       "pre",
				"searchType": "pre",
				"id":         uuid.New().String(),
			},
		}, past48Hours...)
		if i == 0 {
			past48Hours[0]["type"] = "current"
			past48Hours[0]["searchType"] = "current"
		}
	}

	diffMinutes := minutesFromNow(inputDateTime)
	if diffMinutes > 1 {
		// 生成未来48小时的时间数组 288 两天
		diffNumber := diffMinutes
		if diffNumber > 288 {
			diffNumber = 288
		}
		for i := 1; i < diffNumber; i++ {
			dateTime := inputDateTime.Add(time.Duration(i*10) * time.Minute)
			minutes := (dateTime.Minute() / 10) * 10
			dateTime = dateTime.Truncate(time.Hour).Add(time.Duration(minutes) * time.Minute)
			future48Hours = append(future48Hours, map[string]interface{}{
				"date":       formatDateTime(dateTime)["date"],
				"year":       formatDateTime(dateTime)["year"],
				"month":      formatDateTime(dateTime)["month"],
				"day":        formatDateTime(dateTime)["day"],
				"hours":      formatDateTime(dateTime)["hours"],
				"minutes":    formatDateTime(dateTime)["minutes"],
				"seconds":    formatDateTime(dateTime)["seconds"],
				"type":       "pre",
				"searchType": "next",
				"id":         uuid.New().String(),
			})
		}
	} else {
		// 生成未来48小时的时间数组
		for i := 1; i < 48; i++ {
			dateTime := inputDateTime.Add(time.Duration(i) * time.Hour)
			future48Hours = append(future48Hours, map[string]interface{}{
				"date":       formatDateTime(dateTime)["date"],
				"year":       formatDateTime(dateTime)["year"],
				"month":      formatDateTime(dateTime)["month"],
				"day":        formatDateTime(dateTime)["day"],
				"hours":      formatDateTime(dateTime)["hours"],
				"minutes":    formatDateTime(dateTime)["minutes"],
				"seconds":    formatDateTime(dateTime)["seconds"],
				"type":       "next",
				"searchType": "next",
				"id":         uuid.New().String(),
			})
		}
	}

	past48Hours = getViewList(past48Hours, "pre")
	if diffMinutes > 1 {
		future48Hours = getViewList(future48Hours, "pre")
	} else {
		future48Hours = getViewList(future48Hours, "next")
	}

	timeList := append(past48Hours, future48Hours...)
	currentItem := findCurrentItem(timeList)
	currentTime := fmt.Sprintf("%s %s:%s:00", currentItem["date"], currentItem["hours"], currentItem["minutes"])

	return map[string]interface{}{
		"currentTime": currentTime,
		"searchType":  currentItem["searchType"],
		"timeList":    timeList,
	}
}

func findCurrentItem(list []map[string]interface{}) map[string]interface{} {
	for _, item := range list {
		if item["type"] == "current" {
			return item
		}
	}
	return nil
}

//func main() {
//	// 示例调用
//	replayTime := time.Now()
//	result := CreateTimeArray(&replayTime)
//	fmt.Println(MapToTimeList(result["timeList"].([]map[string]interface{})))
//}

// 将 []map[string]interface{} 映射成 []TimeList
func MapToTimeList(timeList []map[string]interface{}) []TimeList {
	var result []TimeList
	for _, item := range timeList {
		yearInt, err := strconv.Atoi(item["year"].(string))
		if err != nil {
			logx.Error(err)
		}
		timeStruct := TimeList{
			Year:       int64(yearInt),
			Month:      item["month"].(string),
			Day:        item["day"].(string),
			Hours:      item["hours"].(string),
			Seconds:    item["seconds"].(string),
			SearchType: item["searchType"].(string),
			Date:       item["date"].(string),
			Minutes:    item["minutes"].(string),
			Type:       item["type"].(string),
			Id:         item["id"].(string),
			RowSpan:    item["rowSpan"].(int),
			LargeW:     item["largeW"].(bool),
		}
		result = append(result, timeStruct)
	}
	return result
}

type TimeList struct {
	Date       string `json:"date"`
	Year       int64  `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	Hours      string `json:"hours"`
	Minutes    string `json:"minutes"`
	Seconds    string `json:"seconds"`
	SearchType string `json:"searchType"`
	Type       string `json:"type"`
	Id         string `json:"id"`
	RowSpan    int    `json:"rowSpan"`
	LargeW     bool   `json:"largeW"`
}
