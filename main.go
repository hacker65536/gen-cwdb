package main

import (
        "bufio"
        "encoding/json"
        "fmt"
        "os"
        "regexp"
)

var (
        metricnames = []string{
                "ActiveTransactions",
                "AuroraBinlogReplicaLag",
                "AuroraReplicaLagMaximum",
                "BinLogDiskUsage",
                "BacktrackChangeRecordsCreationRate",
                "BacktrackChangeRecordsStored",
                "BacktrackWindowActual",
                "BacktrackWindowAlert",
                "BlockedTransactions",
                "CommitLatency",
                "CommitThroughput",
                "CPUCreditBalance",
                "CPUCreditUsage",
                "CPUUtilization",
                "DatabaseConnections",
                "DDLLatency",
                "DDLThroughput",
                "Deadlocks",
                "DeleteLatency",
                "DeleteThroughput",
                "DMLLatency",
                "DMLThroughput",
                "EngineUptime",
                "FreeableMemory",
                "FreeLocalStorage",
                "InsertLatency",
                "InsertThroughput",
                "LoginFailures",
                "NetworkThroughput",
                "NetworkTransmitThroughput",
                "Queries",
                "ResultSetCacheHitRatio",
                "SelectLatency",
                "SelectThroughput",
                "UpdateLatency",
                "UpdateThroughput",
                "VolumeBytesUsed",
                "VolumeReadIOPs",
                "VolumeWriteIOPs",
                "AuroraReplicaLag",
                "AuroraReplicaLagMinimum",
                "NetworkReceiveThroughput",
                "BufferCacheHitRatio",
        }

        dblistfile = "dblist"

        dblist []string
        prefix = "${prefix}"
        region = "${region}"
)

type M map[string]interface{}
type A []interface{}

func chkerr(err error) {
        if err != nil {
                fmt.Println("ErrOutput: ", err)
                os.Exit(1)
        }
}

func init() {
        //read file

        file, err := os.Open(dblistfile)
        chkerr(err)
        defer file.Close()

        sc := bufio.NewScanner(file)
        for i := 1; sc.Scan(); i++ {
                if err := sc.Err(); err != nil {
                        chkerr(err)
                        break
                }

                str := sc.Text()
                reg := regexp.MustCompile(`^\s*?#|^\s*?//`)
                if !reg.MatchString(str) {
                        dblist = append(dblist, str)
                }
        }

}

func main() {

        body := M{
                "widgets": func() A {
                        var a A
                        for _, mn := range metricnames {
                                var m A
                                for _, db := range dblist {

                                        reg := regexp.MustCompile(`aurora`)
                                        var dim string
                                        if reg.MatchString(db) {
                                                dim = "DBClusterIdentifier"
                                        } else {
                                                dim = "DBInstanceIdentifier"
                                        }

                                        metrics := A{
                                                "AWS/RDS", mn, dim, prefix + "-" + db,
                                                M{
                                                        "stat":   "Maximum",
                                                        "period": 60,
                                                },
                                        }
                                        m = append(m, metrics)
                                }

                                prop := M{
                                        "title":   prefix + "-" + mn,
                                        "view":    "timeSeries",
                                        "stacked": false,
                                        "region":  region,
                                        "metrics": m,
                                }
                                b := M{
                                        "type":       "metric",
                                        "x":          0,
                                        "y":          0,
                                        "width":      24,
                                        "height":     6,
                                        "properties": prop,
                                }

                                a = append(a, b)
                        }
                        return a
                }(),
        }

        bytes, err := json.MarshalIndent(body, "", "    ")
        if err == nil {
                jsonstring := string(bytes)
                fmt.Println(jsonstring)
        } else {
                fmt.Println("Err: ", err)
        }

}
