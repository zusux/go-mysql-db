package Db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
)

/*
 @Author 周旭鑫
 @Email zhxx136@qq.com
 @date 2020年6月4日 22:00
 */
var (
	connections *connectionStruct
	tmpSql *sqlTmp
)

func Connect(hostname string,port int,database string,username string, password string ,prefex string,charset string,maxOpenConns int,maxIdleConns int)  {
	_config := &configStruct{
		dbtype:"mysql",
		hostname:hostname,
		port:port,
		database: database,
		username:username,
		password:password,
		dsn:"",
		charset:charset,
		Prefex:prefex,
		Debug:true,
	}
	connections = &connectionStruct{
		Config:_config,
	}
	//连接数据库
	connections.Connection()
	connections.SqlDb.SetMaxOpenConns(maxOpenConns) //设置最大打开
	connections.SqlDb.SetMaxIdleConns(maxIdleConns) //设置最大闲置连接

	//模板sql
	tmpSql = &sqlTmp{
		SelectSql:"SELECT %DISTINCT%%FIELD% FROM %TABLE%%FORCE%%JOIN%%WHERE%%GROUP%%HAVING%%UNION%%ORDER%%LIMIT%%LOCK%",
		InsertSql:"%INSERT% INTO %TABLE% (%FIELD%) VALUES (%DATA%)",
		InsertAllSql : "%INSERT% INTO %TABLE% (%FIELD%) VALUES %DATA%",
		UpdateSql:"UPDATE %TABLE% SET %SET% %JOIN%%WHERE%%ORDER%%LIMIT%",
		DeleteSql:"DELETE FROM %TABLE%  %JOIN%%WHERE%%ORDER%%LIMIT%",
		MaxSql:"SELECT MAX(%FIELD%) as zusux_max FROM %TABLE% %WHERE%",
		MinSql:"SELECT MIN(%FIELD%) as zusux_min FROM %TABLE% %WHERE%",
		CountSql:"SELECT COUNT(*) as zusux_count FROM %TABLE% %WHERE%",
		AvgSql:"SELECT AVG(%FIELD%) as zusux_avg FROM %TABLE% %WHERE%",
		SumSql:"SELECT SUM(%FIELD%) as zusux_sum FROM %TABLE% %WHERE%",
	}
}

type Where struct {
	Field string
	Condition string
	Value interface{}
}

type configStruct struct {
	dbtype string
	hostname string
	port int
	database string
	username string
	password string
	dsn string
	charset string
	Prefex string
	Debug bool
}

type connectionStruct struct {
	Config *configStruct
	SqlDb *sql.DB
	err error
}

func ( c  *connectionStruct) Connection ()  {
	dns := c.parseDsn()
	c.SqlDb ,c.err = sql.Open(c.Config.dbtype,dns)
	if c.err != nil{
		panic(c.err.Error())
	}
}
func (c *connectionStruct) parseDsn () string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		c.Config.username,
		c.Config.password,
		c.Config.hostname,
		c.Config.port,
		c.Config.database,
		c.Config.charset,
	)
	c.Config.dsn = dsn
	return dsn
}

type zdb struct {
	Conn *connectionStruct
	Tmp *sqlTmp
	Build *sqlBuild
}
type sqlTmp struct {
	SelectSql string
	InsertSql string
	InsertAllSql string
	UpdateSql string
	DeleteSql string
	MaxSql string
	MinSql string
	CountSql string
	AvgSql string
	SumSql string
}
type sqlBuild struct {
	Table_ string
	Alias_ string
	Distinct_ string
	Field_ string
	//where
	Where_ map[string]interface{}
	WhereCond_ string
	WhereOr_ map[string]interface{}
	//having
	Having_ map[string]interface{}
	HavingCond_ string
	HavingOr_ map[string]interface{}
	//Data_  map[string]interface{}
	Union_ []string
	UnionType_ string
	Force_ string
	Join_ []string
	Order_ []string
	Group_ []string
	Comment_ string
	Lock_ bool
	Offset_ int64
	Rows_ int64
	Debug_ bool
}
func (b *sqlBuild) Reset(){
	b.Table_ =""
	b.Alias_ =""
	b.Distinct_ =""
	b.Field_ ="*"
	b.Where_ = map[string]interface{}{}
	b.WhereCond_ = " and "
	b.WhereOr_ = map[string]interface{}{}

	b.Having_ = map[string]interface{}{}
	b.HavingCond_ = " and "
	b.HavingOr_ = map[string]interface{}{}

	b.Union_ = []string{}
	b.Force_  = ""
	b.Join_ = []string{}
	b.Order_ = []string{}
	b.Group_ = []string{}
	b.Comment_ = ""
	b.Lock_ = false
	b.Offset_ = 0
	b.Rows_ = 0
	b.Debug_ = false
}


