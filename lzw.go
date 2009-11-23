package lzw

import "crush"

const (
  _LZW_EOD       = 257;
  _LZW_RESET     = 256;
  _LZW_DICSIZE   = 4096;
  _LZW_STARTBITS = 9;
  _LZW_STARTUTOK = 258;
)

type lzwDecoder struct {
  bits   *crush.BitT;
  bc, cp int;
  early  bool;
}

func (lzw *lzwDecoder) reset() {
  lzw.bc = _LZW_STARTBITS;
  lzw.cp = _LZW_STARTUTOK - 1;
}

func newLzwDecoder(s []byte, early bool) (lzw *lzwDecoder) {
  lzw = new(lzwDecoder);
  lzw.bits = crush.NewBits(s);
  lzw.early = early;
  lzw.reset();
  return;
}

func (lzw *lzwDecoder) update() bool {
  if lzw.cp < _LZW_DICSIZE-1 {
    if lzw.early {
      lzw.cp++
    }
    switch lzw.cp {
    case 512:
      lzw.bc = 10
    case 1024:
      lzw.bc = 11
    case 2048:
      lzw.bc = 12
    }
    if !lzw.early {
      lzw.cp++
    }
    return true;
  }
  return false;
}

func (lzw *lzwDecoder) token() (r int) {
  for {
    r = lzw.bits.Get(lzw.bc);
    if r != _LZW_RESET {
      break
    }
    lzw.reset();
  }
  return r;
}

func DecodeToSlice(s []byte, out []byte, early bool) (r int) {
  lzw := newLzwDecoder(s, early);
  dict := make([][]byte, _LZW_DICSIZE);
  for i := 0; i <= 255; i++ {
    dict[i] = []byte{byte(i)}
  }
  for c := lzw.token(); c != _LZW_EOD; c = lzw.token() {
    k := r;
    for i := 0; i < len(dict[c]); i++ {
      out[r] = dict[c][i];
      r++;
    }
    if lzw.update() {
      dict[lzw.cp] = out[k : r+1]
    }
  }
  return;
}

func CalculateLength(s []byte, early bool) (r int) {
  lzw := newLzwDecoder(s, early);
  dict := make([]int, _LZW_DICSIZE);
  for i := 0; i <= 255; i++ {
    dict[i] = 1
  }
  for c := lzw.token(); c != _LZW_EOD; c = lzw.token() {
    r += dict[c];
    if lzw.update() {
      dict[lzw.cp] = dict[c] + 1
    }
  }
  return;
}

func Decode(s []byte, early bool) []byte {
  r := make([]byte, CalculateLength(s, early)+1);
  return r[0:DecodeToSlice(s, r, early)];
}