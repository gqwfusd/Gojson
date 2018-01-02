package main
import (
  "fmt"
  "encoding/json"
  "io/iouitl"
  "os"
  "time"
)

type caseinput struct {
  Domain string `json:"domain"`
  Mode   string `json:"mode"`
  Value  string `josn:"value"`
  Uuid   string `json:"uuid"`
  Methods []int32 'json:"methods"'
  Actions []int32 `json:"actions"`
  

type testcase struct {
	Apiname string                 `json:"apiname"`
	Repeat  int32                  `json:"repeat"`
	Input   caseinput              `json:"input"`
}

type casegroup struct {
	Info []testcase `json:"info"`
}

var (
	//	edge       *sdk.Edge
	centerAddr string
	edgeAddr   string
)

func main{
  runJsonCases()

}
//获取目录下所有的json文件
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 50)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

func runJsonCases() {
	//读取conf/tc目录下的cases
	edge := sdk.NewEdge()
	edge.SetAddr("10.115.25.15:5050")
	edge.Connect()
	var beginTime, endTime time.Time
	files, _ := ListDir("./", "json")
	for i := 0; i < len(files); i++ {
		file, err := ioutil.ReadFile(files[i])
		if err != nil {
			fmt.Println("Failed to open json file '%s': %s\n", err.Error())
			return
		}
		var r casegroup
		//		var a testcase
		err = json.Unmarshal(file, &r)
		if err != nil {
			fmt.Println("err was:", err.Error())
		}

		for _, v := range r.Info {
			fmt.Println(v.Apiname)
			switch v.Apiname {
      case "ConsultAvailableMethod":
				fmt.Println("####################ConsultAvailableMethod####################")
				domain := v.Input["domain"]
				mode := v.Input["mode"]
				value := v.Input["value"]
				uuid := v.Input["uuid"]
        methods := v.Input["methoods"]
        actions := v.Input["actions"]

				beginTime = time.Now()
				resCon, err := edge.ConsultAvailableMethod(domain, mode, value, uuid, methods, actions)
				endTime = time.Now()
				elapsedS := int32(endTime.Sub(beginTime).Nanoseconds()) / 1000000

				if resCon.Status.Code == 0 {
					util.RollLog().Print("|"+"success"+"|", elapsedS, "|", resCon)
					fmt.Println("Passed")
				} else {
					util.RollLog().Print("|"+"failed"+"|", elapsedS, "|", resCon)
					fmt.Println("Fialed")
				}
				fmt.Println(resCon, err)
        			default:
				fmt.Println("not find case input")

			}

		}

	}

}