func NewDb() *zdb{
	return &zdb{
		Conn:connections,
		Tmp:tmpSql,
		Build: &sqlBuild{
			Table_ : "",
			Alias_ : "",
			Distinct_ : "",
			Field_ : "*",
			Where_ : map[string]interface{}{},
			WhereCond_ : "and",
			WhereOr_ : map[string]interface{}{},
			//having
			Having_ : map[string]interface{}{},
			HavingCond_ : "and",
			HavingOr_ : map[string]interface{}{},

			//Data_ :map[string]interface{}{},
			Union_ : []string{},
			Join_ : []string{},
			Order_ : []string{},
			Group_ : []string{},
			Comment_ : "",
			Lock_ : false,
			Offset_ : 0,
			Rows_ : 0,
			Debug_ : false,
		},
	}
}

func (db *zdb) Name (name string) *zdb  {
	db.Build.Table_ = db.Conn.Config.Prefex + name
	return db
}

func (db *zdb) Table (table string) *zdb  {
	db.Build.Table_ = table
	return db
}

func (db *zdb) Alias (alias string) *zdb  {
	db.Build.Alias_ = " "+alias
	return db
}

func (db *zdb) Distinct (distinct bool) *zdb  {
	if distinct {
		db.Build.Distinct_ = "DISTINCT "
	}else{
		db.Build.Distinct_ = ""
	}

	return db
}



func (db *zdb) Field (fields []string) *zdb {
	if len(fields) > 0{
		db.Build.Field_ = strings.Join(fields,",")
	}else{
		db.Build.Field_ = "*"
	}
	return db
}

func (db *zdb) Order (field string,order string) *zdb  {
	o := []string{}
	o = append(o,field)
	o = append(o,order)
	db.Build.Order_ = append(db.Build.Order_,strings.Join(o," "))
	return db
}

func (db *zdb) Debug(debug bool) *zdb{
  db.Build.Debug_ = debug
  return db
}

func (db *zdb) Group (field string) *zdb {
	db.Build.Group_ = append(db.Build.Group_,field)
	return db
}

func (db *zdb) Comment (comment string) *zdb  {
	db.Build.Comment_ = comment
	return db
}

func (db *zdb) Join(table string, condition string, jointype string) *zdb  {
	db.Build.Join_ = append(db.Build.Join_, fmt.Sprintf("%s JOIN %s ON %s",strings.ToUpper(jointype),table,condition))
	return db
}

func (db *zdb) Lock(isLock bool) *zdb {
	db.Build.Lock_ = isLock
	return db
}

func (db *zdb) Force(index string) *zdb  {
	db.Build.Force_ = index
	return db
}

func (db *zdb) Where_(where Where) *zdb  {
	getType := reflect.TypeOf(where.Value)
	//getValue := reflect.ValueOf(option)
	k := getType.Kind()
	switch k {
	case reflect.String, reflect.Int,reflect.Int64 :
		key := fmt.Sprintf("%s %s ?",where.Field,where.Condition)
		db.Build.Where_[key] = where.Value
	case reflect.Slice, reflect.Array:
		arr := reflect.ValueOf(where.Value)
		holders := make([]string,0,arr.Len())
		values := make([]interface{}, 0, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			holders = append(holders,"?")
			s := fmt.Sprintf("%v",arr.Index(i))
			values = append(values, s)
		}
		key := fmt.Sprintf("%s %s (%s)",where.Field,where.Condition,strings.Join(holders,","))
		db.Build.Where_[key] = values
	}
	return db
}

