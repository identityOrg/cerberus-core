package core

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewTokenStoreServiceImpl,
	NewSPStoreServiceImpl,
	NewUserStoreServiceImpl,
	NewScopeClaimStoreServiceImpl,
	wire.Bind(new(ITokenStoreService), new(*TokenStoreServiceImpl)),
	wire.Bind(new(ISPStoreService), new(*SPStoreServiceImpl)),
	wire.Bind(new(IUserStoreService), new(*UserStoreServiceImpl)),
	wire.Bind(new(IScopeClaimStoreService), new(*ScopeClaimStoreServiceImpl)),
)
