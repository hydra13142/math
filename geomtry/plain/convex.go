package plain

const (
	Outter = 0 // 外部
	OnLine = 1 // 边上
	OnApex = 3 // 在顶点上
	Inside = 4 // 内部
)

const (
	Outside = 0 // 直线在多边形外部
	Tangent = 1 // 直线与多边形相切
	Through = 2 // 直线穿越多边形
)

// 顺时针排布点的凸多边形
type Convex []Spot

// 本函数设计为将随机分布的点用尽量小的凸多边形全部包含在内
// 位于边界和内部都算包含在内，本函数用于生成凸多边形
// 会将位于边界上的点（非顶点）和凸多边形内的点丢弃
func NewConvex(ps []Spot) Convex {
	type link struct {
		cont       Spot
		last, next *link
	}
	var (
		p, q *link
		a, b Spot
		s    Segment
	)
	i, l := 0, len(ps)
	if l < 3 {
		return nil
	}
	a = ps[0]
	for i = 1; i < l && ps[i] == a; i++ {
	}
	if i+2 > l {
		return nil
	}
	b = ps[i]
	p = &link{cont: a}
	q = &link{cont: b}
	p.next, p.last = q, q
	q.next, q.last = p, p
	for ; i < l; i++ {
	inner:
		for q = p; ; {
			s.Set(q.cont, q.next.cont)
			switch u := s.Location(ps[i]); {
			case u == WM&HM:
				break inner
			case u&WR == 0:
				t := q
				for ; ; t = t.last {
					s.Set(t.last.cont, t.cont)
					if s.Location(ps[i])&(WL|HB) == 0 {
						break
					}
				}
				for q = q.next; ; q = q.next {
					s.Set(q.cont, q.next.cont)
					if s.Location(ps[i])&(WL|HT) == 0 {
						break
					}
				}
				t.next.last = nil
				q.last.next = nil
				p = new(link)
				p.cont = ps[i]
				t.next, p.last = p, t
				p.next, q.last = q, p
				break inner
			default:
				if q = q.next; q == p {
					break inner
				}
			}
		}
	}
	for i, q = 1, p.next; q != p; i, q = i+1, q.next {
	}
	ans := make(Convex, i)
	for i--; i >= 0; i, p = i-1, q {
		q = p.last
		ans[i] = p.cont
		p.last, p.next = nil, nil
	}
	if len(ans) < 3 {
		return nil
	}
	return ans
}

// 点相对凸多边形的位置
func (c Convex) Location(p Spot) int {
	var (
		s Segment
		t int
	)
	l := len(c)
	if l < 3 {
		return Outter
	}
	for i, j := l-1, 0; j < l; i, j = j, j+1 {
		s.Set(c[i], c[j])
		t = s.Location(p)
		if t&WL != 0 {
			return Outter
		}
		if t == WM|HM {
			if p == c[i] || p == c[j] {
				return OnApex
			}
			return OnLine
		}
	}
	return Inside
}

// 线段相对凸多边形的位置
func (c Convex) Cross(s Segment) (int, []Spot) {
	var (
		p, q, r int
		a       Spot
		x       = make([]Spot, 2)
	)
	inter := func(s, l Segment) Spot {
		k := OuterProduct(AimTo(s.O, l.O), l.K) / OuterProduct(s.K, l.K)
		return Spot{s.O.X + s.K.I*k, s.O.Y + s.K.J*k}
	}
	a = c[len(c)-1]
	p = s.Location(a) & (WM | WL | WR)
	for _, b := range c {
		q = s.Location(b) & (WM | WL | WR)
		switch p {
		case WM:
			switch q {
			case WL:
				x[0] = a
			case WR:
				x[1] = a
			}
		case WL:
			switch q {
			case WM:
				x[1] = b
			case WR:
				x[1] = inter(Segment{a, AimTo(a, b)}, s)
			}
		case WR:
			switch q {
			case WM:
				x[0] = b
			case WL:
				x[0] = inter(Segment{a, AimTo(a, b)}, s)
			}
		}
		a, r = b, r|p
	}
	if r&WL != 0 && r&WR != 0 {
		switch both(&s, &x[0], &x[1]) {
		case 0:
			x = nil
		case 1:
			x = x[:1]
		case 2:
		}
		return Through, x
	}
	if r&WM != 0 || r&WM != 0 {
		switch both(&s, &x[0], &x[1]) {
		case 0:
			x = nil
		case 1:
			x = x[:1]
		case 2:
			if x[0] == x[1] {
				x = x[:1]
			}
		}
		return Tangent, x
	}
	return Outside, nil
}

