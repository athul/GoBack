package main

import (
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
        "github.com/jinzhu/gorm"
        _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func init() {
        //open a db connection
        db, err := gorm.Open("sqlite3", "test.db")
        if err != nil {
                panic("Could not open Database")
        }
        defer db.Close()

        //Migrate the schema
        db.AutoMigrate(&todoModel{})
}
func main() {
        router := gin.Default()

        v1 := router.Group("/api/v1/todos")
        {
                v1.POST("/", createTodo)
                v1.GET("/", fetchAllTodo)
                v1.GET("/:id", fetchSingleTodo)
                v1.PUT("/:id", updateTodo)
                v1.DELETE("/:id", deleteTodo)
        }
        router.Run()
}

type todoModel struct {
        gorm.Model
        ID        uint   `json:"id"`
        Title     string `json:"title"`
        Desc      string `json:"description"`
        Completed int    `json:"completed"`
}
type transformedTodo struct {
        ID        uint   `json:"id"`
        Title     string `json:"title"`
        Completed bool   `json:"completed"`
}

func createTodo(c *gin.Context) {
        completed, _ := strconv.Atoi(c.PostForm("completed"))
        todo := todoModel{Title: c.PostForm("title"), Completed: completed}
        db.Save(&todo)
        c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}
func fetchAllTodo(c *gin.Context) {
        var todos []todoModel
        var _todos []transformedTodo

        db.Find(&todos)

        if len(todos) <= 0 {
                c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
                return
        }

        //transforms the todos for building a good response
        for _, item := range todos {
                completed := false
                if item.Completed == 1 {
                        completed = true
                } else {
                        completed = false
                }
                _todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
        }
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

// fetchSingleTodo fetch a single todo
func fetchSingleTodo(c *gin.Context) {
        var todo todoModel
        todoID := c.Param("id")

        db.First(&todo, todoID)

        if todo.ID == 0 {
                c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
                return
        }

        completed := false
        if todo.Completed == 1 {
                completed = true
        } else {
                completed = false
        }

        _todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// updateTodo update a todo
func updateTodo(c *gin.Context) {
        var todo todoModel
        todoID := c.Param("id")

        db.First(&todo, todoID)

        if todo.ID == 0 {
                c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
                return
        }

        db.Model(&todo).Update("title", c.PostForm("title"))
        completed, _ := strconv.Atoi(c.PostForm("completed"))
        db.Model(&todo).Update("completed", completed)
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

// deleteTodo remove a todo
func deleteTodo(c *gin.Context) {
        var todo todoModel
        todoID := c.Param("id")

        db.First(&todo, todoID)

        if todo.ID == 0 {
                c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
                return
        }

        db.Delete(&todo)
        c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}

