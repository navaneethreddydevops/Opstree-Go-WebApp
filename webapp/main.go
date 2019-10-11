package webapp

import (
    "net/http"
    log "github.com/sirupsen/logrus"
    "os"
    dbcheck "github.com/dimiro1/health/db"
    "github.com/dimiro1/health"
    "github.com/dimiro1/health/redis"
    "gopkg.in/ini.v1"
    "github.com/gorilla/mux"
    "go.elastic.co/apm/module/apmgorilla"
)

func Run() {

	propertyfile := "/etc/conf.d/ot-go-webapp/application.ini"
    if fileExists(propertyfile) {
        vaules, err := ini.Load(propertyfile)
        if err != nil {
            log.Error("No property file found in " + propertyfile)
        }
        redisHost = vaules.Section("redis").Key("REDIS_HOST").String()
        redisPort = vaules.Section("redis").Key("REDIS_PORT").String()
        logStdout()
        log.WithFields(log.Fields{
            "file": propertyfile,
          }).Info("Reading properties from " + propertyfile)
        logFile("info")
        log.WithFields(log.Fields{
          "file": propertyfile,
        }).Info("Reading properties from " + propertyfile)
    } else {
        redisHost = os.Getenv("REDIS_HOST")
        redisPort = os.Getenv("REDIS_PORT")
        logStdout()
        log.WithFields(log.Fields{
            "file": propertyfile,
          }).Info("No property file found, using environment variables")
        logFile("info")
        log.WithFields(log.Fields{
            "file": propertyfile,
          }).Info("No property file found, using environment variables")
    }

    generateLogsFile()
    createDatabaseTable()
    db := dbConn()
    mysql := dbcheck.NewMySQLChecker(db)
    r := mux.NewRouter()
    handler := health.NewHandler()
    handler.AddChecker("MySQL", mysql)
    handler.AddChecker("Redis", redis.NewChecker("tcp", redisHost + ":" + redisPort))
    http.Handle("/health", handler)
    r.HandleFunc("/", Index)
    r.HandleFunc("/show", Show)
    r.HandleFunc("/new", New)
    r.HandleFunc("/edit", Edit)
    r.HandleFunc("/insert", Insert)
    r.HandleFunc("/update", Update)
    r.HandleFunc("/delete", Delete)
    r.Use(apmgorilla.Middleware())
    // http.Handle("/health", handler)
    // http.HandleFunc("/", Index)
    // http.HandleFunc("/show", Show)
    // http.HandleFunc("/new", New)
    // http.HandleFunc("/edit", Edit)
    // http.HandleFunc("/insert", Insert)
    // http.HandleFunc("/update", Update)
    // http.HandleFunc("/delete", Delete)
    go func() {
			http.ListenAndServe(":8080", r)
    }()
    http.ListenAndServe(":8081", nil)
}
