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
	app.Get("/count/requests/:year/:code", jwt.DecodeToken, handler.CountRequestByYear)
	app.Get("/request/count/:status/:code", jwt.DecodeToken, handler.CountRequestByStatusAndCode)
	app.Get("/request/approve/count/:status/:code", jwt.DecodeToken, handler.CountRequestStatusAndApproveByCode)
	app.Put("/request/update/:requestNo/:rev", jwt.DecodeToken, handler.ApproveRequestByNo)
	app.Put("/request/cancel/:requestNo/:rev", jwt.DecodeToken, handler.CancelRequestByReqNo)
	app.Post("/revise/request", jwt.DecodeToken, handler.ReviseRequestOvertime)
	app.Get("/menu/year", jwt.DecodeToken, handler.GetYearMenu)
	app.Get("/menu/month/:year", jwt.DecodeToken, handler.GetMonthMenu)
	app.Get("/user/requests/:code", jwt.DecodeToken, handler.GetRequestByEmpCode)
	app.Get("/summary/request/last/all", jwt.DecodeToken, handler.SummaryLastRevRequestAll)
	app.Get("/summary/request/:reqNo/:rev", jwt.DecodeToken, handler.SummaryLastRevRequestAllByReqNo)
	app.Get("/comment/request/:requestNo/:rev", jwt.DecodeToken, handler.GetApproverCommentByRequestNo)
	app.Get("/lasted/request/:status/:code", jwt.DecodeToken, handler.GetRequestListByStatusApprove)
	app.Get("/details/request/:status/:code", jwt.DecodeToken, handler.GetRequestListByStatusApproveAndCode)
	app.Get("/lasted/user/request/:status/:code", jwt.DecodeToken, handler.GetUserRequestListByStatusPendApprove)
	app.Get("/details/user/request/:status/:code", jwt.DecodeToken, handler.GetUserRequestListByStatusApprove)
	app.Get("/details/old/request/:status/:requestNo/:rev", jwt.DecodeToken, handler.GetDetailOldRequestByStatus)

	app.Get("/pending/request/approver", jwt.DecodeToken, handler.GetApproverPending)
	app.Get("/approve/count/:code", jwt.DecodeToken, handler.CountApproveStatusByCode)
	app.Get("/approve/reqList/:code/:status", jwt.DecodeToken, handler.GetRequestListByCodeAndStatus)

	app.Get("/role", handler.GetAllRole)
	app.Post("/login", auth.LoginDomain)

	// --- Users Route ---

	app.Get("/users/group/:id/:factory", jwt.DecodeToken, handler.GetUserInfoByGroupId)
	app.Get("/users/factory/:factory", jwt.DecodeToken, handler.GetUserByFactory)
	app.Get("/users/ugroup", jwt.DecodeToken, handler.GetUserGroup)
	app.Get("/users/type", jwt.DecodeToken, handler.GetUserType)
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
	app.Post("/factory", jwt.DecodeToken, handler.InsertFactory)
	app.Get("/factory/:group", jwt.DecodeToken, handler.GetAllFactoryByGroup)

	//<---- Workcell ----->
	app.Get("/workcell", jwt.DecodeToken, handler.GetWorkCellByAll)
	app.Post("/workcell", jwt.DecodeToken, handler.InsertWorkcell)
	app.Put("/workcell/:id", jwt.DecodeToken, handler.UpdateWorkcell)
	app.Delete("/workcell/:id", jwt.DecodeToken, handler.DeleteWorkcell)
	app.Get("/workcell/:group", jwt.DecodeToken, handler.GetWorkCellByGroup)
	app.Get("/workcell/factory/:id", jwt.DecodeToken, handler.GetWorkcellByFactory)

	//<---- Group Workcell ----->
	app.Get("/workgroup", jwt.DecodeToken, handler.GetAllGroupWorkcell)

	//<--- Plan ---->
	app.Get("/plan/main", jwt.DecodeToken, handler.GetAllMainPlan)
	app.Post("/plan/main", jwt.DecodeToken, handler.AddMainPlan)
	app.Put("/plan/main/:id", jwt.DecodeToken, handler.UpdateMainPlan)
	app.Get("/plan/workcell/:year/:month/:id", jwt.DecodeToken, handler.GetPlanByWorkcell)
	app.Get("/plan/factory/:year/:month/:id", jwt.DecodeToken, handler.GetPlanByFactory)
	app.Delete("/plan/:id", jwt.DecodeToken, handler.DeletePlan)

	//<---- Plan OB ---->
	app.Get("/ob/plan", jwt.DecodeToken, handler.GetMainPlanOb)
	app.Delete("/ob/plan/:id", jwt.DecodeToken, handler.DeletePlanOB)
	app.Post("/ob/plan", jwt.DecodeToken, handler.AddPlanOB)
	app.Put("/ob/plan/:id", jwt.DecodeToken, handler.UpdatePlanOB)

	// <--- Actual ---->
	app.Post("/actual/overtime", jwt.DecodeToken, handler.SaveActualOvertime)
	app.Get("/actual/overtime", jwt.DecodeToken, handler.GetActualOvertime)
	app.Get("/actual/overtime/:start/:end", jwt.DecodeToken, handler.GetActualByDate)
	app.Get("/actual/summary/compare/plan/:year", jwt.DecodeToken, handler.SummaryActualComparePlan)
	app.Get("/actual/summary/factory/plan/:start/:end/:ugroup", jwt.DecodeToken, handler.SummaryActualByDurationAndFac)
	app.Get("/actual/count/:start/:end/:ugroup", jwt.DecodeToken, handler.GetCountActualOvertime)
	app.Get("/actual/summary/date/:start/:end/:ugroup", jwt.DecodeToken, handler.SummaryActualByDate)
	app.Get("/actual/ot/:start/:end/:ugroup/:fac", jwt.DecodeToken, handler.SummaryActualOvertime)
	app.Get("/actual/factory/:start/:end/:ugroup/:fac", jwt.DecodeToken, handler.SummaryActualOvertimeGroupFac)
	app.Get("/actual/type/:start/:end/:ugroup/:fac", jwt.DecodeToken, handler.SummaryActualOvertimeByType)
	app.Get("/actual/bydate/:start/:end/:ugroup/:fac", jwt.DecodeToken, handler.SummaryActualOvertimeByDate)
	app.Get("/actual/workcell/:requestNo/:rev/:year/:month", jwt.DecodeToken, handler.CalActualByWorkcell)
	app.Get("/actual/cal/:year/:month/:fac", jwt.DecodeToken, handler.CalActualByFactory)
	app.Get("/actual/all/workgroup/:start/:end", jwt.DecodeToken, handler.GetActualCompareWorkgroup)
	app.Get("/actual/group/workcell/:start/:end", jwt.DecodeToken, handler.GetActualCompareGroupWorkCell)
	app.Get("/actual/group/workgroup/:start/:end", jwt.DecodeToken, handler.GetActualCompareGroupWorkGroup)
	app.Delete("/actual/:id", jwt.DecodeToken, handler.DeleteActualById)

	//----  Permission ----
	app.Post("/permission/user", jwt.DecodeTokenAdmin, handler.InsertUserPermission)
	app.Put("/permission/user/:id", jwt.DecodeTokenAdmin, handler.UpdateUserPermission)
	app.Delete("/permission/user/:id", jwt.DecodeTokenAdmin, handler.DeleteUserPermission)

	// ----  Employee ----
	app.Get("/employee", jwt.DecodeTokenAdmin, handler.GetEmployeeAll)
	app.Post("/employee", jwt.DecodeTokenAdmin, handler.InsertEmployee)
	app.Put("/employee/:code", jwt.DecodeTokenAdmin, handler.UpdateEmployee)
	app.Delete("/employee/:code", jwt.DecodeTokenAdmin, handler.DeleteEmployee)
	app.Get("/employee/:code", jwt.DecodeToken, handler.GetEmployeeByCode)

	// Auth
	app.Get("/auth", jwt.CheckToken)

	app.Get("/container", handler.TestApp)
	app.Get("/mail/:mail", handler.TestingMail)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "3005" // Default port
	}

	// Application Listen :PORT
	app.Listen(":" + PORT)

}
