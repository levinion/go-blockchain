package wallet

import (
	"errors"
	"gbc/constant"
	"gbc/util"
	"io/fs"
	"path/filepath"
)

// 钱包地址到别名的映射表
type RefList map[string]string

type IRefList interface {
	Update()
	Store()
	Load()
}

func NewRefList() *RefList {
	ref:=make(RefList)
	return &ref
}

func (r *RefList) Store() {
	util.NewFileDB(constant.RefListPath).Store(r)
}

// 不存在则从本地更新，否则加载并更新
func (r *RefList) Load() *RefList {
	if util.FileIsExit(constant.RefListPath){
		util.NewFileDB(constant.RefListPath).Load(r)
		r.Update()
	}else{
		r.Update()
	}
	return r
}

func (r *RefList) Update() {
	filepath.Walk(constant.WalletPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, ok := (*r)[info.Name()]; !ok {
			(*r)[info.Name()] = "" //若不存在该钱包条目则创建
		}
		return nil
	})
}

func (r *RefList) BindRef(addr string, ref string) {
	(*r)[addr] = ref
}

func (r *RefList) FindRef(ref string) (string, error) {
	for k, v := range *r {
		if v == ref {
			return k, nil
		}
	}
	return "", errors.New("ref not found")
}
