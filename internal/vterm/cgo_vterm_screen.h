#ifndef __CGO_VTERM_SCREEN_H__
#define __CGO_VTERM_SCREEN_H__

#include <stdbool.h>
#include <vterm.h>

typedef struct {
  VTermPos cursor_pos;
  int cursor_visible;
  int cursor_blink;
  int cursor_shape;
} CGoVTermScreenUser;

static int cgo_vterm_screen_user_movecursor(VTermPos pos, VTermPos oldpos,
                                            int visible, void *user) {
  CGoVTermScreenUser *u = user;
  u->cursor_pos = pos;
  return 1;
}

static int cgo_vterm_screen_user_settermprop(VTermProp prop, VTermValue *val,
                                             void *user) {
  CGoVTermScreenUser *u = user;

  switch (prop) {
  case VTERM_PROP_CURSORVISIBLE:
    u->cursor_visible = val->boolean;
    break;
  case VTERM_PROP_CURSORBLINK:
    u->cursor_blink = val->boolean;
    break;
  case VTERM_PROP_CURSORSHAPE:
    u->cursor_shape = val->number;
    break;
  default:
    break;
  }

  return 1;
}

VTermScreenCallbacks cgo_vterm_screen_user_callbacks = {
    .movecursor = &cgo_vterm_screen_user_movecursor,
    .settermprop = &cgo_vterm_screen_user_settermprop,
};

static unsigned int cgo_vterm_screen_attrs_bold(VTermScreenCellAttrs attrs) {
  return attrs.bold;
}

static unsigned int
cgo_vterm_screen_attrs_underline(VTermScreenCellAttrs attrs) {
  return attrs.underline;
}

static unsigned int cgo_vterm_screen_attrs_italic(VTermScreenCellAttrs attrs) {
  return attrs.italic;
}

static unsigned int cgo_vterm_screen_attrs_blink(VTermScreenCellAttrs attrs) {
  return attrs.blink;
}

static unsigned int cgo_vterm_screen_attrs_reverse(VTermScreenCellAttrs attrs) {
  return attrs.reverse;
}

static unsigned int cgo_vterm_screen_attrs_conceal(VTermScreenCellAttrs attrs) {
  return attrs.conceal;
}

static unsigned int cgo_vterm_screen_attrs_strike(VTermScreenCellAttrs attrs) {
  return attrs.strike;
}

static unsigned int cgo_vterm_screen_attrs_font(VTermScreenCellAttrs attrs) {
  return attrs.font;
}

static unsigned int cgo_vterm_screen_attrs_dwl(VTermScreenCellAttrs attrs) {
  return attrs.dwl;
}

static unsigned int cgo_vterm_screen_attrs_dhl(VTermScreenCellAttrs attrs) {
  return attrs.dhl;
}

static unsigned int cgo_vterm_screen_attrs_small(VTermScreenCellAttrs attrs) {
  return attrs.small;
}

static unsigned int
cgo_vterm_screen_attrs_baseline(VTermScreenCellAttrs attrs) {
  return attrs.baseline;
}

#endif
