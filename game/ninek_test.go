package game_test

import (
	"999k_engine/constant"
	"999k_engine/game"
	"999k_engine/model"
	"testing"
)

func TestThreeOfAKind(t *testing.T) {
	t.Run("score=10000002", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			0, 1, 2,
		})
		score := scores[0] + scores[1]
		if score != 10000002 || kind != constant.ThreeOfAKind {
			t.Fail()
		}
	})
	t.Run("b>a", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			49, 50, 51,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			4, 5, 6,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a != 10000051 || b != 10000052 || a > b ||
			akind != constant.ThreeOfAKind ||
			bkind != constant.ThreeOfAKind {
			t.Fail()
		}
	})
	t.Run("a>b", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 1, 2,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			4, 17, 20,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a != 10000002 || bs[0] != 6 || a < b ||
			akind != constant.ThreeOfAKind ||
			bkind != constant.Nothing {
			t.Fail()
		}
	})
}

func TestStraightFlush(t *testing.T) {
	t.Run("8c,9c,10c is correct", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			24, 28, 32,
		})
		score := scores[0] + scores[1]
		if score != 1000032 || kind != constant.StraightFlush {
			t.Error()
		}
	})
	t.Run("5c,7c,8c is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			12, 20, 24,
		})
		if scores[0] != 1000 || kind != constant.Flush {
			t.Error()
		}
	})
	t.Run("Kc,Ad,2d is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			44, 49, 1,
		})
		if scores[0] != 3 || kind != constant.Nothing {
			t.Error()
		}
	})
	t.Run("Ac,2c,3c is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			48, 0, 4,
		})
		if scores[0] != 1000 || scores[1] != 48 || kind != constant.Flush {
			t.Error()
		}
	})
	t.Run("2c,3d,4h is wrong", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			0, 5, 10,
		})
		score := scores[0] + scores[1]
		if score != 10010 || kind != constant.Straight {
			t.Error()
		}
	})
	t.Run("2c,3c,4c < 6c,7c,8c", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			16, 20, 24,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a > b || a != 1000008 || b != 1000024 ||
			akind != constant.StraightFlush ||
			bkind != constant.StraightFlush {
			t.Error()
		}
	})
	t.Run("2c,3c,4c > 5c,7c,8c", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			12, 20, 24,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a < b ||
			akind != constant.StraightFlush ||
			bkind != constant.Flush {
			t.Error()
		}
	})
	t.Run("2c,3c,4c < 2d,3d,4d", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			0, 4, 8,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			1, 5, 9,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a > b ||
			akind != constant.StraightFlush ||
			bkind != constant.StraightFlush {
			t.Error()
		}
	})
}

func TestStraight(t *testing.T) {
	t.Run("2c,3d,4h is collect", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			0, 5, 10,
		})
		score := scores[0] + scores[1]
		if score != 10010 || kind != constant.Straight {
			t.Error()
		}
	})
	t.Run("Jc,Qs,Js", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			36, 43, 39,
		})
		score := scores[0] + scores[1]
		if score != 100043 ||
			kind != constant.Royal {
			t.Fail()
		}
	})
}

func TestRoyal(t *testing.T) {
	t.Run("Qc,Jd,Js is correct", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			40, 37, 39,
		})
		score := scores[0] + scores[1]
		if score != 100040 || kind != constant.Royal {
			t.Fail()
		}
	})
	t.Run("Jc,Jd,Js is not correct because it is three of a kind", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			36, 37, 38,
		})
		score := scores[0] + scores[1]
		if score == 100038 || kind != constant.ThreeOfAKind {
			t.Fail()
		}
	})
	t.Run("Jc,Qs,Js < Kc,Qc,Jh", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			36, 43, 39,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			44, 40, 38,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a > b ||
			akind != constant.Royal ||
			bkind != constant.Royal {
			t.Fail()
		}
	})
	t.Run("Kh,Qs,Js > Kc,Qc,Jh", func(t *testing.T) {
		as, akind := game.NineK{}.Evaluate(model.Cards{
			46, 43, 39,
		})
		bs, bkind := game.NineK{}.Evaluate(model.Cards{
			44, 40, 38,
		})
		a := as[0] + as[1]
		b := bs[0] + bs[1]
		if a < b ||
			akind != constant.Royal ||
			bkind != constant.Royal {
			t.Fail()
		}
	})
}

func TestSum(t *testing.T) {
	t.Run("Qc,10d,1s is nothing but has bonus", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			40, 33, 51,
		})
		if scores[0] != 1 ||
			scores[1] != 51 ||
			kind != constant.Nothing {
			t.Fail()
		}
	})
	t.Run("Jd,Qd,Ah is nothing", func(t *testing.T) {
		scores, kind := game.NineK{}.Evaluate(model.Cards{
			37, 41, 50,
		})
		if scores[0] != 1 ||
			scores[1] != 50 ||
			kind != constant.Nothing {
			t.Fail()
		}
	})
}

func TestFlush(t *testing.T) {
	t.Run("6,2,9 hearts must win 10s,2s,5d", func(t *testing.T) {
		a, akind := game.NineK{}.Evaluate(model.Cards{
			18, 2, 30,
		})
		b, bkind := game.NineK{}.Evaluate(model.Cards{
			4, 12, 16,
		})
		if a[0] != 1000 || akind != constant.Flush {
			t.Error()
		}
		if b[0] != 1000 || bkind != constant.Flush {
			t.Error()
		}
		if a[0]+a[1] < b[0]+b[1] {
			t.Error()
		}
	})
}
