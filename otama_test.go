package otama

import (
    "bytes"
    "fmt"
    "os"
    "regexp"
    "testing"
    "path/filepath"
)

func setup() (*Otama) {
    os.Mkdir("examples/data", 0755)
    os.Chdir("examples")
    return new(Otama)
}

func TestLibOtamaVersion(t *testing.T) {
    r, _ := regexp.Compile("([0-9]+).([0-9]+)*([0-9]+)")
    if r.MatchString(LIBOTAMA_VERSION) != true {
        t.Errorf("fail")
    }
}

func TestOtamaOpen(t *testing.T) {
    o := setup()
    o.Open("test.conf")
}

func TestOtamaClose(t *testing.T) {
    o := setup()
    o.Open("test.conf")
    o.Close()
}

func TestOtamaCreateDatabase(t *testing.T) {
    o := setup()
    o.Open("test.conf")
    o.CreateDatabase()
}

func TestOtamaDropDatabase(t *testing.T) {
    o := setup()
    o.Open("test.conf")
    o.CreateDatabase()
    o.DropDatabase()
}

func create_dataset(t *testing.T, o *Otama, pwd string) (ids []map[string]string) {
    buf := bytes.NewBufferString(pwd)
    buf.WriteString("/image")

    // create database
    filepath.Walk(buf.String(), func (path string, info os.FileInfo, err error) error {
        if info == nil || info.IsDir() {
            return nil
        }

        id, err := o.Insert(path)
        if err != nil {
            t.Error(err)
        }

        ids = append(ids, map[string]string{"id": id, "filename": path})
        return nil
    })

    err := o.Pull()
    if err != nil {
        t.Error(err)
        return nil
    }

    return ids
}

func TestOtamaInsertAndSearch(t *testing.T) {
    o := setup()
    o.Open("test.conf")
    o.CreateDatabase()

    pwd, _ := os.Getwd()
    var _ []map[string]string = create_dataset(t, o, pwd)

    buf := bytes.NewBufferString(pwd)
    buf.WriteString("/image/lena.jpg")
    results, err := o.Search(10, buf.String())
    if err != nil {
        t.Error(err)
    }
    if len(results) == 0 {
        t.Errorf("result not found")
    }
    for result := range results {
        fmt.Println(fmt.Sprintf("key=%s, sim=%0.3f", results[result].Id, results[result].Similarity))
    }
}

func TestOtamaExists(t *testing.T) {
    o := setup()
    o.Open("test.conf")
    o.CreateDatabase()

    pwd, _ := os.Getwd()
    var _ []map[string]string = create_dataset(t, o, pwd)

    ret, _ := o.Exists("3444d453367af67e18dd20f99cdb4d90397a1fa9")  // lena.jpg
    if ret != true {
        t.Errorf("not exist")
    }
}
