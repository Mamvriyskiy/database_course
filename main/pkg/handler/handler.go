package handler

import (
	"net/http"
	"strings"
	"fmt"
	_ "github.com/santosh/gingo/docs"
	"github.com/Mamvriyskiy/database_course/main/logger"
	"github.com/Mamvriyskiy/database_course/main/pkg/service"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/gin-contrib/cors"
	"time"
)

const signingKey = "jaskljfkdfndnznmckmdkaf3124kfdlsf"

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

// Middleware для извлечения данных из JWT и добавления их в контекст запроса.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверить URL запроса
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			// Если URL неначинается с /api, пропустить проверку JWT
			c.Next()
			return
		}

		// Получить токен из заголовка запроса или из куки
		tokenString := c.GetHeader("Authorization")
		fmt.Println("Token:", tokenString)
		var err error
		if tokenString == "" {
			// Если токен не найден в заголовке, попробуйте из куки
			tokenString, err = c.Cookie("jwt")
			if err != nil {
				logger.Log("Error", "c.Cookie(jwt)", "Error", err, "jwt")
			}
		}

		// Проверить, что токен не пустой
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Empty token"})
			c.Abort()
			return
		}

		// Парсинг токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Здесь нужно вернуть ключ для проверки подписи токена.
			// В реальном приложении, возможно, это будет случайный секретный ключ.
			return []byte(signingKey), nil
		})
		
		fmt.Println(err) 
		// Проверить наличие ошибок при парсинге токена
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "detail:": err.Error()})
			c.Abort()
			return
		}

		fmt.Println("+")
		// Добавить данные из токена в контекст запроса
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("+")
			c.Set("userId", claims["userId"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешить CORS заголовки
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")  // Замените на нужный домен
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Обработка OPTIONS запроса
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Продолжить обработку запроса
		c.Next()

		fmt.Println(c.Request.Method)
	}
}

// @title     Gingo Bookstore API
func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()
	fmt.Println("+")

	// router.Use(CORSMiddleware())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Access-Control-Allow-Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// }))
	
	// router.OPTIONS("/*any", func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "*")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
	// 	c.Status(http.StatusOK)
	// })	

	router.Use(func(ctx *gin.Context) {
        fmt.Println("Requested URL:", ctx.Request.URL.String()) // Логируем URL запроса
		fmt.Println("Request Method:", ctx.Request.Method) 
        ctx.Next() // Продолжаем обработку запроса
    })

	router.Use(AuthMiddleware())

	// router.Static("/css", "./templates/css")
	// router.LoadHTMLGlob("templates/*.html")

	router.Static("/docs", "./docs")

	fmt.Println("+++")
    // Настройка Swagger UI
    router.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
        ginSwagger.URL("/docs/swagger.yaml")))

	// app := router.Group("/app")
	// app.GET("/menu", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "menu.html", nil)
	// })

	// app.GET("/home", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "home.html", nil)
	// })

	// app.GET("/access", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "access.html", nil)
	// })

	// app.GET("/device", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "device.html", nil)
	// })
	// auth.GET("/sign-up", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "registr.html", nil)
	// })

	// auth.GET("/sign-in", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "auth.html", nil)
	// })

	// auth.GET("/reset-password", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "send.html", nil)
	// })

	// auth.GET("/verification", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "checkcode.html", nil)
	// })

	// auth.GET("/password", func(ctx *gin.Context) {
	// 	ctx.HTML(http.StatusOK, "changepswrd.html", nil)
	// })
	
	auth := router.Group("/auth")
	auth.POST("/sign-up", h.SignUp)
	auth.POST("/sign-in", h.signIn)
	auth.PUT("/password", h.changePassword)
	auth.POST("/verification", h.checkCode)
	auth.POST("/code", h.code)

	api := router.Group("/api")

	home := api.Group("/homes")
	home.POST("/", h.createHome)
	home.GET("/", h.listHome)
	home.DELETE("/:homeID", h.deleteHome)
	home.PUT("/:homeID", h.updateHome)
	home.GET("/:homeID", h.infoHome)

	home.POST("/:homeID/accesses", h.addUser)
	home.DELETE("/:homeID/accesses/:accessID", h.deleteUser)
	home.GET("/:homeID/accesses", h.getListUserHome)
	home.PUT("/:homeID/accesses/:accessID", h.updateLevel)
	home.GET("/:homeID/accesses/:accessID", h.getInfoAccess)

	home.POST("/:homeID/devices", h.createDevice)
	home.GET("/:homeID/devices", h.getListDevice)
	home.DELETE("/:homeID/devices/:deviceID", h.deleteDevice)
	home.GET("/:homeID/devices/:deviceID", h.getInfoDevice)

	home.POST("/:homeID/devices/:deviceID/status", h.createDeviceHistory)
	home.GET("/:homeID/devices/:deviceID/history", h.getDeviceHistory)

	logger.Log("Info", "", "Create router", nil)

	return router
}
