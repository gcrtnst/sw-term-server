#ifndef __CGO_VTERM_COLOR_H__
#define __CGO_VTERM_COLOR_H__

#include <vterm.h>

static void cgo_vterm_color(VTermColor *col, uint8_t type, uint8_t red,
                            uint8_t green, uint8_t blue, uint8_t idx) {
  col->type = type;
  if (VTERM_COLOR_IS_RGB(col)) {
    col->rgb.red = red;
    col->rgb.green = green;
    col->rgb.blue = blue;
  } else if (VTERM_COLOR_IS_INDEXED(col)) {
    col->indexed.idx = idx;
  }
}

static uint8_t cgo_vterm_color_type(VTermColor *col) { return col->type; }

static uint8_t cgo_vterm_color_red(VTermColor *col) {
  if (!VTERM_COLOR_IS_RGB(col)) {
    return 0;
  }
  return col->rgb.red;
}

static uint8_t cgo_vterm_color_green(VTermColor *col) {
  if (!VTERM_COLOR_IS_RGB(col)) {
    return 0;
  }
  return col->rgb.green;
}

static uint8_t cgo_vterm_color_blue(VTermColor *col) {
  if (!VTERM_COLOR_IS_RGB(col)) {
    return 0;
  }
  return col->rgb.blue;
}

static uint8_t cgo_vterm_color_idx(VTermColor *col) {
  if (!VTERM_COLOR_IS_INDEXED(col)) {
    return 0;
  }
  return col->indexed.idx;
}

#endif
