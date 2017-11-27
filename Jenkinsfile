import java.text.SimpleDateFormat

node('linux-ubuntu-16.04-amd64') {
	checkout scm
	def services = ['evauth', 'evbolt', 'evlog', 'evschedule']
	def tools = ['eve-bintray', 'eve-gen', 'eve-setup']
	def dependencies = ['github.com/boltdb/bolt', 'github.com/gorilla/mux', 'github.com/prometheus/client_golang/prometheus', 'github.com/prometheus/client_golang/prometheus/promhttp', 'github.com/dchest/uniuri', 'github.com/mitchellh/go-ps', 'github.com/axw/gocov/...', 'github.com/AlekSi/gocov-xml','github.com/kless/osutil/user/crypt/sha512_crypt']
	def oses = ['darwin', 'linux', 'windows']
	def archs = ['amd64']
	def version = '0.0.2'
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
	def slackNotificationChannel = "build"
	switch (env.BRANCH_NAME) {
		case "master":
			try {
				withEnv(["GOROOT=${curr}/${build}/${goroot}", "GOPATH=${curr}/${build}/${gopath}", "PATH+GOPATHBIN=${curr}/${build}/${gopath}/bin", "PATH+GOROOTBIN=${curr}/${build}/${goroot}/bin"]){
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
						sh("${gobin} test -coverprofile=dist/coverage.out")
						sh("cd ${curr}/${build}/${gopath}/src/evalgo.org/eve && gocov test | gocov-xml > ${curr}/dist/coverage.xml")
					}
					stage ('Upload CodeCoverage to codecov.io'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'codecov', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("cd ${curr}/${build}/${gopath}/src/evalgo.org/eve && curl -o codecov.sh https://codecov.io/bash && chmod +x codecov.sh && ./codecov.sh -t ${PASSWORD}")
						}
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
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${curr}/${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} ${src}/bin/${tools[t]}/main.go")
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
							sh("${curr}/${dist}/linux-amd64-${version}_eve-gen generate -service ${services[s]} ${use_flags}  -target ${tmp}/${services[s]}_main.go")
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${curr}/${dist}/${oses[o]}-${archs[a]}-${version}_${services[s]}${ext} ${tmp}/${services[s]}_main.go")
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
										sh("curl --user '${USERNAME}:${PASSWORD}' -T ${curr}/${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} https://u162240.your-storagebox.de/eve/backend/${osrename}/")
									}
								}
							}
					}
				}
				stage ('CleanUp EVE ARTIFACTS at BinTray'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'bintray', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("${curr}/${dist}/linux-amd64-0.0.2_eve-bintray evalgo eve-backend core ${version} https://api.bintray.com ${USERNAME} ${PASSWORD}")
						}
				}
				stage ('Deploy EVE ARTIFACTS to BinTray'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'bintray', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("cp -Rf ${dist} "+dateFormat.format(date))
							sh("zip -r "+dateFormat.format(date)+".zip "+dateFormat.format(date))
							sh("curl -v -X PUT --header 'X-Bintray-Package: core' --header 'X-Bintray-Version: ${version}' --header 'X-Bintray-Publish: 1' --header 'X-Bintray-Override: 1' --header 'X-Bintray-Explode: 1' --user '${USERNAME}:${PASSWORD}' -T "+dateFormat.format(date)+".zip 'https://api.bintray.com/content/evalgo/eve-backend/"+dateFormat.format(date)+".zip'")
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
					stage("master :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
		        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} success", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
		    	}
				}
			} catch (Exception e) {
				stage("master error :: post to slack") {
					withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
						slackSend channel: '#build', color: 'danger', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} failed", teamDomain: "${USERNAME}", token: "${PASSWORD}"
					}
				}
				throw e;
			}
			break;
		case "dev":
			try {
				withEnv(["GOROOT=${curr}/${build}/${goroot}", "GOPATH=${curr}/${build}/${gopath}", , "PATH+GOPATHBIN=${curr}/${build}/${gopath}/bin", "PATH+GOROOTBIN=${curr}/${build}/${goroot}/bin"]){
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
						sh("${gobin} test -coverprofile=dist/coverage.out")
						sh("cd ${curr}/${build}/${gopath}/src/evalgo.org/eve && gocov test | gocov-xml > ${curr}/dist/coverage.xml")
					}
					stage ('Upload CodeCoverage to codecov.io'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'codecov', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("cd ${curr}/${build}/${gopath}/src/evalgo.org/eve && curl -o codecov.sh https://codecov.io/bash && chmod +x codecov.sh && ./codecov.sh -t ${PASSWORD}")
						}
					}
					stage("dev :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
		        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} success", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
		    	}
				}
			} catch (Exception e) {
				stage("dev error :: post to slack") {
					withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
						slackSend channel: '#build', color: 'danger', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} failed", teamDomain: "${USERNAME}", token: "${PASSWORD}"
					}
				}
				throw e;
			}
			break;
		case "test":
			try {
				withEnv(["GOROOT=${curr}/${build}/${goroot}", "GOPATH=${curr}/${build}/${gopath}", , "PATH+GOPATHBIN=${curr}/${build}/${gopath}/bin", "PATH+GOROOTBIN=${curr}/${build}/${goroot}/bin"]){
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
						sh("${gobin} test -coverprofile=dist/coverage.out")
						sh("cd ${curr}/${build}/${gopath}/src/evalgo.org/eve && gocov test | gocov-xml > ${curr}/dist/coverage.xml")
					}
					stage ('Upload CodeCoverage to codecov.io'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'codecov', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("cd ${curr}/${build}/${gopath}/src/evalgo.org/eve && curl -o codecov.sh https://codecov.io/bash && chmod +x codecov.sh && ./codecov.sh -t ${PASSWORD}")
						}
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
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${curr}/${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} ${src}/bin/${tools[t]}/main.go")
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
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} ${gobin} build -o ${curr}/${dist}/${oses[o]}-${archs[a]}-${version}_${services[s]}${ext} ${tmp}/${services[s]}_main.go")
								}
							}
						}
					}
				}
				stage ('Archive EVE ARTIFACTS'){
					archiveArtifacts("${dist}/*")
				}
				stage("test :: post to slack") {
					withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
	        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} success", teamDomain: "${USERNAME}", token: "${PASSWORD}"
					}
	    	}
			} catch (Exception e) {
				stage("test error :: post to slack") {
					withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
						slackSend channel: '#build', color: 'danger', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} failed", teamDomain: "${USERNAME}", token: "${PASSWORD}"
					}
				}
				throw e;
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
