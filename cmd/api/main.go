package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/platform/config"
	"go-waste-routes/internal/platform/database"
	"go-waste-routes/internal/platform/jwt"
	"go-waste-routes/internal/platform/logger"
	redisplatform "go-waste-routes/internal/platform/redis"
	"go-waste-routes/internal/repository/memory"
	"go-waste-routes/internal/service"
	httptransport "go-waste-routes/internal/transport/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.Log.Level, cfg.Log.Format)
	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	ctx := context.Background()

	db, err := openDatabase(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	redisClient, err := redisplatform.Open(ctx, cfg.Redis)
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()
	if err := database.Migrate(ctx, db, "migrations"); err != nil {
		panic(err)
	}

	userStore := memory.New[domain.User]()
	roleStore := memory.New[domain.Role]()
	permissionStore := memory.New[domain.Permission]()
	customerStore := memory.New[domain.Customer]()
	wasteCategoryStore := memory.New[domain.WasteCategory]()
	containerStore := memory.New[domain.Container]()
	vehicleStore := memory.New[domain.Vehicle]()
	driverStore := memory.New[domain.Driver]()
	planStore := memory.New[domain.CollectionPlan]()
	planRuleStore := memory.New[domain.PlanRule]()
	routeStore := memory.New[domain.Route]()
	routeStopStore := memory.New[domain.RouteStop]()
	taskStore := memory.New[domain.Task]()
	taskStopStore := memory.New[domain.TaskStop]()
	weighingStore := memory.New[domain.WeighingRecord]()
	abnormalStore := memory.New[domain.WeighingAbnormality]()
	invoiceStore := memory.New[domain.Invoice]()
	invoiceItemStore := memory.New[domain.InvoiceItem]()
	paymentStore := memory.New[domain.PaymentRecord]()
	weighbridgeStore := memory.New[domain.Weighbridge]()
	auditLogStore := memory.New[domain.AuditLog]()
	systemConfigStore := memory.New[domain.SystemConfig]()
	holidayStore := memory.New[domain.HolidayConfig]()

	users := service.NewUserService(userStore)
	auth := service.NewAuthService(users, jwtManager, cfg.Security.LoginMaxAttempts, cfg.Security.LockoutDuration)
	customers := service.NewResourceService(customerStore)
	wasteCategories := service.NewResourceService(wasteCategoryStore)
	containers := service.NewResourceService(containerStore)
	vehicles := service.NewResourceService(vehicleStore)
	drivers := service.NewResourceService(driverStore)
	plans := service.NewResourceService(planStore)
	planRules := service.NewResourceService(planRuleStore)
	routes := service.NewResourceService(routeStore)
	routeStops := service.NewResourceService(routeStopStore)
	tasks := service.NewResourceService(taskStore)
	taskStops := service.NewResourceService(taskStopStore)
	weighings := service.NewResourceService(weighingStore)
	abnormalities := service.NewResourceService(abnormalStore)
	invoices := service.NewResourceService(invoiceStore)
	invoiceItems := service.NewResourceService(invoiceItemStore)
	payments := service.NewResourceService(paymentStore)
	weighbridges := service.NewResourceService(weighbridgeStore)
	auditLogs := service.NewResourceService(auditLogStore)
	systemConfigs := service.NewResourceService(systemConfigStore)
	holidays := service.NewResourceService(holidayStore)
	dashboard := service.NewDashboardService(customers, tasks, routes, weighings, abnormalities, invoices)
	exporter := service.NewExportService(service.ExportService{
		Customers:     customers,
		Routes:        routes,
		Tasks:         tasks,
		Invoices:      invoices,
		Weighings:     weighings,
		Abnormalities: abnormalities,
		Payments:      payments,
		AuditLogs:     auditLogs,
		PlanRules:     planRules,
		RouteStops:    routeStops,
		TaskStops:     taskStops,
		InvoiceItems:  invoiceItems,
		Weighbridges:  weighbridges,
		SystemConfigs: systemConfigs,
		Holidays:      holidays,
	})
	reviews := service.NewReviewService(service.NewWorkflowEngine(), service.NewBillingEngine(service.NewConfigService()), service.NewWeighingEngine())

	seedCatalog(roleStore, permissionStore, users, cfg)

	deps := httptransport.Deps{
		Logger:       log,
		Auth:         auth,
		Users:        users,
		Dashboard:    dashboard,
		Exporter:     exporter,
		Reviews:      reviews,
		Customers:    customers,
		WasteCats:    wasteCategories,
		Containers:   containers,
		Vehicles:     vehicles,
		Drivers:      drivers,
		Plans:        plans,
		PlanRules:    planRules,
		Routes:       routes,
		RouteStops:   routeStops,
		Tasks:        tasks,
		TaskStops:    taskStops,
		Weighings:    weighings,
		Abnormalities: abnormalities,
		Invoices:     invoices,
		InvoiceItems: invoiceItems,
		Payments:     payments,
		Weighbridges: weighbridges,
		Roles:        service.NewResourceService(roleStore),
		Permissions:  service.NewResourceService(permissionStore),
		AuditLogs:    auditLogs,
		SystemConfig: systemConfigs,
		Holidays:     holidays,
	}

	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      httptransport.NewRouter(deps),
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
	}

	go func() {
		log.WithField("addr", server.Addr).Info("server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

func seedCatalog(roleStore *memory.Store[domain.Role], permissionStore *memory.Store[domain.Permission], users *service.UserService, cfg *config.Config) {
	for _, permission := range service.DefaultPermissions {
		_, _ = permissionStore.Create(context.Background(), permission)
	}
	_, _ = roleStore.Create(context.Background(), service.DefaultRole)
	_, _ = roleStore.Create(context.Background(), service.AdminRole)
	admin := domain.User{Username: cfg.Admin.Username, DisplayName: "System Admin", Status: domain.UserStatusEnabled, Roles: []domain.Role{service.AdminRole}, Permissions: service.DefaultPermissions}
	_, _ = users.Create(context.Background(), admin, cfg.Admin.Password, service.DefaultPermissions)
}

func openDatabase(ctx context.Context, cfg *config.Config) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		db, err = database.Open(cfg.DB)
		if err == nil {
			return db, nil
		}
		time.Sleep(time.Second)
	}
	return nil, err
}
