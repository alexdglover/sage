package dependencyregistry

import (
	"log"
	"os"
	"time"

	"github.com/alexdglover/sage/internal/api"
	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DependencyRegistry struct {
	DbConnection               *gorm.DB
	Bootstrapper               *models.Bootstrapper
	AccountRepository          *models.AccountRepository
	BalanceRepository          *models.BalanceRepository
	BudgetRepository           *models.BudgetRepository
	CategoryRepository         *models.CategoryRepository
	ImportSubmissionRepository *models.ImportSubmissionRepository
	TransactionRepository      *models.TransactionRepository

	AccountManager *services.AccountManager
	ImportService  *services.ImportService

	AccountController     *api.AccountController
	BalanceController     *api.BalanceController
	BudgetController      *api.BudgetController
	CategoryController    *api.CategoryController
	ImportController      *api.ImportController
	MainController        *api.MainController
	NetWorthController    *api.NetWorthController
	TransactionController *api.TransactionController
	ApiServer             *api.ApiServer
}

func (dr *DependencyRegistry) GetBootstrapper() *models.Bootstrapper {
	if dr.Bootstrapper == nil {
		dbConnection, err := dr.GetDbConnection()
		if err != nil {
			panic("failed to get database connection")
		}
		dr.Bootstrapper = models.NewBootstrapper(dbConnection)
	}
	return dr.Bootstrapper
}

func (dr *DependencyRegistry) GetDbConnection() (*gorm.DB, error) {
	var err error
	if dr.DbConnection == nil {
		// TODO: Reduce logger verbosity once stable
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				// LogLevel:                  logger.Info, // Log level
				LogLevel:                  logger.Warn, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      false,       // Include params in the SQL log
				Colorful:                  true,        // Enable color
			},
		)

		var sageFilePath string
		if sageFileEnvVar, ok := os.LookupEnv("SAGE_FILE"); ok {
			sageFilePath = sageFileEnvVar
		} else {
			sageFilePath = "sage.db"
		}
		dr.DbConnection, err = gorm.Open(sqlite.Open(sageFilePath), &gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
			Logger: newLogger,
		})
		if err != nil {
			return nil, err
		}
	}
	return dr.DbConnection, nil
}

// Repositories

func (dr *DependencyRegistry) GetAccountRepository() (*models.AccountRepository, error) {
	if dr.AccountRepository == nil {
		dbConnection, err := dr.GetDbConnection()
		if err != nil {
			return nil, err
		}
		dr.AccountRepository = &models.AccountRepository{
			DB: dbConnection,
		}
	}
	return dr.AccountRepository, nil
}

func (dr *DependencyRegistry) GetBalanceRepository() (*models.BalanceRepository, error) {
	if dr.BalanceRepository == nil {
		dbConnection, err := dr.GetDbConnection()
		if err != nil {
			return nil, err
		}
		dr.BalanceRepository = &models.BalanceRepository{
			DB: dbConnection,
		}
	}
	return dr.BalanceRepository, nil
}

func (dr *DependencyRegistry) GetBudgetRepository() (*models.BudgetRepository, error) {
	if dr.BudgetRepository == nil {
		dbConnection, err := dr.GetDbConnection()
		if err != nil {
			return nil, err
		}
		dr.BudgetRepository = &models.BudgetRepository{
			DB: dbConnection,
		}
	}
	return dr.BudgetRepository, nil
}

func (dr *DependencyRegistry) GetCategoryRepository() (*models.CategoryRepository, error) {
	if dr.CategoryRepository == nil {
		dbConnection, err := dr.GetDbConnection()
		if err != nil {
			return nil, err
		}
		dr.CategoryRepository = &models.CategoryRepository{
			DB: dbConnection,
		}
	}
	return dr.CategoryRepository, nil
}

func (dr *DependencyRegistry) GetTransactionRepository() (*models.TransactionRepository, error) {
	if dr.TransactionRepository == nil {
		dbConnection, err := dr.GetDbConnection()
		if err != nil {
			return nil, err
		}
		dr.TransactionRepository = &models.TransactionRepository{
			DB: dbConnection,
		}
	}
	return dr.TransactionRepository, nil
}

func (dr *DependencyRegistry) GetImportSubmissionRepository() (*models.ImportSubmissionRepository, error) {
	if dr.ImportSubmissionRepository == nil {
		dbConnection, err := dr.GetDbConnection()
		if err != nil {
			return nil, err
		}
		dr.ImportSubmissionRepository = &models.ImportSubmissionRepository{
			DB: dbConnection,
		}
	}
	return dr.ImportSubmissionRepository, nil
}

// Services

func (dr *DependencyRegistry) GetAccountManager() (*services.AccountManager, error) {
	if dr.AccountManager == nil {
		accountRepository, err := dr.GetAccountRepository()
		if err != nil {
			return nil, err
		}
		dr.AccountManager = &services.AccountManager{
			AccountRepository: accountRepository,
		}
	}
	return dr.AccountManager, nil
}

func (dr *DependencyRegistry) GetImportService() (*services.ImportService, error) {
	if dr.ImportService == nil {
		accountRepository, err := dr.GetAccountRepository()
		if err != nil {
			return nil, err
		}
		balanceRepository, err := dr.GetBalanceRepository()
		if err != nil {
			return nil, err
		}
		importSubmissionRepository, err := dr.GetImportSubmissionRepository()
		if err != nil {
			return nil, err
		}
		transactionRepository, err := dr.GetTransactionRepository()
		if err != nil {
			return nil, err
		}
		dr.ImportService = &services.ImportService{
			AccountRepository:          accountRepository,
			BalanceRepository:          balanceRepository,
			ImportSubmissionRepository: importSubmissionRepository,
			TransactionRepository:      transactionRepository,
		}
	}
	return dr.ImportService, nil
}

