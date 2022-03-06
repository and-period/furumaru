# -*- mode: Python -*-

if k8s_context() != 'docker-desktop':
  fail('failed to connect kubernetes cluster for local env')

allow_k8s_contexts('docker-desktop')

##################################################
# User Gateway
##################################################
docker_build(
  'marche/user-gateway',
  '.',
  dockerfile='infra/docker/api/gateway/user/Dockerfile'
)

k8s_yaml(helm(
  'infra/helm/gateway',
  name='user-gateway',
  namespace='default',
  values=['./api/config/gateway/user/dev.yaml']
))

k8s_resource(
  'user-gateway',
  port_forwards=['18000:9000', '18001:9001', '18002:9002'],
  labels=['gateway']
)

##################################################
# User Servcer
##################################################
docker_build(
  'marche/user-server',
  '.',
  dockerfile='infra/docker/api/user/server/Dockerfile'
)

k8s_yaml(helm(
  'infra/helm/service',
  name='user-server',
  namespace='default',
  values=['./api/config/user/server/dev.yaml']
))

k8s_resource(
  'user-server',
  port_forwards=['19000:9000', '19001:9001', '19002:9002'],
  labels=['service']
)
