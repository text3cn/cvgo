package kvsKey

const WorkspacePath = "WorkspacePath"
const WorkspaceName = "WorkspaceName"
const AllocatedPort = "AllocatedPort"

func ModuleWebFramework(modName string) string {
	return modName + ".webFramework"
}

func ModuleHttpPort(modName string) string {
	return modName + ".httpPort"
}

func ModuleSwaggerEnable(modName string) string {
	return modName + ".swagger"
}

func ModuleI18n(modName string) string {
	return modName + ".i18n"
}

func ModuleRoute(modName string) string {
	return modName + ".routes"
}

func ModuleSvc(modName string) string {
	return modName + ".services"
}
