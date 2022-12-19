package clock

import "strconv"

type SettingType interface {
    GetName() string
    GetType() string
    GetDisplayValue() string
    GetSaveValue() interface{}
}

type IntRangeSetting struct {
    val  int
    max  int
    min  int
    name string
}

func (i IntRangeSetting) GetName() string {
    return i.name
}

func (i IntRangeSetting) GetType() string {
    return "RangeSetting"
}

func (i IntRangeSetting) GetDisplayValue() string {
    return strconv.Itoa(i.val)
}

func (i IntRangeSetting) GetSaveValue() interface{} {
    return i.val
}

func (i IntRangeSetting) GetValue() int {
    return i.val
}

func (i IntRangeSetting) GetMaxValue() int {
    return i.max
}

func (i IntRangeSetting) GetMinValue() int {
    return i.min
}

func NewIntRangeSetting(name string, def, min, max int) SettingType {
    return &IntRangeSetting{
        name: name,
        val:  def,
        min:  min,
        max:  max,
    }
}
