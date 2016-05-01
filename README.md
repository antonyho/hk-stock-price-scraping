# hk-stock-price-scraping
Hong Kong stock prices data scraping from [Hong Kong Yahoo Finance](https://hk.finance.yahoo.com/)

This tool scraps opening and closing stock prices for all Hong Kong HSI listed companies. Data will be stored to SQLite3 database file named "stockprice.db".

There are 3 tables in the database: "daily", "weekly", "monthly". There store stock prices in different time intervals respectively.


#### Technical information

This tool is written in [Go](https://golang.org/).

##### Dependencies:

- [go-sqlite3](https://github.com/mattn/go-sqlite3)

##### Build:
```
go get github.com/antonyho/hk-stock-price-scraping
```
On Linux/Mac OS X:
```
go build -o saver github.com/antonyho/hk-stock-price-scraping
```
On Windows:
```
go build -o saver.exe github.com/antonyho/hk-stock-price-scraping
```
