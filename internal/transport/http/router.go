package httptransport

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/service"
	handlerpkg "go-waste-routes/internal/transport/http/handler"
	"go-waste-routes/internal/transport/http/middleware"
	"go-waste-routes/internal/transport/http/response"
)

type Deps struct {
	Logger       *logrus.Logger
	Auth         *service.AuthService
	Users        *service.UserService
	Dashboard    *service.DashboardService
	Exporter     *service.ExportService
	Reviews      *service.ReviewService
	Customers    *service.ResourceService[domain.Customer]
	WasteCats    *service.ResourceService[domain.WasteCategory]
	Containers   *service.ResourceService[domain.Container]
	Vehicles     *service.ResourceService[domain.Vehicle]
	Drivers      *service.ResourceService[domain.Driver]
	Plans        *service.ResourceService[domain.CollectionPlan]
	PlanRules    *service.ResourceService[domain.PlanRule]
	Routes       *service.ResourceService[domain.Route]
	RouteStops   *service.ResourceService[domain.RouteStop]
	Tasks        *service.ResourceService[domain.Task]
	TaskStops    *service.ResourceService[domain.TaskStop]
	Weighings    *service.ResourceService[domain.WeighingRecord]
	Abnormalities *service.ResourceService[domain.WeighingAbnormality]
	Invoices     *service.ResourceService[domain.Invoice]
	InvoiceItems *service.ResourceService[domain.InvoiceItem]
	Payments     *service.ResourceService[domain.PaymentRecord]
	Weighbridges *service.ResourceService[domain.Weighbridge]
	Roles        *service.ResourceService[domain.Role]
	Permissions  *service.ResourceService[domain.Permission]
	AuditLogs    *service.ResourceService[domain.AuditLog]
	SystemConfig *service.ResourceService[domain.SystemConfig]
	Holidays     *service.ResourceService[domain.HolidayConfig]
}

