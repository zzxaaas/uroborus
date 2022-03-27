package common

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

// ConcatURL 连接 URL
func ConcatURL(segments ...string) string {
	urls := make([]string, len(segments))
	for i, segment := range segments {
		if i == 0 || i == len(segment)-1 {
			urls[i] = segment
		}
		urls[i] = strings.Trim(segment, "/")
	}
	return strings.Join(urls, "/")
}

// ToStringSlice 从空接口slice到字符串slice
func ToStringSlice(input []interface{}) []string {
	output := make([]string, len(input))
	for i := range input {
		output[i] = input[i].(string)
	}
	return output
}

// ConvertStringSlice 从字符串slice到空接口slice
func ConvertStringSlice(input []string) []interface{} {
	output := make([]interface{}, len(input))
	for i := range input {
		output[i] = input[i]
	}
	return output
}

// NewInt 返回指针
func NewInt(value int) *int {
	return &value
}

// ConvertToPQInt64Array []int to pq.Int64Array
func ConvertToPQInt64Array(input []int) pq.Int64Array {
	var output pq.Int64Array
	for _, in := range input {
		output = append(output, int64(in))
	}
	return output
}

// ConvertUint64ToPQInt64Array []uint64 to pq.Int64Array
func ConvertUint64ToPQInt64Array(input []uint64) pq.Int64Array {
	var output pq.Int64Array
	for _, in := range input {
		output = append(output, int64(in))
	}
	return output
}

// ConvertPQInt64ArrToCommaString -
func ConvertPQInt64ArrToCommaString(input pq.Int64Array) string {
	var strList []string
	for _, it := range input {
		strList = append(strList, strconv.Itoa(int(it)))
	}
	return strings.Join(strList, ",")
}

// ConvertPQInt64ArrToIntSlice -
func ConvertPQInt64ArrToIntSlice(input pq.Int64Array) []int {
	var intList []int
	for _, it := range input {
		intList = append(intList, int(it))
	}
	return intList
}

// ConvertIntSliceToCommaString -
func ConvertIntSliceToCommaString(input []int) string {
	var strList []string
	for _, it := range input {
		strList = append(strList, strconv.Itoa(it))
	}
	return strings.Join(strList, ",")
}

// ConvertUint64ToCommaString -
func ConvertUint64ToCommaString(input []uint64) string {
	var strList []string
	for _, it := range input {
		strList = append(strList, strconv.FormatUint(it, 10))
	}
	return strings.Join(strList, ",")
}

// StringInSlice 字符串是否在Slice中
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// IsStringInMapValue 字符串是否在Map的Value中
func IsStringInMapValue(amb string, m map[string]string) bool {
	for _, v := range m {
		if amb == v {
			return true
		}
	}
	return false
}

// Float64Ptr makes a copy and returns the pointer to a float64.
func Float64Ptr(v float64) *float64 {
	return &v
}

// StringPtr makes a copy and returns the pointer to a string.
func StringPtr(v string) *string {
	return &v
}

// Int32Ptr makes a copy and returns the pointer to a int32.
func Int32Ptr(v int32) *int32 {
	return &v
}

// Int64Ptr makes a copy and returns the pointer to a int64.
func Int64Ptr(v int64) *int64 {
	return &v
}

// IntPtr makes a copy and returns the pointer to a int.
func IntPtr(v int) *int {
	return &v
}

// UintPtr -
func UintPtr(i uint) *uint {
	return &i
}

// TimePtr -
func TimePtr(t time.Time) *time.Time {
	return &t
}

// IsNil 用于interface判断是否为空
// 若为map，slice，则长度为空
func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	if vi.Kind() == reflect.Slice {
		return vi.Len() == 0
	}
	if vi.Kind() == reflect.Map {
		return vi.Len() == 0
	}
	return i == nil
}

// GetAllFields 获取某类型所有的字段
// 匿名字段自动递归解析，非匿名结构体字段不递归解析
// eg : types := common.GetAllFields(reflect.TypeOf(store.SegmentTaskRes{}),make([]string,0))
func GetAllFields(t reflect.Type, rec []string) []string {
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		fi := t.Field(i)
		// 匿名
		if fi.Anonymous && fi.Name == fi.Type.Name() {
			rec = GetAllFields(fi.Type, rec)
		} else {
			rec = append(rec, fi.Name)
		}
	}
	return rec
}

// AlertIntLength 改变Int64的长度
// num 输入数字
// length 截取前length位
func AlertIntLength(num int64, length int) (res int64, err error) {
	fmt.Printf("AlertIntLength input: num: %d, length: %d", num, length)
	if num <= 1000000000 {
		return res, nil
	}
	s := strconv.FormatInt(num, 10)[:length]
	res, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return res, fmt.Errorf("time convert failed: %v", err)
	}
	return res, nil
}

// RemoveRepeatedElement 去除切片中重复的元素
func RemoveRepeatedElement(array []string) []string {
	temMap := make(map[string]interface{})
	for _, e := range array {
		temMap[e] = nil
	}
	res := []string{}
	for k, _ := range temMap {
		res = append(res, k)
	}
	return res
}

// GetFirstDateOfWeek 获取本周周一的日期
func GetFirstDateOfWeek() (weekMonday time.Time) {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekMonday = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int64
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func Md5(file *multipart.FileHeader) string {
	f, err := file.Open()
	if err != nil {
		return ""
	}
	defer f.Close()

	m := md5.New()
	_, err = io.Copy(m, f)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(m.Sum(nil))
}
