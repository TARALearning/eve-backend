import java.text.SimpleDateFormat

node('linux-ubuntu-16.04-amd64') {
	checkout scm
	def services = ['evauth', 'evbolt', 'evlog', 'evschedule', 'evweb']
	def tools = ['eve']
	def dependencies = ['github.com/boltdb/bolt', 'github.com/gorilla/mux', 'github.com/prometheus/client_golang/prometheus', 'github.com/prometheus/client_golang/prometheus/promhttp', 'github.com/dchest/uniuri', 'github.com/mitchellh/go-ps', 'github.com/axw/gocov/...', 'github.com/AlekSi/gocov-xml','github.com/kless/osutil/user/crypt/sha512_crypt']
	def oses = ['darwin', 'linux', 'windows']
	def archs = ['amd64']
	def version = '0.0.3'
	def ext = ''
	def curr = pwd()
	def build = "${curr}/build"
	def dist = "${curr}/dist"
	def go_version = '1.9'
	def api_bintray = "https://api.bintray.com"
	def goroot = "/usr/lib/go-${go_version}"
	def gopath = "${curr}/${build}/go-${go_version}-packages"
	def gorootbin = "${goroot}/bin"
	def gopathbin = "${gopath}/bin"
	def src = "${gopath}/src/evalgo.org/eve"
	def tmp = "tests/tmp"
	def dateFormat = new SimpleDateFormat("yyyy-MM-dd-HH-mm")
	def date = new Date()
	def use_flags = ""
	def use_flags_default = "-use debug"
	def slackNotificationChannel = "build"
	switch (env.BRANCH_NAME) {
		case "master":
			try {
				withEnv(["GOROOT=${goroot}", "GOPATH=${gopath}", "PATH+GOPATHBIN=/${gopath}/bin", "PATH+GOROOTBIN=${goroot}/bin"]){
					stage("master build start :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} starting build...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
					}
					stage ('init go environment'){
						sh("rm -rf ${build} ${dist}")
						sh("mkdir ${build} ${dist}")
						sh("mkdir -p ${gopath}/src/evalgo.org/eve")
						sh("rsync -av --exclude='build' --exclude='dist' ./ ${src}/")
						for(int i = 0; i < dependencies.size(); i++) {
							sh("go get -v ${dependencies[i]}")
						}
					}
					stage("master running tests and codecoverage :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
		        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} running tests and code coverage analysis...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
		    	}
					stage ('run unit tests'){
						sh("cd ${src}/tests && chmod +x gen.ssl.client.crt.sh && ./gen.ssl.client.crt.sh")
						sh("cd ${src} && gocov test -v -race -timeout=5m | gocov-xml > ${dist}/coverage.xml")
					}
					stage ('Upload CodeCoverage to codecov.io'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'codecov', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("cd ${dist} && curl -o codecov.sh https://codecov.io/bash && chmod +x codecov.sh && ./codecov.sh -t ${PASSWORD}")
						}
					}
					stage("master build tools :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
		        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} building tools...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
		    	}
					stage ('build tools'){
						for (int t = 0; t < tools.size(); t++){
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} go build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} ${src}/bin/${tools[t]}/main.go")
								}
							}
						}
					}
					stage("master build services :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
		        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} building services...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
		    	}
					stage ('build services') {
						for (int s = 0; s < services.size(); s++){
							switch("${services[s]}".toString()){
								case "evbolt":
									use_flags = "-use evBoltRoot=. "+use_flags_default
									break;
								default:
									use_flags = use_flags_default
									break;
							}
							sh("eve generate golang -service ${services[s]} ${use_flags} -target ${tmp}/${services[s]}_main.go")
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} go build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${services[s]}${ext} ${tmp}/${services[s]}_main.go")
								}
							}
						}
				}
				stage ('archive artifacts'){
					archiveArtifacts("dist/*")
				}
				stage("master cleanup bintray :: post to slack") {
					withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
						slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} cleanup old bintray versions...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
					}
				}
				stage ('cleanup artifacts at bintray'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'bintray', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("eve delete bintray -subject evalgo -repo eve-backend -rpackage core -version ${version} -url ${api_bintray} -username ${USERNAME} -password ${PASSWORD}")
						}
				}
				stage("master deploy to bintray :: post to slack") {
					withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
						slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} deploy to bintray...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
					}
				}
				stage ('deploy artifacst to bintray'){
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'bintray', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
							sh("cp -Rf ${dist} "+dateFormat.format(date))
							sh("zip -r "+dateFormat.format(date)+".zip "+dateFormat.format(date))
							sh("curl -v -X PUT --header 'X-Bintray-Package: core' --header 'X-Bintray-Version: ${version}' --header 'X-Bintray-Publish: 1' --header 'X-Bintray-Override: 1' --header 'X-Bintray-Explode: 1' --user '${USERNAME}:${PASSWORD}' -T "+dateFormat.format(date)+".zip '"+api_bintray+"/content/evalgo/eve-backend/"+dateFormat.format(date)+".zip'")
							sh("sleep 10")
							sh("eve publish bintray -subject evalgo -repo eve-backend -rpackage core -version ${version} -url ${api_bintray} -username ${USERNAME} -password ${PASSWORD}")
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
		case "test":
			try {
				withEnv(["GOROOT=${goroot}", "GOPATH=${gopath}", , "PATH+GOPATHBIN=${gopath}/bin", "PATH+GOROOTBIN=${goroot}/bin"]){
					stage ('init go environment'){
						sh("rm -rf ${build} ${dist}")
						sh("mkdir ${build} ${dist}")
						sh("mkdir -p ${src}")
						sh("rsync -av --exclude='build' --exclude='dist' ./ ${src}/")
						for(int i = 0; i < dependencies.size(); i++) {
							sh("go get -v ${dependencies[i]}")
						}
					}
					stage ('run unit tests'){
						sh("cd ${src}/tests && chmod +x gen.ssl.client.crt.sh && ./gen.ssl.client.crt.sh")
						sh("cd ${src} && gocov test -v -race -timeout=5m | gocov-xml > ${dist}/coverage.xml")
					}
					stage("test build tools :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
		        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} building tools...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
		    	}
					stage ('build tools'){
						for (int t = 0; t < tools.size(); t++){
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} go build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${tools[t]}${ext} ${src}/bin/${tools[t]}/main.go")
								}
							}
						}
					}
					stage("test build services :: post to slack") {
						withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'slack', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
		        	slackSend channel: '#build', color: 'good', message: "${env.JOB_NAME} ${env.BUILD_NUMBER} building services...", teamDomain: "${USERNAME}", token: "${PASSWORD}"
						}
		    	}
					stage ('build services') {
						for (int s = 0; s < services.size(); s++){
							switch("${services[s]}".toString()){
								case "evbolt":
									use_flags = "-use evBoltRoot=. "+use_flags_default
									break;
								default:
									use_flags = use_flags_default
									break;
							}
							sh("eve generate golang -service ${services[s]} ${use_flags}  -target ${tmp}/${services[s]}_main.go")
							for (int o = 0; o < oses.size(); o++){
								if ("${oses[o]}" == "windows"){
									ext = ".exe"
								}else{
									ext = ""
								}
								for (int a = 0; a < archs.size(); a++){
									sh("GOOS=${oses[o]} GOARCH=${archs[a]} go build -o ${dist}/${oses[o]}-${archs[a]}-${version}_${services[s]}${ext} ${tmp}/${services[s]}_main.go")
								}
							}
						}
					}
				}
				stage ('archive artifacts'){
					archiveArtifacts("dist/*")
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
		case "dev":
			try {
				withEnv(["GOROOT=${goroot}", "GOPATH=${gopath}", , "PATH+GOPATHBIN=${gopath}/bin", "PATH+GOROOTBIN=${goroot}/bin"]){
					stage ('init go environment'){
						sh("rm -rf ${build} ${dist}")
						sh("mkdir ${build} ${dist}")
						sh("mkdir -p ${src}")
						sh("rsync -av --exclude='build' --exclude='dist' ./ ${src}/")
						for(int i = 0; i < dependencies.size(); i++) {
							sh("go get -v ${dependencies[i]}")
						}
					}
					stage ('run unit tests'){
						sh("cd ${src}/tests && chmod +x gen.ssl.client.crt.sh && ./gen.ssl.client.crt.sh")
						sh("cd ${src} && gocov test -v -race -timeout=5m | gocov-xml > ${dist}/coverage.xml")
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
		default:
			stage ('no pipeline'){
				sh('echo "this branch ${env.BRANCH_NAME} has no pipeline stages now!"')
			}
			break;
    }
		cobertura autoUpdateHealth: false, autoUpdateStability: false, coberturaReportFile: 'dist/coverage.xml', conditionalCoverageTargets: '70, 0, 0', failUnhealthy: false, failUnstable: false, lineCoverageTargets: '80, 0, 0', maxNumberOfBuilds: 0, methodCoverageTargets: '80, 0, 0', onlyStable: false, sourceEncoding: 'ASCII'
}