func (db *zdb) Where(field string,condition string,option interface{}) *zdb  {
	getType := reflect.TypeOf(option)
	//getValue := reflect.ValueOf(option)
	k := getType.Kind()
	switch k {
		case reflect.String, reflect.Int,reflect.Int64 :
			key := fmt.Sprintf("%s %s ?",field,condition)
			db.Build.Where_[key] = option
		case reflect.Slice, reflect.Array:
			arr := reflect.ValueOf(option)
			holders := make([]string,0,arr.Len())
			values := make([]interface{}, 0, arr.Len())
			for i := 0; i < arr.Len(); i++ {
				holders = append(holders,"?")
				s := fmt.Sprintf("%v",arr.Index(i))
				values = append(values, s)
			}
			key := fmt.Sprintf("%s %s (%s)",field,condition,strings.Join(holders,","))
			db.Build.Where_[key] = values
	}
	return db
}
func (db *zdb) WhereCond(cond bool) *zdb{
	if cond {
		db.Build.WhereCond_ = "AND"
	}else{
		db.Build.WhereCond_ = "OR"
	}
	return db
}
func (db *zdb) WhereOr(field string,condition string,option interface{}) *zdb  {
	getType := reflect.TypeOf(option)
	//getValue := reflect.ValueOf(option)
	k := getType.Kind()
	switch k {
	case reflect.String, reflect.Int,reflect.Int64 :
		key := fmt.Sprintf("%s %s ?",field,condition)
		db.Build.WhereOr_[key] = option
	case reflect.Slice, reflect.Array:
		arr := reflect.ValueOf(option)
		holders := make([]string,0,arr.Len())
		values := make([]interface{}, 0, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			holders = append(holders,"?")
			values = append(values, arr.Index(i))
		}
		key := fmt.Sprintf("%s %s (%s)",field,condition,strings.Join(holders,","))
		db.Build.WhereOr_[key] = values
	}
	return db
}


func (db *zdb) Having(field string,condition string,option interface{}) *zdb  {
	getType := reflect.TypeOf(option)
	//getValue := reflect.ValueOf(option)
	k := getType.Kind()
	switch k {
	case reflect.String, reflect.Int,reflect.Int64 :
		key := fmt.Sprintf("%s %s ?",field,condition)
		db.Build.Having_[key] = option
	case reflect.Slice, reflect.Array:
		arr := reflect.ValueOf(option)
		holders := make([]string,0,arr.Len())
		values := make([]interface{}, 0, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			holders = append(holders,"?")
			s := fmt.Sprintf("%v",arr.Index(i))
			values = append(values, s)
		}
		key := fmt.Sprintf("%s %s (%s)",field,condition,strings.Join(holders,","))
		db.Build.Having_[key] = values
	}
	return db
}
func (db *zdb) HavingCond(cond bool) *zdb{
	if cond {
		db.Build.HavingCond_ = "AND"
	}else{
		db.Build.HavingCond_ = "OR"
	}
	return db
}
func (db *zdb) HavingOr(field string,condition string,option interface{}) *zdb  {
	getType := reflect.TypeOf(option)
	//getValue := reflect.ValueOf(option)
	k := getType.Kind()
	switch k {
	case reflect.String, reflect.Int,reflect.Int64 :
		key := fmt.Sprintf("%s %s ?",field,condition)
		db.Build.HavingOr_[key] = option
	case reflect.Slice, reflect.Array:
		arr := reflect.ValueOf(option)
		holders := make([]string,0,arr.Len())
		values := make([]interface{}, 0, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			holders = append(holders,"?")
			values = append(values, arr.Index(i))
		}
		key := fmt.Sprintf("%s %s (%s)",field,condition,strings.Join(holders,","))
		db.Build.HavingOr_[key] = values
	}
	return db
}

func (db *zdb) Union(sql string,isAll bool) *zdb  {

	var UnionType string
	if isAll{
		UnionType = "ALL"
	}else{
		UnionType = "DISTINCT"
	}
	db.Build.Union_ = append(db.Build.Union_," UNION "+UnionType+" "+sql)
	return db
}


