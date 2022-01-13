package main

import (
    "fmt"
    "database/sql"
    "log"
    "net/http"
    yaml "gopkg.in/yaml.v2"
	"github.com/go-redis/redis"
	"io/ioutil"
    "github.com/julienschmidt/httprouter"
    _ "github.com/go-sql-driver/mysql"
    "os"
    "path/filepath"
)
type Conf struct {
    Listen struct {
        Port int `yaml:"port"`
    }
    MySQL struct {
        Db string `yaml:"db"`
        Host string `yaml:"host"`
        Port int `yaml:"port"`
        User string `yaml:"user"`
        Pass string `yaml:"pass"`
    }
    Redis struct {
        Host string `yaml:"host"`
        Port int `yaml:"port"`
        Db int `yaml:"db"`
        Pass string `yaml:"pass"`
    }
	Logfile string `yaml:"logfile"`
}
//全局公共变量
var (
    Db *sql.DB
    logger *log.Logger
    client *redis.Client
    cnf Conf
)
//先利用日志文件创建调试输出，然后初始化数据库连接
func init() {
    dir,_:=filepath.Abs(filepath.Dir(os.Args[0]))
    //抓全部的配置信息
    yamlBytes, err := ioutil.ReadFile("config.yml")
    if err!=nil {
        log.Fatalf("无法打开环境配置文件: %v",err)
    }
    yaml.Unmarshal(yamlBytes,&cnf)
    file, err := os.OpenFile(dir+"/"+cnf.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("无法打开日志文件", err)
    }
    logger = log.New(file, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
    client=redis.NewClient(&redis.Options{
        Addr:       fmt.Sprintf("%s:%d",cnf.Redis.Host,cnf.Redis.Port),
        Password:   cnf.Redis.Pass,
        DB:         cnf.Redis.Db,
    })
    _, err = client.Ping().Result()
    if err!=nil {
        logger.Fatalf("redis连接异常：%v\n",err)
    }
    dsn:=fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true", cnf.MySQL.User, cnf.MySQL.Pass, cnf.MySQL.Host, cnf.MySQL.Port, cnf.MySQL.Db)
    //fmt.Println("Data Source Name: ",dsn)
    Db,err=sql.Open("mysql",dsn)
    if err!=nil {
        logger.Fatalf("mysql连接异常：%v\n",err)
    }
}
func main() {
    // handle static assets
    router := httprouter.New()
    router.ServeFiles("/static/*filepath", http.Dir("static"))
    //页面
    router.GET("/", index)
    router.GET("/ea/:Type", equipType)
    router.GET("/pa/:Type", propertyType)
    router.GET("/pra/:Type", prescriptionType)
    router.GET("/qa/:Type", questionType)
    router.GET("/mi/:Type", missionType)
    router.GET("/magic/:Type", magicType)
    router.GET("/stunt/:Type", stuntType)
    router.GET("/up/:Name", upgradeRole)
    router.GET("/ma/common", enemyCommon)
    router.GET("/ma/basic", enemyBasic)
    router.GET("/ma/resistance", enemyResistance)
    router.GET("/ma/skill", enemySkill)
    router.GET("/ma/drop", enemyDrop)
    //依据双层分类，搜寻包括装备、道具、配方在内的绝大部分物品，先做个捡场景物品的试手
    router.GET("/pu/:Type", pickUp)
    router.POST("/pu/:Type", pickUp)
    router.GET("/fm/:cat/:Type", findItem)
    router.POST("/fm/:cat/:Type", findItem)
    logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",cnf.Listen.Port),router))
}