func NewRouter(deps Deps) http.Handler {
	router := chi.NewRouter()
	router.Use(requestIDMiddleware)
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger(deps.Logger))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) { response.OK(w, map[string]any{"status": "ok"}, contextx.RequestID(r.Context())) })
	router.Get("/ready", func(w http.ResponseWriter, r *http.Request) { response.OK(w, map[string]any{"status": "ready"}, contextx.RequestID(r.Context())) })
	router.Get("/api/swagger", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "docs/swagger.yaml") })

	api := chi.NewRouter()
	api.Use(middleware.Idempotency(24 * time.Hour))

	authHandler := &handlerpkg.AuthHandler{Auth: deps.Auth}
	userHandler := &handlerpkg.UserHandler{Users: deps.Users, Auth: deps.Auth}

	api.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout)
	})

	protected := chi.NewRouter()
	protected.Use(middleware.Auth(deps.Auth))
	opsHandler := &handlerpkg.OpsHandler{
		Dashboard:     deps.Dashboard,
		Exporter:      deps.Exporter,
		Reviews:       deps.Reviews,
		Abnormalities: deps.Abnormalities,
		Invoices:      deps.Invoices,
		Payments:      deps.Payments,
	}
	protected.Get("/dashboard", opsHandler.DashboardView)
	protected.Get("/exports/kinds", opsHandler.ExportKinds)
	protected.Get("/exports", opsHandler.Export)
	protected.Get("/users/me", userHandler.Me)
	protected.Put("/users/password", userHandler.ChangePassword)
	protected.Route("/users", func(r chi.Router) {
		r.Use(middleware.RBAC("user.manage"))
		r.Get("/", userHandler.List)
		r.Post("/", userHandler.Create)
		r.Get("/{id}", userHandler.Get)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
		r.Put("/{id}/status", userHandler.Status)
		r.Put("/{id}/password/reset", userHandler.ResetPassword)
	})

	registerResource[domain.Role](protected, "/roles", handlerpkg.NewCRUDHandler[domain.Role](deps.Roles), "role.manage")
	registerResource[domain.Permission](protected, "/permissions", handlerpkg.NewCRUDHandler[domain.Permission](deps.Permissions), "user.manage")
	registerResource[domain.Customer](protected, "/customers", handlerpkg.NewCRUDHandler[domain.Customer](deps.Customers), "customer.manage")
	registerResource[domain.WasteCategory](protected, "/waste-categories", handlerpkg.NewCRUDHandler[domain.WasteCategory](deps.WasteCats), "customer.manage")
	registerResource[domain.Container](protected, "/containers", handlerpkg.NewCRUDHandler[domain.Container](deps.Containers), "customer.manage")
	registerResource[domain.Vehicle](protected, "/vehicles", handlerpkg.NewCRUDHandler[domain.Vehicle](deps.Vehicles), "route.manage")
	registerResource[domain.Driver](protected, "/drivers", handlerpkg.NewCRUDHandler[domain.Driver](deps.Drivers), "route.manage")
	registerResource[domain.CollectionPlan](protected, "/plans", handlerpkg.NewCRUDHandler[domain.CollectionPlan](deps.Plans), "plan.manage")
	registerResource[domain.Route](protected, "/routes", handlerpkg.NewCRUDHandler[domain.Route](deps.Routes), "route.manage")
	registerResource[domain.Task](protected, "/tasks", handlerpkg.NewCRUDHandler[domain.Task](deps.Tasks), "task.manage")
	registerResource[domain.WeighingRecord](protected, "/weighing", handlerpkg.NewCRUDHandler[domain.WeighingRecord](deps.Weighings), "weighing.manage")
	registerResource[domain.WeighingAbnormality](protected, "/abnormalities", handlerpkg.NewCRUDHandler[domain.WeighingAbnormality](deps.Abnormalities), "weighing.manage")
	registerResource[domain.Invoice](protected, "/invoices", handlerpkg.NewCRUDHandler[domain.Invoice](deps.Invoices), "invoice.manage")
	registerResource[domain.PlanRule](protected, "/plan-rules", handlerpkg.NewCRUDHandler[domain.PlanRule](deps.PlanRules), "plan.manage")
	registerResource[domain.RouteStop](protected, "/route-stops", handlerpkg.NewCRUDHandler[domain.RouteStop](deps.RouteStops), "route.manage")
	registerResource[domain.TaskStop](protected, "/task-stops", handlerpkg.NewCRUDHandler[domain.TaskStop](deps.TaskStops), "task.manage")
	registerResource[domain.InvoiceItem](protected, "/invoice-items", handlerpkg.NewCRUDHandler[domain.InvoiceItem](deps.InvoiceItems), "invoice.manage")
	registerResource[domain.PaymentRecord](protected, "/payment-records", handlerpkg.NewCRUDHandler[domain.PaymentRecord](deps.Payments), "invoice.manage")
	registerResource[domain.Weighbridge](protected, "/weighbridges", handlerpkg.NewCRUDHandler[domain.Weighbridge](deps.Weighbridges), "weighing.manage")
	registerResource[domain.AuditLog](protected, "/audit-logs", handlerpkg.NewCRUDHandler[domain.AuditLog](deps.AuditLogs), "user.manage")
	registerResource[domain.SystemConfig](protected, "/system-configs", handlerpkg.NewCRUDHandler[domain.SystemConfig](deps.SystemConfig), "user.manage")
	registerResource[domain.HolidayConfig](protected, "/holiday-configs", handlerpkg.NewCRUDHandler[domain.HolidayConfig](deps.Holidays), "user.manage")

	protected.Put("/abnormalities/{id}/review", opsHandler.ReviewAbnormality)
	protected.Put("/invoices/{id}/confirm", opsHandler.ConfirmInvoice)
	protected.Post("/invoices/{id}/payment", opsHandler.RecordPayment)

	router.Mount("/api/v1", api)
	api.Mount("/", protected)
	return router
}

func registerResource[T any](root chi.Router, path string, crud *handlerpkg.CRUDHandler[T], permission string) {
	root.Route(path, func(r chi.Router) {
		r.Use(middleware.RBAC(permission))
		r.Get("/", crud.List)
		r.Post("/", crud.Create)
		r.Get("/{id}", crud.Get)
		r.Put("/{id}", crud.Update)
		r.Delete("/{id}", crud.Delete)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()
		w.Header().Set("X-Request-ID", requestID)
		ctx := contextx.WithRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
