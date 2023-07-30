PACKAGE_LIST="sysadm,registryctl,infrastructure,agent,apiserver"
DEFAULT_IMAGE_VER="v1.0.1"
BASE_IMG="hb.sysadm.cn/os/centos:centos7.9.2009"
EMAIL="net_use@bzhy.com"
DEFAULT_REGISTRY_URL="hb.sysadm.cn/sysadm/"
DEFAULT_DEPLOY_SERVER="192.53.117.73"
DEFAULT_DEPLOY_SERVER_PORT="2218"
DEFAULT_DEPLOY_DOCKER_PATH="/usr/bin/docker"
DEFAULT_DEPLOY_DOCKERCOMPOSE_PATH="/usr/local/bin/docker-compose"
DEFAULT_DEPLOY_CONFIG_FILE="/data/k8ssysadm/docker-compose.yml"
DEFAULT_DEPLOY_TYPE="k8s"
DOCKER_BIN_PATH="/usr/bin/docker"

function deploy::package(){
  package_name=$1
  version=$2

  echo "${DEFAULT_DEPLOY_DOCKERCOMPOSE_PATH} -f ${DEFAULT_DEPLOY_CONFIG_FILE} down " |ssh -p ${DEFAULT_DEPLOY_SERVER_PORT}  root@${DEFAULT_DEPLOY_SERVER} -q
  if [ $? != 0 ]; then
    echo "stop service error"
    exit 2
  fi

#  imageID=`echo "${DEFAULT_DEPLOY_DOCKER_PATH} images -q ${DEFAULT_REGISTRY_URL}${package_name}:${version}" |ssh -p ${DEFAULT_DEPLOY_SERVER_PORT}  root@${DEFAULT_DEPLOY_SERVER} -q |tail -n 1`
#  if ["X${imageID}" != "X"]; then
  echo "${DEFAULT_DEPLOY_DOCKER_PATH} rmi ${DEFAULT_REGISTRY_URL}${package_name}:${version}" |ssh -p ${DEFAULT_DEPLOY_SERVER_PORT}  root@${DEFAULT_DEPLOY_SERVER} -q
#  fi

  docker push "${DEFAULT_REGISTRY_URL}${package_name}:${version}"
  echo "${DEFAULT_DEPLOY_DOCKERCOMPOSE_PATH} -f ${DEFAULT_DEPLOY_CONFIG_FILE} up -d " |ssh -p ${DEFAULT_DEPLOY_SERVER_PORT}  root@${DEFAULT_DEPLOY_SERVER} -q
  if [ $? != 0 ]; then
    echo "start service error"
    exit 2
  fi

  echo "start service successful"
}
