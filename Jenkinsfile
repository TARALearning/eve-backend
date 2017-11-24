import java.text.SimpleDateFormat

node('linux-ubuntu-16.04-amd64') {
	checkout scm
	def services = ['evauth', 'evbolt', 'evlog', 'evschedule']
	def tools = ['eve-gen', 'eve-setup']
	def dependencies = ['github.com/boltdb/bolt', 'github.com/gorilla/mux', 'github.com/prometheus/client_golang/prometheus', 'github.com/prometheus/client_golang/prometheus/promhttp', 'github.com/dchest/uniuri', 'github.com/mitchellh/go-ps', 'go get github.com/axw/gocov/...', 'go get github.com/AlekSi/gocov-xml']
	def oses = ['darwin', 'linux', 'windows']
	def archs = ['amd64']
	def version = '0.0.1'
	def ext = ''
	def dist = 'dist'
	def curr = pwd()
	def build = 'build'
	def go_version = '1.9.2'
	def go = "go${go_version}.linux-amd64"
	def goroot = "go-${go_version}"
	def gopath = "${go}-${go_version}-packages"
	def gobin = "${curr}/${build}/${goroot}/bin/go"
	def gopathbin = "${build}/${gopath}/bin"
	def src = "${build}/${gopath}/src/evalgo.org/eve"
	def tmp = "tests/tmp"
	def dateFormat = new SimpleDateFormat("yyyy-MM-dd-HH-mm")
	def date = new Date()
	def use_flags = ""
	def use_flags_default = "-use debug"
	switch (env.BRANCH_NAME) {
		case "master":
			withEnv(["GOROOT=${curr}/${build}/${goroot}", "GOPATH=${curr}/${build}/${gopath}"]){
				stage ('Init GO ENV'){
					sh("rm -rf ${build} .goget ${dist}")
					sh("mkdir ${build} ${dist}")
					sh("cd ${build} && wget -q https://redirector.gvt1.com/edgedl/go/${go}.tar.gz")
					sh("cd ${build} && tar xfz ${go}.tar.gz && mv go go-${go_version} && rm -f ${go}.tar.gz")
					sh("mkdir -p ${build}/${gopath}/src/evalgo.org/eve")
					sh("rsync -av --exclude='${build}' ./ ${src}/")
					for(int i = 0; i < dependencies.size(); i++) {
						sh("${gobin} get -v ${dependencies[i]}")
					}
				}
				stage ('Run TESTS'){
					sh("cd tests && chmod +x gen.ssl.client.crt.sh && ./gen.ssl.client.crt.sh")
					sh("${gobin} test -v")
					sh("go test -coverprofile=dist/coverage.out")
					sh("go tool cover -func=dist/coverage.out")
					sh("${gopathbin}/gocov test evalgo.org/eve | ${gopathbin}/gocov-xml > dist/coverage.xml")
				}
				stage ('Build TOOLS'){
					for (int t = 0; t < tools.size(); t++){
						for (int o = 0; o < oses.size(); o++){
							if ("${oses[o]}" == "windows"){
								ext = ".exe"
							}else{
								ext = ""
							}
							for (int a = 0; a < archs.size(); a++){
								sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} ${src}/bin/${tools[t]}/main.go")
							}
						}
					}
				}
				stage ('Build EVE SERVICES') {
					for (int s = 0; s < services.size(); s++){
						switch("${services[s]}".toString()){
							case "evbolt":
								use_flags = "-use evBoltRoot=. "+use_flags_default
								break;
							default:
								use_flags = use_flags_default
								break;
						}
						sh("${dist}/linux-amd64-${version}_eve-gen generate -service ${services[s]} ${use_flags}  -target ${tmp}/${services[s]}_main.go")
						for (int o = 0; o < oses.size(); o++){
							if ("${oses[o]}" == "windows"){
								ext = ".exe"
							}else{
								ext = ""
							}
							for (int a = 0; a < archs.size(); a++){
								sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${services[s]}${ext} ${tmp}/${services[s]}_main.go")
							}
						}
					}
			}
			stage ('Archive EVE ARTIFACTS'){
				archiveArtifacts("${dist}/*")
			}
			stage ('Deploy EVE ARTIFACTS to StorageBox'){
				withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'storagebox', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
						for (int s = 0; s < services.size(); s++){
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									def osrename = "${oses[o]}"
									if ("${oses[o]}" == "darwin"){
										osrename = "macos"
									}
									sh("curl --user '${USERNAME}:${PASSWORD}' -T ${dist}/${oses[o]}-${archs[a]}-${version}_${services[s]}${ext} https://u162240.your-storagebox.de/eve/backend/${osrename}/")
								}
							}
						}
						for (int t = 0; t < tools.size(); t++){
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									def osrename = "${oses[o]}"
									if ("${oses[o]}" == "darwin"){
										osrename = "macos"
									}
									sh("curl --user '${USERNAME}:${PASSWORD}' -T ${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} https://u162240.your-storagebox.de/eve/backend/${osrename}/")
								}
							}
						}
				}
			}
			stage ('Deploy EVE ARTIFACTS to BinTray'){
					withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'bintray', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
						sh("cp -Rf ${dist} "+dateFormat.format(date))
						sh("zip -r "+dateFormat.format(date)+".zip "+dateFormat.format(date))
						sh("curl -v -X PUT --header 'X-Bintray-Package: core' --header 'X-Bintray-Version: 0.0.1' --header 'X-Bintray-Publish: 1' --header 'X-Bintray-Override: 1' --header 'X-Bintray-Explode: 1' --user '${USERNAME}:${PASSWORD}' -T "+dateFormat.format(date)+".zip 'https://api.bintray.com/content/evalgo/eve-backend/"+dateFormat.format(date)+".zip'")
						sh("sleep 20")
						for (int s = 0; s < services.size(); s++){
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									def osrename = "${oses[o]}"
									if ("${oses[o]}" == "darwin"){
										osrename = "macos"
									}
									sh("sleep 10")
									sh("curl -v -X PUT -d '{\"list_in_downloads\":true}' --header 'Content-Type: application/json' --user '${USERNAME}:${PASSWORD}' 'https://api.bintray.com/file_metadata/evalgo/eve-backend/"+dateFormat.format(date)+"%2F${oses[o]}-${archs[a]}-${version}_${services[s]}${ext}'")
								}
							}
						}
						for (int t = 0; t < tools.size(); t++){
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									def osrename = "${oses[o]}"
									if ("${oses[o]}" == "darwin"){
										osrename = "macos"
									}
									sh("sleep 10")
									sh("curl -v -X PUT -d '{\"list_in_downloads\":true}' --header 'Content-Type: application/json' --user '${USERNAME}:${PASSWORD}' 'https://api.bintray.com/file_metadata/evalgo/eve-backend/"+dateFormat.format(date)+"%2F${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext}'")
								}
							}
						}
					}
				}
			}
			break;
		case "dev":
			withEnv(["GOROOT=${curr}/${build}/${goroot}", "GOPATH=${curr}/${build}/${gopath}"]){
				stage ('Init GO ENV'){
					sh("rm -rf ${build} .goget ${dist}")
					sh("mkdir ${build} ${dist}")
					sh("cd ${build} && wget -q https://redirector.gvt1.com/edgedl/go/${go}.tar.gz")
					sh("cd ${build} && tar xfz ${go}.tar.gz && mv go go-${go_version} && rm -f ${go}.tar.gz")
					sh("mkdir -p ${build}/${gopath}/src/evalgo.org/eve")
					sh("rsync -av --exclude='${build}' ./ ${src}/")
					for(int i = 0; i < dependencies.size(); i++) {
						sh("${gobin} get -v ${dependencies[i]}")
					}
				}
				stage ('Run TESTS'){
					sh("cd tests && chmod +x gen.ssl.client.crt.sh && ./gen.ssl.client.crt.sh")
					sh("${gobin} test -v")
					sh("go test -coverprofile=dist/coverage.out")
					sh("go tool cover -func=dist/coverage.out")
					sh("${gopathbin}/gocov test evalgo.org/eve | ${gopathbin}/gocov-xml > dist/coverage.xml")
				}
			}
			break;
		case "test":
			withEnv(["GOROOT=${curr}/${build}/${goroot}", "GOPATH=${curr}/${build}/${gopath}"]){
				stage ('Init GO ENV'){
					sh("rm -rf ${build} .goget ${dist}")
					sh("mkdir ${build} ${dist}")
					sh("cd ${build} && wget -q https://redirector.gvt1.com/edgedl/go/${go}.tar.gz")
					sh("cd ${build} && tar xfz ${go}.tar.gz && mv go go-${go_version} && rm -f ${go}.tar.gz")
					sh("mkdir -p ${build}/${gopath}/src/evalgo.org/eve")
					sh("rsync -av --exclude='${build}' ./ ${src}/")
					for(int i = 0; i < dependencies.size(); i++) {
						sh("${gobin} get -v ${dependencies[i]}")
					}
				}
				stage ('Run TESTS'){
					sh("cd tests && chmod +x gen.ssl.client.crt.sh && ./gen.ssl.client.crt.sh")
					sh("${gobin} test -v")
					sh("go test -coverprofile=dist/coverage.out")
					sh("go tool cover -func=dist/coverage.out")
					sh("${gopathbin}/gocov test evalgo.org/eve | ${gopathbin}/gocov-xml > dist/coverage.xml")
				}
				stage ('Build TOOLS'){
					for (int t = 0; t < tools.size(); t++){
						for (int o = 0; o < oses.size(); o++){
							if ("${oses[o]}" == "windows"){
								ext = ".exe"
							}else{
								ext = ""
							}
							for (int a = 0; a < archs.size(); a++){
								sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} ${src}/bin/${tools[t]}/main.go")
							}
						}
					}
				}
				stage ('Build EVE SERVICES') {
					for (int s = 0; s < services.size(); s++){
						switch("${services[s]}".toString()){
							case "evbolt":
								use_flags = "-use evBoltRoot=. "+use_flags_default
								break;
							default:
								use_flags = use_flags_default
								break;
						}
						sh("${dist}/linux-amd64-${version}_eve-gen generate -service ${services[s]} ${use_flags}  -target ${tmp}/${services[s]}_main.go")
						for (int o = 0; o < oses.size(); o++){
							if ("${oses[o]}" == "windows"){
								ext = ".exe"
							}else{
								ext = ""
							}
							for (int a = 0; a < archs.size(); a++){
								sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${services[s]}${ext} ${tmp}/${services[s]}_main.go")
							}
						}
					}
				}
			}
			stage ('Archive EVE ARTIFACTS'){
				archiveArtifacts("${dist}/*")
			}
			break;
		default:
			stage ('no pipeline'){
				sh('echo "this branch ${env.BRANCH_NAME} has no pipeline stages now!"')
			}
			break;
    }
		cobertura autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: 'dist/coverage.xml', conditionalCoverageTargets: '70, 0, 0', failUnhealthy: false, failUnstable: false, lineCoverageTargets: '80, 0, 0', maxNumberOfBuilds: 0, methodCoverageTargets: '80, 0, 0', onlyStable: false, sourceEncoding: 'ASCII'
}
