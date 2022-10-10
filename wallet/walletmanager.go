package wallet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"goblockchain/constcoe"
	"goblockchain/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 管理一台机器上的所有钱包

// key位地址，value为别名
type RefList map[string]string

// 存储
func (r *RefList) Save() {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(r)
	utils.Handle(err)
	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}

// 更新函数，扫描所有的钱包文件
func (r *RefList) Update() {
	err := filepath.Walk(constcoe.Wallets, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileName := f.Name()
		if strings.Compare(fileName[len(fileName)-4:], ".wlt") == 0 {
			_, ok := (*r)[fileName[:len(fileName)-4]]
			if !ok {
				(*r)[fileName[:len(fileName)-4]] = ""
			}
		}
		return nil
	})
	utils.Handle(err)
}

// 加载已保存的reflist
func LoadRefList() *RefList {
	filename := constcoe.WalletsRefList + "ref_list.data"
	var reflist RefList
	if utils.FileExits(filename) {
		fileContent, err := ioutil.ReadFile(filename)
		utils.Handle(err)
		decoder := gob.NewDecoder(bytes.NewBuffer(fileContent))
		err = decoder.Decode(&reflist)
		utils.Handle(err)
	} else {
		reflist = make(RefList)
		reflist.Update()
	}
	return &reflist
}

// 为了方便操作为钱包设置别名
func (r *RefList) BindRef(address, refname string) {
	(*r)[address] = refname
}

// 通过别名取钱包地址
func (r *RefList) FindRef(refname string) (string, error) {
	temp := ""
	for key, val := range *r {
		if val == refname {
			temp = key
			break
		}
	}
	if temp == "" {
		err := errors.New("the refname is not found")
		return temp, err
	}
	return temp, nil
}
