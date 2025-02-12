package serve

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/robfig/cron"
	// _ "github.com/mattn/go-sqlite3"
)

// 定时任务
func Test_time(t *testing.T) {
	for {
		time.Sleep(time.Millisecond * 5000)
		fmt.Println("Hello TigerwolfC")
	}
}

// 定时任务
func Test_cron(t *testing.T) {
	c := cron.New()
	// 定时5秒，每5秒执行print5
	c.AddFunc("*/5 * * * * *", print5)
	// 定时15秒，每5秒执行print5
	// c.AddFunc("*/15 * * * * *", print15)
	c.Start()
	defer c.Stop()
	select {}
}

func print5() {
	fmt.Println("哈哈哈哈")
}
func Test_exec_shell(t *testing.T) {
	res, err := exec_shell(`docker exec -it mysql /bin/bash -c "mysqldump -uroot -pmypass vpndata > /var/lib/mysql/databackup/vpn1.sql && echo $?"`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	// res1, err := exec_shell("$?")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("这是判断结果" + res1)
}
func Test_exec_ase(t *testing.T) {
	str := "SSS"
	key := make([]byte, 16)

	if len(str) > 16 {
		copy(key, str[:16])
	} else {
		copy(key, str)
		remain := 16 - len(str)
		for i := 0; i < remain; i++ {
			// key = append(key, 'A')
			last := key[len(key)-1]
			if last == 'A' {
				key = append(key, 'B')
			} else {
				key = append(key, 'A')
			}
		}
	}

	println(key)
	println(string(key))

}

// 阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func exec_shell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()

	return out.String(), err
}
