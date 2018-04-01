package strconv

import (
	"errors"
	"strings"
)

var (
	ErrRune   = errors.New("Unrecognised rune")
	ErrEmpty  = errors.New("Empty string")
	ErrSyntax = errors.New("Syntax error")
	ErrRange  = errors.New("Out of range")
	
	cs = "零壹贰叁肆伍陆柒捌玖拾佰仟万亿"
	cb = "零一二三四五六七八九十百千万亿"
	is = []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九", "", "十", "百", "千", "万", "亿"}
	ib = []string{"零", "壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖", "", "拾", "佰", "仟", "万", "亿"}
)

// 对汉字字符串进行一定的翻译，方便后续的处理，并对习惯用法进行修复
func parse(s, list string) ([]byte, error) {
	if s == "" {
		return nil, ErrEmpty
	}
	t := make([]byte, 1, len(s)/3+1)
outer:
	for _, r := range s {
		for j, c := range []rune(list) {
			if r == c {
				t = append(t, byte(j))
				continue outer
			}
		}
		return nil, ErrRune
	}
	if t[1] == 10 {
		t[0] = 1
	} else { // 对应汉字‘十’
		t = t[1:]
	}
	return t, nil
}

// 万以下级别（0 ~ 9999）的汉字到整数转换，这是可重复的最小转换规则单位。
func ctoi1(s []byte) int64 {
	t, k, p, v := 0, 0, 4, false
	for _, c := range s {
		switch c {
		case 10: //‘十’
			if k == 0 || p < 2 {
				return -1
			}
			t, k, p, v = t+k*10, 0, 1, false
		case 11: //‘百’
			if k == 0 || p < 3 {
				return -1
			}
			t, k, p, v = t+k*100, 0, 2, false
		case 12: //‘千’
			if v || k == 0 || p < 4 {
				return -1
			}
			t, k, p, v = t+k*1000, 0, 3, false
		case 0: //‘零’
			if v || k != 0 {
				return -1
			}
			p, v = p-1, true
		default: //‘一’—‘九’
			if k != 0 {
				return -1
			}
			k = int(c)
		}
	}
	if k != 0 {
		if v || p == 1 {
			return int64(t + k)
		}
		return -1
	}
	return int64(t)
}

// 万级别（1e4 ~ 1e8-1）的汉字到整数转换
func ctoi2(s []byte) int64 {
	var h, l int64
	t := 0
	for i, c := range s {
		if c == 13 {
			if t != 0 || i == 0 {
				return -1
			}
			t = i
		}
	}
	if t == 0 {
		return ctoi1(s)
	}
	if h = ctoi1(s[:t]); h < 0 {
		return -1
	}
	if l = ctoi1(s[t+1:]); l < 0 {
		return -1
	}
	return h*1e4 + l
}

// 亿级别（1e8 ~ 1e16-1）的汉字到整数转换
func ctoi3(s []byte) int64 {
	var h, l int64
	t := 0
	for i, c := range s {
		if c == 14 {
			if t != 0 || i == 0 {
				return -1
			}
			t = i
		}
	}
	if t == 0 {
		return ctoi2(s)
	}
	if h = ctoi2(s[:t]); h < 0 {
		return -1
	}
	if l = ctoi2(s[t+1:]); l < 0 {
		return -1
	}
	return h*1e8 + l
}

// 万以下级别（0 ~ 9999）的整数到汉字转换
func itoc1(n int64) []byte {
	var i, j int
	if n == 0 {
		return []byte{0}
	}
	c := make([]byte, 4)
	for i = 0; i < 4; i++ {
		c[3-i] = byte(n % 10)
		n /= 10
	}
	for j = 3; j > 0 && c[j] == 0; j-- {
		c[j] = 10
	}
	for i = 0; i < j && c[i] == 0; i++ {
		c[i] = 10
	}
	if j == 3 && i == 0 && c[1] == 0 && c[2] == 0 {
		c[2] = 10
	}
	if i != 0 {
		i--
		c[i] = 0
	}
	for k := 0; k < 4; k++ {
		c[k]|= byte(3-k)<<4
	}
	return c
}

// 万级别（1e4 ~ 1e8-1）的整数到汉字转换
func itoc2(n int64) []byte {
	h, l := n/1e4, n%1e4
	if h == 0 {
		return itoc1(l)
	}
	H := append(itoc1(h), 14)
	if l == 0 {
		return H
	}
	return append(H, itoc1(l)...)
}

// 亿级别（1e8 ~ 1e16-1）的整数到汉字转换
func itoc3(n int64) []byte {
	h, l := n/1e8, n%1e8
	if h == 0 {
		return itoc2(l)
	}
	H := append(itoc2(h), 15)
	if l == 0 {
		return H
	}
	return append(H, itoc2(l)...)
}

// 将修正过的标记串翻译成汉字数字+单位的数组
func format(s []byte, list []string) []string {
	t := make([]string, len(s))
	for i, c := range s {
		switch k:=c&15; k {
		case 0, 10, 14, 15:
			t[i] = list[k]
		default:
			t[i] = list[k] + list[(c>>4)+10]
		}
	}
	return t
}

// 返回一个小写汉字数字字符串表示的数字和可能的错误
func Ctoi(s string) (int64, error) {
	if t, e := parse(s, cs); e != nil {
		return 0, e
	} else if len(t) > 1 && t[0] == 0 {
		return 0, ErrSyntax
	} else if x := ctoi3(t); x < 0 {
		return 0, ErrSyntax
	} else {
		return x, nil
	}
}

// 返回一个数字的小写汉字数字字符串表示和可能的错误
func Itoc(n int64) (string, error) {
	if n < 0 || n >= 1e16 {
		return "", ErrRange
	}
	s := itoc3(n)
	if len(s) > 1 && s[0]&15 == 0 {
		s = s[1:]
	}
	t := format(s, is)
	if s[0] == 17 {
		t[0] = is[11]
	}
	return strings.Join(t, ""), nil
}


// 返回一个大写汉字数字字符串表示的数字和可能的错误
func Ctoi2(s string) (int64, error) {
	if t, e := parse(s, cb); e != nil {
		return 0, e
	} else if len(t) > 1 && t[0] == 0 {
		return 0, ErrSyntax
	} else if x := ctoi3(t); x < 0 {
		return 0, ErrSyntax
	} else {
		return x, nil
	}
}

// 返回一个数字的大写汉字数字字符串表示和可能的错误
func Itoc2(n int64) (string, error) {
	if n < 0 || n >= 1e16 {
		return "", ErrRange
	}
	s := itoc3(n)
	if len(s) > 1 && s[0]&15 == 0 {
		s = s[1:]
	}
	t := format(s, ib)
	if s[0] == 17 {
		t[0] = ib[11]
	}
	return strings.Join(t, ""), nil
}