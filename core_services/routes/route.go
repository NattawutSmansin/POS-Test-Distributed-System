package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	databases "cor-service/databases"
	_ "cor-service/docs"
	rabbitmq "cor-service/services"

	orderRepo "cor-service/src/order/repositories"
	productRepo "cor-service/src/product/repositories"

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

	connRabbitmq, err := rabbitmq.NewConnectionRabbitMQ()
	if err != nil {
		panic(err)
	}

	r.GET("api/core/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := selectAPIPath(r)
	{

		product := api.Group("/products")
		{
			handler := productRepo.NewRepositoryHandler(databases.DB, connRabbitmq)
			product.POST("", handler.CreateProduct)
			product.PUT("", handler.UpdateProductPrice)
		}

		branch := api.Group("/:branch_id")
		{
			handler := productRepo.NewRepositoryHandler(databases.DB, connRabbitmq)
			branch.GET("products", handler.BranchProduct)
		}

		order := api.Group("/orders")
		{
			handler := orderRepo.NewRepositoryHandler(databases.DB)
			order.GET("", handler.Order)
		}
	}

	fmt.Println("Swagger UI Click ---> http://localhost:3344/api/core/docs/swagger/index.html")
	return r
}

func selectAPIPath(r *gin.Engine) *gin.RouterGroup {
	return r.Group("api/core")
}
