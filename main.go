package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg"
	"web_app/routes"
	"web_app/settings"
	"web_app/validators"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// @title web_app
// @version 1.0
// @description go web framework demo
// @termsOfService http://swagger.io/terms/

// @contact.name qiujun@sina.com
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /

func main() {
	// 加载配置文件
	if err := settings.InitConfig(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	// 初始化日志
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	// 初始化MySQL连接
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		zap.L().Error("init mysql connection failed, err:%v", zap.Error(err))
		return
	}
	defer mysql.Close()
	// 初始化Redis连接
	if err := redis.InitRedis(); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		zap.L().Error("init redis connection failed, err:%v", zap.Error(err))
		return
	}
	defer redis.Close()
	// 4 初始化翻译
	if err := validators.InitTrans("zh"); err != nil {
		zap.L().Error("init validator failed, err:%v", zap.Error(err))
		return
	}
	// 注册自定义验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", pkg.ValidateMobile)
		// 为验证器自定义错误返回
		_ = v.RegisterTranslation("mobile", validators.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0}手机号码不合法", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	// 注册路由
	r := routes.RegisterRouter()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	// 启动服务(优雅关机)
	// 这个goroutine是启动服务的goroutine
	go func() {
		srv.ListenAndServe()
	}()

	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	// 调用Server.Shutdown graceful结束
	//timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}
