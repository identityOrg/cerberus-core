package core

import (
	"github.com/google/wire"
	"github.com/identityOrg/oidcsdk"
)

var ProviderSet = wire.NewSet(
	NewTokenStoreServiceImpl,
	NewSPStoreServiceImpl,
	NewUserStoreServiceImpl,
	NewScopeClaimStoreServiceImpl,
	NewSecretStoreServiceImpl,
	wire.Bind(new(ITokenStoreService), new(*TokenStoreServiceImpl)),
	wire.Bind(new(oidcsdk.ITokenStore), new(*TokenStoreServiceImpl)),
	wire.Bind(new(ISPStoreService), new(*SPStoreServiceImpl)),
	wire.Bind(new(oidcsdk.IClientStore), new(*SPStoreServiceImpl)),
	wire.Bind(new(IUserStoreService), new(*UserStoreServiceImpl)),
	wire.Bind(new(oidcsdk.IUserStore), new(*UserStoreServiceImpl)),
	wire.Bind(new(ISecretStoreService), new(*SecretStoreServiceImpl)),
	wire.Bind(new(oidcsdk.ISecretStore), new(*SecretStoreServiceImpl)),
	wire.Bind(new(IScopeClaimStoreService), new(*ScopeClaimStoreServiceImpl)),
)
