// Copyright 2017 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package prog

func (p *Prog) Clone() *Prog {
	p1 := &Prog{
		Target: p.Target,
		Calls:  make([]*Call, len(p.Calls)),
	}
	newargs := make(map[*ResultArg]*ResultArg)
	for ci, c := range p.Calls {
		c1 := new(Call)
		c1.Meta = c.Meta
		if c.Ret != nil {
			c1.Ret = clone(c.Ret, newargs).(*ResultArg)
		}
		c1.Args = make([]Arg, len(c.Args))
		for ai, arg := range c.Args {
			c1.Args[ai] = clone(arg, newargs)
		}
		p1.Calls[ci] = c1
	}
	p1.debugValidate()
	return p1
}

func (p *Prog) CloneWithUncover() *Prog {
	p1 := &Prog{
		Target:       p.Target,
		Calls:        make([]*Call, len(p.Calls)),
		Uncover:      make(map[int]*Uncover),
		UncoverIdx:   p.UncoverIdx,
		WriteAddress: p.WriteAddress,
		Sig:          p.Sig,
	}
	newargs := make(map[*ResultArg]*ResultArg)
	for ci, c := range p.Calls {
		c1 := new(Call)
		c1.Meta = c.Meta
		if c.Ret != nil {
			c1.Ret = clone(c.Ret, newargs).(*ResultArg)
		}
		c1.Args = make([]Arg, len(c.Args))
		for ai, arg := range c.Args {
			c1.Args[ai] = clone(arg, newargs)
		}
		p1.Calls[ci] = c1
	}
	p1.debugValidate()

	for ui, u := range p.Uncover {
		u1 := new(Uncover)
		u1.UncoveredAddress = u.UncoveredAddress
		u1.Idx = u.Idx

		u1.RelatedAddress = make([]*RelatedAddresses, len(u.RelatedAddress))
		for ai, address := range u.RelatedAddress {
			ra := &RelatedAddresses{
				RelatedAddress: address.RelatedAddress,
				Prio:           address.Prio,
				Repeat:         address.Repeat,
			}
			for _, call := range address.RelatedCalls {
				ra.RelatedCalls = append(ra.RelatedCalls, call)
			}
			for _, prog := range address.RelatedProgs {
				ra.RelatedProgs = append(ra.RelatedProgs, prog.Clone())
			}
			u1.RelatedAddress[ai] = ra
		}
		p1.Uncover[ui] = u1
	}

	return p1
}

func clone(arg Arg, newargs map[*ResultArg]*ResultArg) Arg {
	var arg1 Arg
	switch a := arg.(type) {
	case *ConstArg:
		a1 := new(ConstArg)
		*a1 = *a
		arg1 = a1
	case *PointerArg:
		a1 := new(PointerArg)
		*a1 = *a
		arg1 = a1
		if a.Res != nil {
			a1.Res = clone(a.Res, newargs)
		}
	case *DataArg:
		a1 := new(DataArg)
		*a1 = *a
		a1.data = append([]byte{}, a.data...)
		arg1 = a1
	case *GroupArg:
		a1 := new(GroupArg)
		*a1 = *a
		arg1 = a1
		a1.Inner = make([]Arg, len(a.Inner))
		for i, arg2 := range a.Inner {
			a1.Inner[i] = clone(arg2, newargs)
		}
	case *UnionArg:
		a1 := new(UnionArg)
		*a1 = *a
		arg1 = a1
		a1.Option = clone(a.Option, newargs)
	case *ResultArg:
		a1 := new(ResultArg)
		*a1 = *a
		arg1 = a1
		if a1.Res != nil {
			r := newargs[a1.Res]
			a1.Res = r
			if r.uses == nil {
				r.uses = make(map[*ResultArg]bool)
			}
			r.uses[a1] = true
		}
		a1.uses = nil // filled when we clone the referent
		newargs[a] = a1
	default:
		panic("bad arg kind")
	}
	return arg1
}
