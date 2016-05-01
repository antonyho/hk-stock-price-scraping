package db

import (
    "database/sql"
    "fmt"
    sqlite3 "github.com/mattn/go-sqlite3"
    "log"
)

type DBTool struct {
    filename        string
    dbConn          *sql.DB
}

func NewDBTool(filename string) *DBTool {
    conn, err := sql.Open("sqlite3", filename)
    if err != nil {
        log.Fatal(err)
    }
    dbTool := &DBTool{filename: filename, dbConn: conn}
    if !dbTool.DBExist() {
        dbTool.Initialise()
    }
    
    return dbTool
}

func DefaultDBTool() *DBTool {
    dbFilename := "stockprice.db"
    return NewDBTool(dbFilename)
}

func (dbTool *DBTool) DBExist() (bool) {
    selectDailyTable := "SELECT 1 FROM daily LIMIT 1"
    _, err := dbTool.Query(selectDailyTable)
    if (err != nil) && (err.(sqlite3.Error).Code == 1) {
        return false
    }
    
    return true
}

func (dbTool *DBTool) Initialise() {
    dailyStockPriceTable := "CREATE TABLE daily (stock INTEGER, date TEXT, open REAL, high REAL, low, REAL, close REAL, volume INTEGER)"
    dailyStockPriceUniqueConstraint := "CREATE UNIQUE INDEX cnstr_daily_1 ON daily(stock, date)"
    weeklyStockPriceTable := "CREATE TABLE weekly (stock INTEGER, date TEXT, open REAL, high REAL, low, REAL, close REAL, volume INTEGER)"
    weeklyStockPriceUniqueConstraint := "CREATE UNIQUE INDEX cnstr_weekly_1 ON weekly(stock, date)"
    monthlyStockPriceTable := "CREATE TABLE monthly (stock INTEGER, date TEXT, open REAL, high REAL, low, REAL, close REAL, volume INTEGER)"
    monthlyStockPriceUniqueConstraint := "CREATE UNIQUE INDEX cnstr_monthly_1 ON monthly(stock, date)"
    
    if _, err := dbTool.dbConn.Exec(dailyStockPriceTable); err != nil {
        log.Fatal(err)
    }
    if _, err := dbTool.dbConn.Exec(dailyStockPriceUniqueConstraint); err != nil {
        log.Fatal(err)
    }
    if _, err := dbTool.dbConn.Exec(weeklyStockPriceTable); err != nil {
        log.Fatal(err)
    }
    if _, err := dbTool.dbConn.Exec(weeklyStockPriceUniqueConstraint); err != nil {
        log.Fatal(err)
    }
    if _, err := dbTool.dbConn.Exec(monthlyStockPriceTable); err != nil {
        log.Fatal(err)
    }
    if _, err := dbTool.dbConn.Exec(monthlyStockPriceUniqueConstraint); err != nil {
        log.Fatal(err)
    }
}

func (dbTool *DBTool) Query(queryStr string) (*sql.Rows, error) {
    return dbTool.dbConn.Query(queryStr)
}

func (dbTool *DBTool) Add(table string, stock int, data [][]string) {
    insertStatement := "INSERT INTO %s (stock, date, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?, ?)"
    
    tx, err := dbTool.dbConn.Begin()
    if err != nil {
        log.Fatal(err)
    }
    stmt, err := tx.Prepare(fmt.Sprintf(insertStatement, table))
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    for _, row := range data {
        _, err = stmt.Exec(stock, row[0], row[1], row[2], row[3], row[4], row[5])
        if err != nil {
            if err.(sqlite3.Error).ExtendedCode == sqlite3.ErrConstraint.Extend(8) {
                continue
            }
            log.Fatal(err)
        }
    }
    tx.Commit()
}


func (dbTool *DBTool) Close() {
    dbTool.dbConn.Close()
}