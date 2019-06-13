package ipc

// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/bits/ipc.h.html

const (
	/* Mode bits for `msgget', `semget', and `shmget'.  */
	IPC_CREAT  = 01000 // create key if key does not exist.
	IPC_EXCL   = 02000 // fail if key exists.
	IPC_NOWAIT = 04000 // return error on wait.

	/* Control commands for `msgctl', `semctl', and `shmctl'.  */
	IPC_RMID = 0 // remove identifier.
	IPC_SET  = 1 // set `ipc_perm' options.
	IPC_STAT = 2 // get `ipc_perm' options.
	IPC_INFO = 3 // see ipcs.

	/* Special key values.  */
	IPC_PRIVATE = 0 // private key. NOTE: this value is of type __key_t.
)

// ACCESS_RDWR IPC access permission
const ACCESS_RDWR = 0600
