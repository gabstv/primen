package ggui

import (
	"strconv"
	"strings"
)

type UIAnchor struct {
	offtop, offleft, offwidth, offheight float64
	ptop, pleft, pwidth, pheight         float64
	ltop, lleft, lwidth, lheight         bool
}

func (a *UIAnchor) Locked() bool {
	return a.ltop || a.lleft || a.lwidth || a.lheight
}

//go:generate ecsgen -n UIAnchor -p core -o uianchor_component.go --component-tpl --vars "UUID=D7052642-7A15-44B0-81AA-4B84F49DA9E0"

func parseUIPosition(v string) (pct float64, locked bool, offset float64) {
	v = strings.ReplaceAll(v, " ", "")
	pcti := strings.IndexByte(v, '%')
	if pcti == -1 {
		locked = false
		offset, _ = strconv.ParseFloat(v, 64)
		return
	}
	pct, _ = strconv.ParseFloat(v[:pcti], 64)
	locked = true
	offset, _ = strconv.ParseFloat(v[pcti+1:], 64)
	return
}

func parseUIPosX(v string, a *UIAnchor) float64 {
	pct, locked, offset := parseUIPosition(v)
	if locked {
		a.ltop = true
		a.ptop = pct
		a.offtop = offset
	}
	return offset
}

func parseUIPosY(v string, a *UIAnchor) float64 {
	pct, locked, offset := parseUIPosition(v)
	if locked {
		a.lleft = true
		a.pleft = pct
		a.offleft = offset
	}
	return offset
}

func parseUIPosWidth(v string, a *UIAnchor) float64 {
	pct, locked, offset := parseUIPosition(v)
	if locked {
		a.lwidth = true
		a.pwidth = pct
		a.offwidth = offset
	}
	return offset
}

func parseUIPosHeight(v string, a *UIAnchor) float64 {
	pct, locked, offset := parseUIPosition(v)
	if locked {
		a.lheight = true
		a.pheight = pct
		a.offheight = offset
	}
	return offset
}
