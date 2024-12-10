package main

import (
	"log"
	"os"

	"gitgub.com/Aukawut/ServerOTManagement/auth"
	jwt "gitgub.com/Aukawut/ServerOTManagement/auth/jwt"
	"gitgub.com/Aukawut/ServerOTManagement/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	// Load Env
	err := godotenv.Load()
	if err != nil {
		log.Println("Load .env error!")

	}

	// Enable CORS for all routes
	app.Use(cors.New())

	//<------ Request ----->
	app.Post("/request", handler.RequestOvertime)
	app.Get("/count/request/:code", jwt.DecodeToken, handler.CountRequestByEmpCode)
	app.Get("/count/requests/:year", jwt.DecodeToken, handler.CountRequestByYear)
	app.Put("/request/:requestNo", jwt.DecodeToken, handler.CancelRequestByReqNo)
	app.Post("/rewrite/request", jwt.DecodeToken, handler.RewiteRequestOvertime)
	app.Get("/menu/year", jwt.DecodeToken, handler.GetYearMenu)
	app.Get("/menu/month/:year", jwt.DecodeToken, handler.GetMonthMenu)
	app.Get("/user/requests/:code", jwt.DecodeToken, handler.GetRequestByEmpCode)
	app.Get("/summary/request/last/all", jwt.DecodeToken, handler.SummaryLastRevRequestAll)
	app.Get("/summary/request/lasted/:reqNo", jwt.DecodeToken, handler.SummaryLastRevRequestAllByReqNo)

	app.Get("/pending/request/approver", jwt.DecodeToken, handler.GetApproverPending)
	app.Get("/approve/count/:code", jwt.DecodeToken, handler.CountApproveStatusByCode)
	app.Get("/approve/reqList/:code/:status", jwt.DecodeToken, handler.GetRequestListByCodeAndStatus)

	app.Get("/role", handler.GetAllRole)
	app.Post("/login", auth.LoginDomain)

	// --- Users Route ---
	app.Get("/users/group/:id/:factory", jwt.DecodeToken, handler.GetUserInfoByGroupId)
	app.Get("/users/factory/:factory", jwt.DecodeToken, handler.GetUserByFactory)
	app.Get("/approver/group/:id/:factory", jwt.DecodeToken, handler.GetApproverByGroupId)
	app.Get("/approver", jwt.DecodeToken, handler.GetAllApprover)
	app.Post("/approver", jwt.DecodeTokenAdmin, handler.InsertApprover)
	app.Put("/approver/:id", jwt.DecodeTokenAdmin, handler.UpdateApprover)
	app.Get("/requests/users/:requestNo/:rev", jwt.DecodeToken, handler.GetUserByRequestNoAndRev)

	//---- Group Department -----

	app.Get("/group", jwt.DecodeToken, handler.GetGroupDepartment)
	app.Get("/group/:status", jwt.DecodeToken, handler.GetGroupDepartmentByStatus)

	//<----- Shift -------->
	app.Get("/shift", jwt.DecodeToken, handler.GetAllShift)
	app.Get("/shift/:status", jwt.DecodeToken, handler.GetAllShiftActive)

	// <---- Overtime ---->

	app.Get("/overtime", jwt.DecodeToken, handler.GetOvertimeType)

	// <----- Factory ------>
	app.Get("/factory", jwt.DecodeToken, handler.GetAllFactory)
	app.Get("/factory/:group", jwt.DecodeToken, handler.GetAllFactoryByGroup)

	//<---- Workcell ----->
	app.Get("/workcell/:group", jwt.DecodeToken, handler.GetWorkCellByGroup)

	//<---- Group Workcell ----->
	app.Get("/workgroup", jwt.DecodeToken, handler.GetAllGroupWorkcell)

	//<--- Plan ---->
	app.Get("/plan/main", jwt.DecodeToken, handler.GetAllMainPlan)
	app.Post("/plan/main", jwt.DecodeToken, handler.AddMainPlan)
	app.Put("/plan/main/:id", jwt.DecodeToken, handler.UpdateMainPlan)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "4860" // Default port
	}

	// Application Listen :PORT
	app.Listen(":" + PORT)

}
