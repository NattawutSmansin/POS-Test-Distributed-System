package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	_ "branch-service/docs"
	rabbitmq "branch-service/services"
	

	orderRepo "branch-service/src/order/repositories"
	productRepo "branch-service/src/product/repositories"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	})

	conn, err := rabbitmq.NewConnectionRabbitMQ()
	if err != nil {
		panic(err)
	}

	r.GET("api/branch/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := selectAPIPath(r)
	{
		product := api.Group("/:branch_id")
		{
			handler := productRepo.NewRepositoryHandler()
			product.GET("/products", handler.BranchProduct)
		}

		order := api.Group("/orders")
		{
			handler := orderRepo.NewRepositoryHandler(conn)
			order.POST("", handler.CreateOrder)
		}
	}

	fmt.Println("Swagger UI Click ---> http://localhost:5566/api/branch/docs/swagger/index.html")
	return r
}

func selectAPIPath(r *gin.Engine) *gin.RouterGroup {
	return r.Group("api/branch")
}
