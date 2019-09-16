// Copyright 2019 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package linux

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/ZHYfeng/2018_dependency/03-syzkaller/prog"
)

const (
	USB_DEVICE_ID_MATCH_VENDOR = 1 << iota
	USB_DEVICE_ID_MATCH_PRODUCT
	USB_DEVICE_ID_MATCH_DEV_LO
	USB_DEVICE_ID_MATCH_DEV_HI
	USB_DEVICE_ID_MATCH_DEV_CLASS
	USB_DEVICE_ID_MATCH_DEV_SUBCLASS
	USB_DEVICE_ID_MATCH_DEV_PROTOCOL
	USB_DEVICE_ID_MATCH_INT_CLASS
	USB_DEVICE_ID_MATCH_INT_SUBCLASS
	USB_DEVICE_ID_MATCH_INT_PROTOCOL
	USB_DEVICE_ID_MATCH_INT_NUMBER

	BytesPerUsbID = 17
)

type UsbDeviceID struct {
	MatchFlags         uint16
	IDVendor           uint16
	IDProduct          uint16
	BcdDeviceLo        uint16
	BcdDeviceHi        uint16
	BDeviceClass       uint8
	BDeviceSubClass    uint8
	BDeviceProtocol    uint8
	BInterfaceClass    uint8
	BInterfaceSubClass uint8
	BInterfaceProtocol uint8
	BInterfaceNumber   uint8
}

func (arch *arch) generateUsbDeviceDescriptor(g *prog.Gen, typ0 prog.Type, old prog.Arg) (
	arg prog.Arg, calls []*prog.Call) {

	if old == nil {
		arg = g.GenerateSpecialArg(typ0, &calls)
	} else {
		arg = old
		calls = g.MutateArg(arg)
	}
	if g.Target().ArgContainsAny(arg) {
		return
	}

	totalIds := len(usbIds) / BytesPerUsbID
	idNum := g.Rand().Intn(totalIds)
	base := usbIds[idNum*BytesPerUsbID : (idNum+1)*BytesPerUsbID]

	p := strings.NewReader(base)
	var id UsbDeviceID
	if binary.Read(p, binary.LittleEndian, &id) != nil {
		panic("not enough data to read")
	}

	if (id.MatchFlags & USB_DEVICE_ID_MATCH_VENDOR) == 0 {
		id.IDVendor = uint16(g.Rand().Intn(0xffff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_PRODUCT) == 0 {
		id.IDProduct = uint16(g.Rand().Intn(0xffff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_DEV_LO) == 0 {
		id.BcdDeviceLo = 0x0
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_DEV_HI) == 0 {
		id.BcdDeviceHi = 0xffff
	}
	bcdDevice := id.BcdDeviceLo + uint16(g.Rand().Intn(int(id.BcdDeviceHi-id.BcdDeviceLo)+1))
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_DEV_CLASS) == 0 {
		id.BDeviceClass = uint8(g.Rand().Intn(0xff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_DEV_SUBCLASS) == 0 {
		id.BDeviceSubClass = uint8(g.Rand().Intn(0xff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_DEV_PROTOCOL) == 0 {
		id.BDeviceProtocol = uint8(g.Rand().Intn(0xff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_INT_CLASS) == 0 {
		id.BInterfaceClass = uint8(g.Rand().Intn(0xff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_INT_SUBCLASS) == 0 {
		id.BInterfaceSubClass = uint8(g.Rand().Intn(0xff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_INT_PROTOCOL) == 0 {
		id.BInterfaceProtocol = uint8(g.Rand().Intn(0xff + 1))
	}
	if (id.MatchFlags & USB_DEVICE_ID_MATCH_INT_NUMBER) == 0 {
		id.BInterfaceNumber = uint8(g.Rand().Intn(0xff + 1))
	}

	patchGroupArg(arg, 7, "idVendor", uint64(id.IDVendor))
	patchGroupArg(arg, 8, "idProduct", uint64(id.IDProduct))
	patchGroupArg(arg, 9, "bcdDevice", uint64(bcdDevice))
	patchGroupArg(arg, 3, "bDeviceClass", uint64(id.BDeviceClass))
	patchGroupArg(arg, 4, "bDeviceSubClass", uint64(id.BDeviceSubClass))
	patchGroupArg(arg, 5, "bDeviceProtocol", uint64(id.BDeviceProtocol))

	configArg := arg.(*prog.GroupArg).Inner[14].(*prog.GroupArg).Inner[0]
	interfaceArg := configArg.(*prog.GroupArg).Inner[8].(*prog.GroupArg).Inner[0]

	patchGroupArg(interfaceArg, 5, "bInterfaceClass", uint64(id.BInterfaceClass))
	patchGroupArg(interfaceArg, 6, "bInterfaceSubClass", uint64(id.BInterfaceSubClass))
	patchGroupArg(interfaceArg, 7, "bInterfaceProtocol", uint64(id.BInterfaceProtocol))
	patchGroupArg(interfaceArg, 2, "bInterfaceNumber", uint64(id.BInterfaceNumber))

	return
}

func patchGroupArg(arg prog.Arg, index int, field string, value uint64) {
	fieldArg := arg.(*prog.GroupArg).Inner[index].(*prog.ConstArg)
	if fieldArg.Type().FieldName() != field {
		panic(fmt.Sprintf("bad field, expected %v, found %v", field, fieldArg.Type().FieldName()))
	}
	fieldArg.Val = value
}
