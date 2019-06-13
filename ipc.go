package ipc

// https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/bits/ipc.h.html

const (
	/* Mode bits for `msgget', `semget', and `shmget'.  */
	IPC_CREAT  = 01000 // create key if key does not exist.
	IPC_EXCL   = 02000 // fail if key exists.
	IPC_NOWAIT = 04000 // return error on wait.

	/* Permission flag for `msgget', `semget', and `shmget'.  */
	IPC_R  = 0400 // read
	IPC_W  = 0200 // write
	IPC_RW = 0600 // read and write

	/* Control commands for `msgctl', `semctl', and `shmctl'.  */
	IPC_RMID = 0 // remove identifier.
	IPC_SET  = 1 // set `ipc_perm' options.
	IPC_STAT = 2 // get `ipc_perm' options.
	IPC_INFO = 3 // see ipcs.

	/* Special key values.  */
	IPC_PRIVATE = 0 // private key. NOTE: this value is of type __key_t.
)