// 两个凸多边形的交集凸多边形
func (p Convex) And(q Convex) (r Convex) {
	var (
		u, v, w int
		i, j    int
		x, y    int
		t       bool
		o       Spot
		s       Segment
	)
	i, j = len(p)-1, 0
	if q.Location(p[i]) != Outside {
		r = append(r, p[i])
		x, y = len(q)-1, 0
		t = (u != OnApex)
		goto walk
	}
	for j < len(p) {
		if u = q.Location(p[j]); u != Outside {
			r = append(r, p[j])
			x, y = len(q)-1, 0
			i, j = j, j+1
			t = (u != OnApex)
			goto walk
		}
		s.Set(p[i], p[j])
		u = s.Location(q[x])
		w = 0
		for x, y = len(q)-1, 0; y < len(q); x, y = y, y+1 {
			v = s.Location(q[y])
			if u&WR == 0 && v&WR != 0 {
				kl := AimTo(q[x], q[y])
				sl := AimTo(s.O, q[x])
				k := OuterProduct(sl, kl) / OuterProduct(s.K, kl)
				f := Spot{s.O.X + s.K.I*k, s.O.Y + s.K.J*k}
				g := AimTo(f, s.O)
				if InnerProduct(g, g.Add(s.K)) <= 0 {
					r = append(r, f)
					i, j = j, (j+1)%len(p)
					t = true
					p, i, j, q, x, y = q, x, y, p, i, j
					goto walk
				}
			}
			u, w = v, w|v
		}
		if w&WR == 0 {
			return nil
		}
		i, j = j, j+1
	}
	r = make(Convex, len(q))
	copy(r, q)
	return r
walk:
	m, n := true, j
	for {
		if u = q.Location(p[j]); u != Outside {
			r = append(r, p[j])
			i, j = j, (j+1)%len(p)
			t = (u != OnApex)
		} else {
			s.Set(p[i], p[j])
			u = s.Location(q[x])
			for {
				v = s.Location(q[y])
				if u&WR == 0 && v&WR != 0 {
					if u&WM != 0 {
						o = q[x]
					} else {
						kl := AimTo(q[x], q[y])
						sl := AimTo(s.O, q[x])
						k := OuterProduct(sl, kl) / OuterProduct(s.K, kl)
						o = Spot{s.O.X + s.K.I*k, s.O.Y + s.K.J*k}
					}
					if t {
						r = append(r, o)
					}
					i, j = j, (j+1)%len(p)
					t = true
					p, i, j, q, x, y = q, x, y, p, i, j
					m = !m
					break
				}
				u, x, y = v, y, (y+1)%len(q)
			}
		}
		if m && j == n {
			break
		}
	}
	if r[0] == r[len(r)-1] {
		r = r[:len(r)-1]
	}
	if len(r) < 3 {
		r = nil
	}
	return r
}

// 求多边形c的面积
func (c Convex) Area() (s float64) {
	area := func(A, B, C Spot) float64 {
		return OuterProduct(AimTo(A, C), AimTo(A, B)) / 2
	}
	s = 0
	for i := 2; i < len(c); i++ {
		s += area(c[0], c[i-1], c[i])
	}
	return s
}
