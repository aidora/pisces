#!/usr/bin/env bats

load helpers

function teardown() {
	stop_manager
	stop_docker
}

@test "docker info should return the number of nodes" {
	start_docker 3
	start_manager
	run docker_swarm info
	[ "$status" -eq 0 ]
	[[ "${lines[1]}" == *"Nodes: 3" ]]
}

@test "docker ps -n 3 should return the 3 last containers, including non running one" {
skip
       start_docker 1
       start_manager
       run docker_swarm run -d busybox sleep 42
       run docker_swarm run -d busybox false
       run docker_swarm run -d busybox true
       run docker_swarm ps -a -n 3
       [ "${#lines[@]}" -eq  3 ]
}

@test "docker ps -l should return the last container, including non running one" {
skip
       start_docker 1
       start_manager
       run docker_swarm run -d busybox sleep 42
       run docker_swarm run -d busybox false
       run docker_swarm run -d busybox true
       run docker_swarm ps -l
       [ "${#lines[@]}" -eq  2 ]
       [[ "${lines[1]}" == *"true"* ]]
}
