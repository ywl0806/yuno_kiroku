package utils

import (
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
)

func ConvertStruct(dst any, src any) error {
	err := copier.Copy(dst, src)
	copier.CopyWithOption(dst, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
	if err != nil {
		return err
	}
	return nil
}

func ConvertToInt32(src any) (int32, error) {
	return cast.ToInt32E(src)
}

func ConvertToBool(src any) (bool, error) {
	return cast.ToBoolE(src)
}