func (db *zdb) Limit(offset int64, rows int64) *zdb  {
	if offset > 0{
		db.Build.Offset_ = offset
	}
	if rows > 0{
		db.Build.Rows_ = rows
	}
	return db
}
//执行插入语句
func (db *zdb) Insert (data map[string]interface{},replace bool) (int64,error) {
	sqlStr, binds := db.BuildInsertSql(data,replace)
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(binds...)
	if err != nil{
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil{
		return 0, err
	}
	return id,nil
}


//执行插入语句
func (db *zdb) InsertAll (data []map[string]interface{},replace bool) (int64,error) {
	sqlStr, binds := db.BuildInsertAllSql(data,replace)
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(binds...)
	if err != nil{
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil{
		return 0, err
	}
	return id,nil
}

//构建插入语句多条记录批量插入
func (db *zdb) BuildInsertAllSql(data []map[string]interface{}, replace bool) (sql string, values []interface{}) {
	keys,holders,values  := db.array_map_keys_values(data)
	var typeWords string
	if replace{
		typeWords = "REPLACE"
	}else{
		typeWords = "INSERT"
	}
	//%INSERT% INTO %TABLE% (%FIELD%) VALUES %DATA%
	replacer := strings.NewReplacer(
		"%INSERT%",typeWords,
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%FIELD%",strings.Join(keys,","),
		"%DATA%",strings.Join(holders,","),
	)
	sql  = replacer.Replace(db.Tmp.InsertAllSql)
	return
}

//构建插入语句
func (db *zdb) BuildInsertSql(data map[string]interface{}, replace bool) (sql string, values []interface{}) {
	keys,holders,values  := db.map_keys_values(data)

	var typeWords string
	if replace{
		typeWords = "REPLACE"
	}else{
		typeWords = "INSERT"
	}

	replacer := strings.NewReplacer(
		"%INSERT%",typeWords,
				  "%TABLE%",db.Build.Table_ + db.Build.Alias_,
				  "%FIELD%",strings.Join(keys,","),
				  "%DATA%",strings.Join(holders,","),
	)
	sql  = replacer.Replace(db.Tmp.InsertSql)
	return
}

func (db *zdb) map_keys_values (data map[string]interface{}) ([]string, []string, []interface{}) {
	keys := make([]string, 0, len(data))
	holders := make([]string,0,len(data))
	values := make([]interface{}, 0, len(data))
	for k,v := range data {
		keys = append(keys, k)
		holders = append(holders,"?")
		values = append(values, v)
	}
	return keys,holders,values
}

func (db *zdb) array_map_keys_values (data []map[string]interface{}) ([]string, []string, []interface{}) {
	keys := make([]string, 0, len(data))
	holdersArr := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for index,item := range data {

		holders := make([]string,0,len(item))
		if index == 0{
			for k,v := range item{
				keys = append(keys, k)
				holders = append(holders,"?")
				values = append(values, v)
			}
		}else{
			for _,field := range keys{
				holders = append(holders,"?")
				if v, ok := item[field]; ok {
					values = append(values, v)
				}else {
					values = append(values, "")
				}
			}
		}
		holderString := strings.Join(holders,",")
		holdersArr = append(holdersArr,"("+holderString+")")
	}
	return keys,holdersArr,values
}


func (db *zdb) showDebug(sql string, bings []interface{}){
	fmt.Println("[sql]"+sql+"  [binds]"+fmt.Sprintf("%v",bings))
}

//执行查询语句
func (db *zdb) Query(sqlStr string,binds []interface{}) (*[]map[string]string,error) {
	ret := make([]map[string]string,0)
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return &ret,err
	}
	rows, err := stmt.Query(binds...)

	if err != nil{
		return  &ret,err
	}
	columns,err := rows.Columns()
	if err != nil{
		return  &ret,err
	}
	values := make([]sql.RawBytes,len(columns))
	scanArgs := make([]interface{},len(values))
	for i := range values{
		scanArgs[i] = &values[i]
	}
	for rows.Next(){
		err = rows.Scan(scanArgs...)
		if err != nil{
			return  &ret,err
		}
		var value string
		vmap := make(map[string]string,len(scanArgs))
		for i,col := range values{
			if col == nil{
				value = ""
			}else{
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret,vmap)
	}
	return &ret,nil
}

//执行查询语句
func (db *zdb) Select() (*[]map[string]string,error) {
	ret := make([]map[string]string,0)
	sqlStr, binds := db.BuildSelectSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return &ret,err
	}
	rows, err := stmt.Query(binds...)
	if err != nil{
		return  &ret,err
	}
	columns,err := rows.Columns()
	if err != nil{
		return  &ret,err
	}
	values := make([]sql.RawBytes,len(columns))
	scanArgs := make([]interface{},len(values))
	for i := range values{
		scanArgs[i] = &values[i]
	}
	for rows.Next(){
		err = rows.Scan(scanArgs...)
		if err != nil{
			return  &ret,err
		}
		var value string
		vmap := make(map[string]string,len(scanArgs))
		for i,col := range values{
			if col == nil{
				value = ""
			}else{
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret,vmap)
	}
	return &ret,nil
}


//执行查询find语句
func (db *zdb) Find() (*map[string]string,error) {
	ret := map[string]string{}
	sqlStr, binds := db.Limit(0,1).BuildSelectSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return &ret,err
	}
	rows, err := stmt.Query(binds...)
	if err != nil{
		return  &ret,err
	}
	columns,err := rows.Columns()
	if err != nil{
		return  &ret,err
	}
	values := make([]sql.RawBytes,len(columns))
	scanArgs := make([]interface{},len(values))
	for i := range values{
		scanArgs[i] = &values[i]
	}
	vmap := make(map[string]string,len(scanArgs))
	for rows.Next(){
		err = rows.Scan(scanArgs...)
		if err != nil{
			return  &ret,err
		}
		var value string
		for i,col := range values{
			if col == nil{
				value = ""
			}else{
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = vmap
		break
	}
	return &ret,nil

}

//执行查询Value语句
func (db *zdb) Value(field string) (string,error) {
	var ret string
	sqlStr, binds := db.Limit(0,1).BuildSelectSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return ret,err
	}
	rows, err := stmt.Query(binds...)
	if err != nil{
		return  ret,err
	}
	columns,err := rows.Columns()
	if err != nil{
		return  ret,err
	}
	values := make([]sql.RawBytes,len(columns))
	scanArgs := make([]interface{},len(values))
	for i := range values{
		scanArgs[i] = &values[i]
	}
	for rows.Next(){
		err = rows.Scan(scanArgs...)
		if err != nil{
			return  ret,err
		}
		var value string
		for i,col := range values{
			if col == nil{
				value = ""
			}else{
				value = string(col)
			}
			if columns[i] == field{
				ret = value
				break
			}
		}
		break
	}
	return ret,nil
}


//执行查询count语句
func (db *zdb) Count() (int64,error) {
	var zusuxCount int64
	sqlStr, binds := db.BuildCountSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return zusuxCount,err
	}
	row := stmt.QueryRow(binds...)
	if err := row.Scan(&zusuxCount); err !=nil{
		return 0,err
	}
	return zusuxCount,nil
}

//执行查询max语句
func (db *zdb) Max(field string) (float64,error) {
	var zusuxMax float64
	sqlStr, binds := db.Field([]string{field}).BuildMaxSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return zusuxMax,err
	}
	row := stmt.QueryRow(binds...)
	if err := row.Scan(&zusuxMax); err !=nil{
		return 0,err
	}
	return zusuxMax,nil
}