// APIs

func (dr *DependencyRegistry) GetAccountController() (*api.AccountController, error) {
	if dr.AccountController == nil {
		accountRepository, err := dr.GetAccountRepository()
		if err != nil {
			return nil, err
		}
		balanceRepository, err := dr.GetBalanceRepository()
		if err != nil {
			return nil, err
		}
		dr.AccountController = &api.AccountController{
			AccountRepository: accountRepository,
			BalanceRepository: balanceRepository,
		}
	}
	return dr.AccountController, nil
}

func (dr *DependencyRegistry) GetBalanceController() (*api.BalanceController, error) {
	if dr.BalanceController == nil {
		accountRepository, err := dr.GetAccountRepository()
		if err != nil {
			return nil, err
		}
		balanceRepository, err := dr.GetBalanceRepository()
		if err != nil {
			return nil, err
		}
		dr.BalanceController = &api.BalanceController{
			AccountRepository: accountRepository,
			BalanceRepository: balanceRepository,
		}
	}
	return dr.BalanceController, nil
}

func (dr *DependencyRegistry) GetBudgetController() (*api.BudgetController, error) {
	if dr.BudgetController == nil {
		budgetRepository, err := dr.GetBudgetRepository()
		if err != nil {
			return nil, err
		}
		categoryRepository, err := dr.GetCategoryRepository()
		if err != nil {
			return nil, err
		}
		dr.BudgetController = &api.BudgetController{
			BudgetRepository:   budgetRepository,
			CategoryRepository: categoryRepository,
		}
	}
	return dr.BudgetController, nil
}

func (dr *DependencyRegistry) GetCategoryController() (*api.CategoryController, error) {
	if dr.CategoryController == nil {
		balanceRepository, err := dr.GetBalanceRepository()
		if err != nil {
			return nil, err
		}
		categoryRepository, err := dr.GetCategoryRepository()
		if err != nil {
			return nil, err
		}
		dr.CategoryController = &api.CategoryController{
			BalanceRepository:  balanceRepository,
			CategoryRepository: categoryRepository,
		}
	}
	return dr.CategoryController, nil
}

func (dr *DependencyRegistry) GetImportController() (*api.ImportController, error) {
	if dr.ImportController == nil {
		accountManager, err := dr.GetAccountManager()
		if err != nil {
			return nil, err
		}
		importService, err := dr.GetImportService()
		if err != nil {
			return nil, err
		}
		transactionRepository, err := dr.GetTransactionRepository()
		if err != nil {
			return nil, err
		}
		dr.ImportController = &api.ImportController{
			AccountManager:        accountManager,
			ImportService:         importService,
			TransactionRepository: transactionRepository,
		}
	}
	return dr.ImportController, nil
}

func (dr *DependencyRegistry) GetMainController() (*api.MainController, error) {
	if dr.MainController == nil {
		dr.MainController = &api.MainController{}
	}
	return dr.MainController, nil
}

func (dr *DependencyRegistry) GetNetWorthController() (*api.NetWorthController, error) {
	if dr.NetWorthController == nil {
		balanceRepository, err := dr.GetBalanceRepository()
		if err != nil {
			return nil, err
		}
		dr.NetWorthController = &api.NetWorthController{
			BalanceRepository: balanceRepository,
		}
	}
	return dr.NetWorthController, nil
}

func (dr *DependencyRegistry) GetTransactionController() (*api.TransactionController, error) {
	if dr.TransactionController == nil {
		accountRepository, err := dr.GetAccountRepository()
		if err != nil {
			return nil, err
		}
		transactionRepository, err := dr.GetTransactionRepository()
		if err != nil {
			return nil, err
		}
		dr.TransactionController = &api.TransactionController{
			AccountRepository:     accountRepository,
			TransactionRepository: transactionRepository,
		}
	}
	return dr.TransactionController, nil
}

func (dr *DependencyRegistry) GetApiServer() (*api.ApiServer, error) {
	if dr.ApiServer == nil {
		accountController, err := dr.GetAccountController()
		if err != nil {
			return nil, err
		}
		balanceController, err := dr.GetBalanceController()
		if err != nil {
			return nil, err
		}
		budgetController, err := dr.GetBudgetController()
		if err != nil {
			return nil, err
		}
		categoryController, err := dr.GetCategoryController()
		if err != nil {
			return nil, err
		}
		importController, err := dr.GetImportController()
		if err != nil {
			return nil, err
		}
		mainController, err := dr.GetMainController()
		if err != nil {
			return nil, err
		}
		netWorthController, err := dr.GetNetWorthController()
		if err != nil {
			return nil, err
		}
		transactionController, err := dr.GetTransactionController()
		if err != nil {
			return nil, err
		}
		dr.ApiServer = &api.ApiServer{
			AccountController:     accountController,
			BalanceController:     balanceController,
			BudgetController:      budgetController,
			CategoryController:    categoryController,
			ImportController:      importController,
			MainController:        mainController,
			NetWorthController:    netWorthController,
			TransactionController: transactionController,
		}
	}
	return dr.ApiServer, nil
}
