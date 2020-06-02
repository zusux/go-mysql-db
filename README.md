# go-mysql-db

### 连接数据库
```
  db := Db.NewDb("127.0.0.1",3306,"zusux","root","123456","") 
  db.Conn.Connection()
``` 
  
  
### 插入数据
`id, err := db.Table("user").Insert(map[string]interface{}{"username":"aaa","nickname":"sssd","password":"ddd"},false)`
  
### 更新数据
```
  rows,err := db.Table("user").
		Where("password","=","eee").
		Where("username","in",[]interface{}{"889"}).
		WhereCond(false).
		WhereOr("username","=","880").
		Limit(0,10).
		Order("id","desc").
		Debug(true).
		Update(map[string]interface{}{"username":"889","nickname":"ccc","password":"eee"})
```
### 删除数据
``` 
    rows,err := db.Table("user").
		Where("password","=","eee").
		Where("username","in",[]interface{}{"889"}).
		WhereCond(false).
		WhereOr("username","=","880").
		Limit(0,10).
		Order("id","desc").
		Debug(true).
		Delete()
```    
### 查询数据
 
#### 查询多条记录
```
	all,err := db.Table("book").
		Alias("b").
		Distinct(true).
		Force("idx").
		Field([]string{"book_name"}).
		Where("project_id","=","1").
		Where("book_name","=","工程").
		Join("project p","p.project_id = b.project_id","left").
		Join("user u","u.user_id = b.user_id","left").
		Order("sort","asc").
		Order("book_name","desc").
		Union("select * from book",false).
		Group("project_id").
		Group("book_name").
		Having("num","=","3").
		Having("name","=","测试").
		Lock(true).
		Limit(0,10).
		Select()
```   
#### 查询一条记录
`
record,err := db.Table("book").Where("project_id","=","1").Where("book_name","=","工程").Order("sort","asc").Order("book_name","desc").Find() `
		
#### 查询单个字段
```
  value,err := db.Table("book").
		Where("project_id","=","1").
		Where("book_name","=","工程").
		Order("sort","asc").
		Order("book_name","desc").
		Value("author")		
  ```  
  ### 聚合函数 
   
  #### Count
  ``` 
  count,err :=db.Table("book").Where("project_id","=","1").Count() 
  ```
  
  #### Max
  ```
  max,err :=db.Debug(true).Table("book").Where("project_id","=","1").Max("book_id")
  ```
  
  #### Min
  ```
  min,err :=db.Debug(true).Table("book").Where("project_id","=","1").Min("book_id")
  ```
  
  #### Avg
  ```
  avg,err :=db.Debug(true).Table("book").Where("project_id","=","1").Avg("number")
  ```
  
  #### Sum
 ```
 avg,err :=db.Debug(true).Table("book").Where("project_id","=","1").Sum("number")
 ```
