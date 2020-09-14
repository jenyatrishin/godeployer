package deployerFactory

import "dep2go/deployer"

func GetDeployer () deployer.AbstractDeployer {
	return &deployer.SshDeployer{}
}
