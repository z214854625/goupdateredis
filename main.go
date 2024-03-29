package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	// 打开或创建文件，如果文件不存在将会创建一个新的
	// os.O_APPEND: 如果文件存在，移动到文件末尾
	// os.O_CREATE: 创建文件，如果文件不存在
	// os.O_WRONLY: 以写入模式打开文件
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Cannot create/open file", err)
		return
	}
	defer file.Close()

	Log(file, "业务执行中!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	StartRedis(file)
	Log(file, "业务执行结束!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	// InstallSignal(func() {
	// 	fmt.Println("bye")
	// 	os.Exit(0)
	// })
}

// 写入内容到文件
func Log(writer *os.File, args ...interface{}) {
	if writer == nil {
		fmt.Println("writer nil")
		return
	}
	str := MixJoin(args)
	//拼接时间，格式化为字符串，使用Go的时间布局格式
	now := time.Now()
	layout := "2006-01-02 15:04:05"
	formattedTime := now.Format(layout)
	str = formattedTime + " " + str
	//输出到console
	fmt.Println(str)
	//写入文本
	str = str + "\n"
	_, err := writer.WriteString(str)
	if err != nil {
		fmt.Println("Cannot write to file", err)
		return
	}
}

// MixJoin 使用fmt.Sprintf将不同类型参数拼接成字符串
func MixJoin(args ...interface{}) string {
	var parts []string
	for _, arg := range args {
		part := fmt.Sprintf("%v", arg) // 使用%v格式化任意类型的值
		parts = append(parts, part)
	}
	return strings.Join(parts, " ") // 使用空格作为分隔符
}

func StartRedis(writer *os.File) {
	//解析命令行参数
	if len(os.Args) != 5 {
		Log(writer, "No argument provided.")
		return
	}
	strAddr := os.Args[2]
	strPwd := os.Args[4]
	Log(writer, "redis param=", strAddr, strPwd)
	// 连接到Redis
	rdsCli := redis.NewClient(&redis.Options{
		Addr:     strAddr, // Redis地址 127.0.0.1:6379
		Password: strPwd,  // 密码（无密码则为空）
		DB:       0,       // 使用的数据库
	})
	if rdsCli == nil {
		Log(writer, "rdsCli nil ", strAddr, strPwd)
		return
	}
	_, err1 := rdsCli.Ping().Result()
	if err1 != nil {
		Log(writer, "ping error=", err1, strAddr, strPwd)
		return
	}
	//redis测试
	//TestRedis(writer, rdsCli)
	//执行业务
	Execute(writer, rdsCli)
}

func Execute(writer *os.File, rdsCli *redis.Client) {
	//1.删除video的所有录像
	keys, err := rdsCli.Keys("video*").Result() // 获取所有键，"*"为通配符
	if err != nil {
		Log(writer, "Execute keys cmd error1!", err)
		return
	}
	delNum := 0
	for _, v1 := range keys {
		Log(writer, "del key=", v1)
		_, err0 := rdsCli.Del(v1).Result()
		if err0 != nil {
			Log(writer, "del key failed!", err0, v1)
			continue
		}
		delNum++
	}
	//处理太空掠夺
	keys1, err1 := rdsCli.Keys("tkld:records*").Result()
	if err1 != nil {
		Log(writer, "Execute keys cmd error2!", err1)
		return
	}
	delNum1 := 0
	for _, v1 := range keys1 {
		Log(writer, "del key=", v1)
		_, err0 := rdsCli.Del(v1).Result()
		if err0 != nil {
			Log(writer, "Execute del key failed!", err0, v1)
			continue
		}
		delNum1++
	}
	//处理太空掠夺
	keys2, err2 := rdsCli.Keys("arena:records:world*").Result()
	if err2 != nil {
		Log(writer, "Execute keys cmd error2!", err2)
		return
	}
	delNum2 := 0
	for _, v1 := range keys2 {
		Log(writer, "del key=", v1)
		_, err0 := rdsCli.Del(v1).Result()
		if err0 != nil {
			Log(writer, "Execute del key failed!", err0, v1)
			continue
		}
		delNum2++
	}
	Log(writer, "video keyNum=", len(keys), "delNum= ", delNum)
	Log(writer, "tkld keyNum=", len(keys1), "delNum= ", delNum1)
	Log(writer, "arena keyNum=", len(keys2), "delNum= ", delNum2)
}

func TestRedis(writer *os.File, rdsCli *redis.Client) {
	//写入key
	for i := 1; i <= 100; i++ {
		str := strconv.Itoa(i)
		key := "testkey" + str
		rdsCli.Set(key, str, 0)
	}
	//读取并删除key
	keys, err := rdsCli.Keys("testkey*").Result() // 获取所有键，"*"为通配符
	if err != nil {
		Log(writer, "TestRedis Error getting keys:", err)
		return
	}
	delNum := 0
	for k1, v1 := range keys {
		Log(writer, "k1=", k1, "val=", v1)
		_, err2 := rdsCli.Del(v1).Result()
		if err2 != nil {
			Log(writer, "TestRedis del key failed!", err2)
			continue
		}
		delNum++
	}
	Log(writer, "keyNum=", len(keys), "delNum= ", delNum)
}
