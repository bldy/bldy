package gvisor

import (
	"bytes"
	"encoding/json"

	"bldy.build/bldy/src/namespace/gvisor/boot"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"gvisor.dev/gvisor/runsc/boot/platforms"
)

var specTemplate = []byte(`{
	"ociVersion": "1.0.0",
	"process": {
		"terminal": true,
		"user": {
			"uid": 0,
			"gid": 0
		},
		"args": [
			"sh"
		],
		"env": [
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"TERM=xterm"
		],
		"cwd": "/",
		"capabilities": {
			"bounding": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"effective": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"inheritable": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"permitted": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			],
			"ambient": [
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE"
			]
		},
		"rlimits": [
			{
				"type": "RLIMIT_NOFILE",
				"hard": 1024,
				"soft": 1024
			}
		]
	},
	"root": {
		"path": "rootfs",
		"readonly": true
	},
	"hostname": "runsc",
	"mounts": [
		{
			"destination": "/proc",
			"type": "proc",
			"source": "proc"
		},
		{
			"destination": "/dev",
			"type": "tmpfs",
			"source": "tmpfs",
			"options": []
		},
		{
			"destination": "/sys",
			"type": "sysfs",
			"source": "sysfs",
			"options": [
				"nosuid",
				"noexec",
				"nodev",
				"ro"
			]
		}
	],
	"linux": {
		"namespaces": [
			{
				"type": "pid"
			},
			{
				"type": "network"
			},
			{
				"type": "ipc"
			},
			{
				"type": "uts"
			},
			{
				"type": "mount"
			}
		]
	}
}`)

func defaultArgs() boot.Args {
	spec := &specs.Spec{}
	dec := json.NewDecoder(bytes.NewBuffer(specTemplate))
	dec.Decode(spec)

	return boot.Args{
		ID:   "testID",
		Spec: spec,
		Conf: &boot.Config{
			Debug:      true,
			LogFormat:  "text",
			LogPackets: true,
			Network:    boot.NetworkNone,
			Strace:     true,
			FileAccess: boot.FileAccessExclusive,
			TestOnlyAllowRunAsCurrentUserWithoutChroot: true,
			NumNetworkChannels:                         1,
			Platform:                                   platforms.Ptrace,
		},
		NumCPU:  1,
		Console: false,
		//	ControllerFD: p.ns.fd,
		//	GoferFDs:     []int{sandEnd},
		//	StdioFDs:     p.ns.stdio,
	}
}
