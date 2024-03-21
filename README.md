## 使用方式

### install
`go get github.com/hu-1996/gormx`

### Init
```go
DB, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{
    SkipDefaultTransaction: true,
    PrepareStmt:            true,
    Logger:                 logger.Default.LogMode(logger.Info),
})
if err != nil {
    os.Exit(1)
}

// Init
gormx.Init(DB)
```

### struct
```go
package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

// User
type User struct {
    gorm.Model
    Name  string `gorm:"column:name;not null;index:idx_name,unique"`
    Age   int    `gorm:"column:age;not null"`
    Roles Roles  `gorm:"column:roles;type:text;"`
}

type Roles []string

func (i *Roles) Scan(src interface{}) error {
    return json.Unmarshal(src.([]byte), i)
}

func (i Roles) Value() (driver.Value, error) {
    return json.Marshal(i)
}

type UserVO struct {
    Name  string `gorm:"column:name;not null;index:idx_name,unique"`
    Age   int    `gorm:"column:age;not null"`
    Roles Roles  `gorm:"column:roles;type:text;"`
}

// implemented ConvertInterface
func (u *User) Convert() interface{} {
    return &UserVO{
    Name:  u.Name,
    Age:   u.Age,
    Roles: u.Roles,
    }
}

```

### SelectById
```go
user, err := gormx.SelectById[mysql.User](1)
if err != nil {
    panic(err)
}
fmt.Printf("user: %+v", user)

// SELECT * FROM `users` WHERE `users`.`id` = 1 AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` DESC LIMIT 1
```

### SelectConvertById
```go
userVo, err := gormx.SelectConvertById[mysql.User, mysql.UserVO](1)
if err != nil {
    panic(err)
}
fmt.Printf("user: %+v", userVo)

// SELECT * FROM `users` WHERE `users`.`id` = 1 AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` DESC LIMIT 1
```

### SelectByIds
```go
users, err := gormx.SelectByIds[mysql.User]([]uint{1, 2})
if err != nil {
    panic(err)
}
for i := range users {
    fmt.Printf("user: %+v \n ", users[i])
}

// SELECT * FROM `users` WHERE `users`.`id` IN (1,2) AND `users`.`deleted_at` IS NULL
```

### SelectConvertByIds
```go
userVos, err := gormx.SelectConvertByIds[mysql.User, mysql.UserVO]([]uint{1, 2})
if err != nil {
    panic(err)
}
for i := range userVos {
    fmt.Printf("user: %+v \n ", userVos[i])
}

// SELECT * FROM `users` WHERE `users`.`id` IN (1,2) AND `users`.`deleted_at` IS NULL
```

### SelectOne
```go
user, err := gormx.SelectOne[mysql.User]("name = ?", "test1")
if err != nil {
    panic(err)
}
fmt.Printf("user: %+v ", user)

// SELECT * FROM `users` WHERE name = 'test1' AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` DESC LIMIT 1
```

### SelectOneConvert
```go
userVo, err := gormx.SelectOneConvert[mysql.User, mysql.UserVO]("name = ?", "test1")
if err != nil {
    panic(err)
}
fmt.Printf("user: %+v", userVo)

// SELECT * FROM `users` WHERE name = 'test1' AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` DESC LIMIT 1
```

### SelectList
```go
records, err := gormx.SelectList[mysql.User](nil, "name like ? AND age = ?", "%test%", 28)
if err != nil {
    panic(err)
}
for i := range records {
    fmt.Printf("user: %+v \n ", records[i])
}

// SELECT * FROM `users` WHERE (name like '%test%' AND age = 28) AND `users`.`deleted_at` IS NULL
```

### SelectListConvert
```go
records, err := gormx.SelectListConvert[mysql.User, mysql.UserVO](nil, "name like ?", "%test%")
if err != nil {
    panic(err)
}
for i := range records {
    fmt.Printf("user: %+v \n ", records[i])
}

// SELECT * FROM `users` WHERE (name like '%test%' AND age = 28) AND `users`.`deleted_at` IS NULL
```

### SelectPage
```go
records, total, err := gormx.SelectPage[mysql.User](1, 10, "id desc", "name like ?", "%test%")
if err != nil {
    panic(err)
}
fmt.Printf("total: %d \n", total)
for i := range records {
    fmt.Printf("user: %+v \n ", records[i])
}

