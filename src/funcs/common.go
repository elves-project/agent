package funcs

import (
	"time"
	"unsafe"
	"encoding/json"
	"strconv"
	"sort"
	"crypto/md5"
	"fmt"
	"io"
	"path/filepath"
	"os"
	"io/ioutil"
	"archive/zip"
	"net/http"
	"errors"
)

func GetMD5(str string) string{
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr:= md5Ctx.Sum(nil)
	return fmt.Sprintf("%x", cipherStr)
}

func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func BytesString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func JsonToString(j string)string{
	r,_ := json.Marshal(j)
	return BytesString(r)
}

func Sign(uri string , pastr map[string]string ,Authkey string)( string,string){
	keys := make([]string, len(pastr))
	i := 0
	for k, _ := range pastr {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	str :=""
	for _, k := range keys {
		str = str+k+"="+pastr[k]
		if i !=1 {
			str = str+"&"
		}
		i--
	}
	has := md5.Sum([]byte(uri+str+Authkey))
	//seelog.Debug(uri+str+Authkey)
	md5str := fmt.Sprintf("%x", has)
	return md5str,str

}

func StrToInt(str string) int{
	i,_:=strconv.Atoi(str)
	return i
}

func Unzip(src, dest string) error {

	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		if f.FileInfo().IsDir() {
			path := filepath.Join(dest, f.Name)
			os.MkdirAll(path, f.Mode())
		} else {
			buf := make([]byte, f.UncompressedSize)
			_, err = io.ReadFull(rc, buf)
			if err != nil {
				return err
			}
			path := filepath.Join(dest, f.Name)
			err := ioutil.WriteFile(path, buf, f.Mode())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func Download(url string,fileDir string,fileName string) error {
	if err := os.MkdirAll(fileDir, 0700); err != nil {
		return err
	}else{
		res, err := http.Get(url)
		if res.StatusCode == 200{
			if err != nil {
				return err
			}
			f, err := os.Create(fileDir+"/"+fileName)
			defer f.Close()
			if err != nil {
				return err
			}
			io.Copy(f, res.Body)
		}else{
			return errors.New("download file "+fileName+" fail StatusCode:"+strconv.Itoa(res.StatusCode))
		}
	}
	return nil
}