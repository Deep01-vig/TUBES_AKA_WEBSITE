package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

var data []int

func binaryIterative(arr []int, target int) int {
	low, high := 0, len(arr)-1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

func binaryRecursive(arr []int, target, low, high int) int {
	if low > high {
		return -1
	}
	mid := (low + high) / 2
	if arr[mid] == target {
		return mid
	} else if arr[mid] < target {
		return binaryRecursive(arr, target, mid+1, high)
	} else {
		return binaryRecursive(arr, target, low, mid-1)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, _ := r.FormFile("file")
	defer file.Close()

	reader := csv.NewReader(file)
	data = []int{}

	reader.Read() // skip header
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		val, _ := strconv.Atoi(row[0])
		data = append(data, val)
	}
	w.Write([]byte("Data uploaded"))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	target, _ := strconv.Atoi(r.URL.Query().Get("target"))

	start := time.Now()
	binaryIterative(data, target)
	iterTime := time.Since(start).Microseconds()

	start = time.Now()
	binaryRecursive(data, target, 0, len(data)-1)
	recTime := time.Since(start).Microseconds()

	resp := map[string]interface{}{
		"iterative": iterTime,
		"recursive": recTime,
		"size":      len(data),
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", nil)
}
