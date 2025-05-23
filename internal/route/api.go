package route

import (
	address "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/address/handler/rest"
	banner "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/banner/handler/rest"
	cart "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/cart/handler/rest"
	city "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/city/handler/rest"
	dashboard "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/dashboard/handler/rest"
	expenses "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/expenses/handler/rest"
	application "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/list_application/handler/rest"
	notification "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/notification/handler/rest"
	order "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/order/handler/rest"
	product "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product/handler/rest"
	productCategory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_category/handler/rest"
	productSubCategory "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/product_sub_category/handler/rest"
	province "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/province/handler/rest"
	role "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/role/handler/rest"
	shipping "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/shipping/handler/rest"
	subDistrict "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/sub_district/handler/rest"
	user "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/user/handler/rest"
	voucher "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/internal/module/voucher/handler/rest"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-user-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	var (
		userAPI       = app.Group("/api/users")
		superadminAPI = app.Group("/api/superadmin")
		publicAPI     = app.Group("/api")
	)

	user.NewUserHandler().UserRoute(userAPI, superadminAPI)
	role.NewRoleHandler().RoleRoute(userAPI, superadminAPI)
	application.NewApplicationHandler().ApplicationRoute(superadminAPI)
	productCategory.NewProductCategoryHandler().ProductCategoryRoute(publicAPI, superadminAPI)
	productSubCategory.NewProductSubCategoryHandler().ProductSubCategoryRoute(publicAPI, superadminAPI)
	product.NewProductHandler().ProductRoute(publicAPI, superadminAPI)
	voucher.NewVoucherHandler().VoucherRoute(superadminAPI, userAPI)
	banner.NewBannerHandler().BannerRoute(publicAPI, superadminAPI)
	province.NewProvinceHandler().ProvinceRoute(publicAPI)
	city.NewCityHandler().CityRoute(publicAPI)
	subDistrict.NewSubDistrict().SubDistrictRoute(publicAPI)
	address.NewAddressHandler().AddressRoute(userAPI)
	cart.NewCartHandler().CartRoute(userAPI)
	shipping.NewShippingHandler().ShippingRoute(userAPI)
	order.NewOrderHandler().OrderRoute(userAPI, superadminAPI)
	notification.NewNotificationHandler().NotificationRoute(superadminAPI)
	expenses.NewExpensesHandler().ExpensesRoute(superadminAPI)
	dashboard.NewDashboardHandler().DashboardRoute(superadminAPI)

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Info().
			Str("url", c.OriginalURL()).
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found"))
	})
}
