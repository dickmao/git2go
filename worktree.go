package git

/*
#include <git2.h>

extern int _go_git_worktree_add(git_worktree **out, git_repository *repo, const char *name, const char *worktree);

*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Worktree struct {
	ptr  *C.git_worktree
	repo *Repository
}

func newWorktreeFromC(ptr *C.git_worktree, repo *Repository) *Worktree {
	idx := &Worktree{ptr, repo}
	runtime.SetFinalizer(idx, (*Worktree).Free)
	return idx
}

func ExistingWorktree(repo *Repository) (*Worktree, error) {
	var ptr *C.git_worktree

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := C.git_worktree_open_from_repository(&ptr, repo.ptr); err < 0 {
		return nil, MakeGitError(err)
	}

	return newWorktreeFromC(ptr, repo), nil
}

func AddWorktree(repo *Repository, name string, worktree string) (*Worktree, error) {
	var ptr *C.git_worktree

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cworktree := C.CString(worktree)
	defer C.free(unsafe.Pointer(cworktree))

	err := C._go_git_worktree_add(&ptr, repo.ptr, cname, cworktree)
	if err < 0 {
		return nil, MakeGitError(err)
	}
	return newWorktreeFromC(ptr, repo), nil
}

// Path returns the worktree's path on disk or an empty string if it
// exists only in memory.
func (v *Worktree) Path() string {
	ret := C.GoString(C.git_worktree_path(v.ptr))
	runtime.KeepAlive(v)
	return ret
}

func (v *Worktree) Free() {
	runtime.SetFinalizer(v, nil)
	C.git_worktree_free(v.ptr)
}
