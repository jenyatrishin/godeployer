package deployerFactory

import "../../deployer"

func GetDeployer () deployer.AbstractDeployer {
	return &deployer.SshDeployer{}
}
