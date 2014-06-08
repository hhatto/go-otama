package otama

/*
#cgo CFLAGS: -c -g -O3 -Wall
#cgo LDFLAGS: -lotama
#include "otama.h"
*/
import "C"

import (
    "fmt"
    "bytes"
    "errors"
)


const (
    LIBOTAMA_VERSION = C.OTAMA_VERSION
)

type Otama struct {
    otama *C.otama_t
}

type OtamaResult struct {
    Id string
    Similarity float64
}

type OtamaFeatureRaw struct {
    base Otama
    raw *C.otama_feature_raw_t
}


func get_status_message(ret C.otama_status_t) (string) {
    return C.GoString(C.otama_status_message(ret))
}

func variant2goobj(v *C.otama_variant_t) (float64) {
    switch C.otama_variant_type(v) {
    case C.OTAMA_VARIANT_TYPE_INT:
        return float64(C.otama_variant_to_int(v))
    case C.OTAMA_VARIANT_TYPE_FLOAT:
        return float64(C.otama_variant_to_float(v))
    case C.OTAMA_VARIANT_TYPE_HASH:
        var keys *C.otama_variant_t = C.otama_variant_hash_keys(v)
        return float64(variant2goobj(C.otama_variant_hash_at2(v, C.otama_variant_array_at(keys, 0))))
    case C.OTAMA_VARIANT_TYPE_STRING:
        fmt.Println("string")
        break
    case C.OTAMA_VARIANT_TYPE_ARRAY:
        fmt.Println("array")
        break
    case C.OTAMA_VARIANT_TYPE_NULL:
        fmt.Println("null")
        break
    default:
        fmt.Println("other")
    }
    return 1
}

func goobj2variant(o map[string]string, v C.otama_variant_t) {
}

func make_results(raw_results *C.otama_result_t) ([]OtamaResult) {
    var hexid[C.OTAMA_ID_HEXSTR_LEN] C.char
    var result_num = int(C.otama_result_count(raw_results))
    var results []OtamaResult
    var value *C.otama_variant_t

    if result_num != 0 {
        results = make([]OtamaResult, result_num)
    }

    for i := 0; i < result_num; i++ {
        value = C.otama_result_value(raw_results, C.int(i))
        C.otama_id_bin2hexstr(&hexid[0], C.otama_result_id(raw_results, C.int(i)))

        results[i] = OtamaResult{Id: C.GoStringN(&hexid[0], C.OTAMA_ID_HEXSTR_LEN),
                                 Similarity: variant2goobj(value)}
    }

    return results
}

func (o *Otama) Open(config string) (err error) {
    ret := C.otama_open(&o.otama, C.CString(config))
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("otama_open: ")
        buf.WriteString(get_status_message(ret))
        err = errors.New(buf.String())
    }

    return err
}

func (o *Otama) Close() {
    C.otama_close(&o.otama)
}

func (o *Otama) CreateDatabase() (err error) {
    ret := C.otama_create_database(o.otama)
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("otama_create_database: ")
        buf.WriteString(get_status_message(ret))
        err = errors.New(buf.String())
    }

    return err
}

func (o *Otama) DropDatabase() (err error) {
    ret := C.otama_drop_database(o.otama)
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("otama_drop_database: ")
        buf.WriteString(get_status_message(ret))
        err = errors.New(buf.String())
    }

    return err
}

func (o *Otama) Insert(filename string) (id string, err error) {
    var _id C.otama_id_t
    var hexid[C.OTAMA_ID_HEXSTR_LEN] C.char

    ret := C.otama_insert_file(o.otama, &_id, C.CString(filename))
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("otama_insert_file: ")
        buf.WriteString(get_status_message(ret))
        err = errors.New(buf.String())
    }

    C.otama_id_bin2hexstr(&hexid[0], &_id)
    return C.GoStringN(&hexid[0], C.OTAMA_ID_HEXSTR_LEN), err
}

func (o *Otama) Pull() (err error) {
    ret := C.otama_pull(o.otama)
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("otama_pull: ")
        buf.WriteString(get_status_message(ret))
        return errors.New(buf.String())
    }

    return nil
}

func (o *Otama) Search(num int, filename string) (results []OtamaResult, err error) {
    var otama_results *C.otama_result_t

    ret := C.otama_search_file(o.otama, &otama_results, C.int(num), C.CString(filename))
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("otama_search_file: ")
        buf.WriteString(get_status_message(ret))
        err = errors.New(buf.String())
    }

    results = make_results(otama_results)
    C.otama_result_free(&otama_results);

    return results, err
}

func (o *Otama) Exists(id string) (r bool, err error) {
    var otama_id C.otama_id_t
    var ret C.otama_status_t
    var result C.int

    ret = C.otama_id_hexstr2bin(&otama_id, C.CString(id));
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("Exists otama_id_hexstr2bin: ")
        buf.WriteString(get_status_message(ret))
        err = errors.New(buf.String())
        return false, err
    }

    ret = C.otama_exists(o.otama, &result, &otama_id)
    if ret != C.OTAMA_STATUS_OK {
        buf := bytes.NewBufferString("Exists otama_exists: ")
        buf.WriteString(get_status_message(ret))
        err = errors.New(buf.String())
        return false, err
    }

    if int(result) == 0 { return false, err }
    return true, err
}
