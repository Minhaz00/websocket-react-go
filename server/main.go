// main.go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type TaskRequest struct {
    TaskID string `json:"task_id"`
}

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    e.GET("/get-task", handleWebSocket)

    log.Println("Server started at :8080")
    e.Logger.Fatal(e.Start(":8080"))
}

func handleWebSocket(c echo.Context) error {
    ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        log.Println(err)
        return err
    }
    defer ws.Close()

    for {
        _, msg, err := ws.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }
        log.Printf("Received: %s", msg)

        var taskRequest TaskRequest
        if err := json.Unmarshal(msg, &taskRequest); err != nil {
            log.Println(err)
            break
        }

        // Simulating task retrieval
        response := map[string]string{
            "task_id": taskRequest.TaskID,
            "details": "Task details for ID: " + taskRequest.TaskID,
        }
        responseJSON, err := json.Marshal(response)
        if err != nil {
            log.Println(err)
            break
        }

        if err := ws.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
            log.Println(err)
            break
        }
    }

    return nil
}