// SELECT * FROM `users` WHERE name like '%test%' AND `users`.`deleted_at` IS NULL ORDER BY id desc LIMIT 10
```

### SelectPageConvert
```go
records, total, err := gormx.SelectPageConvert[mysql.User, mysql.UserVO](1, 10, "id desc", "name like ?", "%test%")
if err != nil {
    panic(err)
}
fmt.Printf("total: %d \n", total)
for i := range records {
    fmt.Printf("user: %+v \n ", records[i])
}
// SELECT * FROM `users` WHERE name like '%test%' AND `users`.`deleted_at` IS NULL ORDER BY id desc LIMIT 10
```

### Insert
```go
user := mysql.User{
    Name:  "test",
    Age:   18,
    Roles: mysql.Roles{"admin", "user"},
}
_, err := gormx.Insert(&user)
if err != nil {
    panic(err)
}

// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`roles`) VALUES ('2024-03-21 18:13:15.635','2024-03-21 18:13:15.635',NULL,'test',18,'["admin","user"]')
```

### InsertBatches
```go
user1 := mysql.User{
    Name:  "test1",
    Age:   18,
    Roles: mysql.Roles{"admin", "user"},
}
user2 := mysql.User{
    Name:  "test2",
    Age:   18,
    Roles: mysql.Roles{"admin", "user"},
}
_, err := gormx.InsertBatches([]*mysql.User{&user1, &user2})
if err != nil {
    panic(err)
}

// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`roles`) VALUES ('2024-03-21 18:13:35.086','2024-03-21 18:13:35.086',NULL,'test1',18,'["admin","user"]'),('2024-03-21 18:13:35.086','2024-03-21 18:13:35.086',NULL,'test2',18,'["admin","user"]')
```

### Update
```go
user := mysql.User{
    Name:  "test4",
    Age:   38,
    Roles: mysql.Roles{"admin", "user", "guest"},
}
user.ID = 1
user.CreatedAt = time.Now()
_, err := gormx.Update(&user)
if err != nil {
    panic(err)
}

// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`roles`,`id`) VALUES ('2024-03-21 18:13:52.41','2024-03-21 18:13:52.41',NULL,'test4',38,'["admin","user","guest"]',1) ON DUPLICATE KEY UPDATE `updated_at`='2024-03-21 18:13:52.41',`deleted_at`=VALUES(`deleted_at`),`name`=VALUES(`name`),`age`=VALUES(`age`),`roles`=VALUES(`roles`)
```

### UpdateBatches
```go
user1 := mysql.User{
    Name:  "test1",
    Age:   58,
    Roles: mysql.Roles{"admin", "user"},
}
user1.ID = 2
user1.CreatedAt = time.Now()
user2 := mysql.User{
    Name:  "test2",
    Age:   10,
    Roles: mysql.Roles{"admin", "user"},
}
user2.ID = 3
user2.CreatedAt = time.Now()
_, err := gormx.UpdateBatches([]*mysql.User{&user1, &user2})
if err != nil {
    panic(err)
}

// INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`name`,`age`,`roles`,`id`) VALUES ('2024-03-21 18:14:04.676','2024-03-21 18:14:04.676',NULL,'test1',58,'["admin","user"]',2),('2024-03-21 18:14:04.676','2024-03-21 18:14:04.676',NULL,'test2',10,'["admin","user"]',3) ON DUPLICATE KEY UPDATE `updated_at`='2024-03-21 18:14:04.676',`deleted_at`=VALUES(`deleted_at`),`name`=VALUES(`name`),`age`=VALUES(`age`),`roles`=VALUES(`roles`)
```

### Updates
```go
user := mysql.User{
    Age: 100,
}
_, err := gormx.Updates(&user, "name = ?", "test1")
if err != nil {
    panic(err)
}

// UPDATE `users` SET `updated_at`='2024-03-21 18:14:24.633',`age`=100 WHERE name = 'test1' AND `users`.`deleted_at` IS NULL
```

### DeleteById
```go
_, err := gormx.DeleteById[mysql.User](1)
if err != nil {
    panic(err)
}

// UPDATE `users` SET `deleted_at`='2024-03-21 18:14:55.965' WHERE `users`.`id` = 1 AND `users`.`deleted_at` IS NULL
```

### DeleteByIds
```go
_, err := gormx.DeleteByIds[mysql.User]([]uint{2, 3})
if err != nil {
    panic(err)
}

// UPDATE `users` SET `deleted_at`='2024-03-21 18:15:44.54' WHERE `users`.`id` IN (2,3) AND `users`.`deleted_at` IS NULL
```

### Delete
```go
ra, err := gormx.Delete[mysql.User]("name like ?", "%test%")
if err != nil {
    panic(err)
}
fmt.Printf("delete rows: %d", ra)

// UPDATE `users` SET `deleted_at`='2024-03-21 18:16:04.765' WHERE name like '%test%' AND `users`.`deleted_at` IS NULL
// delete rows: 3
```