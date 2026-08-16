package application

import (
	"strings"
	"testing"
)

// DIをしていない悪い例。NewDiceUseCase()は引数が0個で、Mockを渡す場所が
// そもそも存在しない。UseCase(ビジネスロジック)のテストのはずなのに、
// 直接本物のDBに繋いでテストするしかない。比較: good/application/dice_usecase_test.go

func TestDiceUseCase_IsBig_5の目はtrue(t *testing.T) {
	useCase := NewDiceUseCase() // ← 引数なし。DIできないのでMockを渡しようがない

	got, err := useCase.IsBig(1) // seed data: id=1, pips=5
	if err != nil {
		t.Skipf("postgres is required for this test (run `docker compose up -d` first): %v", err)
	}
	if !got {
		t.Errorf("IsBig() = false, want true")
	}
}

func TestDiceUseCase_IsBig_2の目はfalse(t *testing.T) {
	useCase := NewDiceUseCase()

	got, err := useCase.IsBig(2) // seed data: id=2, pips=2
	if err != nil {
		t.Skipf("postgres is required for this test (run `docker compose up -d` first): %v", err)
	}
	if got {
		t.Errorf("IsBig() = true, want false")
	}
}

func TestDiceUseCase_IsBig_存在しないIDはエラー(t *testing.T) {
	useCase := NewDiceUseCase()

	_, err := useCase.IsBig(9999) // does not exist
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	// A connection failure also returns a non-nil err, so "err == nil" alone
	// cannot tell "not found" apart from "postgres is unreachable". Without
	// this check, this test would falsely PASS with postgres stopped.
	if !strings.Contains(err.Error(), "not found") {
		t.Skipf("postgres is required for this test (run `docker compose up -d` first): %v", err)
	}
}
