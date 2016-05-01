package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "net/http"
    "time"
    db "github.com/antonyho/hk-stock-price-scraping/db"
)

func main() {
    db := db.DefaultDBTool()
    
    dailyStockPriceURL := "http://real-chart.finance.yahoo.com/table.csv?s=%04d.HK&a=0&b=4&c=%4d&d=4&e=1&f=%04d&g=d&ignore=.csv"
    weeklyStockPriceURL := "http://real-chart.finance.yahoo.com/table.csv?s=%04d.HK&a=00&b=4&c=%04d&d=04&e=1&f=%04d&g=w&ignore=.csv"
    monthlyStockPriceURL := "http://real-chart.finance.yahoo.com/table.csv?s=%04d.HK&a=00&b=4&c=%04d&d=04&e=1&f=%04d&g=m&ignore=.csv"
    
    spinner := []byte{'-', '\\', '|', '/'}
    
    fmt.Println()
    startTime := time.Now()
    currentYear := startTime.Year()
    for stockNo := 1; stockNo <= 9999; stockNo++ {
        fmt.Printf("\r%cWorking Stock:%5d", spinner[stockNo % 4], stockNo)
        
        // Get daily stock price
        dailyStockPrices := getStockPrices(stockNo, dailyStockPriceURL, 2000, currentYear)
        if dailyStockPrices == nil {
            continue
        }
        
        // CSV format:
        // Date,Open,High,Low,Close,Volume,Adj Close
        db.Add("daily", stockNo, dailyStockPrices[1:])
        
        // Get weekly stock price
        monthlyStockPrices := getStockPrices(stockNo, weeklyStockPriceURL, 2000, currentYear)
        if monthlyStockPrices == nil {
            continue
        }
        
        // CSV format:
        // Date,Open,High,Low,Close,Volume,Adj Close
        db.Add("weekly", stockNo, monthlyStockPrices[1:])
        
        // Get monthly stock price
        yearlyStockPrices := getStockPrices(stockNo, monthlyStockPriceURL, 2000, currentYear)
        if yearlyStockPrices == nil {
            continue
        }
        
        // CSV format:
        // Date,Open,High,Low,Close,Volume,Adj Close
        db.Add("monthly", stockNo, yearlyStockPrices[1:])
    }
    
    db.Close()
    
    fmt.Println()
    fmt.Printf("Done (Elapsed Time: %v)\n", time.Since(startTime))
}

func getStockPrices(stockNo int, url string, startYear int, endYear int) ([][]string) {
    stockPriceResp, err := http.Get(fmt.Sprintf(url, stockNo, startYear, endYear))
    if (err != nil) || (stockPriceResp.StatusCode != 200) {
        if err != nil {
            log.Println(err)
        }
        stockPriceResp.Body.Close()
        return nil
    }
    csvReader := csv.NewReader(stockPriceResp.Body)
    stockPrices, err := csvReader.ReadAll();
    stockPriceResp.Body.Close()
    if err != nil {
        log.Println(err)
        return nil
    }
    
    return stockPrices
}