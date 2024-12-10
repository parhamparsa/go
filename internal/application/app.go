package application

type AppCore struct {
	UserService     UserServiceInterface
	AccountService  AccountServiceInterface
	TransferService TransferServiceInterface
}

func NewAppCore(
	UserService UserServiceInterface,
	AccountService AccountServiceInterface,
	TransferService TransferServiceInterface,
) AppCore {
	return AppCore{
		UserService:     UserService,
		AccountService:  AccountService,
		TransferService: TransferService,
	}
}

func (a AppCore) GetUserService() UserServiceInterface {
	return a.UserService
}

func (a AppCore) GetAccountService() AccountServiceInterface {
	return a.AccountService
}

func (a AppCore) GetTransferService() TransferServiceInterface {
	return a.TransferService
}

var _ App = &AppCore{}
