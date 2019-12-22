from invoke import task

GO_VERSION = "1.13"
GO_BIN = f"/usr/lib/go-{GO_VERSION}/bin/go"

DOCKER_TAG = "latest"
DOCKER_IMAGE = f"docker.pkg.github.com/jrabbit/spongebob_exe/sponge:{DOCKER_TAG}"

@task
def build(c, tag=True):
    "build project and make docker image"
    c.run(f"{GO_BIN} build")
    if tag:
        c.run(f"docker build -t {DOCKER_IMAGE} .")

@task
def release(c):
    "push the docker image"
    c.run(f"docker push {DOCKER_IMAGE}")
