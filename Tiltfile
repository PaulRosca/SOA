update_settings ( max_parallel_updates = 3 , k8s_upsert_timeout_secs = 180 , suppress_unused_image_warnings = None ) 

k8s_yaml('mysql-secret.yaml')
k8s_yaml('mysql-storage.yaml')
k8s_yaml('mysql-deployment.yaml')
k8s_yaml('./gateway/deployment.yaml')
k8s_yaml('./catalog/deployment.yaml')

k8s_resource('mysql', port_forwards=3306)

load('ext://helm_resource', 'helm_resource', 'helm_repo')
helm_repo('openfaas-repo', 'https://openfaas.github.io/faas-netes/')
helm_resource('openfaas', 'openfaas/openfaas', resource_deps=['openfaas-repo'],namespace='openfaas')

apply_cmd = """
echo $(kubectl -n openfaas get secret basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode) | faas-cli login -u admin --password-stdin 1>&2 &&
faas-cli up -f stack.yml 1>&2
"""

delete_cmd = 'faas-cli delete -f stack.yml'

k8s_custom_deploy(
    'go-faas',
    apply_cmd=apply_cmd,
    delete_cmd=delete_cmd,
    # apply_cmd will be re-executed whenever these files change
    deps=['./auth/']
)

docker_build('api-gateway', './gateway')
docker_build('catalog', './catalog')

k8s_resource('api-gateway', port_forwards=9000)

local_resource(name='openfaas-gateway-port-forward', serve_cmd='kubectl port-forward svc/gateway -n openfaas 8080:8080', resource_deps=['openfaas'])
local_resource(name='local-registry', cmd='docker stop local-registry;docker rm local-registry;docker run -d -p 6000:5000 --restart=always --name local-registry registry:2.7.0')
