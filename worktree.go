package git

/*
#include <git2.h>

*/
import "C"
import (
	"runtime"
	"unsafe"
)

type Worktree struct {
	ptr  *C.git_worktree
	Repo *Repository
}

type WorktreeAddOptions struct {
	version uint
	lock    int
	ref     *Reference
}

func NewWorktreeAddOptions(version uint, lock int, ref *Reference) (*WorktreeAddOptions, error) {
	return &WorktreeAddOptions{version, lock, ref}, nil
}

func newWorktreeFromC(ptr *C.git_worktree, repo *Repository) *Worktree {
	idx := &Worktree{ptr, repo}
	runtime.SetFinalizer(idx, (*Worktree).Free)
	return idx
}

func (repo *Repository) ListWorktrees() ([]string, error) {
	var r C.git_strarray

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_worktree_list(&r, repo.ptr)
	runtime.KeepAlive(repo)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}
	defer C.git_strarray_free(&r)
	worktrees := makeStringsFromCStrings(r.strings, int(r.count))
	return worktrees, nil
}

func (repo *Repository) NewWorktreeFromSubrepository() (*Worktree, error) {
	var ptr *C.git_worktree

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ecode := C.git_worktree_open_from_repository(&ptr, repo.ptr)
	runtime.KeepAlive(repo)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}
	return newWorktreeFromC(ptr, repo), nil
}

func (repo *Repository) LookupWorktree(name string) (*Worktree, error) {
	var ptr *C.git_worktree

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	ecode := C.git_worktree_lookup(&ptr, repo.ptr, cname)
	runtime.KeepAlive(repo)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
	}
	return newWorktreeFromC(ptr, repo), nil
}

func (repo *Repository) AddWorktree(name string, destdir string, options *WorktreeAddOptions) (*Worktree, error) {
	var ptr *C.git_worktree

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cdestdir := C.CString(destdir)
	defer C.free(unsafe.Pointer(cdestdir))

	ecode := C.git_worktree_add(&ptr, repo.ptr, cname, cdestdir, &C.git_worktree_add_options{C.uint(options.version), C.int(options.lock), options.ref.ptr})
	runtime.KeepAlive(repo)
	if ecode < 0 {
		return nil, MakeGitError(ecode)
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
