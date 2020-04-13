from invoke import task

GO_VERSION = "1.13"
GO_BIN = f"/usr/lib/go-{GO_VERSION}/bin/go"

DOCKER_TAG = "latest"
DOCKER_IMAGE = f"jrabbit/spongebob_exe:{DOCKER_TAG}"

@task
def build(c, fmt=True, tag=True):
    "build project and make docker image"
    if fmt:
        c.run(f"{GO_BIN} fmt")
    c.run(f"{GO_BIN} build")
    if tag:
        c.run(f"docker build -t {DOCKER_IMAGE} .")

@task
def deps(c):
    c.run("docker pull gcr.io/distroless/base")
    #c.run("docker pull golang")

@task
def release(c):
    "push the docker image"
    c.run(f"docker push {DOCKER_IMAGE}")

@task
def prod(c, kill=False, kubectl="sudo k3s kubectl", load_secret=False):
    if load_secret:
        c.run(f"{kubectl} create secret generic sponge-sec-discord --from-file=discord.toml")
    if kill:
        c.run(f"{kubectl} delete pod static-spongebob")
    c.run(f"{kubectl} apply -f spongebob_pod.yaml")

@task
def pull(c):
    c.run(f"docker pull {DOCKER_IMAGE}")