//执行查询min语句
func (db *zdb) Min(field string) (float64,error) {
	var zusuxMin float64
	sqlStr, binds := db.Field([]string{field}).BuildMinSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return zusuxMin,err
	}
	row := stmt.QueryRow(binds...)
	if err := row.Scan(&zusuxMin); err !=nil{
		return 0,err
	}
	return zusuxMin,nil
}


//执行查询sum语句
func (db *zdb) Sum(field string) (float64,error) {
	var zusuxSum float64
	sqlStr, binds := db.Field([]string{field}).BuildSumSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return zusuxSum,err
	}
	row := stmt.QueryRow(binds...)
	if err := row.Scan(&zusuxSum); err !=nil{
		return 0,err
	}
	return zusuxSum,nil
}

//执行查询sum语句
func (db *zdb) Avg(field string) (float64,error) {
	var zusuxAvg float64
	sqlStr, binds := db.Field([]string{field}).BuildAvgSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return zusuxAvg,err
	}
	row := stmt.QueryRow(binds...)
	if err := row.Scan(&zusuxAvg); err !=nil{
		return 0,err
	}
	return zusuxAvg,nil
}


func (db *zdb) BuildSelectSql() (string,[]interface{})  {
	var binds []interface{}
	//force index
	var forceStr string
	if len(db.Build.Force_) >0{
		forceStr = " FORCE INDEX("+db.Build.Force_+")"
	}else{
		forceStr = ""
	}
	// where and 语句  binds
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereArr = append(whereArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereArr = append(whereArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句  //binds
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ " ( "+ strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+strings.Join(whereOrArr," OR ") +" )"
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}
	// group
	var groupStr string
	if len(db.Build.Group_) > 0{
		groupStr = " GROUP BY "+ strings.Join(db.Build.Group_," , ")
	}else{
		groupStr = ""
	}
	// having and 语句 bings
	havingArr := make([]string,0,len(db.Build.Having_))
	for kk,vv := range db.Build.Having_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			havingArr = append(havingArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			havingArr = append(havingArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句 binds
	havingOrArr := make([]string,0,len(db.Build.HavingOr_))
	for kk,vv := range db.Build.HavingOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			havingOrArr = append(havingOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			havingOrArr = append(havingOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var havingStr string
	if len(havingArr) > 0 &&  len(havingOrArr) >0{
		havingStr = " HAVING "+ " ( " + strings.Join(havingArr," AND ")+ " ) " + db.Build.HavingCond_ + " ( " + strings.Join(havingOrArr," OR ") + ")"
	}else if len(havingArr) > 0{
		havingStr = " HAVING "+ strings.Join(havingArr," AND ")
	}else if len(havingOrArr) > 0{
		havingStr = " HAVING "+  strings.Join(havingOrArr," OR ")
	}else{
		havingStr = ""
	}
	// union
	//order str
	var orderStr string
	if len(db.Build.Order_)>0{
		orderStr = " ORDER BY "+ strings.Join(db.Build.Order_ , " , ")
	}else{
		orderStr = ""
	}
	//limit str  binds
	var limitStr string
	if db.Build.Offset_ >0 && db.Build.Rows_ >0 {
		limitStr = " LIMIT ?,?"
		binds = append(binds,db.Build.Offset_)
		binds = append(binds,db.Build.Rows_)
	}else if db.Build.Rows_>0 {
		limitStr = " LIMIT ?"
		binds = append(binds,db.Build.Rows_)
	}else{
		limitStr = ""
	}
	//lock
	var lockStr string
	if db.Build.Lock_ {
		lockStr = " FOR UPDATE"
	}else{
		lockStr = ""
	}
	//SELECT %DISTINCT% %FIELD% FROM %TABLE% %FORCE% %JOIN% %WHERE% %GROUP% %HAVING% %UNION% %ORDER% %LIMIT% %LOCK%
	replacer := strings.NewReplacer(
		"%DISTINCT%",db.Build.Distinct_,
		"%FIELD%",db.Build.Field_,
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%FORCE%",forceStr,
		"%JOIN%", " "+strings.Join(db.Build.Join_," "),
		"%WHERE%", whereStr, //binds
		"%GROUP%",groupStr,
		"%HAVING%",havingStr, //binds
		"%UNION%",strings.Join(db.Build.Union_," "),
		"%ORDER%",orderStr,
		"%LIMIT%",limitStr, //binds
		"%LOCK%",lockStr,
	)
	sql  := replacer.Replace(db.Tmp.SelectSql)
	return sql,binds
}

//总数
func (db *zdb) BuildCountSql() (string,[]interface{})  {
	var binds []interface{}
	// where and 语句  binds
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereArr = append(whereArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereArr = append(whereArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句  //binds
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ " ( "+ strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+strings.Join(whereOrArr," OR ") +" )"
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}
	//SELECT COUNT(*) as zusux_count FROM %TABLE% %WHERE%
	replacer := strings.NewReplacer(
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%WHERE%", whereStr, //binds
	)
	sql  := replacer.Replace(db.Tmp.CountSql)
	return sql,binds
}

//最大值
func (db *zdb) BuildMaxSql() (string,[]interface{})  {
	var binds []interface{}
	// where and 语句  binds
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereArr = append(whereArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereArr = append(whereArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句  //binds
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ " ( "+ strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+strings.Join(whereOrArr," OR ") +" )"
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}
	//SELECT MAX(%FIELD%) as zusux_max FROM %TABLE% %WHERE%
	replacer := strings.NewReplacer(
		"%FIELD%",db.Build.Field_,
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%WHERE%", whereStr, //binds
	)
	sql  := replacer.Replace(db.Tmp.MaxSql)
	return sql,binds
}
//最小值
func (db *zdb) BuildMinSql() (string,[]interface{})  {
	var binds []interface{}
	// where and 语句  binds
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereArr = append(whereArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereArr = append(whereArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句  //binds
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ " ( "+ strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+strings.Join(whereOrArr," OR ") +" )"
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}
	//SELECT MIN(%FIELD%) as zusux_min FROM %TABLE% %WHERE%
	replacer := strings.NewReplacer(
		"%FIELD%",db.Build.Field_,
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%WHERE%", whereStr, //binds
	)
	sql  := replacer.Replace(db.Tmp.MinSql)
	return sql,binds
}

//求平均值
func (db *zdb) BuildAvgSql() (string,[]interface{})  {
	var binds []interface{}
	// where and 语句  binds
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereArr = append(whereArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereArr = append(whereArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句  //binds
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ " ( "+ strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+strings.Join(whereOrArr," OR ") +" )"
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}
	//SELECT AVG(%FIELD%) as zusux_avg FROM %TABLE% %WHERE%
	replacer := strings.NewReplacer(
		"%FIELD%",db.Build.Field_,
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%WHERE%", whereStr, //binds
	)
	sql  := replacer.Replace(db.Tmp.AvgSql)
	return sql,binds
}

//求和
func (db *zdb) BuildSumSql() (string,[]interface{})  {
	var binds []interface{}
	// where and 语句  binds
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereArr = append(whereArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereArr = append(whereArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句  //binds
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ " ( "+ strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+strings.Join(whereOrArr," OR ") +" )"
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}
	//SELECT SUM(%FIELD%) as zusux_sum FROM %TABLE% %WHERE%
	replacer := strings.NewReplacer(
		"%FIELD%",db.Build.Field_,
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%WHERE%", whereStr, //binds
	)
	sql  := replacer.Replace(db.Tmp.SumSql)
	return sql,binds
}

//执行语句
func (db *zdb) Execute(sqlStr string,binds []interface{}) (int64,error) {
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return 0,err
	}
	res, err := stmt.Exec(binds...)
	if err != nil{
		return  0,err
	}
	affectRows ,err := res.RowsAffected()
	return  affectRows ,err
}

//执行更新语句
func (db *zdb) Update(data map[string]interface{}) (int64,error) {

	sqlStr, binds := db.BuildUpdateSql(data)
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return 0,err
	}
	res, err := stmt.Exec(binds...)
	if err != nil{
		return  0,err
	}

	affectRows ,err := res.RowsAffected()
	return  affectRows ,err
}

func (db *zdb) BuildUpdateSql(data map[string]interface{}) (string,[]interface{})  {
	//fmt.Println(data)
	var binds []interface{}
	//set 语句
	sets := make([]string,0,len(data))
	for k,v := range data{
		getType := reflect.TypeOf(v)
		kind := getType.Kind()
		switch kind {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
				sets = append(sets,fmt.Sprintf("%s = ?",k))
				binds = append(binds,v)
			case reflect.Array, reflect.Slice :
				sets = append(sets,fmt.Sprintf("%s = ?",k))
				b, e := json.Marshal(v)
				if e != nil{
					panic(e.Error())
				}
				binds = append(binds,b)
		}
	}
	// where and 语句
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			  whereArr = append(whereArr,kk)
			  binds = append(binds,vv)
		  case  reflect.Slice, reflect.Array :
			  whereArr = append(whereArr,kk)
			  arr := reflect.ValueOf(vv)
			  for i := 0; i < arr.Len(); i++ {
				  v := fmt.Sprintf("%v", arr.Index(i))
				  binds = append(binds, v)
			  }
		}
	}
	// where or 语句
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ "( " + strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+ strings.Join(whereOrArr," OR ") +" ) "
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}

	//order str
	var orderStr string
	if len(db.Build.Order_)>0{
		orderStr = " ORDER BY "+ strings.Join(db.Build.Order_ , ",")
	}else{
		orderStr = " "
	}

	//limit str
	var limitStr string
	if db.Build.Offset_ >0 && db.Build.Rows_ >0 {
		limitStr = " LIMIT ?,?"
		binds = append(binds,db.Build.Offset_)
		binds = append(binds,db.Build.Rows_)
	}else if db.Build.Rows_>0 {
		limitStr = " LIMIT ?"
		binds = append(binds,db.Build.Rows_)
	}else{
		limitStr = " "
	}

	replacer := strings.NewReplacer(
		"%TABLE%",db.Build.Table_+ db.Build.Alias_,
		"%SET%",  strings.Join(sets,","), //binds
		"%JOIN%",strings.Join(db.Build.Join_," "),
		"%WHERE%", whereStr, //binds
		"%ORDER%",orderStr,
		"%LIMIT%",limitStr,  //binds
	)
	sql  := replacer.Replace(db.Tmp.UpdateSql)

	return sql,binds
}


//执行删除语句
func (db *zdb) Delete() (int64,error) {

	sqlStr, binds := db.BuildDeleteSql()
	if db.Build.Debug_ {
		db.showDebug(sqlStr,binds)
	}
	db.Build.Reset()
	stmt, err := db.Conn.SqlDb.Prepare(sqlStr)
	if err != nil {
		return 0,err
	}
	res, err := stmt.Exec(binds...)
	if err != nil{
		return  0,err
	}

	affectRows ,err := res.RowsAffected()
	return  affectRows ,err
}
// 构建删除语句
func (db *zdb) BuildDeleteSql() (string,[]interface{})  {
	//fmt.Println(data)
	var binds []interface{}
	// where and 语句
	whereArr := make([]string,0,len(db.Build.Where_))
	for kk,vv := range db.Build.Where_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereArr = append(whereArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereArr = append(whereArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	// where or 语句
	whereOrArr := make([]string,0,len(db.Build.WhereOr_))
	for kk,vv := range db.Build.WhereOr_ {
		getType := reflect.TypeOf(vv)
		switch getType.Kind() {
		case reflect.String, reflect.Uint,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Float32,reflect.Float64 :
			whereOrArr = append(whereOrArr,kk)
			binds = append(binds,vv)
		case  reflect.Slice, reflect.Array :
			whereOrArr = append(whereOrArr,kk)
			arr := reflect.ValueOf(vv)
			for i := 0; i < arr.Len(); i++ {
				v := fmt.Sprintf("%v", arr.Index(i))
				binds = append(binds, v)
			}
		}
	}
	var whereStr string
	if len(whereArr) > 0 &&  len(whereOrArr) >0{
		whereStr = " WHERE "+ " ( "+ strings.Join(whereArr," AND ") +" ) " + db.Build.WhereCond_ + " ( "+ strings.Join(whereOrArr," OR ") +" )"
	}else if len(whereArr) > 0{
		whereStr = " WHERE "+ strings.Join(whereArr," AND ")
	}else if len(whereOrArr) > 0{
		whereStr = " WHERE "+  strings.Join(whereOrArr," OR ")
	}else{
		whereStr = ""
	}

	//order str
	var orderStr string
	if len(db.Build.Order_)>0{
		orderStr = " ORDER BY "+ strings.Join(db.Build.Order_ , ",")
	}else{
		orderStr = ""
	}
	//limit str
	var limitStr string
	if  db.Build.Rows_ >0 {
		limitStr = " LIMIT  ?"
		binds = append(binds,db.Build.Rows_)
	}else{
		limitStr = ""
	}
	/// %TABLE% %USING% %JOIN% %WHERE% %ORDER% %LIMIT% %LOCK%
	replacer := strings.NewReplacer(
		"%TABLE%",db.Build.Table_ + db.Build.Alias_,
		"%JOIN%",strings.Join(db.Build.Join_," "),
		"%WHERE%", whereStr,
		"%ORDER%",orderStr,
		"%LIMIT%",limitStr,
	)
	sql  := replacer.Replace(db.Tmp.DeleteSql)

	return sql,binds
}

func (db *zdb) GetDb() *sql.DB {
	return db.Conn.SqlDb
}

/*
事务操作
func transaction {
	tx, err := db.GetDb.Begin()
    if err !=nil{
		fmt.Println(err.Error())
		return
	}
    stmt, err1 := tx.Prepare("INSERT INTO userinfo (username, departname, created) VALUES (?, ?, ?)")
    if err1 != nil{
		fmt.Println(err.Error())
		return
	}
    _, err2 := stmt.Exec("test", "测试", "2016-06-20")
    if err2 != nil{
		fmt.Println(err.Error())
		return
	}
    //err3 := tx.Commit()
    err3 := tx.Rollback()
    if err3 != nil{
		fmt.Println(err.Error())
		return
	}
}
 */