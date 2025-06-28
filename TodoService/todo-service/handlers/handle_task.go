package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"todo-service/database"
	"todo-service/metrics"
	"todo-service/model"

	"github.com/gin-gonic/gin"
)

func getTaskID(c *gin.Context) (int, error) {

	idStr := c.Param("id")
	if idStr == "" {
		return 0, fmt.Errorf("缺失任务ID")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("无效任务ID")
	}

	return id, nil
}

func CreateTask(c *gin.Context) {
	startTime := time.Now()

	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		metrics.DBOpCounter.WithLabelValues("create_task", "invalid").Inc()
		return
	}

	// 使用带超时的上下文
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	database.DB.Exec("use testdb")

	// 执行数据库操作
	result, err := database.DB.ExecContext(
		ctx,
		"INSERT INTO todo_tasks (title, description, creator, assignee) VALUES (?, ?, ?, ?)",
		req.Title, req.Description, req.Creator, req.Assignee,
	)
	if err != nil {
		log.Println("创建任务失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		metrics.DBOpCounter.WithLabelValues("create_task", "error").Inc()
		return
	}

	// 获取新任务ID
	id, _ := result.LastInsertId()

	// 查询新创建的任务
	newTask := model.Task{
		ID:          int(id),
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Creator:     req.Creator,
		Assignee:    req.Assignee,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// 记录指标
	duration := time.Since(startTime).Seconds()
	metrics.DBOpDuration.WithLabelValues("create_task").Observe(duration)
	metrics.DBOpCounter.WithLabelValues("create_task", "success").Inc()

	c.JSON(http.StatusCreated, newTask)
}

// 开始任务
func StartTask(c *gin.Context) {
	id, err := getTaskID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	database.DB.Exec("use testdb")

	result, err := database.DB.ExecContext(
		ctx,
		"UPDATE todo_tasks SET status = 'in_progress', started_at = NOW() WHERE id = ? AND status IN ('pending', 'in_progress')",
		id,
	)

	if err != nil {
		log.Printf("开始任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开始任务失败"})
		metrics.DBOpCounter.WithLabelValues("start_task", "error").Inc()
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务无法开始或不存在"})
		metrics.DBOpCounter.WithLabelValues("start_task", "invalid").Inc()
		return
	}

	duration := time.Since(startTime).Seconds()
	metrics.DBOpDuration.WithLabelValues("start_task").Observe(duration)
	metrics.DBOpCounter.WithLabelValues("start_task", "success").Inc()

	c.Status(http.StatusNoContent)
}

// 完成任务
func CompleteTask(c *gin.Context) {
	id, err := getTaskID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	database.DB.Exec("use testdb")

	result, err := database.DB.ExecContext(ctx,
		"UPDATE todo_tasks SET status = 'completed', completed_at = NOW() WHERE id = ? AND status IN ('pending', 'in_progress', 'completed')",
		id,
	)

	if err != nil {
		log.Printf("完成任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "完成任务失败"})
		metrics.DBOpCounter.WithLabelValues("complete_task", "error").Inc()
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务无法完成或不存在"})
		metrics.DBOpCounter.WithLabelValues("complete_task", "invalid").Inc()
		return
	}

	duration := time.Since(startTime).Seconds()
	metrics.DBOpDuration.WithLabelValues("complete_task").Observe(duration)
	metrics.DBOpCounter.WithLabelValues("complete_task", "success").Inc()

	c.Status(http.StatusNoContent)
}

// 删除任务
func DeleteTask(c *gin.Context) {
	id, err := getTaskID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	database.DB.Exec("use testdb")

	result, err := database.DB.ExecContext(ctx,
		"UPDATE todo_tasks SET status = 'deleted', deleted_at = NOW() WHERE id = ? AND status != 'deleted'",
		id,
	)

	if err != nil {
		log.Printf("删除任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除任务失败"})
		metrics.DBOpCounter.WithLabelValues("delete_task", "error").Inc()
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务已删除或不存在"})
		metrics.DBOpCounter.WithLabelValues("delete_task", "invalid").Inc()
		return
	}

	duration := time.Since(startTime).Seconds()
	metrics.DBOpDuration.WithLabelValues("delete_task").Observe(duration)
	metrics.DBOpCounter.WithLabelValues("delete_task", "success").Inc()

	c.Status(http.StatusNoContent)
}

// 获取任务列表
func ListTasks(c *gin.Context) {
	status := c.Query("status")

	startTime := time.Now()
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	database.DB.Exec("use testdb")

	var rows *sql.Rows
	var err error

	if status == "" {
		rows, err = database.DB.QueryContext(ctx,
			"SELECT id, title, description, creator, assignee, status, created_at, started_at, completed_at, deleted_at "+
				"FROM todo_tasks WHERE status != 'deleted' ORDER BY id DESC")
	} else {
		rows, err = database.DB.QueryContext(ctx,
			"SELECT id, title, description, creator, assignee, status, created_at, started_at, completed_at, deleted_at "+
				"FROM todo_tasks WHERE status = ? ORDER BY id DESC", status)
	}

	if err != nil {
		log.Printf("获取任务列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务失败"})
		metrics.DBOpCounter.WithLabelValues("list_tasks", "error").Inc()
		return
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Creator, &task.Assignee,
			&task.Status, &task.CreatedAt, &task.StartedAt, &task.CompletedAt, &task.DeletedAt,
		)
		if err != nil {
			log.Printf("扫描任务错误: %v", err)
			continue
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Printf("遍历任务错误: %v", err)
	}

	duration := time.Since(startTime).Seconds()
	metrics.DBOpDuration.WithLabelValues("list_tasks").Observe(duration)
	metrics.DBOpCounter.WithLabelValues("list_tasks", "success").Inc()

	c.JSON(http.StatusOK, tasks)
}
