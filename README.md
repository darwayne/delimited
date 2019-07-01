# delimited

delimited is library that simplifies dealing with various delimited file types.
  - e.g tsv, csv, space delimited and more
  - it can handle very large delimited file sizes since it parses one row at a time
  
### Example Reader Usage
```go
exampleStr := `1,2,3,4
5,6,7,8
`
exampleReader := strings.NewReader(exampleStr)
reader := delimited.NewReader(exampleReader, delimited.CommaReaderOpt())

ctx := context.Background()
err := reader.EachRow(ctx, func(row []){
	fmt.Println(row)
})
```